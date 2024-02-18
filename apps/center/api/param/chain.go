package param

import "encoding/json"

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

func UnmarshalChainAllDataParams(data []byte) (ChainAllDataParams, error) {
	var r ChainAllDataParams
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ChainAllDataParams) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ChainAllDataParams struct {
	ChainID    string `json:"chainId"`
	Nodes      []Node `json:"nodes"`
	Edges      []Edge `json:"edges"`
	OriginData string `json:"originData"`
}

type Edge struct {
	ID       string `json:"id"`
	SourceID string `json:"sourceId"`
	TargetID string `json:"targetId"`
	Type     int    `json:"type"`
}

type Node struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Shortcut    Shortcut `json:"shortcut"`
}

type Shortcut struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	CreateTime  string `json:"createTime"`
	Timeout     int64  `json:"timeout"`
	JustRun     bool   `json:"justRun"`
	Payload     string `json:"payload"`
	ProbeID     string `json:"probeId"`
}
