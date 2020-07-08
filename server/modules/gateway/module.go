package gateway

import (
	"dante/core/module"
	basemodule "dante/core/module/base"
	"dante/core/module/gateway"
	. "dante/core/msg"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	"strings"
	"time"
)

var NewModule = func() module.Module {
	mod := &Gateway{
		gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3"}},
	}
	mod.Basemodule.Handler = mod.Handler
	return mod
}

type Gateway struct {
	gateway.Gate
}

func (g *Gateway) Handler(conn commconn.CommConn) {
	var Addr string
	var dconn commconn.CommConn
	var err error
	var buff []byte

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
			conn.WriteMsg(g.ResultPackege(1, err.(error).Error(), nil))
		}
		conn.Close()
	}()

	for {
		buff, err = conn.ReadMsg()
		if err != nil {
			break
		}
		fmt.Println(string(buff))
		msg := &Msg{}
		err = json.Unmarshal(buff, msg)
		if err != nil {
			panic(errors.New("网关数据包格式有误：" + err.Error()))
		}
		fmt.Println(msg)

		Addr, err = getIP(msg.Id)
		if err != nil {
			panic(err)
		}

		dconn = ynet.NewTcpclient(Addr)

		CallModule(conn, dconn, buff)

		time.Sleep(1 * time.Millisecond)
	}
}

func CallModule(conn, dconn commconn.CommConn, body []byte) {
	defer dconn.Close()

	err := dconn.WriteMsg(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	buff, err := dconn.ReadMsg()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = conn.WriteMsg(buff)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getIP(moduletype string) (ip string, err error) {
	//lock.Lock()
	//defer lock.Unlock()
	//if len(trueList) < 1 {
	//	return "", false
	//}
	//ip := trueList[0]
	//trueList = append(trueList[1:], ip)
	if moduletype == "Login" {
		ip = "localhost:9201"
	} else {
		return "", errors.New("未定义的模块类型[" + moduletype + "]")
	}

	return ip, nil
}
