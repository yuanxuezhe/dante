package register

import (
	"encoding/json"
	"gitee.com/yuanxuezhe/dante/conf"
	"gitee.com/yuanxuezhe/dante/log"
	"gitee.com/yuanxuezhe/dante/module/base"
	. "gitee.com/yuanxuezhe/dante/msg"
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

// TCP连接回调函数
func (m *BaseRegister) Handler(conn commconn.CommConn) {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		//if err := recover(); err != nil {
		//	if err.(error).Error() == "EOF" {
		//		return
		//	}
		//	if strings.Contains(err.(error).Error(), "An existing connection was forcibly closed by the remote host") {
		//		conn.WriteMsg(ResultPackege(m.ModuleType, 1, "connection was closed!["+conn.LocalAddr().String()+"==》"+conn.RemoteAddr().String()+"]", nil))
		//		return
		//	}
		//	if strings.Contains(err.(error).Error(), "use of closed network connection") {
		//		return
		//	}
		//	conn.WriteMsg(ResultPackege(m.ModuleType, 1, err.(error).Error(), nil))
		//}
		_ = conn.Close()
	}()

	//var err error
	//for {
	buff, err := conn.ReadMsg()
	if err != nil {
		panic(err)
	}

	if m.ModuleType != "Gateway" {
		log.Release("Params:%s", buff)
	}

	// 解析收到的消息
	msg := Msg{}
	err = json.Unmarshal(buff, &msg)
	if err != nil {
		panic(err)
	}

	// 若为注册消息，直接忽略
	if msg.Id != "Register" {
		return
	}

	_ = conn.WriteMsg(ResultPackege(msg.Id, msg.Id, 0, "注册成功！", nil))

	m.ReadChan <- buff
	//}
}

func (m *BaseRegister) CircleRegisterBeats() {
	for {
		_ = m.RegisterBeats()
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

		err = conn.WriteMsg(PackageMsg("RegisterList", string(jsons)))

		if err != nil {
			delete(m.Modules, k)
		}
	}

	return nil
}
