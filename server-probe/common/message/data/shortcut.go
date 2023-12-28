package data

import "time"

type OneLineShortcutCreate struct {
	Name        string
	Command     string
	Timeout     time.Duration
	JustRun     bool
	Description string
}

type OneLineShortcutCreateResp struct {
	Ok         bool
	Error      string
	ShortcutId string
	CreateTime time.Time
}

type ScriptShortcutCreate struct {
	Name        string
	Path        string
	Timeout     time.Duration
	JustRun     bool
	Description string
}

type ScriptShortcutCreateResp struct {
	Ok         bool
	Error      string
	ShortcutId string
	CreateTime time.Time
}

const (
	ONE_LINE_SHORTCUT = 1
	SCRIPT_SHORTCUT   = 2
)

type ShortcutDelete struct {
	Type int    `json:"type"`
	Id   string `json:"id"`
}

type ShortcutDeleteResp struct {
	OK bool `json:"ok"`
}
