package chain

import (
	"config-manager/center/model"
	"config-manager/common/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

type CRUDService struct {
	Db *gorm.DB
}

func (c *CRUDService) SaveChain(ch *model.Chain) error {
	return c.Db.Create(ch).Error
}

func (c *CRUDService) SaveNode(n *model.Node) error {
	return c.Db.Create(n).Error
}

type ConnectTwoNodesParams struct {
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
	ChainId  string `json:"chainId"`
	Type     int    `json:"type"`
}

// LinkNode Chain中连接两个节点。
//
// ConnectNode接受一个ConnectNodesParams结构体作为参数，该结构体包含要连接的源节点和目标节点的ID，以及其他必要的信息，如链ID和边类型。
// 该函数检查数据库中是否存在源节点和目标节点，并在找不到任何一个节点时返回错误。
// 它还检查是否已经存在具有相同源ID和类型的边，如果存在则返回错误。
// 如果所有检查都通过，函数将在数据库中创建一个新的边，并返回创建过程中发生的任何错误。
//
// 该函数不直接返回任何内容，但如果任何数据库操作失败，它会返回一个错误。
func (c *CRUDService) LinkNode(p ConnectTwoNodesParams) error {
	s := model.Node{}
	c1 := c.Db.First(&s, "id = ?", p.SourceId).RowsAffected
	t := model.Node{}
	c2 := c.Db.First(&t, "id = ?", p.TargetId).RowsAffected

	if c1 == 0 || c2 == 0 {
		return errors.New("Node with Id = " + p.SourceId + " or Node with Id = " + p.TargetId + " not found")
	}

	if s.ChainId != p.ChainId || t.ChainId != p.ChainId {
		return errors.New("Node with Id = " + p.SourceId + " and Node with Id = " + p.TargetId + " not in the same chain")
	}

	//检测相同类型的NODE是否已经存在
	err := c.Db.First(&model.Edge{}, "source_id = ? and type = ?", p.SourceId, p.Type).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Edge with SourceId = " + p.SourceId + " and Type = " + string(rune(p.Type)) + " already exists")
	}

	edge := model.Edge{
		Id:         utils.UUID(),
		SourceId:   p.SourceId,
		TargetId:   p.TargetId,
		CreateTime: time.Now(),
		ChainId:    p.ChainId,
		Type:       p.Type,
	}

	return c.Db.Create(&edge).Error
}

// BindShortcut 绑定快捷方式到节点。
//
// 接收节点ID和快捷方式ID作为参数。
// 如果找不到快捷方式或节点，或者节点已经绑定了快捷方式，则返回错误。
func (c *CRUDService) BindShortcut(nodeId string, shortcutId string) error {

	err := c.Db.First(&model.Shortcut{}, "id = ?", shortcutId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("ScriptShortcut with Id = " + shortcutId + " not found")
	}
	err = c.Db.First(&model.Node{}, "id = ?", nodeId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Node with Id = " + nodeId + " not found")
	}

	err = c.Db.First(&model.ShortcutNodeBinding{}, "node_id = ?", nodeId).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Node with Id = " + nodeId + " already binded")
	}

	oneLineToNode := model.ShortcutNodeBinding{
		Id:         utils.UUID(),
		NodeId:     nodeId,
		ShortcutId: shortcutId,
	}
	err = c.Db.Create(&oneLineToNode).Error
	if err != nil {
		return err
	}
	return nil
}

// SetRoot 设置节点为根节点。
// SetRoot接受节点ID作为参数，并将该节点设置为根节点。
func (c *CRUDService) SetRoot(nodeId string) error {

	node := model.Node{}

	rows := c.Db.First(&node, "id = ?", nodeId).RowsAffected

	if rows == 0 {
		return errors.New("Node with Id = " + nodeId + " not found")
	}

	//set other nodes in the same chain to false
	tx := c.Db.Begin()
	tx.Model(&model.Node{}).Where("chain_id = ?", node.ChainId).Update("root", false)
	node.Root = true
	tx.Save(&node)
	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

type ChainInfo struct {
	Chain model.Chain  `json:"chain"`
	Nodes []NodeVO     `json:"nodes"`
	Edges []model.Edge `json:"edges"`
}

type NodeVO struct {
	model.Node `json:"node"`
	Shortcut   *model.Shortcut `json:"shortcut"`
}

// ChainInfo 获取链信息。
// ChainInfo接受链ID作为参数，并返回一个ChainInfo结构体，该结构体包含链的信息，包括节点和边。
// 该函数首先获取链，然后获取链中的所有节点和边。
// 对于每个节点，它还获取与该节点绑定的快捷方式。
// 最后，它将所有信息组合到一个ChainInfo结构体中，并返回该结构体。
func (c *CRUDService) ChainInfo(chainId string) (ChainInfo, error) {

	chainInfo := ChainInfo{}

	//获取chain
	chain := model.Chain{
		Id: chainId,
	}
	err := c.Db.First(&chain).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return chainInfo, err
	}

	//获取nodes
	var nodes []model.Node
	err = c.Db.Find(&nodes, "chain_id = ?", chainId).Error

	//获取所有shortcut
	var nodeVOs []NodeVO

	for i := range nodes {
		//获取shortcut
		var shortcut model.Shortcut
		tx := c.Db.Table("shortcuts").Select("shortcuts.*").Joins("inner join shortcut_node_bindings on shortcuts.id = shortcut_node_bindings.shortcut_id").Where("shortcut_node_bindings.node_id = ?", nodes[i].Id).First(&shortcut)

		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			nodeVOs = append(nodeVOs, NodeVO{
				Node:     nodes[i],
				Shortcut: nil,
			})
			continue
		}

		nodeVOs = append(nodeVOs, NodeVO{
			Node:     nodes[i],
			Shortcut: &shortcut,
		})
	}

	//获取edges
	var edges []model.Edge
	err = c.Db.Find(&edges, "chain_id = ?", chainId).Error

	chainInfo.Chain = chain
	chainInfo.Nodes = nodeVOs
	chainInfo.Edges = edges

	return chainInfo, nil
}

func (c *CRUDService) ChainRoot(chainId string) (model.Node, error) {
	var node model.Node
	err := c.Db.First(&node, "chain_id = ? and root = ?", chainId, true).Error
	return node, err
}

// ChainList 获取所有链。
// ChainList不接受任何参数，并返回一个Chain结构体的切片。
func (c *CRUDService) ChainList() ([]model.Chain, error) {
	var chains []model.Chain
	err := c.Db.Find(&chains).Error
	return chains, err
}

func (c *CRUDService) GetNodesByChainId(chainId string) ([]model.Node, error) {
	var nodes []model.Node
	err := c.Db.Find(&nodes, "chain_id = ?", chainId).Error
	return nodes, err
}

func (c *CRUDService) GetEdgesByChainId(chainId string) ([]model.Edge, error) {
	var edges []model.Edge
	err := c.Db.Find(&edges, "chain_id = ?", chainId).Error
	return edges, err
}
