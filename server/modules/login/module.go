package login

import (
	"dante/core/module"
	base "dante/core/module/base"
	. "dante/core/msg"
	"dante/server/tables"
	"dante/server/util/snogenerator"
	"encoding/json"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"sync"
	"time"
)

const (
	LOGIN_TYPE_REGISTER = 0
	LOGIN_TYPE_LOGIN    = 1
	LOGIN_TYPE_LOGOUT   = 2
)

type Logininfo struct {
	Type    int    `json:"type"`    // 登录类型 0、注册 1、登录 2、登出
	Account string `json:"account"` // 账号 userid/phone num/email
	Phone   int    `json:"phone"`   // 手机号码
	Email   string `json:"email"`   // 邮箱
	Passwd  string `json:"passwd"`  // 密码
}

var NewModule = func() module.Module {
	mod := &Login{Basemodule: base.Basemodule{ModuleType: "Login", ModuleVersion: "1.2.4"}}
	mod.Handler = mod.handler
	return mod
}

type Login struct {
	base.Basemodule
	rw sync.RWMutex
}

func (m *Login) handler(conn net.Conn) {
	//defer conn.Close()

	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) //这里的err其实就是panic传入的内容，bug
			network.SendMsg(conn, m.ResultPackege(1, err.(error).Error()))
		}
		conn.Close()
	}()
	//var err error
	for {
		buff, err := network.ReadMsg(conn)
		if err != nil {
			panic(err)
		}
		// 解析收到的消息
		msg := &Msg{}
		json.Unmarshal(buff, msg)

		if err != nil {
			panic(err)
		}

		// 若为注册消息，直接忽略
		if msg.Id == "Register" {
			fmt.Println("注册消息不处理")
			continue
		}

		// 解析获取登录信息
		loginInfo := Logininfo{}
		err = json.Unmarshal(buff, &loginInfo)
		if err != nil {
			panic(err)
		}

		userinfo := tables.Userinfo{}

		userinfo.Phone = loginInfo.Phone
		userinfo.Email = loginInfo.Email
		userinfo.Passwd = loginInfo.Passwd
		//SetParam()
		// Ckeck params
		err = m.CheckParams(loginInfo.Type, &userinfo)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		err = m.ManageUserinfo(loginInfo.Type, &userinfo)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		userinfo.QueryByKey()
		fmt.Println(userinfo)

		network.SendMsg(conn, m.ResultPackege(0, userinfo))
		time.Sleep(1 * time.Millisecond)
	}
}

// Check params
func (m *Login) CheckParams(Type int, userinfo *tables.Userinfo) error {
	var err error
	if Type == LOGIN_TYPE_REGISTER {
		err = userinfo.CheckAvailable_Phone()
		if err != nil {
			return err
		}
	} else if Type == LOGIN_TYPE_LOGIN {

	} else if Type == LOGIN_TYPE_LOGOUT {

	}
	return nil
}

func (m *Login) ManageUserinfo(Type int, userinfo *tables.Userinfo) (err error) {
	if Type == LOGIN_TYPE_REGISTER {
		m.rw.Lock()
		userinfo.Userid = snogenerator.NewUserid()
		// 用户编号从1111开始
		if userinfo.Userid < 11111 {
			userinfo.Userid = 11111
		}
		userinfo.Insert()
		m.rw.Unlock()
	} else if Type == LOGIN_TYPE_LOGIN {
		userinfo, err = userinfo.CheckAccountExist()
		if err != nil {
			return err
		}

		fmt.Printf("Login successful : %v \n", userinfo)
	} else if Type == LOGIN_TYPE_LOGOUT {

	}
	return nil
}
