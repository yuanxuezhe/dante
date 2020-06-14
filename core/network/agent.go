//package network
//
//type Agent123123 interface {
//	Run()
//	OnClose()
//}

package network

import (
	"dante/core/log"
	"dante/core/module"
	"net"
	"reflect"
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

type Agent struct {
	Conn Conn
	Mod  *module.Module
	//UserData interface{}
}

func (a *Agent) Run() {
	for {
		data, err := a.Conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		log.Release("recive msg: %s", data)
	}
}

func (a *Agent) OnClose() {
}

func (a *Agent) WriteMsg(msg interface{}) {
	//if a.gate.Processor != nil {
	//	data, err := a.gate.Processor.Marshal(msg)
	//	if err != nil {
	//		log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
	//		return
	//	}
	err := a.Conn.WriteMsg(msg.([]byte))

	if err != nil {
		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}
}

func (a *Agent) LocalAddr() net.Addr {
	return a.Conn.LocalAddr()
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.Conn.RemoteAddr()
}

func (a *Agent) Close() {
	a.Conn.Close()
}

func (a *Agent) Destroy() {
	a.Conn.Destroy()
}

//func (a *Agent) UserData() interface{} {
//	return a.userData
//}

//func (a *Agent) SetUserData(data interface{}) {
//	a.userData = data
//}
