package gateway

import (
	"dante/core/module"
	basemodule "dante/core/module/base"
	"dante/core/module/gateway"
	"net"
)

var NewModule = func() module.Module {
	mod := &Gateway{
		gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3", Handler: Handler}},
	}
	return mod
}

type Gateway struct {
	gateway.Gate
}

func Handler(conn net.Conn) {

}
