package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/api/vo"
	"apps/center/global"
	"apps/center/model"
	"apps/center/service"
	"apps/common/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SSHHandlerLoad(e *gin.RouterGroup) {

	sshRouters := e.Group("/ssh")

	sshRouters.GET("/group/content", GetSSHInfo)

	sshRouters.PUT("/info/save", SaveSSHInfo)

	sshRouters.PUT("/group/save", NewGroup)

	sshRouters.POST("/info/update", UpdateSSHInfo)

	sshRouters.DELETE("/info/delete", DeleteSSHInfo)

	sshRouters.POST("/group/update", UpdateGroupInfo)

	sshRouters.GET("/group/tree", GroupTree)

	sshRouters.DELETE("/group/delete", DeleteGroup)

}

// GetSSHInfo 获取组内信息（包括子组，和ssh信息）
// GET
// @Param groupId 组Id
// @Response []vo.GroupContent
func GetSSHInfo(c *gin.Context) {

	id, ok := c.GetQuery("groupId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	// response
	var groupContent []vo.GroupContent

	//检查 groupId 所对应的 ssh_group 是否存在

	// 检查是否存在对应 group
	tx := global.DB.Where("id = ?", id).Find(&model.SSHGroup{})
	rows := tx.RowsAffected
	if rows == 0 {
		c.JSON(http.StatusOK, response.StringFail("could not find ssh_group that id = "+id))
		return
	}

	var sshInfo []model.SSHInfo
	var sshGroup []model.SSHGroup

	global.DB.Where("group_id = ?", id).Find(&sshInfo)
	global.DB.Where("parent = ?", id).Find(&sshGroup)

	for _, group := range sshGroup {
		groupContent = append(groupContent, vo.GroupContent{
			Type:    vo.GROUP,
			Payload: group,
		})
	}

	for _, info := range sshInfo {
		info.Password = "******"
		groupContent = append(groupContent, vo.GroupContent{
			Type:    vo.INFO,
			Payload: info,
		})
	}

	c.JSON(http.StatusOK, response.Success(groupContent))
}

// SaveSSHInfo 保存ssh信息
// Post
// @Body param.SSHInfoParams
// @Response bool
func SaveSSHInfo(c *gin.Context) {

	params := param.SSHInfoParams{}
	err := c.ShouldBindJSON(&params)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	// 检查是否存在对应 group
	tx := global.DB.Where("id = ?", params.GroupId).Find(&model.SSHGroup{})
	rows := tx.RowsAffected
	if rows == 0 {
		c.JSON(http.StatusOK, response.StringFail("could not find ssh_group that id = "+params.GroupId))
		return
	}

	err = global.DB.Save(&model.SSHInfo{
		Id:       utils.UUID(),
		Username: params.Username,
		AddrPort: params.AddrPort,
		Password: params.Password,
		GroupId:  &params.GroupId,
	}).Error

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(true))
}

// NewGroup 创建组

func NewGroup(c *gin.Context) {

	name, n := c.GetQuery("name")
	parentId, p := c.GetQuery("parentId")

	if !(n && p) {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	// 检查是否存在对应 group
	tx := global.DB.Where("id = ?", parentId).Find(&model.SSHGroup{})
	rows := tx.RowsAffected
	if rows == 0 {
		c.JSON(http.StatusOK, response.StringFail("could not find ssh_group that id = "+parentId))
		return
	}

	group := model.SSHGroup{
		Id:     utils.UUID(),
		Name:   name,
		Parent: &parentId,
	}

	err := global.DB.Save(&group).Error

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(true))
}

// UpdateGroupInfo 更新组信息
// POST
// @Body model.SSHGroup
// @Response true or false
func UpdateGroupInfo(c *gin.Context) {

	var groupInfo model.SSHGroup

	err := c.ShouldBindJSON(&groupInfo)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = global.DB.Save(&groupInfo).Error

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(true))
}

func UpdateSSHInfo(c *gin.Context) {

	info := model.SSHInfo{}
	err := c.ShouldBindJSON(&info)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	// 检查是否存在对应 group
	tx := global.DB.Where("id = ?", *info.GroupId).Find(&model.SSHGroup{})
	rows := tx.RowsAffected
	if rows == 0 {
		c.JSON(http.StatusOK, response.StringFail("could not find ssh_group that id = "+*info.GroupId))
		return
	}

	err = global.DB.Updates(&info).Error

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(true))
}

func DeleteSSHInfo(c *gin.Context) {

	id, ok := c.GetQuery("infoId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err := global.DB.Where("id = ?", id).Delete(&model.SSHInfo{}).Error

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(true))

}

func GroupTree(c *gin.Context) {

	rootId, ok := c.GetQuery("rootId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	var root model.SSHGroup

	r := service.SSHCRUD.Db.Where("id = ?", rootId).Find(&root).RowsAffected

	if r == 0 {
		c.JSON(http.StatusOK, response.StringFail("could not find ssh_group that id = "+rootId))
		return
	}

	tree := service.SSHCRUD.GroupTree(rootId)

	rootItem := vo.GroupTreeItem{
		Id:       root.Id,
		Name:     root.Name,
		Children: tree,
	}

	c.JSON(http.StatusOK, response.Success([]vo.GroupTreeItem{rootItem}))
}

// DeleteGroup 删除一个组（SSH）里的所有内容
// DELETE
// @Param groupId 组Id
// @Response bool 是否删除成功
func DeleteGroup(c *gin.Context) {

	groupId, ok := c.GetQuery("groupId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err := service.SSHCRUD.DeleteGroup(groupId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(true))
}
