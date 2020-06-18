package login

import (
	"dante/core/module"
	base "dante/core/module/base"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"time"
)

var NewModule = func() module.Module {
	mod := &Login{base.Basemodule{ModuleType: "Login", ModuleVersion: "1.2.4", Handler: Handler}}
	return mod
}

type Login struct {
	base.Basemodule
}

//
//func (m *Login) Run(closeSig chan bool) {
//	 tcpServer := NewTcpserver(m.TcpAddr, Handler)
//
//	if tcpServer != nil {
//		tcpServer.Start()
//		log.Release("Module[%-10s|%-10s] start successful:[%s]", m.GetId(), m.Version(), m.TcpAddr)
//	}
//
//	<-closeSig
//
//	if tcpServer != nil {
//		tcpServer.Close()
//	}
//}

func Handler(conn net.Conn) {
	for {
		buff, err := network.ReadMsg(conn)
		if err != nil {
			break
		}

		fmt.Println(string(buff))
		network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		time.Sleep(1 * time.Millisecond)
	}
}
