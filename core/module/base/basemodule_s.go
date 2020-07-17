package base

import (
	. "dante/core/conf"
	"dante/core/log"
	. "dante/core/msg"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	web "gitee.com/yuanxuezhe/ynet/http"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	"strings"
	"time"

	//network "gitee.com/yuanxuezhe/ynet/tcp"
	_ "github.com/go-sql-driver/mysql"
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
	//conn          net.Conn
	registerflag bool
	DoWork       func([]byte) ([]byte, error) `json:"-"`
	ConnMang     bool
	Conns        map[commconn.CommConn]string
	ReadChan     chan []byte
	WriteChan    chan []byte
	Count        int
}

// 取模块ID
func (m *Basemodule) GetId() string {
	return m.ModuleId
}

// 取模块类型
func (m *Basemodule) GetType() string {
	//Very important, it needs to correspond to the Module configuration in the configuration file
	return m.ModuleType
}

// 取模块版本
func (m *Basemodule) Version() string {
	//You can understand the code version during monitoring
	return m.ModuleVersion
}

// 模块初始化
func (m *Basemodule) OnInit() {

}

// 运行模块
func (m *Basemodule) Run(closeSig chan bool) {
	var tcpServer *tcp.TCPServer
	var wsServer *web.WSServer

	// tcp
	if len(m.TcpAddr) > 0 {
		tcpServer = &tcp.TCPServer{
			Addr:            m.TcpAddr,
			MaxConnNum:      1000000,
			PendingWriteNum: 1000,
			Callback:        m.Handler,
		}
	}

	// web
	if len(m.WsAddr) > 0 {
		wsServer = &web.WSServer{
			Addr:            m.WsAddr,
			MaxConnNum:      1000000,
			PendingWriteNum: 1000,
			HTTPTimeout:     5 * time.Second,
			Callback:        m.Handler,
		}
	}

	go m.DealReadChan()

	if tcpServer != nil {
		tcpServer.Start()
		log.Release("Module[%-10s|%-10s] start tcpServer successful:[%s]", m.GetId(), m.Version(), m.TcpAddr)
	}

	if wsServer != nil {
		wsServer.Start()
		log.Release("Module[%-10s|%-10s] start wsServer successful:[%s]", m.GetId(), m.Version(), m.WsAddr)
	}

	// 关闭系统
	<-closeSig

	if tcpServer != nil {
		tcpServer.Close()
	}

	if wsServer != nil {
		wsServer.Close()
	}
}

// 关闭
func (m *Basemodule) OnDestroy() {

}

// 设置模块参数
func (m *Basemodule) SetPorperty(moduleSettings *ModuleSettings) (err error) {
	//m.init()
	m.ModuleId = moduleSettings.Id
	m.Count = 0
	m.ConnMang = false
	m.ReadChan = make(chan []byte, 2)
	m.WriteChan = make(chan []byte, 2)

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

// 注册模块到注册中心
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
		conn := ynet.NewTcpclient(Conf.RegisterCentor)

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

// TCP连接回调函数
func (m *Basemodule) Handler(conn commconn.CommConn) {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			if err.(error).Error() == "EOF" {
				return
			}
			if strings.Contains(err.(error).Error(), "use of closed network connection") {
				return
			}
			//fmt.Println(err) //这里的err其实就是panic传入的内容，bug
			//log.Error(err.(error).Error())
			conn.WriteMsg(ResultPackege(m.ModuleType, 1, err.(error).Error(), nil))
		}
		conn.Close()
	}()

	//var err error
	for {
		buff, err := conn.ReadMsg()
		if err != nil {
			panic(err)
		}
		// 解析收到的消息
		msg := Msg{}
		json.Unmarshal(buff, &msg)
		if err != nil {
			panic(err)
		}

		// 若为注册消息，直接忽略
		if msg.Id == "Register" {
			conn.WriteMsg(ResultPackege(msg.Id, 0, "注册成功！", nil))
			continue
		}

		m.ReadChan <- buff
		m.Count = m.Count + 1
		fmt.Println(m.Count)
	}
}

func (m *Basemodule) DealReadChan() {
	for {
		fmt.Println("len:", len(m.ReadChan))
		select {
		case ri := <-m.ReadChan:
			fmt.Println("ri", m.ModuleType, string(ri))
			buff, err := m.DoWork(ri)
			if err != nil {
				m.WriteChan <- []byte(ResultPackege(m.ModuleType, 1, err.(error).Error(), nil))
			} else {
				m.WriteChan <- buff
			}
		}
	}
}
