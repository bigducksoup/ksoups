package handler

import (
	"apps/center/global"
	"apps/center/model"
	"github.com/gin-gonic/gin"
)

func OnlineNode(c *gin.Context) {

	var nodes []model.ProbeInfo
	global.DB.Find(&nodes)
	c.JSON(200, nodes)

}
