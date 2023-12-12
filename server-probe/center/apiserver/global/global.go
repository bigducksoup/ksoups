package global

import "log"

var Logger *log.Logger

func init() {
	Logger = log.Default()
}
