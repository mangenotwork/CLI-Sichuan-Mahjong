package enum

import "runtime"

type Command string


// 游戏的交互指令包
const (
	HeartPacket Command = "Heart_Packet" // 心跳包
	LoginPacket Command = "Login_Packet" // 登录包
	RegisterPacket Command = "Register_Packet" // 注册包
	RoomListPacket Command = "RoomList_Packet" // 获取房间列表
	CreatRoomPacket Command = "CreatRoom_Packet" // 创建房间
	InToRoomPacket Command = "InToRoom_Packet" // 进入房间
	OutRoomPacket Command = "OutRoom_Packet" // 退出房间
	RoomInfoPacket Command = "RoomInfo_Packet" // 房间信息
	ChatPacket Command = "Chat_Packet" // 聊天
	RefreshRoomListPacket Command = "RefreshRoomList_Packet" // 提醒刷新列表
	GameReadyPacket Command = "GameReady_Packet" // 准备游戏
	GameSayPacket Command = "GameSay_Packet" // 发起聊天
	GameOffPacket Command = "GameOff_Packet" // 取消准备
	StartGamePacket Command = "StartGame_Packet" // 游戏开始指令
)

const (
	SYS_TYPE = runtime.GOOS
)


var StateMap = map[int]string{
	0:"未开始",
	1:"游戏中",
}