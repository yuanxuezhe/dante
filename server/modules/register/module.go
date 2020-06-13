package register

import (
	"dante/core/module"
	base "dante/core/module/Basemodule"
)

var NewModule = func() module.Module {
	mod := &Register{base.Basemodule{ModuleType: "Register", ModuleVersion: "1.0.0"}}
	return mod
}

type Register struct {
	base.Basemodule
}
