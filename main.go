package main

import (
	"dante/core"
	_ "dante/core/conf"
	"dante/server/modules/gateway"
	"dante/server/modules/login"
	"dante/server/modules/register"
	_ "dante/server/util/pool"
)

func main() {
	core.AddMod("Gateway", gateway.NewModule)
	core.AddMod("Register", register.NewModule)
	core.AddMod("Login", login.NewModule)

	core.Run()
}
