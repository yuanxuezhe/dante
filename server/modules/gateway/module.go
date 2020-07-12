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
)

var NewModule = func() module.Module {
	mod := &Gateway{
		gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3"}},
	}
	mod.Basemodule.DoWork = mod.DoWork
	return mod
}

type Gateway struct {
	gateway.Gate
}

func (g *Gateway) DoWork(buff []byte) ([]byte, error) {
	var Addr string
	var dconn commconn.CommConn
	var err error

	module := "Error"
	fmt.Println("GGGGGGGGGGGGGG:", string(buff))
	msg := &Msg{}
	err = json.Unmarshal(buff, msg)
	if err != nil {
		panic(errors.New("Error data format：" + err.Error()))
	}

	module = msg.Id
	if module == "Heart" {
		return g.ResultPackege(module, 0, "Heart beats!", nil), nil
	}

	Addr, err = getIP(module)
	if err != nil {
		panic(err)
	}

	dconn = ynet.NewTcpclient(Addr)

	return CallModule(dconn, []byte(msg.Body))
}

func CallModule(dconn commconn.CommConn, body []byte) ([]byte, error) {
	defer dconn.Close()

	err := dconn.WriteMsg(body)
	if err != nil {
		return nil, err
	}
	buff, err := dconn.ReadMsg()
	if err != nil {
		return nil, err
	}

	return buff, err
}

func getIP(moduletype string) (ip string, err error) {
	if moduletype == "Login" {
		ip = "localhost:9201"
	} else if moduletype == "Goods" {
		ip = "localhost:9301"
	} else {
		return "", errors.New("Undefined moudle:[" + moduletype + "]")
	}

	return ip, nil
}
