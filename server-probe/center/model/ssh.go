package model

type SSHInfo struct {
	Id       string  `gorm:"primaryKey" json:"id"`
	AddrPort string  `json:"addrPort"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	GroupId  *string `json:"groupId"`
}

type SSHGroup struct {
	Id     string  `gorm:"primaryKey" json:"id"`
	Name   string  `json:"name"`
	Parent *string `json:"parent"`
}
