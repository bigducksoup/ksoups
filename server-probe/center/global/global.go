package global

import (
	"config-manager/center/config"
	"gorm.io/gorm"
	"log"
)

var Logger Log
var DB *gorm.DB
var Conf config.Config = config.Config{}

func init() {
	logger := log.Default()
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("[center]")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	Logger = Log{
		logger:     logger,
		ErrPrefix:  "[ERROR]",
		WarnPrefix: "[WARN]",
		InfoPrefix: "[INFO]",
	}
}

type Log struct {
	logger     *log.Logger
	ErrPrefix  string
	WarnPrefix string
	InfoPrefix string
}

func (l *Log) Error() *Log {
	l.logger.SetPrefix(l.ErrPrefix)
	return l
}

func (l *Log) Warn() *Log {
	l.logger.SetPrefix(l.WarnPrefix)
	return l
}
func (l *Log) Info() *Log {
	l.logger.SetPrefix(l.InfoPrefix)
	return l
}

func (l *Log) Println(v ...any) {
	l.logger.Println(v)
}
