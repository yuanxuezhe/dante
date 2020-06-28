package gateway

import (
	"dante/core/module"
	basemodule "dante/core/module/base"
	"dante/core/module/gateway"
	. "dante/core/msg"
	"encoding/json"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"

	"net"
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

func Handler(conn net.Conn) {
	var Addr string
	var dconn net.Conn
	var ok bool
	var err error
	var buff []byte
	for {
		buff, err = network.ReadMsg(conn)
		if err != nil {
			break
		}

		msg := &Msg{}
		err = json.Unmarshal(buff, msg)
		if err != nil {
			network.SendMsg(conn, []byte("错误的数据包格式"))
			//panic("错误的数据包格式")
			continue
		}
		if msg.Id == "Login" {
			//fmt.Println("Recv msg: login : ", msg.Body)
			Addr, ok = getIP()
			if !ok {
				continue
			}
			dconn, err = net.Dial("tcp", Addr)
			if err != nil {
				fmt.Printf("连接%v失败:%v\n", Addr, err)
				return
			}
		} else {
			network.SendMsg(conn, []byte("错误的接口"))
			//panic("错误的接口")
			continue
		}

		CallModule(conn, dconn, msg.Body)

		time.Sleep(1 * time.Millisecond)
	}
}

func CallModule(conn, dconn net.Conn, body string) {
	defer dconn.Close()
	network.SendMsg(dconn, []byte(body))

	buff, _ := network.ReadMsg(dconn)
	err := network.SendMsg(conn, buff)
	if err != nil {
		fmt.Printf(err.Error())
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
