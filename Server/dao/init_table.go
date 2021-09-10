package dao

import "github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/models"

func InitTable(){
	new(models.User).CreateTable()
}