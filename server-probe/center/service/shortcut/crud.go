package shortcut

import (
	"config-manager/center/model"
	"config-manager/common/utils"

	"gorm.io/gorm"
)

type CRUDService struct {
	Db *gorm.DB
}

func (c *CRUDService) SaveShortcut(sc *model.Shortcut) error {
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
