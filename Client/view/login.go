package view

import (
	"fmt"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
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


func GameRoomInit(){
	utils.Cle()

	fmt.Println(`_______________________________________________
房间ID : 1    | 房间名 : aaa    | 人数: 2/4

| aaa    未准备
| bbb    未准备

[房间聊天]



===================== 操作 ==========================
- 输入 q 退出房间
- 输入 ok 准备游戏
- 输入 say 进行聊天
`)

}
