package login

import (
	"dante/core/module"
	base "dante/core/module/base"
	"dante/server/tables"
	"dante/server/util/snogenerator"
	"encoding/json"
	"errors"
	"sync"
)

const (
	LOGIN_TYPE_REGISTER = 0
	LOGIN_TYPE_LOGIN    = 1
	LOGIN_TYPE_LOGOUT   = 2
)

type Logininfo struct {
	Type    int    `json:"type"` // 登录类型 0、注册 1、登录 2、登出
	Account string `json:"account"`
	Userid  int    `json:"userid"` // userid
	Phone   int    `json:"phone"`  // phone
	Email   string `json:"email"`  // email
	Passwd  string `json:"passwd"` // 密码
}

var NewModule = func() module.Module {
	mod := &LoginManage{Basemodule: base.Basemodule{ModuleType: "Login", ModuleVersion: "1.2.4"}}
	mod.Basemodule.DoWork = mod.DoWork
	return mod
}

type LoginManage struct {
	base.Basemodule
	rw sync.RWMutex
}

func (m *LoginManage) DoWork(buff []byte) ([]byte, error) {
	// 解析获取登录信息
	var err error
	// 解析收到的消息
	loginInfo := Logininfo{}
	err = json.Unmarshal(buff, &loginInfo)
	if err != nil {
		panic(err)
	}

	userinfo := tables.Userinfo{}
	userinfo.Userid = loginInfo.Userid
	userinfo.Phone = loginInfo.Phone
	userinfo.Email = loginInfo.Email
	userinfo.Passwd = loginInfo.Passwd

	err = m.CheckParams(loginInfo.Type, &userinfo)
	if err != nil {
		return nil, err
	}
	err = m.ManageUserinfo(loginInfo.Type, &userinfo)
	if err != nil {
		return nil, err
	}

	err = userinfo.QueryByKey()
	if err != nil {
		return nil, err
	}

	return m.ResultPackege(m.ModuleType, 0, m.SetMsgSucc(loginInfo.Type), userinfo), nil
}

// Check params
func (m *LoginManage) CheckParams(Type int, userinfo *tables.Userinfo) error {
	var err error
	if Type == LOGIN_TYPE_REGISTER {
		err = userinfo.CheckAvailable_Phone()
		if err != nil {
			return err
		}

		if userinfo.Passwd == "" {
			return errors.New("Passwd cannot be null!")
		}
		//err = userinfo.CheckAvailable_Email()
		//if err != nil {
		//	return err
		//}
	} else if Type == LOGIN_TYPE_LOGIN {
		if userinfo.Userid == 0 && userinfo.Phone == 0 && userinfo.Email == "" {
			return errors.New("Account cannot be null!")
		}
		if userinfo.Passwd == "" {
			return errors.New("Passwd cannot be null!")
		}
	} else if Type == LOGIN_TYPE_LOGOUT {

	}
	return nil
}

func (m *LoginManage) ManageUserinfo(Type int, userinfo *tables.Userinfo) (err error) {
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
			return err
		}
	} else if Type == LOGIN_TYPE_LOGOUT {

	}
	return nil
}

// Type 操作类型
func (m *LoginManage) SetMsgSucc(Type int) (msg string) {
	if Type == LOGIN_TYPE_REGISTER {
		msg = " Register successful!"
	} else if Type == LOGIN_TYPE_LOGIN {
		msg = " Login successful!"
	} else if Type == LOGIN_TYPE_LOGOUT {
		msg = " Logout successful!"
	}
	return
}
