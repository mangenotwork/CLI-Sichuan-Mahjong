package tcpsrc

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/game"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
)


func Run(){

	//用于重连
Reconnection:

	host := "192.168.0.9:14444"
	hawkServer, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Printf("hawk server [%s] resolve error: [%s]", host, err.Error())
		time.Sleep(1 * time.Second)
		goto Reconnection
	}

	//连接服务器
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	if err != nil {
		log.Printf("connect to hawk server error: [%s]", err.Error())
		time.Sleep(1 * time.Second)
		goto Reconnection
	}
	log.Println("[连接成功] 连接服务器成功")

	//创建客户端实例
	client := &models.TcpClient{
		Connection: connection,
		HawkServer: hawkServer,
		StopChan:   make(chan struct{}),
		CmdChan: make(chan *entity.TransfeData),
	}

	//启动接收
	go func(conn *models.TcpClient){
		for{
			recv := make([]byte, 1024)
			for {
				n, err := conn.Connection.Read(recv)
				if err != nil{
					if err == io.EOF {
						log.Println(conn.Addr(), " 断开了连接!")
						conn.Close()
						return
					}
				}
				if n > 0 && n < 1025 {
					conn.CmdChan <- entity.TransfeDataDecoder(recv[:n])
				}
			}
		}
	}(client)

	// 发送心跳
	go func(conn *models.TcpClient){
		i := 0
		heartBeatTick := time.Tick(10 * time.Second)
		for {
			select {
			case <-heartBeatTick:
				heartBeat := entity.NewTransfeData(enum.HeartPacket, "", i)
				if _, err := conn.Send(heartBeat); err != nil {
					models.RConn <- true
					return
				}
				i++
			case <-conn.StopChan:
				return
			}
		}
	}(client)

	// 登录界面
	go game.Login(client)


	for {
		select {
		case a := <- models.RConn:
			log.Println("global.RConn = ", a)
			goto Reconnection
		}
	}

	//等待退出
	<-client.StopChan

}