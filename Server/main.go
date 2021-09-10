/*
	CLI-四川麻将-血战到底 服务端
 */
package main

import (
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/dao"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/tcpsrc"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/db"
)

func init(){
	// 初始话数据库
	db.InitMysqlDB("root", "root123", "192.168.0.197", "3306", "test")
	dao.InitTable()
}

func main() {
	tcpsrc.Run()
}



