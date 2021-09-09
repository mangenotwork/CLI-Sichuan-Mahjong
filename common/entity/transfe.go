package entity

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
)

func init(){
	gob.Register(User{})
}

type TransfeData struct {
	Cmd enum.Command  // 指令
	Timestamp int64
	Token string // 识别客户端身份
	Data interface{} // 传输的数据
	Message string // 传输消息
	Code int // 传输code
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