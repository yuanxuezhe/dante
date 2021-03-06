package register

import (
	"encoding/json"
	"gitee.com/yuanxuezhe/dante/comm"
	"gitee.com/yuanxuezhe/dante/conf"
	"gitee.com/yuanxuezhe/dante/log"
	"gitee.com/yuanxuezhe/dante/module/base"
	. "gitee.com/yuanxuezhe/dante/msg"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	web "gitee.com/yuanxuezhe/ynet/websocket"
	//"reflect"
	"sync"
	"time"
)

type BaseRegister struct {
	base.Basemodule
	registerConnsRWMutex sync.RWMutex
	//registerConns map[string]commconn.CommConn
	registerConns comm.DMap
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

// SetProperty 设置模块参数
func (m *BaseRegister) SetProperty(moduleSettings *conf.ModuleSettings) (err error) {
	//m.registerCons = make(map[string]common.CommConn, 50)
	m.registerConns.SetCapacity(50)
	m.registerConns.SetDescribe(m.ModuleId + "  " + "m.registerConns  ")
	err = m.Basemodule.SetProperty(moduleSettings)
	if err != nil {
		return err
	}
	return nil
}

// Handler TCP连接回调函数
func (m *BaseRegister) Handler(conn commconn.CommConn) {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		//if err := recover(); err != nil {
		//	if err.(error).Error() == "EOF" {
		//		return
		//	}
		//	if strings.Contains(err.(error).Error(), "An existing connection was forcibly closed by the remote host") {
		//		conn.WriteMsg(ResultPackage(m.ModuleType, 1, "connection was closed!["+conn.LocalAddr().String()+"==》"+conn.RemoteAddr().String()+"]", nil))
		//		return
		//	}
		//	if strings.Contains(err.(error).Error(), "use of closed network connection") {
		//		return
		//	}
		//	conn.WriteMsg(ResultPackage(m.ModuleType, 1, err.(error).Error(), nil))
		//}
		conn.Close()
	}()

	//var err error
	//for {
	buff, err := conn.ReadMsg()
	if err != nil {
		panic(err)
	}

	if m.ModuleType != "Gateway" {
		log.LogPrint(log.LEVEL_RELEASE, "[%-10s]Params:%s", m.ModuleId, buff)
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
	jsons, err := m.Modules.GetJsonFromMap()
	if err != nil {
		return err
	}

	m.Modules.Range(func(key string, values interface{}) bool {
		var moduleInfos = values.([]base.ModuleInfo)

		for k, value := range values.([]base.ModuleInfo) {
			conn, err := m.registerConns.Get(value.TcpAddr)
			if err != nil {
				conn = ynet.NewTcpclient(value.TcpAddr)
				m.registerConns.Set(value.TcpAddr, conn)
			}

			err = conn.(commconn.CommConn).WriteMsg(PackageMsg("RegisterList", string(jsons)))

			if err != nil {
				values = append(moduleInfos[:k], moduleInfos[k+1:]...)
			}
		}
		//m.Modules.Set(key, moduleInfos)
		return true
	})

	return nil
}
