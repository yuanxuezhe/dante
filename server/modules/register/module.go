package register

import (
	"encoding/json"
	"gitee.com/yuanxuezhe/dante/core/module"
	"gitee.com/yuanxuezhe/dante/core/module/base"
	"gitee.com/yuanxuezhe/dante/core/module/register"
	"net"
)

var NewModule = func() module.Module {
	mod := &Register{}
	mod.ModuleType = "Register"
	mod.ModuleVersion = "1.0.0"
	mod.Basemodule.DoWork = mod.DoWork
	return mod
}

type ConnSet map[net.Conn]struct{}

type Register struct {
	register.BaseRegister
}

//var MapRegister map[string]base.Basemodule

func (m *Register) init() {
	//m.Modules = make(map[string]base.ModuleInfo, 50)
	m.Modules = make(map[string]base.ModuleInfo, 50)
	//MapRegister = make(map[string]base.Basemodule, 1000)
}

func (m *Register) DoWork(buff []byte) ([]byte, error) {
	var err error
	// 解析消息体
	moduleInfo := base.ModuleInfo{}
	err = json.Unmarshal(buff, &moduleInfo) //转换成JSON返回的是byte[]
	if err != nil {
		return nil, err
	}

	m.Modules[moduleInfo.ModuleId] = moduleInfo

	// 创建注册连接
	m.RegisterBeats()

	return []byte("Register successful:" + string(buff)), nil
}
