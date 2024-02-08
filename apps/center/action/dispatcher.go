package action

import (
	"apps/center/model"
	"apps/common/utils"
	"errors"
	"time"
)

type DispatchStatus int8

const (
	Ready DispatchStatus = iota
	Abort
	Done
)

type Dispatcher struct {
	Id string
	*Runner
	// pre node output
	Out        string
	Ok         bool
	Status     DispatchStatus
	CurNode    *model.Node
	PreNode    *model.Node
	ChainId    string
	CreateTime time.Time
	shortcuts  map[string]model.Shortcut
	nodes      map[string]model.Node
	edges      []model.Edge
	bindings   map[string]model.ShortcutNodeBinding
}

type ChainExecParams struct {
	Shortcuts []model.Shortcut
	Nodes     []model.Node
	Edges     []model.Edge
	Bindings  []model.ShortcutNodeBinding
	ChainId   string
}

var (
	ErrNoMoreNode = errors.New("no more node")
	ErrAbort      = errors.New("abort")
)

func (d *Dispatcher) Next() error {

	if d.Status == Done {
		return ErrNoMoreNode
	}

	if d.Status == Abort {
		return ErrAbort
	}

	// 如果当前应该执行节点为空，直接返回
	if d.CurNode == nil {
		d.Status = Done
		return ErrNoMoreNode
	}

	binding, ok := d.bindings[d.CurNode.Id]

	if !ok {
		d.Ok = true
		d.Out = "no binding shortcut"
	} else {
		sc := d.shortcuts[binding.ShortcutId]
		out, ok := d.Run(sc)
		d.Ok = ok
		d.Out = out
	}

	d.PreNode = d.CurNode

	// 如果执行失败，ok == false，找到失败对应的边，跳转到失败的节点
	if !d.Ok {
		//find fail edge
		node, ok := d.FailedThenNode()

		if ok {
			d.CurNode = node
			return nil
		}

		// 如果没有找到失败对应的边，直接返回
		d.CurNode = nil
		d.Status = Done
		return nil
	}

	// 如果执行成功，ok == true，找到成功对应的边，跳转到成功的节点
	node, ok := d.SuccessThenNode()

	if ok {
		d.CurNode = node
		return nil
	}

	// 如果没有找到成功对应的边，直接返回
	d.CurNode = nil
	d.Status = Done
	return nil
}

func (d *Dispatcher) Abort() {
	d.Status = Abort
}

func (d *Dispatcher) CurShortcut() (*model.Shortcut, bool) {

	if d.CurNode == nil {
		return nil, false
	}

	binding, ok := d.bindings[d.CurNode.Id]

	if !ok {
		return nil, false
	}

	sc := d.shortcuts[binding.ShortcutId]
	return &sc, true
}

func (d *Dispatcher) CurNodeId() (*string, bool) {
	if d.CurNode == nil {
		return nil, false
	}
	return &d.CurNode.Id, true
}

func (d *Dispatcher) GetCurNode() (*model.Node, bool) {
	return d.CurNode, d.CurNode != nil
}

func (d *Dispatcher) SuccessThenNode() (*model.Node, bool) {

	if d.CurNode == nil {
		return nil, false
	}

	for i := range d.edges {
		if d.edges[i].SourceId == d.CurNode.Id && d.edges[i].Type == model.SUCCESS {
			node := d.nodes[d.edges[i].TargetId]
			return &node, true
		}
	}
	return nil, false
}

func (d *Dispatcher) FailedThenNode() (*model.Node, bool) {

	if d.CurNode == nil {
		return nil, false
	}

	for i := range d.edges {
		if d.edges[i].SourceId == d.CurNode.Id && d.edges[i].Type == model.FAILED {
			node := d.nodes[d.edges[i].TargetId]
			return &node, true
		}
	}
	return nil, false
}

func NewDispatcher(p ChainExecParams) (*Dispatcher, error) {

	scMap, err := utils.SliceToMap[string, model.Shortcut](p.Shortcuts, "Id")
	if err != nil {
		return nil, err
	}

	nodeMap, err := utils.SliceToMap[string, model.Node](p.Nodes, "Id")
	if err != nil {
		return nil, err
	}

	bindingMap, err := utils.SliceToMap[string, model.ShortcutNodeBinding](p.Bindings, "NodeId")

	if err != nil {
		return nil, err
	}

	var root *model.Node

	for i := range p.Nodes {
		if p.Nodes[i].Root {
			root = &p.Nodes[i]
			break
		}
	}

	if root == nil {
		return nil, errors.New("no root node")
	}

	return &Dispatcher{
		Id:         utils.UUID(),
		Runner:     &Runner{},
		CurNode:    root,
		shortcuts:  scMap,
		nodes:      nodeMap,
		edges:      p.Edges,
		bindings:   bindingMap,
		ChainId:    p.ChainId,
		CreateTime: time.Now(),
		Status:     Ready,
	}, nil
}
