package param

type NodeCreateParams struct {
	Name        string `json:"name" binding:"required"`
	ChainId     string `json:"chainId"`
	Description string `json:"description" binding:"required"`
}

type ChainCreateParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type BindShortcutToNodeParams struct {
	ShortcutId string `json:"shortcutId" binding:"required"`
	NodeId     string `json:"nodeId" binding:"required"`
}

type ConnectNodesParams struct {
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
	ChainId  string `json:"chainId"`
	Type     int    `json:"type"`
}
