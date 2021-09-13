/*
	游戏服务逻辑
 */

package tcpsrc

import (
	"fmt"
	"log"
	"net"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/dao"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
)

// 服务逻辑
func Handler(conn net.Conn, b []byte){
	data := entity.TransfeDataDecoder(b)
	log.Println(data.Cmd, data.Timestamp, data.Token, data.Data)

	switch data.Cmd {

	case enum.HeartPacket:
		log.Println("收到心跳包: ", conn.RemoteAddr().String())

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
			_, _ = conn.Write(t.Byte())
			Add(user.Name, &conn)
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
			_, _ = conn.Write(t.Byte())
		}

	case enum.RoomListPacket:
		// 获取房间列表, 房间列表翻页
		log.Println("获取房间列表, 房间列表翻页")
		pg := data.Data.(int)

		var j, z = (pg-1)*10, pg*10
		if j < 0 {
			j = 0
		}
		if j > len(entity.RoomList) {
			j = len(entity.RoomList)
		}
		if z > len(entity.RoomList) {
			z = len(entity.RoomList)
		}

		listData := make([]*entity.RoomShow, 0)
		for _, v := range entity.RoomList[j: z] {
			listData = append(listData, &entity.RoomShow{
				Id: v.Id,
				Name: v.Name,
				State: enum.StateMap[v.State],
				Num: fmt.Sprintf("%d 人", len(v.User)),
			})
		}
		t := &entity.TransfeData{
			Cmd:     enum.CreatRoomPacket,
			Token:   "",
			Code:    1,
			Data:    listData,
			Message: "获取列表成功",
		}
		log.Println("t = ", t)
		_, _ = conn.Write(t.Byte())

	case enum.CreatRoomPacket:
		// 创建房间
		log.Println("创建房间", data.Data)
		room := &entity.Room{
			Id : len(entity.RoomList) +1,
			Name : data.Data.(string),
			User : make([]*net.Conn, 0),
			State : 0, // 状态 0 未开始  1 开始   2结算
		}
		room.User = append(room.User, &conn)
		entity.RoomList = append(entity.RoomList, room)
		entity.RoomMap[room.Id] = room
		// 创建成功
		t := &entity.TransfeData{
			Cmd:     enum.CreatRoomPacket,
			Token:   "",
			Code:    1,
			Data:    true,
			Message: "创建成功",
		}
		_, _ = conn.Write(t.Byte())

	case enum.InToRoomPacket:
		// 进入房间
		log.Println("进入房间", data.Data)
		t := &entity.TransfeData{
			Cmd:     enum.InToRoomPacket,
			Token:   "",
			Code:    1,
			Data:    true,
			Message: "进入房间成功",
		}
		if _, ok := entity.RoomMap[data.Data.(int)]; ok {
			entity.RoomMap[data.Data.(int)].User = append(entity.RoomMap[data.Data.(int)].User, &conn)
		}else{
			t.Data = false
			t.Message = "进入失败，房间ID不存在"
		}
		_, _ = conn.Write(t.Byte())

	case enum.OutRoomPacket:
		//退出房间

			// 准备游戏

			// 游戏开始

			// 发牌

			// 下一名玩家

			// 摸牌

			// 打牌

			// 牌型判定

			// 输赢判定

	}


}

