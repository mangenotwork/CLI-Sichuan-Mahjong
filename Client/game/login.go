package game

import (
	"fmt"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/view"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
	"log"
	"time"
)

func Login(c *models.TcpClient) {
	for {
		utils.Cle()
		fmt.Println(view.LoginTitle)
		var str string
		fmt.Print("请输入: ")
		fmt.Scanln(&str)

		if str == "login" {
			fmt.Print("--------> 登录游戏 \n")
			var user string
			var password string
			fmt.Print("请输入账号: ")
			fmt.Scanln(&user)
			fmt.Print("请输入密码: ")
			fmt.Scanln(&password)
			log.Println("输入 : ", user, password)
			userData := entity.User{
				Name: user,
				Password: password,
			}
			if _, err := c.Send(entity.NewTransfeData(enum.LoginPacket, "", userData)); err != nil {
				models.RConn <- true
				return
			}
			isLogin := <- c.CmdChan
			log.Print("isLogin = ", isLogin)
			if isLogin.Cmd == enum.LoginPacket{
				if isLogin.Data.(bool) == true {
					log.Print("登录成功")
					// 进入游戏界面
					go Home()
					return
				}else{
					log.Print("登录失败")
					time.Sleep(1*time.Second)
					continue
				}
			}
		}

		if str == "reg" {
			fmt.Print("--------> 注册账号 \n")
			var user string
			var password string
			fmt.Print("请输入账号: ")
			fmt.Scanln(&user)
			fmt.Print("请输入密码: ")
			fmt.Scanln(&password)
			log.Println("输入 : ", user, password)
			userData := entity.User{
				Name: user,
				Password: password,
			}
			if _, err := c.Send(entity.NewTransfeData(enum.RegisterPacket, "", userData)); err != nil {
				models.RConn <- true
				return
			}
			isRegister := <- c.CmdChan
			log.Print("isLogin = ", isRegister)
			if isRegister.Cmd == enum.RegisterPacket{
				if isRegister.Data.(bool) == true {
					log.Print("注册成功")
					time.Sleep(1*time.Second)
					continue
				}else{
					log.Print("注册失败")
					time.Sleep(1*time.Second)
					continue
				}
			}
		}
		log.Println("输入未知指令 : ", str)
		time.Sleep(1*time.Second)

	}
}

// 游戏大厅
func Home(){
	utils.Cle()
	fmt.Print(`
================== 游戏大厅 ==================
	
`)
}