package shortcut

import (
	"apps/center/global"
	"apps/center/model"
	v2 "apps/center/service/chain/v2"
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"encoding/json"
	"gorm.io/gorm"
	"slices"
)

type CRUDService struct {
	Db           *gorm.DB
	GraphService *v2.GraphService
}

func (c *CRUDService) SaveShortcut(sc *model.Shortcut) error {

	if sc.Type == model.SCRIPT {
		absPath, err := c.CreateScriptFile(sc.ProbeId, sc.Name, sc.Payload)
		if err != nil {
			return err
		}
		sc.Payload = *absPath
	}

	return c.Db.Create(sc).Error
}

func (c *CRUDService) ListShortcuts(probeId string) ([]model.Shortcut, error) {

	var scs []model.Shortcut

	err := c.Db.Where("probe_id = ?", probeId).Find(&scs).Error

	if err != nil {
		return nil, err
	}

	return scs, nil

}

func (c *CRUDService) UpdateShortcut(sc *model.Shortcut) error {

	var binding []model.ShortcutNodeBinding
	c.Db.Where("shortcut_id = ?", sc.Id).Find(&binding)

	if len(binding) > 0 {

		var nodeIds []string
		for _, nodeBinding := range binding {
			nodeIds = append(nodeIds, nodeBinding.NodeId)
		}

		chainId := binding[0].ChainId
		graphData, err := c.GraphService.GetGraphData(chainId)

		if err != nil {
			return err
		}

		// update data.proto(shortcut)
		for _, nodeId := range nodeIds {
			for _, cell := range graphData {
				if cell.Data != nil && cell.ID == nodeId {
					cell.Data.Proto = model.Proto{
						ID:          sc.Id,
						Name:        sc.Name,
						Description: sc.Description,
						Type:        int64(sc.Type),
						CreateTime:  sc.CreateTime.String(),
						Timeout:     sc.Timeout,
						JustRun:     sc.JustRun,
						Payload:     sc.Payload,
						ProbeID:     sc.ProbeId,
					}
				}
			}
		}

		err = c.GraphService.UpdateGraphData(chainId, graphData)

		if err != nil {
			return err
		}

	}

	return c.Db.Save(sc).Error
}

func (c *CRUDService) ShortcutGroup() (map[string][]model.Shortcut, error) {

	var scs []model.Shortcut

	err := c.Db.Find(&scs).Error

	if err != nil {
		return nil, err
	}

	groups, err := utils.SliceGroupBy[string, model.Shortcut](scs, "ProbeId")

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *CRUDService) RemoveShortcut(id string) error {

	err := c.Db.Where("id = ?", id).Delete(&model.Shortcut{}).Error

	if err != nil {
		return err
	}

	var binding []model.ShortcutNodeBinding
	err = c.Db.Where("shortcut_id = ?", id).Find(&binding).Error

	if err != nil {
		return err
	}

	if len(binding) > 0 {

		var nodeIds []string
		for _, nodeBinding := range binding {
			nodeIds = append(nodeIds, nodeBinding.NodeId)
		}

		chainId := binding[0].ChainId

		c.Db.Delete(&model.Node{}, nodeIds)
		graphData, err := c.GraphService.GetGraphData(chainId)

		if err != nil {
			return err
		}

		for _, nodeId := range nodeIds {
			graphData = slices.DeleteFunc(graphData, func(element model.Element) bool {
				// delete node
				if element.Data != nil && element.ID == nodeId {
					return true
				}
				// delete edge related to node
				if element.Data == nil && (element.Source.Cell == nodeId || element.Target.Cell == nodeId) {
					return true
				}
				return false
			})
		}

		err = c.GraphService.UpdateGraphData(chainId, graphData)

		if err != nil {
			return err
		}

	}

	return nil
}

func (c *CRUDService) CreateScriptFile(probeId string, name string, content string) (absPath *string, err error) {

	d := data.CreateScript{
		Name:    name,
		Content: content,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(probeId, d, message.CREATE_SCRIPT)

	if err != nil {
		return nil, err
	}

	resp := data.CreateScriptResp{}
	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.AbsPath, nil

}
