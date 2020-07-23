package gateway

import (
	"encoding/json"
	"errors"
	"gitee.com/yuanxuezhe/dante/core/log"
	"gitee.com/yuanxuezhe/dante/core/module"
	basemodule "gitee.com/yuanxuezhe/dante/core/module/base"
	"gitee.com/yuanxuezhe/dante/core/module/gateway"
	. "gitee.com/yuanxuezhe/dante/core/msg"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
)

var NewModule = func() module.Module {
	mod := &Gateway{
		Gate: gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3"}},
	}
	mod.Basemodule.DoWork = mod.DoWork
	return mod
}

type Gateway struct {
	gateway.Gate
}

func (g *Gateway) DoWork(buff []byte) ([]byte, error) {
	var dconn commconn.CommConn
	var err error

	module := "Error"

	msg := &Msg{}
	err = json.Unmarshal(buff, msg)
	if err != nil {
		return nil, errors.New("Error data format：" + err.Error())
	}
	module = msg.Id
	if module == "Heart" {
		return ResultPackege(module, 0, "Heart beats!", nil), nil
	}

	times := 0

reconnect:

	dconn, err = g.GetModuleConn(module)
	if err != nil {
		return nil, err
	}

	res, err := g.CallModule(dconn, buff)
	if err != nil {
		times = times + 1
		if times <= 10 {
			//delete(g.ModlueConns, Addr)
			log.Release("Reconnect %d times......", times)
			goto reconnect
		}
		return nil, err
	}

	return res, nil
}

func (g *Gateway) CallModule(dconn commconn.CommConn, body []byte) ([]byte, error) {
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
