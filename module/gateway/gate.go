package gateway

import (
	"encoding/json"

	"gitee.com/yuanxuezhe/dante/conf"
	"gitee.com/yuanxuezhe/dante/log"
	"gitee.com/yuanxuezhe/dante/module/base"
	. "gitee.com/yuanxuezhe/dante/msg"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	web "gitee.com/yuanxuezhe/ynet/websocket"

	//"dante/core/network"
	"fmt"
	"time"
)

type Gate struct {
	base.Basemodule
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32

	// websocket
	HTTPTimeout time.Duration

	CertFile string
	KeyFile  string
}

func (m *Gate) SetProperty(moduleSettings *conf.ModuleSettings) (err error) {
	m.ModuleId = moduleSettings.Id
	m.ReadChan = make(chan []byte, 1000000)
	m.WriteChan = make(chan []byte, 1000000)
	m.Modules.SetCapacity(100)
	m.Conns.SetCapacity(1000000)
	m.ModlueConns.SetCapacity(100)

	m.Modules.SetDescribe(m.ModuleId + "  " + "m.Modules  ")
	m.Conns.SetDescribe(m.ModuleId + "  " + "m.Conns  ")
	m.ModlueConns.SetDescribe(m.ModuleId + "  " + "m.ModlueConns  ")

	if moduleSettings.Settings["TCPAddr"] != nil {
		if value, ok := moduleSettings.Settings["TCPAddr"].(string); ok {
			m.TcpAddr = value
		} else {
			err = fmt.Errorf("ModuleId:%s 参数[TCPAddr]设置有误", moduleSettings.Id)
			return
		}
	}

	if moduleSettings.Settings["WSAddr"] != nil {
		if value, ok := moduleSettings.Settings["WSAddr"].(string); ok {
			m.WsAddr = value
		} else {
			err = fmt.Errorf("ModuleId:%s 参数[TCPAddr]设置有误", moduleSettings.Id)
			return
		}
	}

	m.Registerflag = false
	// 注册标志存在，并且为true时，才发送注册消息
	if v, ok := moduleSettings.Settings["Register"].(bool); ok {
		if v == true {
			m.Registerflag = true
		}
	}

	return
}

// 运行模块
func (m *Gate) Run(closeSig chan bool) {
	var tcpServer *tcp.TCPServer
	var wsServer *web.WSServer

	// tcp
	if len(m.TcpAddr) > 0 {
		tcpServer = &tcp.TCPServer{
			Addr:            m.TcpAddr,
			MaxConnNum:      1000000,
			PendingWriteNum: 1000,
			Callback:        m.Handler,
		}
	}

	// web
	if len(m.WsAddr) > 0 {
		wsServer = &web.WSServer{
			Addr:            m.WsAddr,
			MaxConnNum:      1000000,
			PendingWriteNum: 1000,
			HTTPTimeout:     5 * time.Second,
			Callback:        m.Handler,
		}
	}

	go m.DealReadChan()

	go m.DealWriteChan()

	if tcpServer != nil {
		tcpServer.Start()
		log.LogPrint(log.LEVEL_RELEASE, "Module[%-10s|%-10s] start tcpServer successful:[%s]", m.GetId(), m.Version(), m.TcpAddr)
	}

	if wsServer != nil {
		wsServer.Start()
		log.LogPrint(log.LEVEL_RELEASE, "Module[%-10s|%-10s] start wsServer successful:[%s]", m.GetId(), m.Version(), m.WsAddr)
	}

	// 关闭系统
	<-closeSig

	if tcpServer != nil {
		tcpServer.Close()
	}

	if wsServer != nil {
		wsServer.Close()
	}
}

func (m *Gate) DealReadChan() {
	for {
		select {
		case ri := <-m.ReadChan:
			m.Work(ri)
		}
	}
}

func (m *Gate) Work(msgs []byte) {
	// 解析收到的消息
	dataBuf := make(map[string]interface{})
	err := json.Unmarshal(msgs, &dataBuf)
	if err != nil {
		panic(err)
	}

	Addr, _ := dataBuf["addr"].(string)

	var buff []byte

	buff, err = m.DoWork(msgs)

	if err != nil {
		buff = ResultPackege(m.ModuleType, m.ModuleId, 1, err.(error).Error(), nil)
	}
	m.WriteChan <- ResultIpPackege(Addr, buff)
}

func (m *Gate) DealWriteChan() {
	for {
		select {
		case ri := <-m.WriteChan:
			res := ResultWithIp{}
			err := json.Unmarshal(ri, &res)
			if err != nil {
				continue
			}

			conn, err := m.Conns.Get(res.Ip)
			if err != nil {
				continue
			}

			log.LogPrint(log.LEVEL_DEBUG, "[%-10s][%s ==> %s] %s", m.ModuleId, conn.(commconn.CommConn).LocalAddr().String(), conn.(commconn.CommConn).RemoteAddr().String(), string(res.Results))

			conn.(commconn.CommConn).WriteMsg(res.Results)
		}
	}
}
