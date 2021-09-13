package game

import (
	"fmt"
	"log"
	"time"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/view"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
)

// Login 游戏登录注册
func Login(c *models.TcpClient) {
	for {
		utils.Cle()
		fmt.Println(view.LoginTitle)
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
				models.RConn <- true
				return
			}
			isLogin := <- c.CmdChan
			log.Print("isLogin = ", isLogin)
			if isLogin.Cmd == enum.LoginPacket{
				if isLogin.Data.(bool) == true {
					log.Print("登录成功")
					c.Token = isLogin.Token
					// 进入游戏界面
					go Home(c)
					return
				}else{
					log.Print("登录失败")
					time.Sleep(2*time.Second)
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
				models.RConn <- true
				return
			}
			isRegister := <- c.CmdChan
			log.Print("isLogin = ", isRegister)
			if isRegister.Cmd == enum.RegisterPacket{
				if isRegister.Data.(bool) == true {
					log.Print("注册成功")
					time.Sleep(2*time.Second)
					continue
				}else{
					log.Print("注册失败")
					time.Sleep(2*time.Second)
					continue
				}
			}
		}
		log.Println("输入未知指令 : ", str)
		time.Sleep(2*time.Second)
	}
}


// Home 游戏大厅
func Home(c *models.TcpClient){
	var pg = 1

	// 指令
	go func(c *models.TcpClient){
		for {
			// 获取房间列表
			if _, err := c.Send(entity.NewTransfeData(enum.RoomListPacket, "", pg)); err != nil {
				models.RConn <- true
				return
			}
			roomList := <- c.CmdChan
			log.Println("房间列表: ", roomList)
			roomStr := ""
			for _, v := range roomList.Data.([]*entity.RoomShow) {
				roomStr = roomStr + fmt.Sprintf("ID: %d \t房间名: %s \t房间状态: %s \t房间人数: %s \n", v.Id, v.Name, v.State, v.Num)
			}

			view.HomeView(c.Token, roomStr, pg)

			var str string
			fmt.Print("请输入: ")
			fmt.Scanln(&str)

			if str == "up" {
				// 获取房间列表 - 上一页
				pg++
				continue
			}

			if str == "down" {
				// 获取房间列表 - 下一页
				pg--
				if pg < 1 {
					pg = 1
				}
				continue
			}

			if str == "add" {
				var name string
				fmt.Print("请输入房间名: ")
				fmt.Scanln(&name)
				// 创建房间
				if _, err := c.Send(entity.NewTransfeData(enum.CreatRoomPacket, "", name)); err != nil {
					models.RConn <- true
					return
				}
				res := <- c.CmdChan
				if res.Cmd == enum.CreatRoomPacket && res.Data.(bool) {
					// 创建成功进入游戏房间
					go GameRoom(c)
					return
				} else {
					log.Println("创建房间失败！！！")
					time.Sleep(2*time.Second)
					continue
				}
			}

			if str == "to" {
				// 进入游戏房间
				var roomId int
				fmt.Print("请输入房间ID: ")
				fmt.Scanln(&roomId)
				if _, err := c.Send(entity.NewTransfeData(enum.InToRoomPacket, "", roomId)); err != nil {
					models.RConn <- true
					return
				}
				res := <- c.CmdChan
				if res.Cmd == enum.InToRoomPacket && res.Data.(bool) {
					// 创建成功进入游戏房间
					go GameRoom(c)
					return
				} else {
					log.Println("进入房间失败！", res.Message)
					time.Sleep(2*time.Second)
					continue
				}
			}

		}
	}(c)

}


// GameRoom 游戏房间
func GameRoom(c *models.TcpClient){
	gameChan := make(chan string)
	endChan := make(chan int)


	// 获取房间当前信息
	view.GameRoomInit()

	// 输入操作
	go func(){
		for{
			select {

			case <- endChan:
				return

			default:
				var str string
				fmt.Scanln(&str)
				log.Println("输入 : ", str)
				gameChan <- str
			}
		}
	}()

	// 游戏逻辑
	go func() {
		for {
			select {
			case cmd := <-gameChan:
				log.Println(cmd)

				//退出房间
				if cmd == "q"{
					endChan <- 0
					return
				}
			}
		}
	}()


	for  {
		select {
		case <- endChan:
			// 退出房间
			go Home(c)
			return
		}
	}
}