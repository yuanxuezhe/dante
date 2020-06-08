/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package login

import (
	"dante/core/conf"
	"dante/core/log"
	"dante/core/module"
	"dante/core/network"
	"fmt"
	"math/rand"
	"time"
)

var Module = func() module.Module {
	mod := new(Login)
	return mod
}

type Login struct {
	module.Module
	TCPAddr string
}

func (m *Login) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Login"
}
func (m *Login) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *Login) OnInit() {

}

func (m *Login) Run(closeSig chan bool) {
}

func (m *Login) OnDestroy() {

}
