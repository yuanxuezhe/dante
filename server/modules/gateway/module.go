package gateway

import (
	"dante/core/module"
	basemodule "dante/core/module/base"
	"dante/core/module/gateway"
	. "dante/core/msg"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	"time"
)

var NewModule = func() module.Module {
	mod := &Gateway{
		gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3", Handler: Handler}},
	}
	return mod
}

type Gateway struct {
	gateway.Gate
}

func Handler(conn commconn.CommConn) {
	var Addr string
	var dconn commconn.CommConn
	var ok bool
	var err error
	var buff []byte
	for {
		buff, err = conn.ReadMsg()
		if err != nil {
			break
		}
		fmt.Println(string(buff))
		msg := &Msg{}
		err = json.Unmarshal(buff, msg)
		if err != nil {
			conn.WriteMsg([]byte("错误的数据包格式"))
			//panic("错误的数据包格式")
			continue
		}
		fmt.Println(msg)
		if msg.Id == "Login" {
			Addr, ok = getIP()
			if !ok {
				continue
			}
			dconn = ynet.NewTcpclient(Addr)
		} else {
			conn.WriteMsg([]byte("错误的接口"))
			continue
		}

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

func getIP() (string, bool) {
	//lock.Lock()
	//defer lock.Unlock()
	//if len(trueList) < 1 {
	//	return "", false
	//}
	//ip := trueList[0]
	//trueList = append(trueList[1:], ip)
	ip := "localhost:9201"
	return ip, true
}
