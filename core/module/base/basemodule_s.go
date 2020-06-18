package base

import (
	. "dante/core/conf"
	"dante/core/log"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
)

// 发送注册信息
type ModuleInfo struct {
	ModuleId      string // 模块名称
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	TcpAddr       string
	Status        int
}

type Basemodule struct {
	ModuleId      string // 模块名称
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	registduring  int    // 注册心跳、断开时间间隔
	TcpAddr       string
	WsAddr        string
	conn          net.Conn
	registerflag  bool
	Handler       func(conn net.Conn) `json:"-"`
}

func (m *Basemodule) GetId() string {
	return m.ModuleId //+ "  " + m.TcpAddr + "  " + m.WsAddr
}

func (m *Basemodule) GetType() string {
	//Very important, it needs to correspond to the Module configuration in the configuration file
	return m.ModuleType
}
func (m *Basemodule) Version() string {
	//You can understand the code version during monitoring
	return m.ModuleVersion
}
func (m *Basemodule) OnInit() {

}

func (m *Basemodule) Run(closeSig chan bool) {
	tcpServer := ynet.NewTcpserver(m.TcpAddr, m.Handler)

	if tcpServer != nil {
		tcpServer.Start()
		log.Release("Module[%-10s|%-10s] start successful:[%s]", m.GetId(), m.Version(), m.TcpAddr)
	}

	<-closeSig

	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (m *Basemodule) OnDestroy() {

}

func (m *Basemodule) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	m.ModuleId = moduleSettings.Id

	if moduleSettings.Settings["TCPAddr"] != nil {
		if value, ok := moduleSettings.Settings["TCPAddr"].(string); ok {
			m.TcpAddr = value
		} else {
			err = fmt.Errorf("ModuleId:%s 参数[TCPAddr]设置有误", moduleSettings.Id)
			return
		}
	}

	if moduleSettings.Settings["WSAddr"] != nil {
		if value, ok := moduleSettings.Settings["WSAddr"].(string); ok {
			m.WsAddr = value
		} else {
			err = fmt.Errorf("ModuleId:%s 参数[TCPAddr]设置有误", moduleSettings.Id)
			return
		}
	}

	if value, ok := moduleSettings.Settings["Registduring"].(float64); ok {
		m.registduring = int(value)
	} else {
		err = fmt.Errorf("ModuleId:%s 参数[RegistBeatingduring]设置有误", moduleSettings.Id)
		return
	}

	m.registerflag = false
	// 注册标志存在，并且为true时，才发送注册消息
	if v, ok := moduleSettings.Settings["Register"].(bool); ok {
		if v == true {
			m.registerflag = true
		}
	}

	return
}

func (m *Basemodule) Register(closeSig chan bool) {
	// 注册标志存在，并且为true时，才发送注册消息
	if !m.registerflag {
		return
	}

	moduleInfo := &ModuleInfo{
		ModuleId:      m.ModuleId,
		ModuleType:    m.ModuleType,
		ModuleVersion: m.ModuleVersion,
		TcpAddr:       m.TcpAddr,
		Status:        0, // 0 表示注册
	}
	jsons, errs := json.Marshal(moduleInfo) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
		return
	}

	for {
		conn, err := net.Dial(Conf.RegisterProtocol, Conf.RegisterCentor)
		if err != nil {
			log.Error("Module[%-10s|%-10s] register failes", m.GetId(), m.Version())
			continue
		}

		// 发送注册消息
		err = network.SendMsg(conn, jsons)
		if err != nil {
			fmt.Printf("Module[%-10s|%-10s] register failes:%s", err)
			conn.Close()
			continue
		}

		// 接收注册中心应答
		buff, err := network.ReadMsg(conn.(net.Conn))
		if err != nil {
			fmt.Printf("%s", err)
			conn.Close()
			continue
		} else {
			fmt.Println("首次注册应答  ", conn.LocalAddr(), "==>", conn.RemoteAddr(), "    ", string(buff))
			conn.Close()
			break
		}
	}
}
