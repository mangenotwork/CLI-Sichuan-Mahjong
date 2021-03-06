package game

import (
	"fmt"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/view"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"log"
	"strings"
	"time"
)

var pg = 1 // 房间列表页数

// Home 游戏大厅
func Home(c *models.TcpClient){

	// 指令
	go HomeInput(c)

	for {
		select {
		case res := <- c.CmdChan :

			// 获取房间列表
			if res.Cmd == enum.RoomListPacket {
				roomStr := ""
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
				if res.Code == 0 {
					// 创建成功进入游戏房间
					log.Println("创建成功进入游戏房间")
					go GameRoom(c, res.Data.(entity.RoomInfo))
					return
				} else {
					log.Println("创建房间失败！！！")
					time.Sleep(2*time.Second)
					go HomeInput(c)
				}
			}

			// 进入游戏房间
			if res.Cmd == enum.InToRoomPacket {
				if res.Code == 0 {
					// 进入游戏房间
					log.Println("进入游戏房间")
					go GameRoom(c, res.Data.(entity.RoomInfo))
					return
				} else {
					log.Println("进入房间失败！", res.Message)
					time.Sleep(2*time.Second)
					go HomeInput(c)
				}
			}

		}
	}

}

// 首页交互指令
func HomeInput(c *models.TcpClient){
	for {

		// 获取房间列表
		if _, err := c.Send(entity.NewTransfeData(enum.RoomListPacket, "", pg)); err != nil {
			models.RConn <- true
			return
		}

		var str string
		fmt.Scanln(&str)
		log.Println("Home 输入 --> ", str)

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
			return
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
			return
		}

	}
}


// GameRoom 游戏房间
func GameRoom(c *models.TcpClient,info entity.RoomInfo){
	endChan := make(chan int) // 退出房间的chan
	ChatMsgStr := make([]string, 0)

	// 获取房间当前信息
	view.GameRoomInit(info, strings.Join(ChatMsgStr, "\n"))

	// 输入操作
	go func(ch chan int){
		for{
			var gameInput string
			fmt.Scanln(&gameInput)
			log.Println("GameRoom 输入 : ", gameInput)

			// 退出房间
			if gameInput == "q"{
				log.Println("退出房间")
				if _, err := c.Send(entity.NewTransfeData(enum.OutRoomPacket, c.Token, info.Id)); err != nil {
					models.RConn <- true
					return
				}
				ch <- 0
				return
			}

			// 准备游戏
			if gameInput == "ok"{
				if _, err := c.Send(entity.NewTransfeData(enum.GameReadyPacket, c.Token, info.Id)); err != nil {
					models.RConn <- true
					return
				}
			}

			// 发起聊天
			if gameInput == "say"{
				var msg string
				fmt.Scanln(&msg)
				log.Println("发起聊天 输入 : ", msg)
				if _, err := c.Send(entity.NewTransfeData(enum.GameSayPacket, c.Token, entity.ChatSend{
					RoomId: info.Id,
					Mag: msg,
				})); err != nil {
					models.RConn <- true
					return
				}
			}

			// 取消准备
			if gameInput == "off"{
				if _, err := c.Send(entity.NewTransfeData(enum.GameOffPacket, c.Token, info.Id)); err != nil {
					models.RConn <- true
					return
				}
			}

		}
	}(endChan)

	// 游戏逻辑
	for {
		select {
		case <- endChan:
			log.Println("退出房间")
			go Home(c)
			close(endChan)
			return

		case rse := <- c.CmdChan:
			// 聊天消息
			log.Println(rse)

			// 房间消息
			if rse.Cmd == enum.ChatPacket {
				if len(ChatMsgStr) > 5 {
					ChatMsgStr = ChatMsgStr[1:len(ChatMsgStr)-1]
				}
				ChatMsgStr = append(ChatMsgStr, fmt.Sprintf("%s %s : %s ", time.Now().Format("2006-01-02 15:04:05"),
					rse.Data.(entity.ChatData).From, rse.Data.(entity.ChatData).Mag))
			}

			// 游戏房间信息
			if rse.Cmd == enum.RoomInfoPacket {
				info = rse.Data.(entity.RoomInfo)
			}

			// 游戏开始
			if rse.Cmd == enum.StartGamePacket {
				log.Println("游戏开始")
				for i:=3;i>0;i--{
					time.Sleep(1*time.Second)
					if len(ChatMsgStr) > 5 {
						ChatMsgStr = ChatMsgStr[1:len(ChatMsgStr)-1]
					}
					ChatMsgStr = append(ChatMsgStr, fmt.Sprintf("%s 游戏即将开始倒计时 : %d s", time.Now().Format("2006-01-02 15:04:05"), i))
					view.GameRoomInit(info, strings.Join(ChatMsgStr, "\n"))
				}
				//渲染游戏-> 发牌界面
				Game()
			}

			// 界面渲染
			view.GameRoomInit(info, strings.Join(ChatMsgStr, "\n"))
		}
	}


}


