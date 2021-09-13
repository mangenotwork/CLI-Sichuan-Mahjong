package tcpsrc

import (
	"net"
	"sync"
)

// 存放所有客户端
var AllClient sync.Map

func Add(userName string, conn *net.Conn) {
	AllClient.Store(userName, conn)
}

func Get(userName string) *net.Conn {
	val, ok := AllClient.Load(userName)
	if ok {
		return val.(*net.Conn)
	}
	return nil
}

func Out(userName string) {
	AllClient.Delete(userName)
}



