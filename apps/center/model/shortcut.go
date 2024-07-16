package model

import (
	"time"
)

const (
	ONE_LINE = iota
	SCRIPT
)

type ScriptRunState int

const (
	// start but script return err
	SCRIPT_RUN_ERR ScriptRunState = iota
	// try run but not start
	SCRIPT_RUN_FAIL
	// successfully run and got stdout
	SCRIPT_RUN_SUCCESS
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
	Args    string `json:"args"`
}

type ShortcutExecLog struct {
	// Id
	Id string `json:"id" gorm:"primaryKey;size:20"`
	// script id
	ShortcutId string `json:"shortcutId"`
	StdOut     string
	StdErr     string
	//执行状态
	State ScriptRunState `json:"state"`
	//执行时间
	ExecuteTime time.Time `json:"executeTime"`
	RunByChain  bool      `json:"runByChain"`
	ChainId     *string
	NodeId      *string
}
