package base

import (
	. "dante/core/conf"
	"dante/core/log"
	"dante/core/network"
	"fmt"
	"strings"
	"time"
)

type Basemodule struct {
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	Registduring  int    // 注册心跳、断开时间间隔
	ModuleId      string // 模块名称
	TcpAddr       string
	WsAddr        string
}

func (m *Basemodule) GetId() string {
	return m.ModuleId //+ "  " + m.TcpAddr + "  " + m.WsAddr
}

func (m *Basemodule) GetType() string {
	//Very important, it needs to correspond to the Module configuration in the configuration file
	return m.ModuleType
}
func (m *Basemodule) Version() string {
	//You can understand the code version during monitoring
	return m.ModuleVersion
}
func (m *Basemodule) OnInit() {

}

func (m *Basemodule) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if m.WsAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = m.WsAddr
		//wsServer.MaxConnNum = gate.MaxConnNum
		//wsServer.PendingWriteNum = gate.PendingWriteNum
		//wsServer.MaxMsgLen = gate.MaxMsgLen
		//wsServer.HTTPTimeout = gate.HTTPTimeout
		//wsServer.CertFile = gate.CertFile
		//wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			agent := &agent{conn: conn, mod: m}
			//if gate.AgentChanRPC != nil {
			//	gate.AgentChanRPC.Go("NewAgent", a)
			//}
			return agent
		}
	}
	var tcpServer *network.TCPServer
	if m.TcpAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = m.TcpAddr
		//tcpServer.MaxConnNum = gate.MaxConnNum
		//tcpServer.PendingWriteNum = gate.PendingWriteNum
		//tcpServer.LenMsgLen = gate.LenMsgLen
		//tcpServer.MaxMsgLen = gate.MaxMsgLen
		//tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			agent := &agent{conn: conn, mod: m}
			return agent
		}
	}

	var info string

	if tcpServer != nil {
		tcpServer.Start()
		info += " TcpAddr:" + m.TcpAddr
	}

	if wsServer != nil {
		wsServer.Start()
		info += " WsAddr:" + m.WsAddr
	}

	log.Release("Module[%-10s|%-10s] start successful :%s", m.GetId(), m.Version(), info)

	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (m *Basemodule) OnDestroy() {

}

func (m *Basemodule) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	m.ModuleId = moduleSettings.Id

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

	if value, ok := moduleSettings.Settings["Registduring"].(float64); ok {
		m.Registduring = int(value)
	} else {
		err = fmt.Errorf("ModuleId:%s 参数[RegistBeatingduring]设置有误", moduleSettings.Id)
		return
	}

	return
}

func (m *Basemodule) Register(closeSig chan bool) {
	if strings.ToUpper(Conf.RegisterProtocol) == "TCP" {
		tcpClient := new(network.TCPClient)
		tcpClient.Addr = Conf.RegisterCentor
		tcpClient.ConnectInterval = time.Duration(m.Registduring) * time.Second
		tcpClient.NewAgent = func(conn *network.TCPConn) network.Agent {
			agent := &agent{conn: conn, mod: m}
			return agent
		}
	} else if strings.ToUpper(Conf.RegisterProtocol) == "HTTP" {
		tcpClient := new(network.WSClient)
		tcpClient.Addr = Conf.RegisterCentor
		tcpClient.ConnectInterval = time.Duration(m.Registduring) * time.Second
		tcpClient.NewAgent = func(conn *network.WSConn) network.Agent {
			agent := &agent{conn: conn, mod: m}
			return agent
		}
	}

	client := &network.Client{
		Protocol:      Conf.RegisterProtocol,
		Addr:          Conf.RegisterCentor,
		During:        time.Duration(m.Registduring) * time.Second,
		AutoReconnect: true,
		NewAgent: func(conn *network.Conn) network.Agent {
			agent := &agent{
				conn: *conn,
				mod:  m,
			}

			return agent
		},
	}

	client.Init()

	client.Connect()

	client.Sendmsg("1111123232323")
	<-closeSig

	client.Close()
}
