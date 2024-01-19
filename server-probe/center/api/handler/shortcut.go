package handler

import (
	"config-manager/center/api/param"
	"config-manager/center/api/response"
	"config-manager/center/global"
	"config-manager/center/model"
	"config-manager/center/service"
	"config-manager/common/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

func ShortcutGroup(c *gin.Context) {

	groups, err := service.ShortcutCRUD.ShortcutGroup()

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(groups))

}

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

func checkProbeId(probeId string) bool {
	//检测probe是否存在
	var count int64
	global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probeId).Count(&count)
	if count < 1 {
		return false
	}

	return true
}
