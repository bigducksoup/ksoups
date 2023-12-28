package global

import (
	"gorm.io/gorm"
	"log"
)

var Logger *log.Logger
var DB *gorm.DB

func init() {
	Logger = log.Default()
}
