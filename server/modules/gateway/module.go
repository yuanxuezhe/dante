package gateway

import (
	"dante/core/log"
	"dante/core/module"
	basemodule "dante/core/module/base"
	"dante/core/module/gateway"
	. "dante/core/msg"
	"encoding/json"
	"errors"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
)

var NewModule = func() module.Module {
	mod := &Gateway{
		Gate: gateway.Gate{Basemodule: basemodule.Basemodule{ModuleType: "Gateway", ModuleVersion: "1.3.3"}},
	}
	mod.Basemodule.DoWork = mod.DoWork
	mod.modlueConns = make(map[string]commconn.CommConn, 100)
	return mod
}

type Gateway struct {
	gateway.Gate
	modlueConns map[string]commconn.CommConn
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

	Addr, err := g.getIP(module)
	if err != nil {
		return nil, err
	}

	dconn = g.getModuleConn(Addr)

	res, err := g.CallModule(dconn, buff)
	if err != nil {
		times = times + 1
		if times <= 10 {
			delete(g.modlueConns, Addr)
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

func (g *Gateway) getIP(moduletype string) (ip string, err error) {
	if moduletype == "Login" {
		ip = "192.168.0.3:9201"
	} else if moduletype == "Goods" {
		ip = "192.168.0.3:9301"
	} else {
		return "", errors.New("Undefined moudle:[" + moduletype + "]")
	}

	return ip, nil
}

func (g *Gateway) getModuleConn(ip string) (conn commconn.CommConn) {
	if conn, ok := g.modlueConns[ip]; ok {
		return conn
	} else {
		conns := ynet.NewTcpclient(ip)
		g.modlueConns[ip] = conns
		return g.modlueConns[ip]
	}
}

//
//func (g *Gateway) processMsg() {
//	for {
//		select {
//		//case <-g.ReadChan:
//		//	return
//		case msg := <-g.ReadChan:
//			fmt.Println("channel", msg)
//		}
//	}
//}
