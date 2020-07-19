package gateway

import (
	. "dante/core/conf"
	"dante/core/log"
	base "dante/core/module/base"
	. "dante/core/msg"
	"encoding/json"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	web "gitee.com/yuanxuezhe/ynet/http"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	"strings"

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

	m.ConnMang = true
	m.Conns = make(map[string]commconn.CommConn, 1000000)
	m.ReadChan = make(chan []byte, 1000000)
	m.WriteChan = make(chan []byte, 1000000)
	m.Count = make(chan int, 100)

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

	for k := 0; k < 2; k++ {
		go m.DealReadChan()
	}

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

// TCP连接回调函数
func (m *Gate) Handler(conn commconn.CommConn) {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			if err.(error).Error() == "EOF" {
				return
			}
			if strings.Contains(err.(error).Error(), "use of closed network connection") {
				return
			}
			//fmt.Println(err) //这里的err其实就是panic传入的内容，bug
			//log.Error(err.(error).Error())
			conn.WriteMsg(ResultPackege(m.ModuleType, 1, err.(error).Error(), nil))
		}
		conn.Close()
	}()

	RemoteAddr := conn.RemoteAddr().String()
	// 保存 TCP
	m.Conns[RemoteAddr] = conn
	//var err error
	for {
		buff, err := conn.ReadMsg()
		if err != nil {
			panic(err)
		}
		if m.ModuleType != "Gateway" {
			log.Release("Params:%s", buff)
		}

		// 解析收到的消息
		msg := Msg{}
		json.Unmarshal(buff, &msg)
		if err != nil {
			panic(err)
		}

		// 若为注册消息，直接忽略
		if msg.Id == "Register" {
			conn.WriteMsg(ResultPackege(msg.Id, 0, "注册成功！", nil))
			continue
		}

		msg.Addr = RemoteAddr

		buff, err = json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		m.ReadChan <- buff
	}
}

func (m *Gate) DealReadChan() {
	//startT := time.Now()		//计算当前时间
	//Num := 0
	//var Speed float64
	for {
		select {
		case ri := <-m.ReadChan:
			//Num  = Num + 1
			m.Work(ri)
			//rs := time.Since(startT).Seconds()
			//if rs > 1 {
			//	Speed = float64(Num) / rs
			//	Num = 0
			//	startT = time.Now()
			//}
			//tc := time.Since(startT)
			//if tc.Seconds() < 1 {
			//
			//}
			//fmt.Printf("[%s]Speed = %v\n", m.ModuleId,Speed)
		}
	}
}

func (m *Gate) Work(msgs []byte) {
	// 解析收到的消息
	msg := Msg{}
	err := json.Unmarshal(msgs, &msg)
	if err != nil {
		panic(err)
	}

	var buff []byte

	buff, err = m.DoWork(msgs)

	if err != nil {
		buff = ResultPackege(m.ModuleType, 1, err.(error).Error(), nil)
	}
	m.WriteChan <- ResultIpPackege(msg.Addr, buff)
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
				//log.Release("[%8s][%s ==> %s] %s",m.ModuleId,conn.LocalAddr().String(),  conn.RemoteAddr().String(),string(res.Results))
				conn.WriteMsg(res.Results)
			}
		}
	}
}
