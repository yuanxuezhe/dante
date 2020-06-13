package core

import (
	. "dante/core/conf"
	"dante/core/log"
	"dante/core/module"
	"os"
	"os/signal"
	"strings"
)

func AddMod(tag string, newmodule func() module.Module) {
	module.AddModule(tag, newmodule)
}

func Run() {

	log.Release("Dante %v starting up", version)
	// 按配置注册模块
	mods := strings.Split(Conf.Registermodules, ",")
	for _, mod := range mods {
		module.Register(mod)
	}
	module.Init()

	module.RegisterCentor()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Dante closing down (signal: %v)", sig)
	module.Destroy()
}
