package login

import (
	"dante/core/module"
	base "dante/core/module/base"
	. "dante/core/msg"
	"dante/server/tables"
	"encoding/json"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"time"
)

const (
	LOGIN_TYPE_REGISTER = 0
	LOGIN_TYPE_LOGIN    = 1
	LOGIN_TYPE_LOGOUT   = 2
)

type Logininfo struct {
	Type    int    `json:"type"`    // 登录类型 0、注册 1、登录 2、登出
	Userid  string `json:"userid"`  // 用户名
	Account string `json:"account"` // 账号 userid/phone num/email
	Phone   int    `json:"phone"`   // 手机号码
	Email   string `json:"email"`   // 邮箱
	Passwd  string `json:"passwd"`  // 密码
}

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

		// 解析收到的消息
		msg := &Msg{}
		json.Unmarshal(buff, msg)
		// 若为注册消息，直接忽略
		if msg.Id == "Register" {
			continue
		}

		// 解析获取登录信息
		loginInfo := Logininfo{}
		err = json.Unmarshal(buff, &loginInfo)
		if err != nil {
			break
		}

		userinfo := tables.Userinfo{}

		if loginInfo.Type == LOGIN_TYPE_REGISTER {
			fmt.Println("LOGIN_TYPE_REGISTER")

		} else if loginInfo.Type == LOGIN_TYPE_LOGIN {
			fmt.Println("LOGIN_TYPE_LOGIN")
		} else if loginInfo.Type == LOGIN_TYPE_LOGOUT {
			fmt.Println("LOGIN_TYPE_LOGOUT")
		}

		userinfo.Userid = "1001"
		//userinfo.Insert()
		userinfo.QueryByKey()
		fmt.Println(userinfo)
		network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		time.Sleep(1 * time.Millisecond)
	}
}
