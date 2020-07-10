package login

import (
	"dante/core/module"
	base "dante/core/module/base"
	. "dante/core/msg"
	"dante/server/tables"
	"dante/server/util/snogenerator"
	"encoding/json"
	"fmt"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"strings"
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

func (m *Login) handler(conn commconn.CommConn) {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			if err.(error).Error() == "EOF" {
				return
			}
			if strings.Contains(err.(error).Error(), "use of closed network connection") {
				return
			}
			//fmt.Println(err) //这里的err其实就是panic传入的内容，bug
			//log.Error(err.(error).Error())
			conn.WriteMsg(m.ResultPackege(m.ModuleType, 1, err.(error).Error(), nil))
		}
		conn.Close()
	}()
	//var err error
	for {
		buff, err := conn.(*network.TCPConn).ReadMsg()
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
			conn.WriteMsg(m.ResultPackege("Register", 0, "注册成功！", nil))
			continue
		}

		// 解析获取登录信息
		loginInfo := Logininfo{}
		err = json.Unmarshal([]byte(msg.Body), &loginInfo)
		if err != nil {
			panic(err)
		}

		userinfo := tables.Userinfo{}

		userinfo.Phone = loginInfo.Phone
		userinfo.Email = loginInfo.Email
		userinfo.Passwd = loginInfo.Passwd

		err = m.CheckParams(loginInfo.Type, &userinfo)
		if err != nil {
			panic(err)
		}
		err = m.ManageUserinfo(loginInfo.Type, &userinfo)
		if err != nil {
			panic(err)
		}

		err = userinfo.QueryByKey()
		if err != nil {
			panic(err)
		}

		conn.WriteMsg(m.ResultPackege(m.ModuleType, 0, m.SetMsgSucc(loginInfo.Type), userinfo))
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
		err = userinfo.Insert()
		m.rw.Unlock()
	} else if Type == LOGIN_TYPE_LOGIN {
		userinfo, err = userinfo.CheckAccountExist()
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else if Type == LOGIN_TYPE_LOGOUT {

	}
	return nil
}

// Type 操作类型
func (m *Login) SetMsgSucc(Type int) (msg string) {
	if Type == LOGIN_TYPE_REGISTER {
		msg = " Register successful!"
	} else if Type == LOGIN_TYPE_LOGIN {
		msg = " Login successful!"
	} else if Type == LOGIN_TYPE_LOGOUT {
		msg = " Logout successful!"
	}
	return
}
