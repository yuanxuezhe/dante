package main

import (
	"dante/core"
	_ "dante/core/conf"
	//"leafserver/src/server/conf"
	//"leafserver/src/server/game"
	//"leafserver/src/server/gate"
	//"leafserver/src/server/login"
)

func main() {
	//AddMod("Register", )
	core.Run()
}
