package db

import (
	"config-manager/center/global"
	"config-manager/common/model/node"
	"config-manager/common/model/shortcut"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		SkipDefaultTransaction: false,
	})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&shortcut.OneLineShortcut{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&shortcut.ScriptShortcut{})
	if err != nil {
		panic(err)
	}
	//err = db.AutoMigrate(&trigger.Trigger{})
	//if err != nil {
	//	panic(err)
	//}
	err = db.AutoMigrate(&node.Node{})
	if err != nil {
		panic(err)
	}
	global.DB = db

}
