package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/view"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/entity"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
)

// Login 游戏登录注册
func Login(c *models.TcpClient) {
	for {
		utils.Cle()
		fmt.Println(view.LoginTitle)

		fmt.Print("请输入: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		str := scanner.Text()

		if str == "login" {
			fmt.Print("--------> 登录游戏 \n")
			var user string
			var password string
			fmt.Print("请输入账号: ")
			scanner.Scan()
			user = scanner.Text()
			fmt.Print("请输入密码: ")
			scanner.Scan()
			password = scanner.Text()
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
					c.Token = isLogin.Token
					// 进入游戏界面
					go Home(c)
					return
				}else{
					log.Print("登录失败")
					time.Sleep(2*time.Second)
					continue
				}
			}
		}

		if str == "reg" {
			fmt.Print("--------> 注册账号 \n")
			var user string
			var password string
			fmt.Print("请输入账号: ")
			scanner.Scan()
			user = scanner.Text()
			fmt.Print("请输入密码: ")
			scanner.Scan()
			password = scanner.Text()
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
					time.Sleep(2*time.Second)
					continue
				}else{
					log.Print("注册失败")
					time.Sleep(2*time.Second)
					continue
				}
			}
		}

		if str == "game" {
			Game()
		}

		log.Println("输入未知指令 : ", str)
		time.Sleep(2*time.Second)
	}
}
