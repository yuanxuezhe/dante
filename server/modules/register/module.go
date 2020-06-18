package register

import (
	"dante/core/log"
	"dante/core/module"
	base "dante/core/module/base"
	"dante/core/msg"
	"encoding/json"
	"fmt"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"time"
)

var NewModule = func() module.Module {
	mod := &Register{}
	mod.ModuleType = "Register"
	mod.ModuleVersion = "1.0.0"
	mod.Handler = mod.handler
	return mod
}

type ConnSet map[net.Conn]struct{}

type Register struct {
	base.Basemodule
	conns ConnSet
}

//var MapRegister map[string]base.Basemodule

func (m *Register) init() {
	m.conns = make(ConnSet)
	//MapRegister = make(map[string]base.Basemodule, 1000)
}

func (m *Register) handler(conn net.Conn) {
	for {
		// 接受首次注册消息
		buff, err := network.ReadMsg(conn)
		if err != nil {
			break
		}
		// 发送应答
		network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))
		// 解析消息体
		moduleInfo := base.ModuleInfo{}
		errs := json.Unmarshal(buff, &moduleInfo) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
			return
		}

		// 创建注册连接
		go m.CreateRegisterBeats(moduleInfo)
		//if _, ok := MapRegister[m.ModuleId]; !ok {
		//	MapRegister[m.ModuleId] = m
		//}

		//fmt.Println(MapRegister)

		time.Sleep(1 * time.Millisecond)
	}
}

func (m *Register) CreateRegisterBeats(moduleInfo base.ModuleInfo) {
	moduleInfos := &base.ModuleInfo{
		TcpAddr: m.TcpAddr,
		Status:  0, // 0 表示注册
	}

	jsons, err := json.Marshal(moduleInfos) //转换成JSON返回的是byte[]
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var conn net.Conn
reconnect:

	conn, err = net.Dial("tcp", moduleInfo.TcpAddr)
	if err != nil {
		log.Error("CreateRegisterBeats Module[%-10s|%-10s] register failes: %v  reconnect in 1 s", moduleInfo.ModuleId, moduleInfo.ModuleVersion, err)
		time.Sleep(1 * time.Second)
		goto reconnect
	}

	//go Read(conn)
	//go Write(conn)
	// 发送注册消息
resend:
	err = network.SendMsg(conn, msg.PackageMsg("Register", string(jsons)))
	if err != nil {
		fmt.Printf("Module[%-10s|%-10s] register sendmsg failes:%s", err)
		conn.Close()
		goto reconnect
	} else {
		time.Sleep(10 * time.Second)
		goto resend
	}

}

//
//func  Read(conn net.Conn) {
//	reader := bufio.NewReader(c.TCPConn)
//	for {
//		lineBytes, err := reader.ReadBytes('\n')
//		if err != nil {
//			log.Println("startread read bytes error ", err)
//			break
//		}
//		len:=len(lineBytes)
//		line:=lineBytes[:len-1]
//		log.Println("start read from client ",string(line))
//		go c.HandleMsg(line)
//	}
//}
//func  Write(conn net.Conn) {
//	log.Println("write groutine is waiting")
//	defer log.Println("write groutine exit")
//	for {
//		select {
//		case data := <-c.MsgChan:
//			if _, err := c.TCPConn.Write(data); err != nil {
//				log.Println("startwrite conn write error ", err)
//				return
//			}
//			log.Println("start write from server ",string(data))
//		case <-c.ExitChan:
//			return
//		}
//	}
//}
