package model

import "time"

type Chain struct {
	Id          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	CreateTime  time.Time `json:"createTime"`
	Description string    `json:"description"`
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
}

type DispatchStatus int8

const (
	DispatchStatusRunning DispatchStatus = iota
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
}

type NodeExecLog struct {
	Id         string    `gorm:"primaryKey" json:"id"`
	DispatchId string    `json:"dispatchId"`
	NodeId     string    `json:"nodeId"`
	Ok         bool      `json:"ok"`
	Out        string    `json:"out"`
	CreateTime time.Time `json:"createTime"`
}

type NodeExecDetail struct {
	CreateTime   time.Time `json:"createTime"`
	Ok           bool      `json:"ok"`
	Out          string    `json:"out"`
	NodeName     string    `json:"nodeName"`
	ShortcutName *string   `json:"shortcutName"`
	Payload      *string   `json:"payload"`
	ShortcutType *int      `json:"shortcutType"`
}
