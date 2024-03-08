package global

import (
	"apps/center/config"
	"apps/center/server/core"
	"gorm.io/gorm"
	"log"
)

var CenterServer *core.CenterServer
var DB *gorm.DB
var Conf config.Config = config.Config{}

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("[center]")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
}
