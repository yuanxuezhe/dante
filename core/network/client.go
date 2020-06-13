//package network
//
//import (
//	"fmt"
//	"strings"
//	"time"
//)
//
package network

import (
	//"dante/core/log"
	"fmt"
	//"net"
	"strings"
	//"sync"
	"time"
)

var tcpClient *TCPClient
var wsClient *WSClient

type Client struct {
	Protocol      string
	Addr          string
	During        time.Duration
	AutoReconnect bool
	Msg           string
	NewAgent      func(conn *Conn) Agent
}

func (c *Client) Run() (err error) {
	err = c.Init()
	if err != nil {
		return
	}

	c.Connect()
	return
}
func (c *Client) Init() (err error) {
	if strings.ToUpper(c.Protocol) == "TCP" {
		tcpClient = new(TCPClient)
		tcpClient.Addr = c.Addr
		tcpClient.ConnectInterval = c.During
		tcpClient.AutoReconnect = c.AutoReconnect
		tcpClient.init()
	} else if strings.ToUpper(c.Protocol) == "HTTP" {
		wsClient = new(WSClient)
		wsClient.Addr = c.Addr
		wsClient.ConnectInterval = c.During
		tcpClient.AutoReconnect = c.AutoReconnect
		wsClient.init()
	} else {
		err = fmt.Errorf("协议[%s]配置错误，支持TCP/HTTP")
	}
	return
}

func (c *Client) Connect() (err error) {
	if strings.ToUpper(c.Protocol) == "TCP" {
		tcpClient.connect()
	} else if strings.ToUpper(c.Protocol) == "HTTP" {
		wsClient.connect()
	}
	return
}

func (c *Client) Close() {
	if strings.ToUpper(c.Protocol) == "TCP" {
		tcpClient.Close()
	} else if strings.ToUpper(c.Protocol) == "HTTP" {
		wsClient.Close()
	}
}

func (c *Client) SendMsg() (err error) {
	if strings.ToUpper(c.Protocol) == "TCP" {
		tcpClient.connect()
	} else if strings.ToUpper(c.Protocol) == "HTTP" {
		wsClient.connect()
	}
	return
}
