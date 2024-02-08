package shortcut

import (
	"apps/center/model"
	"apps/center/server"
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"encoding/json"

	"gorm.io/gorm"
)

type CRUDService struct {
	Db *gorm.DB
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

	tx := c.Db.Begin()
	tx.Where("id = ?", id).Delete(&model.Shortcut{})
	tx.Where("shortcut_id = ?", id).Delete(&model.ShortcutNodeBinding{})
	return tx.Commit().Error
}

func (c *CRUDService) CreateScriptFile(probeId string, name string, content string) (absPath *string, err error) {

	d := data.CreateScript{
		Name:    name,
		Content: content,
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probeId, d, message.CREATE_SCRIPT)

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
