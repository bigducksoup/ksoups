package db

import (
	"config-manager/center/global"
	"config-manager/center/model"
	"config-manager/common/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserInfo struct {
	Id       string `json:"id"`
	Account  string `json:"account" gorm:"unique"`
	Password string `json:"password"`
}

func InitDB() {

	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	// 自动生成表
	err = db.AutoMigrate(
		&UserInfo{},
		&model.Shortcut{},
		&model.ShortcutExecLog{},
		&model.Chain{},
		&model.DispatchLog{},
		&model.NodeExecLog{},
		&model.Node{},
		&model.Edge{},
		&model.ShortcutNodeBinding{},
		&model.ProbeInfo{},
	)
	if err != nil {
		panic(err)
	}

	// 初始化数据，设置全部探针为offline状态
	db.Model(&model.ProbeInfo{}).Where("online = ?", true).Update("online", false)

	// 初始化账号密码
	db.Create(&UserInfo{
		Id:       utils.UUID(),
		Account:  global.Conf.Account,
		Password: utils.Md5([]byte(global.Conf.Password)),
	})

	// 初始化全局变量
	global.DB = db
}
