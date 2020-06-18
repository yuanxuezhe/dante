package register

import (
	"dante/core/module"
	base "dante/core/module/base"
	"encoding/json"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"time"
)

var NewModule = func() module.Module {
	mod := &Register{base.Basemodule{ModuleType: "Register", ModuleVersion: "1.0.0", Handler: Handler}}

	return mod
}

type Register struct {
	base.Basemodule
}

var MapRegister map[string]base.Basemodule

func init() {
	MapRegister = make(map[string]base.Basemodule, 1000)
}

func Handler(conn net.Conn) {
	for {
		buff, err := network.ReadMsg(conn)
		if err != nil {
			break
		}

		m := base.Basemodule{}
		errs := json.Unmarshal(buff, &m) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
			return
		}

		MapRegister[m.ModuleType] = m

		fmt.Println(MapRegister)
		network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		time.Sleep(1 * time.Millisecond)
	}
}
