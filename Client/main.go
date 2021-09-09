/*
	CLI-四川麻将-血战到底 客户端
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
)

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

// 重连
var RConn = make(chan bool)

func main(){

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
	client := &TcpClient{
		Connection: connection,
		HawkServer: hawkServer,
		StopChan:   make(chan struct{}),
		CmdChan: make(chan *entity.TransfeData),
	}

	//启动接收
	go func(conn *TcpClient){
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
	go func(conn *TcpClient){
		i := 0
		heartBeatTick := time.Tick(10 * time.Second)
		for {
			select {
			case <-heartBeatTick:
				heartBeat := entity.NewTransfeData(enum.HeartPacket, "", i)
				if _, err := conn.Send(heartBeat); err != nil {
					RConn <- true
					return
				}
				i++
			case <-conn.StopChan:
				return
			}
		}
	}(client)

	// 登录界面
	go Login(client)


	for {
		select {
		case a := <- RConn:
			log.Println("global.RConn = ", a)
			goto Reconnection
		}
	}

	//等待退出
	<-client.StopChan


}

func Login(c *TcpClient) {
	for {
		utils.Cle()
		fmt.Println(`
============ 四川麻将 血战到底 ==============
输入 login : 登录游戏
输入 reg : 注册账号
`)
		var str string
		fmt.Print("请输入: ")
		fmt.Scanln(&str)

		if str == "login" {
			fmt.Print("--------> 登录游戏 \n")
			var user string
			var password string
			fmt.Print("请输入账号: ")
			fmt.Scanln(&user)
			fmt.Print("请输入密码: ")
			fmt.Scanln(&password)
			log.Println("输入 : ", user, password)
			userData := entity.User{
				Name: user,
				Password: password,
			}
			if _, err := c.Send(entity.NewTransfeData(enum.LoginPacket, "", userData)); err != nil {
				RConn <- true
				return
			}
			isLogin := <- c.CmdChan
			log.Print("isLogin = ", isLogin)
			if isLogin.Cmd == enum.LoginPacket{
				if isLogin.Data.(bool) == true {
					log.Print("登录成功")
					// 进入游戏界面
					go Home()
					return
				}else{
					log.Print("登录失败")
					time.Sleep(1*time.Second)
					continue
				}
			}
		}

		if str == "reg" {
			fmt.Print("--------> 注册账号 \n")
			var user string
			var password string
			fmt.Print("请输入账号: ")
			fmt.Scanln(&user)
			fmt.Print("请输入密码: ")
			fmt.Scanln(&password)
			log.Println("输入 : ", user, password)
			userData := entity.User{
				Name: user,
				Password: password,
			}
			if _, err := c.Send(entity.NewTransfeData(enum.RegisterPacket, "", userData)); err != nil {
				RConn <- true
				return
			}
			isRegister := <- c.CmdChan
			log.Print("isLogin = ", isRegister)
			if isRegister.Cmd == enum.RegisterPacket{
				if isRegister.Data.(bool) == true {
					log.Print("注册成功")
					time.Sleep(1*time.Second)
					continue
				}else{
					log.Print("注册失败")
					time.Sleep(1*time.Second)
					continue
				}
			}
		}
		log.Println("输入未知指令 : ", str)
		time.Sleep(1*time.Second)

	}
}

// 游戏大厅
func Home(){
	utils.Cle()
	fmt.Print(`
================== 游戏大厅 ==================

`)
}