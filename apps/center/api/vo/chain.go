package vo

import "apps/center/model"

type ChainInfo struct {
	Chain model.Chain  `json:"chain"`
	Nodes []NodeVO     `json:"nodes"`
	Edges []model.Edge `json:"edges"`
}

type NodeVO struct {
	model.Node      `json:"node"`
	Shortcut        *model.Shortcut `json:"shortcut"`
	SuccessThenId   *string         `json:"successThenId"`
	SuccessThenName *string         `json:"successThenName"`
	FailThenId      *string         `json:"failThenId"`
	FailThenName    *string         `json:"failThenName"`
}

type NodeDetail struct {
	Node     *model.Node     `json:"node"`
	Shortcut *model.Shortcut `json:"shortcut"`
}

// DispatchStatus 调度情况汇总，包含调度日志，当前节点，成功节点，失败节点，前置节点
type DispatchStatus struct {
	DispatchLog *model.DispatchLog   `json:"dispatchLog"`
	CurNode     *NodeDetail          `json:"curNode"`
	SuccessThen *NodeDetail          `json:"successThen"`
	FailThen    *NodeDetail          `json:"failThen"`
	PreNodes    []*model.NodeExecLog `json:"preNodes"`
}
