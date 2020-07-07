package base

import (
	. "dante/core/conf"
	"dante/core/log"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	web "gitee.com/yuanxuezhe/ynet/http"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"

	//network "gitee.com/yuanxuezhe/ynet/tcp"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"reflect"
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
	Registduring  int    // 注册心跳、断开时间间隔
	TcpAddr       string
	WsAddr        string
	conn          net.Conn
	registerflag  bool
	Handler       func(conn commconn.CommConn) `json:"-"`
	//Mysqlpool     *yconnpool.ConnPool
}

type Result struct {
	Status string // 状态
	Code   int    // 错误码
	Msg    string // 消息
	Data   string // 结果
}

//func (m *Basemodule) init() {
//	var err error
//	m.Mysqlpool, err = yconnpool.NewConnPool(func() (yconnpool.ConnRes, error) {
//		return sql.Open("mysql", "root:1@tcp(192.168.3.25:3306)/dante?parseTime=true")
//	}, 100, time.Second*100)
//	if err != nil {
//		panic(err)
//	}
//}

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
	var tcpServer *tcp.TCPServer
	var wsServer *web.WSServer
	if len(m.TcpAddr) > 0 {
		tcpServer = ynet.NewTcpserver(m.TcpAddr, m.Handler)
	}
	if len(m.WsAddr) > 0 {
		wsServer = ynet.NewWsserver(m.WsAddr, m.Handler)
	}

	//wsServer := ynet.NewTcpserver(m.TcpAddr, m.Handler)
	if tcpServer != nil {
		tcpServer.Start()
		log.Release("Module[%-10s|%-10s] start tcpServer successful:[%s]", m.GetId(), m.Version(), m.TcpAddr)
	}

	if wsServer != nil {
		wsServer.Start()
		log.Release("Module[%-10s|%-10s] start wsServer successful:[%s]", m.GetId(), m.Version(), m.WsAddr)
	}

	<-closeSig

	if tcpServer != nil {
		tcpServer.Close()
	}

	if wsServer != nil {
		wsServer.Close()
	}
}

func (m *Basemodule) OnDestroy() {

}

func (m *Basemodule) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	//m.init()
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
		m.Registduring = int(value)
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
	jsons, err := json.Marshal(moduleInfo) //转换成JSON返回的是byte[]
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		//conn, err := net.Dial(Conf.RegisterProtocol, Conf.RegisterCentor)
		conn := ynet.NewTcpclient(Conf.RegisterCentor)
		//if err != nil {
		//	log.Error("Module[%-10s|%-10s] register failes", m.GetId(), m.Version())
		//	continue
		//}

		// 发送注册消息
		err = conn.WriteMsg(jsons)
		if err != nil {
			fmt.Printf("Module[%-10s|%-10s] register failes:%s", err)
			conn.Close()
			continue
		}

		// 接收注册中心应答
		buff, err := conn.ReadMsg()
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

func (m *Basemodule) ResultPackege(code int, msg string, data interface{}) []byte {
	result := &Result{}
	if code == 0 {
		result.Status = "ok"
	} else {
		result.Status = "err"
	}

	result.Code = code

	result.Msg = msg

	data_type := reflect.TypeOf(data)
	if data_type != nil {
		if data_type.Kind().String() == "struct" {
			buff, _ := json.Marshal(data)
			result.Data = string(buff)
		}
	}

	buff, _ := json.Marshal(result)

	if result.Status == "ok" {
		log.Release(string(buff))
	} else {
		log.Error(string(buff))
	}

	return buff
}
