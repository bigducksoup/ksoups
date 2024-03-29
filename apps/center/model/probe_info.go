package model

import (
	"time"
)

type ProbeInfo struct {
	Id      string    `gorm:"primaryKey" json:"id"`
	Name    string    `gorm:"size:20" json:"name"`
	Address string    `gorm:"size:20" json:"addr"`
	KeyId   uint      `json:"keyId"`
	RegTime time.Time `json:"regTime"`
	Online  bool      `json:"online"`
}

type RegisterKey struct {
	BaseModel
	Name         string `json:"name"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	PublicKeyMd5 string `json:"public_key_md5"`
}

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
