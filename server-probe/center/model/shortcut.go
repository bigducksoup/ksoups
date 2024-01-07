package model

import (
	"time"
)

const (
	ONE_LINE = iota
	SCRIPT
)

// Shortcut is  a shortcut just execute one line command
type Shortcut struct {
	Id          string    `json:"id" gorm:"primaryKey;size:20"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        int       `json:"type"`
	CreateTime  time.Time `json:"createTime"`
	//超时时间
	Timeout int64 `json:"timeout"`
	//仅运行，不处理结果
	JustRun bool   `json:"justRun"`
	Payload string `json:"payload"`
	ProbeId string `json:"probeId"`
}

type ShortcutExecLog struct {
	Id         string    `json:"id" gorm:"primaryKey;size:20"`
	ShortcutId string    `json:"shortcutId"`
	CreateTime time.Time `json:"createTime"`
	//执行结果
	Out string `json:"out"`
	//执行状态
	OK bool `json:"ok"`
	//执行时间
	ExecuteTime time.Time `json:"executeTime"`
}
