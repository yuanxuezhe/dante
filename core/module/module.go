package module

import (
	"dante/core/conf"
	"dante/core/log"
	"runtime"
	"sync"
)

type Module interface {
	Version() string //模块版本
	GetType() string //模块类型
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
}

type module struct {
	mi       Module
	closeSig chan bool
	wg       sync.WaitGroup
}

var mods []*module
var mpmods map[string]Module

func init() {
	mpmods = make(map[string]Module)
}

func AddModule(tag string, mi Module) {
	mpmods[tag] = mi
}

func Register(mod string) {
	if mpmods[mod] == nil {
		log.Fatal("模块[%s]不存在", mod)
		return
	}
	m := new(module)
	m.mi = mpmods[mod]
	m.closeSig = make(chan bool, 1)

	mods = append(mods, m)
}

func Init() {
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit()
	}

	for i := 0; i < len(mods); i++ {
		m := mods[i]
		m.wg.Add(1)
		go run(m)
	}
}

func Destroy() {
	for i := len(mods) - 1; i >= 0; i-- {
		m := mods[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func run(m *module) {
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *module) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	m.mi.OnDestroy()
}
