package global

import (
	"config-manager/center/config"
	"gorm.io/gorm"
	"log"
)

var Logger *log.Logger
var DB *gorm.DB
var Conf config.Config = config.Config{}

func init() {
	Logger = log.Default()
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("[center]")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
}
