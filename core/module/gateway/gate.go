package gateway

import (
	. "dante/core/conf"
	base "dante/core/module/base"
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
	m.ReadChan = make(chan []byte, 2)
	m.WriteChan = make(chan []byte, 2)

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
