package tcpsrc

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
	"io"
	"log"
	"net"
)

// tcp
type TcpServer struct {
	Listener   *net.TCPListener
	HawkServer *net.TCPAddr
}

// 运行服务
func Run() {
	//类似于初始化套接字，绑定端口
	hawkServer, err := net.ResolveTCPAddr("tcp", "0.0.0.0:14444")
	utils.PanicErr(err)

	//侦听
	listen, err := net.ListenTCP("tcp", hawkServer)
	utils.PanicErr(err)

	//关闭
	defer listen.Close()

	tcpServer := &TcpServer{
		Listener:   listen,
		HawkServer: hawkServer,
	}
	log.Println("start Master TCP server successful.")

	//接收请求
	for {

		//来自客户端的连接
		conn, err := tcpServer.Listener.Accept()
		if err != nil {
			log.Println("[连接失败]:", err.Error())
			continue
		}
		log.Println("[连接成功]: ", conn.RemoteAddr().String(), conn)

		go func(conn net.Conn){
			recv := make([]byte, 1024)
			for {

				//err = conn.SetReadDeadline(time.Now().Add(2*time.Second)) // timeout
				//if err != nil {
				//	log.Println("setReadDeadline failed:", err)
				//}

				n, err := conn.Read(recv)
				log.Println(n, err)
				if err != nil{
					if err == io.EOF {
						log.Println(conn.RemoteAddr().String(), " 断开了连接!")
						conn.Close()
						return
					}
				}
				if n > 0 && n < 1025 {
					Handler(conn, recv[:n])
				}
			}
		}(conn)

	}
}
