package register

import (
	"dante/core/log"
	"dante/core/module"
	base "dante/core/module/base"
	"dante/core/msg"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	"net"
	"time"
)

var NewModule = func() module.Module {
	mod := &Register{}
	mod.ModuleType = "Register"
	mod.ModuleVersion = "1.0.0"
	mod.Basemodule.DoWork = mod.DoWork
	return mod
}

type ConnSet map[net.Conn]struct{}

type Register struct {
	base.Basemodule
	modules map[string]base.ModuleInfo
}

//var MapRegister map[string]base.Basemodule

func (m *Register) init() {
	m.modules = make(map[string]base.ModuleInfo)
	//MapRegister = make(map[string]base.Basemodule, 1000)
}

func (m *Register) DoWork(buff []byte) ([]byte, error) {
	var err error
	// 解析消息体
	moduleInfo := base.ModuleInfo{}
	err = json.Unmarshal(buff, &moduleInfo) //转换成JSON返回的是byte[]
	if err != nil {
		return nil, err
	}

	// 创建注册连接
	go m.CreateRegisterBeats(moduleInfo)
	return []byte("Register successful:" + string(buff)), nil
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

	var conn commconn.CommConn
reconnect:

	//conn, err = net.Dial("tcp", moduleInfo.TcpAddr)
	conn = ynet.NewTcpclient(moduleInfo.TcpAddr)
	if err != nil {
		//if _, ok := m.modules[moduleInfo.ModuleId]; ok {
		//	delete(m.modules,moduleInfo.ModuleId)
		//}
		log.Error("CreateRegisterBeats Module[%-10s|%-10s] register failes: %v  reconnect in 1 s", moduleInfo.ModuleId, moduleInfo.ModuleVersion, err)
		time.Sleep(1 * time.Second)
		goto reconnect
	}

	//go Read(conn)
	//go Write(conn)
	// 发送注册消息
resend:

	err = conn.WriteMsg(msg.PackageMsg("Register", string(jsons)))
	if err != nil {
		//if _, ok := m.modules[moduleInfo.ModuleId]; ok {
		//	delete(m.modules,moduleInfo.ModuleId)
		//}
		fmt.Printf("Module[%-10s|%-10s] register sendmsg failes:%s", err)
		conn.Close()
		goto reconnect
	} else {
		//if _, ok := m.modules[moduleInfo.ModuleId]; !ok {
		//	m.modules[moduleInfo.ModuleId] = moduleInfo
		//}
		time.Sleep(time.Duration(m.Registduring) * time.Second)
		goto resend
	}
}
