package data

import "time"

type ShortcutRun struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
	//超时时间
	Timeout time.Duration `json:"timeout"`
	//仅运行，不处理结果
	JustRun  bool   `json:"justRun"`
	Payload  string `json:"payload"`
	RealTime bool   `json:"realTime"`
}

type ShortcutRunResp struct {
	Ok     bool
	Err    string
	StdOut string
	StdErr string
}

type RealTimeShortcutRunResp struct {
	Ok    bool
	Err   string
	RunId string
}

type RealTimeShortcutOutPut struct {
	// 0 : stdout
	// 1 : stderr
	Type    int    `json:"type"`
	Payload string `json:"payload"`
	RunId   string `json:"runId"`
}

type CreateScript struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CreateScriptResp struct {
	Name    string `json:"name"`
	AbsPath string `json:"absPath"`
}

const (
	ONE_LINE = iota
	SCRIPT
)
