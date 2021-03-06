/*
	游戏服务逻辑
 */

package tcpsrc

import (
	"fmt"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/dao"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
	"log"
	"net"
)

// 服务逻辑
func Handler(client *ClientUser, b []byte){
	data := entity.TransfeDataDecoder(b)
	log.Println(data.Cmd, data.Timestamp, data.Token, data.Data)

	switch data.Cmd {

	case enum.HeartPacket:
		log.Println("收到心跳包: ", client.Conn.RemoteAddr().String())

	case enum.LoginPacket:
		// 登录
		if v, ok := data.Data.(entity.User); ok{
			log.Println("登录信息: ", v.Name, v.Password)
			user, _ := dao.User().WhereName(v.Name).Get()

			t := &entity.TransfeData{
				Cmd : enum.LoginPacket,
				Token : v.Name,
				Data : true,
				Code : 1,
				Message: "登录成功",
			}

			if user.Id == 0 || user.Password != utils.MD5(v.Password) {
				t.Token = ""
				t.Data = false
				t.Message = "账号或密码错误"
			}
			_, _ = client.Conn.Write(t.Byte())
			client.Token = v.Name
			Add(user.Name, &client.Conn)
		}

	case enum.RegisterPacket:
		// 注册
		if v, ok := data.Data.(entity.User); ok {
			log.Println("登录信息: ", v.Name, v.Password)
			t := &entity.TransfeData{
				Cmd:     enum.RegisterPacket,
				Token:   "",
				Code:    1,
				Data:    true,
				Message: "注册成功",
			}

			if user, _ := dao.User().WhereName(v.Name).Get(); user.Id != 0 {
				t.Data = false
				t.Message = "账号已经存在"
			} else if err := dao.User().Create(models.User{
				Name:     v.Name,
				Password: utils.MD5(v.Password),
			}); err != nil {
				t.Data = false
				t.Message = "注册失败:" + err.Error()
			}
			_, _ = client.Conn.Write(t.Byte())
		}

	case enum.RoomListPacket:
		// 获取房间列表, 房间列表翻页
		log.Println("获取房间列表, 房间列表翻页")
		pg := data.Data.(int)

		var j, z = (pg-1)*10, pg*10
		if j < 0 {
			j = 0
		}
		if j > len(RoomList) {
			j = len(RoomList)
		}
		if z > len(RoomList) {
			z = len(RoomList)
		}
		listData := make([]*entity.RoomShow, 0)
		for _, v := range RoomList[j: z] {
			listData = append(listData, &entity.RoomShow{
				Id: v.Id,
				Name: v.Name,
				State: enum.StateMap[v.State],
				Num: fmt.Sprintf("%d 人", len(v.User)),
			})
		}
		t := &entity.TransfeData{
			Cmd:     enum.RoomListPacket,
			Token:   "",
			Code:    1,
			Data:    listData,
			Message: "获取列表成功",
		}
		log.Println("t = ", t)
		_, _ = client.Conn.Write(t.Byte())

	case enum.CreatRoomPacket:
		// 创建房间
		log.Println("创建房间", data.Data)
		room := &Room{
			Id : len(RoomList) +1,
			Name : data.Data.(string),
			User : make([]*ClientUser, 0),
			State : 0, // 状态 0 未开始  1 开始   2结算
			Ready: make(map[string]bool),
		}
		room.User = append(room.User, client)
		room.Ready[client.Token] = false
		RoomList = append(RoomList, room)
		RoomMap[room.Id] = room
		client.NowRoom = room.Id

		// 房间信息
		roomInfo := entity.RoomInfo{
			Id : room.Id,
			Name : room.Name,
			State : room.State,
			Num : fmt.Sprintf("%d /4 人", len(room.User)),
			Ready : room.Ready,
		}

		// 创建成功返回
		t := &entity.TransfeData{
			Cmd:     enum.CreatRoomPacket,
			Token:   "",
			Code:    0,
			Data: roomInfo,
			Message: "创建成功",
		}
		_, _ = client.Conn.Write(t.Byte())

		//广播所有服务器用户更新
		AllClient.Range(func(k interface{}, v interface{}) bool {
			conn :=  v.(*net.Conn)
			ref := &entity.TransfeData{
				Cmd:     enum.RefreshRoomListPacket,
				Code:    0,
				Data:    true,
			}
			(*conn).Write(ref.Byte())
			return true
		})

	case enum.InToRoomPacket:
		// 进入房间
		log.Println("进入房间", data.Data)
		t := &entity.TransfeData{
			Cmd:     enum.InToRoomPacket,
			Token:   "",
			Code:    0,
			Message: "进入房间成功",
		}
		if room, ok := RoomMap[data.Data.(int)]; ok {
			if len(room.User) >= 4 {
				t.Code = 1
				t.Data = nil
				t.Message = "进入失败，房间满员"
			}else{
				// 进入房间
				client.NowRoom = room.Id
				room.User = append(room.User, client)
				room.Ready[client.Token] = false
				// 下发系统消息
				room.Chat(entity.ChatData{
					From: "[系统]",
					Mag: fmt.Sprintf("%s 进入了房间", client.Token),
				})
				// 房间信息
				roomInfo := entity.RoomInfo{
					Id : room.Id,
					Name : room.Name,
					State : room.State,
					Num : fmt.Sprintf("%d /4 人", len(room.User)),
					Ready : room.Ready,
				}
				t.Data = roomInfo
			}
		}else{
			t.Code = 1
			t.Data = nil
			t.Message = "进入失败，房间ID不存在"
		}
		_, _ = client.Conn.Write(t.Byte())

		// 广播给房间所有人更新
		RoomMap[data.Data.(int)].SendInfo()

	case enum.OutRoomPacket:
		//退出房间
		log.Println("退出房间", data.Data)
		if room, ok := RoomMap[data.Data.(int)]; ok {
			room.Out(client)
			room.Chat(entity.ChatData{
				From: "[系统]",
				Mag: fmt.Sprintf("%s 退出了房间", client.Token),
			})
			room.SendInfo()
		}

	case enum.GameReadyPacket:
		// 准备游戏
		if _, ok := RoomMap[data.Data.(int)]; !ok {
			return
		}
		RoomMap[data.Data.(int)].ReadyLock.Lock()
		RoomMap[data.Data.(int)].Ready[client.Token] = true
		RoomMap[data.Data.(int)].ReadyLock.Unlock()

		// 下发消息
		RoomMap[data.Data.(int)].Chat(entity.ChatData{
			From: "[系统]",
			Mag: fmt.Sprintf("%s 准备了游戏", client.Token),
		})
		// 广播房间信息
		RoomMap[data.Data.(int)].SendInfo()

		// 游戏开始触发条件
		isStart := true
		for _,start := range RoomMap[data.Data.(int)].Ready {
			if !start {
				isStart = false
			}
		}
		// 4名玩家都准备就下发游戏开始指令
		if isStart && len(RoomMap[data.Data.(int)].User) == 4{
			t := &entity.TransfeData{
				Cmd:     enum.StartGamePacket,
				Message: "游戏即将开始",
			}
			log.Println(RoomMap[data.Data.(int)].User)
			for _, user := range RoomMap[data.Data.(int)].User {
				user.IsStart = true
				user.Conn.Write(t.Byte())
			}
		}

	case enum.GameSayPacket:
		// 发起聊天
		msgData := data.Data.(entity.ChatSend)
		if _, ok := RoomMap[msgData.RoomId]; !ok {
			return
		}
		// 下发消息
		RoomMap[msgData.RoomId].Chat(entity.ChatData{
			From: client.Token,
			Mag: msgData.Mag,
		})

	case enum.GameOffPacket:
		// 取消准备
		if _, ok := RoomMap[data.Data.(int)]; !ok {
			return
		}
		RoomMap[data.Data.(int)].ReadyLock.Lock()
		RoomMap[data.Data.(int)].Ready[client.Token] = false
		RoomMap[data.Data.(int)].ReadyLock.Unlock()

		// 下发消息
		RoomMap[data.Data.(int)].Chat(entity.ChatData{
			From: "[系统]",
			Mag: fmt.Sprintf("%s 取消了准备", client.Token),
		})
		// 广播房间信息
		RoomMap[data.Data.(int)].SendInfo()

			// 游戏开始

			// 发牌

			// 下一名玩家

			// 摸牌

			// 打牌

			// 牌型判定

			// 输赢判定

	}


}

