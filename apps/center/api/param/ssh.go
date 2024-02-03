package param

type SSHInfoParams struct {
	GroupId  string `json:"groupId"`
	AddrPort string `json:"addrPort"`
	Username string `json:"username"`
	Password string `json:"password"`
}
