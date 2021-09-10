package dao

import (
	"errors"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Server/models"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/db"
)

func User() UserDaoInterface{
	return &userDao{}
}

type UserDaoInterface interface {
	WhereId(id int64) *userDao
	WhereName(name string) *userDao
	Create(user models.User) error
	IsHave(userName string) bool
	Get() (models.User, error)
}

type userDao struct {
	id int64	// 查询字段id
	name string	// 查询字段name
}

func (d *userDao) WhereId(id int64) *userDao {
	d.id = id
	return d
}

func (d *userDao) WhereName(name string) *userDao {
	d.name = name
	return d
}

func (d *userDao) Get() (models.User, error) {
	var (
		data models.User
		err error
		ok = false
		dbConn = db.MysqlDB
	)
	dbConn = dbConn.Table(data.TableName())

	if d.id != 0 {
		ok = true
		dbConn = dbConn.Where("id=?", d.id)
	}

	if len(d.name) > 0 {
		ok = true
		dbConn = dbConn.Where("user_name=?", d.name)
	}

	if ok {
		err = dbConn.First(&data).Error

	} else {
		err = errors.New("参数不够")
	}

	return data, err
}

func (d *userDao) IsHave(userName string) bool {
	d.name = userName
	data, _ := d.Get()
	if data.Id > 0 {
		return true
	}
	return false
}

func (d *userDao) Create(user models.User) error {
	var dbConn = db.MysqlDB
	return dbConn.Table(user.TableName()).Create(&user).Error
}

