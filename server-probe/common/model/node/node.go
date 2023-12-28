package node

import "time"

type Node struct {
	Id      string    `gorm:"primaryKey" json:"id"`
	Name    string    `gorm:"size:20" json:"name"`
	Address string    `gorm:"size:20" json:"addr"`
	RegTime time.Time `json:"regTime"`
	Online  bool      `json:"online"`
}