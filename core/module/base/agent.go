package base

import (
	"dante/core/log"
	"dante/core/network"
	"net"
)

//type Agent interface {
//	WriteMsg(msg interface{})
//	LocalAddr() net.Addr
//	RemoteAddr() net.Addr
//	Close()
//	Destroy()
//	UserData() interface{}
//	SetUserData(data interface{})
//}

type agent struct {
	conn     network.Conn
	mod      *Basemodule
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
	//if a.gate.Processor != nil {
	//	data, err := a.gate.Processor.Marshal(msg)
	//	if err != nil {
	//		log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
	//		return
	//	}
	//	err = a.conn.WriteMsg(data...)
	//	if err != nil {
	//		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	//	}
	//}
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
