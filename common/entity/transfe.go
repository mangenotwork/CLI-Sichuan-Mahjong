package entity

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
)

// 注册gob类型
func init(){
	gob.Register(User{})
	gob.Register([]*RoomShow{})
	gob.Register(ChatData{})
}

// 传输数据结构
type TransfeData struct {
	Cmd enum.Command  // 指令
	Timestamp int64
	Token string // 识别客户端身份
	Data interface{} // 传输的数据
	Message string // 传输消息
	Code int // 传输code
}

func (t *TransfeData) Byte() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}

func NewTransfeData(cmd enum.Command, token string, data interface{}) []byte {
	tra := &TransfeData{
		Cmd: cmd,
		Timestamp: time.Now().Unix(),
		Token: token,
		Data: data,
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(tra)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}

func TransfeDataDecoder(data []byte) *TransfeData {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	tra := &TransfeData{}
	err := dec.Decode(&tra)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return tra
}

// User
type User struct {
	Name string
	Password string
}

// 客户端的显示
type RoomShow struct {
	Id int
	Name string
	State string // 0:未开始   1:游戏中
	Num string // 多少人
}

// 聊天
type ChatData struct {
	From string // 使用token
	Mag string // 内容
}