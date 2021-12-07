package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	. "gitee.com/yuanxuezhe/dante/conf"
	"gitee.com/yuanxuezhe/dante/log"
	. "gitee.com/yuanxuezhe/dante/msg"
	"gitee.com/yuanxuezhe/ynet"
	commconn "gitee.com/yuanxuezhe/ynet/Conn"
	tcp "gitee.com/yuanxuezhe/ynet/tcp"
	web "gitee.com/yuanxuezhe/ynet/websocket"

	_ "github.com/go-sql-driver/mysql"
)

// 发送注册信息
type ModuleInfo struct {
	ModuleId      string // 模块名称
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	TcpAddr       string // 注册地址
	//Status        int
}

// 模块基类
type Basemodule struct {
	sync.RWMutex
	ModuleId      string // 模块名称
	ModuleType    string // 模块类型
	ModuleVersion string // 模块版本号
	Registduring  int    // 注册心跳、断开时间间隔
	TcpAddr       string // TCP连接地址
	WsAddr        string // WEB连接地址

	Registerflag bool                         // 注册标志
	DoWork       func([]byte) ([]byte, error) // 回调函数
	ReadChan     chan []byte                  // 读入队列
	WriteChan    chan []byte                  // 写出队列
	Modules      map[string]ModuleInfo        // 记录注册信息
	Conns        map[string]commconn.CommConn // 客户端连接
	ModlueConns  map[string]commconn.CommConn // 记录模块连接
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

	go m.DealWriteChan()

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
	m.ModuleId = moduleSettings.Id
	m.ReadChan = make(chan []byte, 20)
	m.WriteChan = make(chan []byte, 20)
	m.Conns = make(map[string]commconn.CommConn, 500)
	m.Modules = make(map[string]ModuleInfo, 50)
	m.ModlueConns = make(map[string]commconn.CommConn, 100)

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

	m.Registerflag = false
	// 注册标志存在，并且为true时，才发送注册消息
	if v, ok := moduleSettings.Settings["Register"].(bool); ok {
		if v == true {
			m.Registerflag = true
		}
	}

	return
}

// 注册模块到注册中心
func (m *Basemodule) Register(closeSig chan bool) {
	// 注册标志存在，并且为true时，才发送注册消息
	if !m.Registerflag {
		return
	}

	moduleInfo := &ModuleInfo{
		ModuleId:      m.ModuleId,
		ModuleType:    m.ModuleType,
		ModuleVersion: m.ModuleVersion,
		TcpAddr:       m.TcpAddr,
		//Status:        0, // 0 表示注册
	}

	jsons, err := json.Marshal(moduleInfo) //转换成JSON返回的是byte[]
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsons = PackageMsg("Register", string(jsons))

	for {
		conn := ynet.NewTcpclient(Conf.RegisterCentor)

		// 发送注册消息
		err = conn.WriteMsg(jsons)
		if err != nil {
			fmt.Printf("Module[%-10s|%-10s] register failes:%s", err)
			_ = conn.Close()
			continue
		}

		// 接收注册中心应答
		buff, err := conn.ReadMsg()
		if err != nil {
			fmt.Printf("%s", err)
			conn.Close()
			continue
		} else {
			log.Release("[%-10s]%s", m.ModuleId, "首次注册应答  "+conn.LocalAddr().String()+"==>", conn.RemoteAddr().String()+"    "+string(buff))
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
			if strings.Contains(err.(error).Error(), "An existing connection was forcibly closed by the remote host") {
				conn.WriteMsg(ResultPackege(m.ModuleType, m.ModuleId, 1, "connection was closed!["+conn.LocalAddr().String()+"==》"+conn.RemoteAddr().String()+"]", nil))
				m.Lock()
				delete(m.Conns, conn.RemoteAddr().String())
				m.Unlock()
				return
			}
			if strings.Contains(err.(error).Error(), "use of closed network connection") {
				return
			}
			conn.WriteMsg(ResultPackege(m.ModuleType, m.ModuleId, 1, err.(error).Error(), nil))
		}
		conn.Close()
	}()

	RemoteAddr := conn.RemoteAddr().String()
	// 保存 TCP
	m.Lock()
	m.Conns[RemoteAddr] = conn
	m.Unlock()
	//var err error
	for {
		buff, err := conn.ReadMsg()
		if err != nil {
			panic(err)
		}

		if m.ModuleType != "Gateway" {
			log.Release("[%-10s]Params:%s", m.ModuleId, buff)
		}

		// 解析收到的消息
		msg := Msg{}
		json.Unmarshal(buff, &msg)
		if err != nil {
			panic(err)
		}

		// Register list refresh
		if msg.Id == "RegisterList" {
			json.Unmarshal([]byte(msg.Body), &m.Modules)
			log.Release("[%-10s]%s", m.ModuleId, "Register list refresh successful!")
			continue
		}

		msg.Addr = RemoteAddr

		buff, err = json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		m.ReadChan <- buff
	}
}

func (m *Basemodule) DealReadChan() {
	for {
		select {
		case ri := <-m.ReadChan:
			go m.Work(ri)
		}
	}
}

func (m *Basemodule) Work(msgs []byte) {
	// 解析收到的消息
	msg := Msg{}
	err := json.Unmarshal(msgs, &msg)
	if err != nil {
		panic(err)
	}

	var buff []byte

	buff, err = m.DoWork([]byte(msg.Body))

	if err != nil {
		buff = ResultPackege(m.ModuleType, m.ModuleId, 1, err.(error).Error(), nil)
	}
	m.WriteChan <- ResultIpPackege(msg.Addr, buff)
}

func (m *Basemodule) DealWriteChan() {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		m.DealWriteChan()
	}()
	for {
		select {
		case ri := <-m.WriteChan:
			res := ResultWithIp{}
			err := json.Unmarshal(ri, &res)
			if err != nil {
				continue
			}

			if conn, ok := m.Conns[res.Ip]; ok {
				conn.WriteMsg(res.Results)
			}
		}
	}
}

func (m *Basemodule) GetModuleConn(moduletype string) (commconn.CommConn, error) {
	var ip string
	for _, value := range m.Modules {
		if value.ModuleType == moduletype {
			ip = value.TcpAddr
			break
		}
	}

	if len(ip) == 0 {
		return nil, errors.New("Undefined module:" + moduletype)
	}
	if conn, ok := m.ModlueConns[ip]; ok {
		return conn, nil
	} else {
		conns := ynet.NewTcpclient(ip)
		m.ModlueConns[ip] = conns
		return m.ModlueConns[ip], nil
	}
}
