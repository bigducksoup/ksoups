package vo

// GroupContent 一个组（SSH）里的内容
// Payload model.SSHInfo || model.SSHGroup
type GroupContent struct {
	Type    int `json:"type"`
	Payload any `json:"payload"`
}

const (
	GROUP = iota
	INFO
)

type GroupTreeItem struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	Children []GroupTreeItem `json:"children"`
}
