package enum

import "runtime"

type Command string

const (
	HeartPacket Command = "Heart_Packet" // 心跳包
	LoginPacket Command = "Login_Packet" // 登录包
	RegisterPacket Command = "Register_Packet" // 注册包

)

const (
	SYS_TYPE = runtime.GOOS
)