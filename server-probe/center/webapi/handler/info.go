package handler

import (
	"config-manager/center/global"
	"config-manager/center/model"
	"github.com/gin-gonic/gin"
)

func OnlineNode(c *gin.Context) {

	var nodes []model.ProbeInfo
	global.DB.Find(&nodes)
	c.JSON(200, nodes)

}
