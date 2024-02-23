package action

import (
	"apps/center/model"
	"errors"
)

// Seeker seek for next shortcut(s) to run
type Seeker interface {
	SeekNext(preResults []*NodeRunResult) ([]*NodeRunJob, error)
	seekRoot() (*NodeRunJob, error)
}

type NormalSeeker struct {
	nodes     []model.Node
	edges     []model.Edge
	shortcuts []model.Shortcut
	bindings  []model.ShortcutNodeBinding
}

func NewNormalSeeker(nodes []model.Node, edges []model.Edge, shortcuts []model.Shortcut, bindings []model.ShortcutNodeBinding) *NormalSeeker {
	return &NormalSeeker{
		nodes:     nodes,
		edges:     edges,
		shortcuts: shortcuts,
		bindings:  bindings,
	}
}

// SeekNext if empty return []
func (n *NormalSeeker) SeekNext(preResults []*NodeRunResult) ([]*NodeRunJob, error) {

	var next []*NodeRunJob

	// if n.pre is nil, means root has not been found
	if len(preResults) == 0 {
		// find root
		root, err := n.seekRoot()
		if err != nil {
			return nil, err
		}
		return append(next, root), nil
	}

	var targetIds []string
	// find next by edge, source is n.pre
	for _, edge := range n.edges {
		for _, preResult := range preResults {

			// fail
			t := model.FAILED
			//success
			if preResult.RunResult.Ok {
				t = model.SUCCESS
			}
			if edge.SourceId == preResult.Node.Id && edge.Type == t {
				targetIds = append(targetIds, edge.TargetId)
			}
		}
	}

	if len(targetIds) == 0 {
		return next, nil
	}

	for _, targetId := range targetIds {
		node, err := n.findNode(targetId)
		if err != nil {
			continue
		}
		shortcut, err := n.findShortcutByNodeId(targetId)
		if err == nil {
			next = append(next, &NodeRunJob{
				Node:     node,
				Shortcut: shortcut,
			})
		}
	}

	return next, nil
}

func (n *NormalSeeker) seekRoot() (*NodeRunJob, error) {

	for _, node := range n.nodes {
		if node.Root {
			shortcut, err := n.findShortcutByNodeId(node.Id)
			if err != nil {
				return nil, err
			}
			return &NodeRunJob{
				Node:     &node,
				Shortcut: shortcut,
			}, nil
		}
	}
	return nil, errors.New("could not find a root")
}

func (n *NormalSeeker) findShortcutByNodeId(id string) (*model.Shortcut, error) {

	var shortcutId *string

	for _, binding := range n.bindings {
		if binding.NodeId == id {
			shortcutId = &binding.ShortcutId
			break
		}
	}

	if shortcutId == nil {
		return nil, errors.New("unable to find a shortcut that is binded to node")
	}

	for _, shortcut := range n.shortcuts {
		if shortcut.Id == *shortcutId {
			return &shortcut, nil
		}
	}

	return nil, errors.New("there is a shortcut to node binding, but no shortcut could be found")

}

func (n *NormalSeeker) findNode(id string) (*model.Node, error) {

	for _, node := range n.nodes {
		if node.Id == id {
			return &node, nil
		}
	}

	return nil, errors.New("no such node with id = " + id)

}
