package network

import (
	"dante/core/log"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	sync.Mutex
	Addr            string
	ConnNum         int
	ConnectInterval time.Duration
	PendingWriteNum int
	AutoReconnect   bool
	NewAgent        func(*TCPConn) Agent
	Agent           *Agent
	conns           ConnSet
	//wg              sync.WaitGroup
	closeFlag bool

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	msgParser    *MsgParser
}

func (client *TCPClient) Start() {
	client.init()

	for i := 0; i < client.ConnNum; i++ {
		//client.wg.Add(1)
		go client.Connect()
	}
}

func (client *TCPClient) Start1() {
	client.init()

	//client.wg.Add(1)
	client.Connect()
}

func (client *TCPClient) init() {
	client.Lock()
	defer client.Unlock()

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		log.Release("invalid ConnNum, reset to %v", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		log.Release("invalid ConnectInterval, reset to %v", client.ConnectInterval)
	}

	if client.Agent == nil {
		log.Fatal("Agent must not be nil")
	}
	if client.conns != nil {
		log.Fatal("client is running")
	}

	client.conns = make(ConnSet)
	client.closeFlag = false

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(client.LenMsgLen, client.MinMsgLen, client.MaxMsgLen)
	msgParser.SetByteOrder(client.LittleEndian)
	client.msgParser = msgParser
}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.closeFlag {
			return conn
		}

		log.Release("connect to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) Connect() {
	//defer client.wg.Done()

	//reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return
	}
	//client.conns[conn] = struct{}{}
	//client.Unlock()

	tcpConn := newTCPConn(conn, client.PendingWriteNum, client.msgParser)

	//data := []byte(`{
	//        "YRequest": {
	//            "type": "register"
	//            "modid": "Login"
	//        }
	//    }`)

	//agent := client.NewAgent(tcpConn)
	//agent.Run()
	client.Agent.Conn = tcpConn

	// cleanup
	//tcpConn.Close()

	//client.Agent.Conn.WriteMsg(data)

	//client.Lock()
	//delete(client.conns, conn)
	//client.Unlock()
	//agent.OnClose()

	//if client.AutoReconnect {
	//	time.Sleep(client.ConnectInterval)
	//	goto reconnect
	//}
}

func (client *TCPClient) Close() {
	client.Lock()
	client.closeFlag = true
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

	//client.wg.Wait()
}
