package main

import (
	"gitee.com/yuanxuezhe/dante/core"
	_ "gitee.com/yuanxuezhe/dante/core/conf"
	"gitee.com/yuanxuezhe/dante/server/modules/gateway"
	"gitee.com/yuanxuezhe/dante/server/modules/goods"
	"gitee.com/yuanxuezhe/dante/server/modules/login"
	"gitee.com/yuanxuezhe/dante/server/modules/register"
	_ "gitee.com/yuanxuezhe/dante/server/util/pool"
)

func main() {
	core.AddMod("Gateway", gateway.NewModule)
	core.AddMod("Register", register.NewModule)
	core.AddMod("Login", login.NewModule)
	core.AddMod("Goods", goods.NewModule)

	core.Run()
}
