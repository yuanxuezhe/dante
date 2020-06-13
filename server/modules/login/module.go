package login

import (
	"dante/core/module"
	base "dante/core/module/base"
)

var NewModule = func() module.Module {
	mod := &Login{base.Basemodule{ModuleType: "Login", ModuleVersion: "1.2.4"}}
	return mod
}

type Login struct {
	base.Basemodule
}
