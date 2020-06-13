package gateway

import (
	"dante/core/module"
	basemodule "dante/core/module/base"
	"dante/core/module/gateway"
)

var NewModule = func() module.Module {
	mod := &Gateway{
		gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3"}},
	}
	return mod
}

type Gateway struct {
	gateway.Gate
}
