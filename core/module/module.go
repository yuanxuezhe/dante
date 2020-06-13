package module

import (
	. "dante/core/conf"
	"dante/core/log"
	"fmt"
	"sync"
)

type Module interface {
	GetId() string
	Version() string //module version
	GetType() string //module type
	OnInit()
	OnDestroy()
	SetPorperty(*ModuleSettings) error
	Run(closeSig chan bool)
}

type module struct {
	mi       Module
	closeSig chan bool
	wg       sync.WaitGroup
}

// all module collections
var mods []*module

// Create module function map
var mpmods map[string]func() Module

func init() {
	mpmods = make(map[string]func() Module)
}

func AddModule(tag string, newModule func() Module) {
	mpmods[tag] = newModule
}

func Register(mod string) {
	if mpmods[mod] == nil {
		log.Fatal("模块[%s]不存在", mod)
		return
	}
	if Conf.Module[mod] == nil {
		log.Fatal("模块[%s]配置信息不存在", mod)
		return
	}

	for _, moduleSettings := range Conf.Module[mod] {
		m := new(module)
		Model := mpmods[mod]()
		err := Model.SetPorperty(moduleSettings)
		if err != nil {
			fmt.Println(err)
			continue
		}

		m.mi = Model
		m.closeSig = make(chan bool, 1)

		mods = append(mods, m)
	}
}

func Init() {
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit()
	}

	for i := 0; i < len(mods); i++ {
		m := mods[i]
		m.wg.Add(1)
		go run(m)

		log.Release("%s:%s:%s", m.mi.GetId(), m.mi.GetType(), m.mi.Version())
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
	m.mi.OnDestroy()
}
