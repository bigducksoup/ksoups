package action

import (
	"config-manager/center/model"
	"config-manager/common/utils"
	"errors"
)

type Dispatcher struct {
	Id string
	*Runner
	Out       string
	Ok        bool
	curNode   *model.Node
	shortcuts map[string]model.Shortcut
	nodes     map[string]model.Node
	edges     []model.Edge
	bindings  map[string]model.ShortcutNodeBinding
}

type ChainExecParams struct {
	Shortcuts []model.Shortcut
	Nodes     []model.Node
	Edges     []model.Edge
	Bindings  []model.ShortcutNodeBinding
}

var (
	ErrNoMoreNode = errors.New("no more node")
)

func (d *Dispatcher) Next() error {

	// 如果当前应该执行节点为空，直接返回
	if d.curNode == nil {
		return ErrNoMoreNode
	}

	binding, ok := d.bindings[d.curNode.Id]

	if !ok {
		d.Ok = false
		d.Out = "no binding shortcut"
		return nil
	}

	sc := d.shortcuts[binding.ShortcutId]

	out, ok := d.Run(sc)
	d.Ok = ok
	d.Out = out

	// 如果执行失败，ok == false，找到失败对应的边，跳转到失败的节点
	if !ok {
		//find fail edge
		for i := range d.edges {
			if d.edges[i].SourceId == d.curNode.Id && d.edges[i].Type == model.FAILED {
				node := d.nodes[d.edges[i].TargetId]
				d.curNode = &node

				return nil
			}
		}
		// 如果没有找到失败对应的边，直接返回
		d.curNode = nil
		return nil
	}

	// 如果执行成功，ok == true，找到成功对应的边，跳转到成功的节点
	for i := range d.edges {
		if d.edges[i].SourceId == d.curNode.Id && d.edges[i].Type == model.SUCCESS {
			node := d.nodes[d.edges[i].TargetId]
			d.curNode = &node
			return nil
		}
	}
	// 如果没有找到成功对应的边，直接返回
	d.curNode = nil
	return nil
}

func (d *Dispatcher) CurShortcut() (*model.Shortcut, bool) {

	if d.curNode == nil {
		return nil, false
	}

	binding, ok := d.bindings[d.curNode.Id]

	if !ok {
		return nil, false
	}

	sc := d.shortcuts[binding.ShortcutId]
	return &sc, true
}

func (d *Dispatcher) GetCurNodeId() (*string, bool) {
	if d.curNode == nil {
		return nil, false
	}
	return &d.curNode.Id, true
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

	var root model.Node

	for i := range p.Nodes {
		if p.Nodes[i].Root {
			root = p.Nodes[i]
			break
		}
	}

	return &Dispatcher{
		Id:        utils.UUID(),
		Runner:    &Runner{},
		curNode:   &root,
		shortcuts: scMap,
		nodes:     nodeMap,
		edges:     p.Edges,
		bindings:  bindingMap,
	}, nil
}
