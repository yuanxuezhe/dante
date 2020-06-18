package login

import (
	"dante/core/module"
	base "dante/core/module/base"
	. "dante/core/msg"
	"encoding/json"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"time"
)

var NewModule = func() module.Module {
	mod := &Login{base.Basemodule{ModuleType: "Login", ModuleVersion: "1.2.4"}}
	mod.Handler = mod.handler
	return mod
}

type Login struct {
	base.Basemodule
}

func (m *Login) handler(conn net.Conn) {
	for {
		buff, err := network.ReadMsg(conn)
		if err != nil {
			break
		}

		msg := &Msg{}
		json.Unmarshal(buff, msg)
		if msg.Id == "Register" {
			fmt.Println("Recv register heats    ", m.ModuleId, "       ", string(buff))
			continue
		}
		fmt.Println(string(buff))
		network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		time.Sleep(1 * time.Millisecond)
	}
}
