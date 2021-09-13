package tcpsrc

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"net"
	"sync"
)

// 客户端
type ClientUser struct{
	Conn net.Conn
	Token string
}

// 存放所有客户端
var AllClient sync.Map

func Add(userName string, conn *net.Conn) {
	AllClient.Store(userName, conn)
}

func Get(userName string) *net.Conn {
	val, ok := AllClient.Load(userName)
	if ok {
		return val.(*net.Conn)
	}
	return nil
}

func Out(userName string) {
	AllClient.Delete(userName)
}


// 游戏房间列表
var RoomList = make([]*Room, 0)
var RoomMap = make(map[int]*Room, 0)

// 游戏房间
type Room struct {
	Id int
	Name string
	User []*ClientUser // token对应conn
	State int // 状态
}

// 房间里聊天
func (r *Room) Chat(msg entity.ChatData){
	for _, v := range r.User{
		t := &entity.TransfeData{
			Cmd:     enum.ChatPacket,
			Token:   "",
			Code:    1,
			Data:    msg,
			Message: "",
		}
		_, _ = v.Conn.Write(t.Byte())
	}
}