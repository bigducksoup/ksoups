package ssh

import (
	"apps/center/api/vo"
	"apps/center/model"
	"gorm.io/gorm"
)

type CRUDService struct {
	Db *gorm.DB
}

// GroupTree 递归查询一个组（SSH）里的子组
func (c *CRUDService) GroupTree(rootId string) []vo.GroupTreeItem {

	var group []model.SSHGroup
	c.Db.Where("parent = ?", rootId).Find(&group)

	var groupTree []vo.GroupTreeItem

	if len(group) == 0 {
		return groupTree
	}

	for _, v := range group {
		groupTree = append(groupTree, vo.GroupTreeItem{
			Id:       v.Id,
			Name:     v.Name,
			Children: c.GroupTree(v.Id),
		})
	}

	return groupTree
}

// DeleteGroup 递归删除一个组（SSH）里的所有内容
func (c *CRUDService) DeleteGroup(groupId string) error {

	var childGroup []model.SSHGroup

	err := c.Db.Where("parent = ?", groupId).Find(&childGroup).Error

	if err != nil {
		return err
	}

	for _, v := range childGroup {
		if err := c.DeleteGroup(v.Id); err != nil {
			return err
		}
	}

	err = c.Db.Where("id = ?", groupId).Delete(&model.SSHGroup{}).Error

	if err != nil {
		return err
	}

	err = c.Db.Where("group_id = ?", groupId).Delete(&model.SSHInfo{}).Error

	return err
}
