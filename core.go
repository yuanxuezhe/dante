package dante

import (
	"danteserver/server/pool"
	"fmt"
	. "gitee.com/yuanxuezhe/dante/conf"
	"gitee.com/yuanxuezhe/dante/log"
	"gitee.com/yuanxuezhe/dante/module"
	"gitee.com/yuanxuezhe/dante/public"
	logs "log"
	"os"
	"os/signal"
	"strings"
)

func AddMod(tag string, newmodule func() module.Module) {
	module.AddModule(tag, newmodule)
}

func Run() {
	defaultLogPath := fmt.Sprintf("%s/%s", public.ApplicationRoot, Conf.Log["LogPath"].(string))
	//fmt.Println(defaultLogPath)
	// 定义日志配置
	if Conf.Log["LogLevel"].(string) != "" {
		logger, err := log.New(Conf.Log["LogLevel"].(string), Conf.Log["PrintLevel"].(string), defaultLogPath, logs.Ldate|logs.Lmicroseconds, Conf.Log["Console"].(bool))
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	// 初始化连接池
	pool.InitConnpoll()

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
