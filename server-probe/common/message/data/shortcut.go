package data

import "time"

type ShortcutRun struct {
	Type int `json:"type"`
	//超时时间
	Timeout time.Duration `json:"timeout"`
	//仅运行，不处理结果
	JustRun bool   `json:"justRun"`
	Payload string `json:"payload"`
}

type ShortcutRunResp struct {
	Ok  bool
	Err string
	Out string
}

const (
	ONE_LINE = iota
	SCRIPT
)
