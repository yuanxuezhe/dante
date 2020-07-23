package gateway

import (
	"encoding/json"
	. "gitee.com/yuanxuezhe/dante/core/conf"
	"gitee.com/yuanxuezhe/dante/core/log"
	base "gitee.com/yuanxuezhe/dante/core/module/base"
	. "gitee.com/yuanxuezhe/dante/core/msg"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	web "gitee.com/yuanxuezhe/ynet/http"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	//"dante/core/network"
	"fmt"
	"time"
)

type Gate struct {
	base.Basemodule
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	//Processor       network.Processor
	//AgentChanRPC    *chanrpc.Server

	// websocket
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	LenMsgLen    int
	LittleEndian bool
}

func (m *Gate) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	m.ModuleId = moduleSettings.Id

	m.Conns = make(map[string]commconn.CommConn, 1000000)
	m.ReadChan = make(chan []byte, 1000000)
	m.WriteChan = make(chan []byte, 1000000)
	m.ModlueConns = make(map[string]commconn.CommConn, 100)

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
		log.Release("Module[%-10s|%-10s] start tcpServer successful:[%s]", m.GetId(), m.Version(), m.TcpAddr)
	}

	if wsServer != nil {
		wsServer.Start()
		log.Release("Module[%-10s|%-10s] start wsServer successful:[%s]", m.GetId(), m.Version(), m.WsAddr)
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

//// TCP连接回调函数
//func (m *Gate) Handler(conn commconn.CommConn) {
//	// 将客户端远程地址作为连接的key，保存TCP连接，供返回值调用
//	RemoteAddr := conn.RemoteAddr().String()
//	m.Conns[RemoteAddr] = conn
//	defer func() { //必须要先声明defer，否则不能捕获到panic异常
//		if err := recover(); err != nil {
//			if err.(error).Error() == "EOF" {
//				return
//			}
//			if strings.Contains(err.(error).Error(), "use of closed network connection") {
//				return
//			}
//			m.WriteChan <- ResultIpPackege(RemoteAddr, ResultPackege(m.ModuleType, 1, err.(error).Error(), nil))
//		}
//		conn.Close()
//	}()
//
//	for {
//		buff, err := conn.ReadMsg()
//		if err != nil {
//			panic(err)
//		}
//		// 解析收到的消息
//		msg := Msg{}
//		json.Unmarshal(buff, &msg)
//		if err != nil {
//			panic(err)
//		}
//
//		msg.Addr = RemoteAddr
//
//		buff, err = json.Marshal(msg)
//		if err != nil {
//			panic(err)
//		}
//
//		// 若为注册消息，直接忽略
//		if msg.Id == "Register" {
//			conn.WriteMsg(ResultPackege(msg.Id, 0, "注册成功！", nil))
//			//continue
//		}
//		if msg.Id == "RegisterList" {
//			json.Unmarshal([]byte(msg.Body), &m.Modules)
//			fmt.Println(m.ModuleId," Register info:",m.Modules)
//			continue
//		}
//
//		m.ReadChan <- buff
//	}
//}

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
		buff = ResultPackege(m.ModuleType, 1, err.(error).Error(), nil)
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

			if conn, ok := m.Conns[res.Ip]; ok {
				log.Release("[%8s][%s ==> %s] %s", m.ModuleId, conn.LocalAddr().String(), conn.RemoteAddr().String(), string(res.Results))
				conn.WriteMsg(res.Results)
			}
		}
	}
}
