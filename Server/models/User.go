package models

import (
	"log"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/db"
)


type User struct {
	Id int64 `gorm:"primary_key;column:id;size:11" json:"id"`
	Name string `gorm:"column:user_name;size:20" json:"user_name"`
	Password string `gorm:"column:password;size:32" json:"password"`
}

//TableName  默认获取table name
func (*User) TableName() string {
	return "tbl_user"
}

//CreateTable 创建表
func (u *User) CreateTable(){
	if !db.MysqlDB.HasTable(u.TableName()) {
		log.Println("CreateTable User")
		db.MysqlDB.CreateTable(&User{})
	}
}