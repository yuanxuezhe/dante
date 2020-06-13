package basemodule

import (
	. "dante/core/conf"
	"fmt"
)

type Basemodule struct {
	ModuleType    string  // 模块类型
	ModuleVersion string  // 模块版本号
	Registduring  float64 // 注册心跳、断开时间间隔
	ModuleId      string  // 模块名称
	TCPAddr       string  // 监听地址

}

func (m *Basemodule) GetId() string {
	return m.ModuleId + "  " + m.TCPAddr
}

func (m *Basemodule) GetType() string {
	//Very important, it needs to correspond to the Module configuration in the configuration file
	return m.ModuleType
}
func (m *Basemodule) Version() string {
	//You can understand the code version during monitoring
	return m.ModuleVersion
}
func (m *Basemodule) OnInit() {

}

func (m *Basemodule) Run(closeSig chan bool) {
	//if m.ModuleType == "Register" {
	//
	//}

}

func (m *Basemodule) OnDestroy() {

}

func (m *Basemodule) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	m.ModuleId = moduleSettings.Id

	if value, ok := moduleSettings.Settings["TCPAddr"].(string); ok {
		m.TCPAddr = value
	} else {
		err = fmt.Errorf("ModuleId:%s 参数[TCPAddr]设置有误", moduleSettings.Id)
		return
	}

	if value, ok := moduleSettings.Settings["Registduring"].(float64); ok {
		m.Registduring = value
	} else {
		err = fmt.Errorf("ModuleId:%s 参数[RegistBeatingduring]设置有误", moduleSettings.Id)
		return
	}

	return
}
