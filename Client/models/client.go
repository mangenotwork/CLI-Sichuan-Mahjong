package models

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"net"
)

// 重连
var RConn = make(chan bool)

type TcpClient struct {
	Connection *net.TCPConn
	HawkServer *net.TCPAddr
	StopChan   chan struct{}
	CmdChan chan *entity.TransfeData
}

func (c *TcpClient) Send(b []byte) (int, error) {
	return c.Connection.Write(b)
}

func (c *TcpClient) Read(b []byte) (int, error) {
	return c.Connection.Read(b)
}

func (c *TcpClient) Addr() string {
	return c.Connection.RemoteAddr().String()
}

func (c *TcpClient) Close(){
	c.Connection.Close()
	RConn <- true
}
