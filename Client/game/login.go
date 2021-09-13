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
	var pg = 1 // 房间列表页数
	roomStr := "" // 房间列表
	departure := make(chan int)

	// 获取房间列表
	if _, err := c.Send(entity.NewTransfeData(enum.RoomListPacket, "", pg)); err != nil {
		models.RConn <- true
		return
	}

	// 指令
	go func(c *models.TcpClient){
		for {

			select {
			case <-departure:
				return

			default:
				var str string
				fmt.Print("请输入: ")
				fmt.Scanln(&str)
				log.Println("输入 --> ", str)

				if str == "up" {
					// 获取房间列表 - 上一页
					pg++
					if _, err := c.Send(entity.NewTransfeData(enum.RoomListPacket, "", pg)); err != nil {
						models.RConn <- true
						return
					}
					continue
				}

				if str == "down" {
					// 获取房间列表 - 下一页
					pg--
					if pg < 1 {
						pg = 1
					}
					if _, err := c.Send(entity.NewTransfeData(enum.RoomListPacket, "", pg)); err != nil {
						models.RConn <- true
						return
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
				}

			}

		}
	}(c)


	for {
		select {
			case res := <- c.CmdChan :

				// 获取房间列表
				if res.Cmd == enum.RoomListPacket {
					for _, v := range res.Data.([]*entity.RoomShow) {
						roomStr = roomStr + fmt.Sprintf("ID: %d \t房间名: %s \t房间状态: %s \t房间人数: %s \n", v.Id, v.Name, v.State, v.Num)
					}
					view.HomeView(c.Token, roomStr, pg)
				}

				// 刷新列表
				if res.Cmd == enum.RefreshRoomListPacket {
					if _, err := c.Send(entity.NewTransfeData(enum.RoomListPacket, "", pg)); err != nil {
						models.RConn <- true
						return
					}
				}

				// 创建房间后的应答
				if res.Cmd == enum.CreatRoomPacket {
					if res.Data.(bool) {
						// 创建成功进入游戏房间
						go GameRoom(c)
						departure <- 0
						return
					} else {
						log.Println("创建房间失败！！！")
					}
				}

				// 进入游戏房间
				if res.Cmd == enum.InToRoomPacket {
					if res.Data.(bool) {
						// 创建成功进入游戏房间
						go GameRoom(c)
						departure <- 0
						return
					} else {
						log.Println("进入房间失败！", res.Message)
					}
				}


		}
	}

}


// GameRoom 游戏房间
func GameRoom(c *models.TcpClient){
	gameChan := make(chan string)
	endChan := make(chan int)
	ChatMsgStr := ""

	// 获取房间当前信息
	view.GameRoomInit(ChatMsgStr)

	// 输入操作
	go func(){
		for{
			var str string
			fmt.Scanln(&str)
			log.Println("输入 : ", str)
			if str == "q"{
				if _, err := c.Send(entity.NewTransfeData(enum.OutRoomPacket, c.Token, "")); err != nil {
					models.RConn <- true
					return
				}
				endChan <- 0
				return
			}
		}
	}()

	// 游戏逻辑
	for {
		select {
		case <- endChan:
			go Home(c)
			close(endChan)
			close(gameChan)
			return

		case cmd := <-gameChan:
			//来自操作交互
			log.Println(cmd)

		case rse := <- c.CmdChan:
			// 聊天消息
			log.Println(rse)
			if rse.Cmd == enum.ChatPacket {
				ChatMsgStr = ChatMsgStr + fmt.Sprintf("%s : %s", rse.Data.(entity.ChatData).From,
					rse.Data.(entity.ChatData).Mag)
			}
			view.GameRoomInit(ChatMsgStr)

		}
	}


}