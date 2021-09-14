package view

import (
	"fmt"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
	"log"
)


const LoginTitle = `
============ 四川麻将 血战到底 ==============
输入 login : 登录游戏
输入 reg : 注册账号
`


func HomeView(user, roomList string, pg int){

	utils.Cle()

	if roomList == "" {
		roomList = "没有房间了\n"
	}

	fmt.Print(`
================== 游戏大厅 ==================
	用户: `+user+`           

- 输入 up 上一页； 输入 down 下一页
- 输入 add+空格+房名 创建房间
- 输入 to+空格+房间编号 进入房间

———————————— 第 `+fmt.Sprintf("%d",pg)+` 页 ————————————————

`)
	// room列表
	fmt.Print(roomList)
	fmt.Print("\n_____________________________________________\n\n")

}


func GameRoomInit(roomInfo entity.RoomInfo, chatMsg string){
	utils.Cle()
	log.Println("roomInfo = ", roomInfo)
	readyStr := ""
	for k,v := range roomInfo.Ready {
		if v {
			readyStr = readyStr + fmt.Sprintf("| %s \t已准备\n", k)
		}else{
			readyStr = readyStr + fmt.Sprintf("| %s \t未准备\n", k)
		}

	}

	fmt.Println(`_______________________________________________
房间ID : `+fmt.Sprintf("%d",roomInfo.Id)+`    | 房间名 : `+roomInfo.Name+`   | 人数: `+roomInfo.Num+`   | 状态: `+enum.StateMap[roomInfo.State]+`

`+readyStr+`

[房间聊天]
`+chatMsg+`


===================== 操作 ==========================
- 输入 q 退出房间
- 输入 ok 准备游戏
- 输入 say 进行聊天
- 输入 off 取消准备
`)

}
