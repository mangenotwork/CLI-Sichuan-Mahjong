/*
	游戏服务逻辑
 */

package tcpsrc

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/dao"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
	"log"
	"net"
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
				Token : "",
				Data : true,
				Code : 1,
				Message: "登录成功",
			}

			if user.Id == 0 || user.Password != utils.MD5(v.Password) {
				t.Data = false
				t.Message = "账号或密码错误"
			}
			_, _ = conn.Write(t.Byte())
			Add(user.Name, &conn)
		}

	case enum.RegisterPacket:
		// 注册
		if v, ok := data.Data.(entity.User); ok{
			log.Println("登录信息: ", v.Name, v.Password)
			t := &entity.TransfeData{
				Cmd : enum.RegisterPacket,
				Token : "",
				Code : 1,
				Data : true,
				Message: "注册成功",
			}

			if user, _ := dao.User().WhereName(v.Name).Get(); user.Id != 0{
				t.Data = false
				t.Message = "账号已经存在"
			}else if err := dao.User().Create(models.User{
				Name: v.Name,
				Password: utils.MD5(v.Password),
			}); err != nil {
				t.Data = false
				t.Message = "注册失败:"+err.Error()
			}
			_, _ = conn.Write(t.Byte())

			// 获取房间列表

			// 房间列表翻页

			// 创建房间

			// 进入房间

		}
	}


}