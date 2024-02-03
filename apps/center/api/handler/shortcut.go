package handler

import (
	"config-manager/center/api/param"
	"config-manager/center/api/response"
	"config-manager/center/global"
	"config-manager/center/model"
	"config-manager/center/service"
	"config-manager/common/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// ShortcutCreate 创建快捷方式
// Params: probeId 探针id，name 快捷方式名称，timeout 超时时间，justRun 是否只运行一次，payload 脚本内容，type 脚本类型，description 描述
// Response: 200，“success”， true
// Err: 200, “fail”， err
func ShortcutCreate(c *gin.Context) {
	p := param.CreateShortcutParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	ok := checkProbeId(p.ProbeId)

	if !ok {
		c.JSON(http.StatusOK, response.StringFail("can not find a probe with id = "+p.ProbeId))
		return
	}

	sc := model.Shortcut{}
	sc.Id = utils.UUID()
	sc.ProbeId = p.ProbeId
	sc.Name = p.Name
	sc.Timeout = p.Timeout
	sc.JustRun = p.JustRun
	sc.Payload = p.Payload
	sc.CreateTime = time.Now()
	sc.Type = p.Type
	sc.Description = p.Description

	err = service.ShortcutCRUD.SaveShortcut(&sc)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success[bool](true))
}

// ListShortcuts 列出指定探针的所有快捷方式
// Params: probeId 探针id
// Response: 200，“success”， []model.Shortcut
// Err: 200, “fail”， err
func ListShortcuts(c *gin.Context) {

	probeId, b := c.GetQuery("probeId")

	if !b {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	ok := checkProbeId(probeId)

	if !ok {
		c.JSON(http.StatusOK, response.StringFail("can not find a probe with id = "+probeId))
		return
	}

	shortcuts, err := service.ShortcutCRUD.ListShortcuts(probeId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(shortcuts))

}

// ShortcutGroup 快捷方式分组汇总
// Params: 无
// Response: 200，“success”， map[string][]model.Shortcut(key:探针id，value:快捷方式数组)
func ShortcutGroup(c *gin.Context) {

	groups, err := service.ShortcutCRUD.ShortcutGroup()

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(groups))

}

// RunShortcut 运行快捷方式
// Params: shortcutId 快捷方式id
// Response: 200，“success”， {ok: true, out: out}
// Err: 200, “fail”， {ok: false, out: err}
func RunShortcut(c *gin.Context) {

	scId, ok := c.GetQuery("shortcutId")
	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	out, err := service.ShortcutRUN.Run(scId)

	if err != nil {
		c.JSON(http.StatusOK, response.Success(gin.H{
			"ok":  false,
			"out": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"ok":  true,
		"out": out,
	}))
}

// ShortcutRunHistory 快捷方式运行历史
// Params: shortcutId 快捷方式id
// Response: 200，“success”， []model.ShortcutExecLog
// Err: 200, “fail”， err
func ShortcutRunHistory(c *gin.Context) {

	shortcutId, ok := c.GetQuery("shortcutId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	logs, err := service.ShortcutRUN.RunHistory(shortcutId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(logs))
}

// DeleteShortcut 删除快捷方式
// Params: shortcutId 快捷方式id
// Response: 200，“success”， true
// Err: 200, “fail”， err
func DeleteShortcut(c *gin.Context) {

	id, ok := c.GetQuery("shortcutId")

	if !(ok) {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err := service.ShortcutCRUD.RemoveShortcut(id)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

// UpdateShortcut 更新快捷方式
// Params: shortcutId 快捷方式id
// Response: 200，“success”， true
// Err: 200, “fail”， err
func UpdateShortcut(c *gin.Context) {
	var shortcut = model.Shortcut{}

	err := c.ShouldBindJSON(&shortcut)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ShortcutCRUD.UpdateShortcut(&shortcut)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())

}

func checkProbeId(probeId string) bool {
	//检测probe是否存在
	var count int64
	global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probeId).Count(&count)
	if count < 1 {
		return false
	}

	return true
}
