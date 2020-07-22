package register

import (
	"dante/core/conf"
	"dante/core/log"
	base "dante/core/module/base"
	"dante/core/msg"
	"encoding/json"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	web "gitee.com/yuanxuezhe/ynet/http"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	"time"
)

type BaseRegister struct {
	base.Basemodule
	registerConns map[string]commconn.CommConn
}

// 运行模块
func (m *BaseRegister) Run(closeSig chan bool) {
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

	go m.CircleRegisterBeats()

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

// 设置模块参数
func (m *BaseRegister) SetPorperty(moduleSettings *conf.ModuleSettings) (err error) {
	m.registerConns = make(map[string]commconn.CommConn, 50)
	err = m.Basemodule.SetPorperty(moduleSettings)
	if err != nil {
		return err
	}
	return nil
}

func (m *BaseRegister) CircleRegisterBeats() {
	for {
		m.RegisterBeats()
		time.Sleep(time.Duration(m.Registduring) * time.Second)
	}
}

func (m *BaseRegister) RegisterBeats() error {
	jsons, err := json.Marshal(m.Modules) //转换成JSON返回的是byte[]
	if err != nil {
		return err
	}

	var conn commconn.CommConn

	for k, value := range m.Modules {
		if c, ok := m.registerConns[value.TcpAddr]; ok {
			conn = c
		} else {
			conn = ynet.NewTcpclient(value.TcpAddr)
			m.registerConns[value.TcpAddr] = conn
		}

		err = conn.WriteMsg(msg.PackageMsg("RegisterList", string(jsons)))

		if err != nil {
			delete(m.Modules, k)
		}
	}

	return nil
}
