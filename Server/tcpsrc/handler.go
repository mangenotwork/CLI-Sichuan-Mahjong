/*
	游戏服务逻辑
 */

package tcpsrc

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
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
		log.Println("登录")
		if v, ok := data.Data.(entity.User); ok{
			log.Println("登录信息: ", v.Name, v.Password)
			_, _ = conn.Write(entity.NewTransfeData(enum.LoginPacket,"", true))
		}
		// 创建新的客户端连接实例
		// 存储客户端连接

	}


}