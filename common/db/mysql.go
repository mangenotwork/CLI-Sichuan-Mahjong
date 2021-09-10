package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

var MysqlDB *gorm.DB

func InitMysqlDB(user, pass, host, port, dbname string) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname) + "?charset=utf8mb4&parseTime=true&loc=Local&timeout=5s"
	log.Println("连接数据 = ", dsn)
	MysqlDB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	MysqlDB.LogMode(true)
	MysqlDB.DB().SetMaxIdleConns(10)
	MysqlDB.DB().SetMaxOpenConns(20)
}

func GetMysqlDB() *gorm.DB {
	return MysqlDB
}