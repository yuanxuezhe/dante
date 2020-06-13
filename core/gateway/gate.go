package gateway

import (
	. "dante/core/conf"
	"dante/core/log"
	base "dante/core/module/Basemodule"
	"dante/core/network"
	"fmt"
	"net"
	"reflect"
	"time"
)

type Gate struct {
	base.Basemodule
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	//AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	LenMsgLen    int
	LittleEndian bool

	/////////////
	//TCPAddr string				// 监听地址
	//ModuleType string			// 模块类型
	//ModuleVersion string		// 模块版本号
	//Registduring float64		// 注册心跳、断开时间间隔
	//ModuleId string				// 模块名称
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			agent := &agent{
				conn: conn,
				gate: gate,
			}

			return agent
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

func (gate *Gate) GetId() string {
	return gate.ModuleId + "  " + gate.TCPAddr
}

func (gate *Gate) GetType() string {
	//Very important, it needs to correspond to the Module configuration in the configuration file
	return gate.ModuleType
}
func (gate *Gate) Version() string {
	//You can understand the code version during monitoring
	return gate.ModuleVersion
}

func (gate *Gate) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	gate.ModuleId = moduleSettings.Id

	if value, ok := moduleSettings.Settings["TCPAddr"].(string); ok {
		gate.TCPAddr = value
	} else {
		err = fmt.Errorf("ModuleId:%s 参数[TCPAddr]设置有误", moduleSettings.Id)
		return
	}

	return
}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		log.Release("recive msg: %s", data)
	}
}

func (a *agent) OnClose() {
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
