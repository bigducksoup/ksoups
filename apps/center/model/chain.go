package model

import "time"

type Chain struct {
	Id          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	CreateTime  time.Time `json:"createTime"`
	Description string    `json:"description"`
	OriginData  string    `json:"originData"`
}

// Edge type
const (
	SUCCESS = iota
	FAILED
)

type Node struct {
	Id          string `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	ChainId     string `json:"chainId"`
	Description string `json:"description"`
	Root        bool   `json:"root"`
}

type Edge struct {
	Id         string    `gorm:"primaryKey" json:"id"`
	SourceId   string    `json:"sourceId"`
	TargetId   string    `json:"targetId"`
	CreateTime time.Time `json:"createTime"`
	ChainId    string    `json:"chainId"`
	Type       int       `json:"type"`
}

type ShortcutNodeBinding struct {
	Id         string `gorm:"primaryKey"`
	NodeId     string `json:"nodeId"`
	ShortcutId string `json:"shortcutId"`
	ChainId    string `json:"chainId"`
}

type DispatchStatus int8

const (
	DispatchStatusReady DispatchStatus = iota
	DispatchStatusRunning
	DispatchStatusDone
	DispatchStatusFailed
	DispatchStatusAborted
)

type DispatchLog struct {
	Id         string         `gorm:"primaryKey" json:"id"`
	ChainId    string         `json:"chainId"`
	CreateTime time.Time      `json:"createTime"`
	Status     DispatchStatus `json:"status"`
	Done       bool           `json:"done"`
	ChainData  *string        `json:"chainData,omitempty"`
}

type NodeExecLog struct {
	Id           string    `gorm:"primaryKey" json:"id"`
	ChainId      string    `json:"chainId"`
	DispatchId   string    `json:"dispatchId"`
	NodeId       string    `json:"nodeId"`
	NodeName     *string   `json:"nodeName"`
	ShortcutId   *string   `json:"shortcutId"`
	ShortcutName *string   `json:"shortcutName"`
	ProbeId      *string   `json:"probeId"`
	ShortcutType *int      `json:"shortcutType"`
	Payload      *string   `json:"payload"`
	Ok           bool      `json:"ok"`
	StdOut       string    `json:"stdOut"`
	StdErr       string    `json:"stdErr"`
	CreateTime   time.Time `json:"createTime"`
}

type NodeExecDetail struct {
	CreateTime   time.Time `json:"createTime"`
	Ok           bool      `json:"ok"`
	Out          string    `json:"out"`
	NodeName     string    `json:"nodeName"`
	ShortcutName *string   `json:"shortcutName"`
	Payload      *string   `json:"payload"`
	ProbeId      *string   `json:"probeId"`
	ShortcutType *int      `json:"shortcutType"`
}
