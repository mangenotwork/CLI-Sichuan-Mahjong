package tcpsrc

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"net"
	"sync"
	"fmt"
)

// 客户端
type ClientUser struct{
	Conn net.Conn
	Token string
	NowRoom int // 当前所在的房间
	IsStart bool // 是否在游戏进行中
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
	Ready map[string]bool //房间用户准备情况
	ReadyLock sync.Mutex
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
	return
}

// 广播房间信息
func (r *Room) SendInfo(){
	// 广播房间信息
	r.ReadyLock.Lock()
	defer r.ReadyLock.Unlock()

	t := &entity.TransfeData{
		Cmd:     enum.RoomInfoPacket,
		Data: entity.RoomInfo{
			Id : r.Id,
			Name : r.Name,
			State : r.State,
			Num : fmt.Sprintf("%d /4 人", len(r.User)),
			Ready : r.Ready,
		},
	}
	for _, v := range r.User {
		v.Conn.Write(t.Byte())
	}
	return
}

// 退出房间
func (r *Room) Out(client *ClientUser){
	for i, c := range r.User{
		if c == client{
			r.User = append(r.User[:i], r.User[i+1:]...)
		}
	}
	delete(r.Ready, client.Token)
	client.NowRoom = 0
}