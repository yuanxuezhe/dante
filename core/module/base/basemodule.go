package base

import (
	. "dante/core/conf"
	"dante/core/log"
	"dante/core/network"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Basemodule struct {
	ModuleId      string // 模块名称
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	registduring  int    // 注册心跳、断开时间间隔
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
			agent := network.Agent{Conn: conn}
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
		tcpServer.LenMsgLen = 4
		tcpServer.MaxMsgLen = 1000000
		tcpServer.LittleEndian = false
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			agent := network.Agent{Conn: conn}
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
		m.registduring = int(value)
	} else {
		err = fmt.Errorf("ModuleId:%s 参数[RegistBeatingduring]设置有误", moduleSettings.Id)
		return
	}

	return
}

func (m *Basemodule) Register(closeSig chan bool) {
	agent := &network.Agent{}
	//o := new(network.TCPConn)
	//a.conn = o
	if strings.ToUpper(Conf.RegisterProtocol) == "TCP" {
		tcpClient := new(network.TCPClient)
		tcpClient.Addr = Conf.RegisterCentor
		tcpClient.ConnectInterval = time.Duration(m.registduring) * time.Second
		tcpClient.PendingWriteNum = 100
		tcpClient.LenMsgLen = 4
		tcpClient.MinMsgLen = 0
		tcpClient.MaxMsgLen = 100000
		tcpClient.LittleEndian = false

		tcpClient.Agent = agent

		tcpClient.Start1()
	} else if strings.ToUpper(Conf.RegisterProtocol) == "HTTP" {
		wsClient := new(network.WSClient)
		wsClient.Addr = Conf.RegisterCentor
		wsClient.ConnectInterval = time.Duration(m.registduring) * time.Second

		//wsClient.Agent = agent
		wsClient.Connect()
	}

	jsons, errs := json.Marshal(m) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}

	agent.WriteMsg(jsons)
	<-closeSig

	agent.Conn.Close()
}
