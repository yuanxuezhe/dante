package base

import (
	. "dante/core/conf"
	"dante/core/log"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	network "gitee.com/yuanxuezhe/ynet/tcp"
	"net"
	"time"
)

type Basemodule struct {
	ModuleId      string // 模块名称
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	registduring  int    // 注册心跳、断开时间间隔
	TcpAddr       string
	WsAddr        string
	conn          net.Conn
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

	return
}

func (m *Basemodule) Register(closeSig chan bool) {
	if m.ModuleType == "Register" || m.ModuleType == "Gateway" {
		return
	}

	conn, err := net.Dial(Conf.RegisterProtocol, Conf.RegisterCentor)
	if err != nil {
		log.Error("Module[%-10s|%-10s] register failes", m.GetId(), m.Version())
		return
	}
	m.conn = conn

	jsons, errs := json.Marshal(m) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
		return
	}

	for {
		// 发送注册消息
		err = network.SendMsg(m.conn, jsons)
		if err != nil {
			fmt.Printf("Module[%-10s|%-10s] register failes:%s", err)
			break
		}

		// 接收注册中心应答
		buff, err := network.ReadMsg(conn.(net.Conn))
		if err != nil {
			fmt.Printf("%s", err)
		}
		//fmt.Println(conn.(net.Conn).LocalAddr(), "==>", conn.(net.Conn).RemoteAddr(), "    ", string(buff))
		fmt.Println(m.conn.LocalAddr(), "==>", m.conn.RemoteAddr(), "    ", string(buff))

		time.Sleep(time.Duration(m.registduring) * time.Second)
	}
}

func Handler(conn net.Conn) {
	for {
		buff, err := network.ReadMsg(conn)
		if err != nil {
			break
		}

		fmt.Println(string(buff))
		network.SendMsg(conn, []byte("Hello,Recv msg:"+string(buff)))

		time.Sleep(1 * time.Millisecond)
	}
}
