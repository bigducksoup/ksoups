package handler

import (
	"apps/center/api/response"
	"apps/center/global"
	"apps/center/model"
	"github.com/gin-gonic/gin"
)

func OnlineNode(c *gin.Context) {

	var nodes []model.ProbeInfo
	global.DB.Find(&nodes)
	c.JSON(200, response.Success(nodes))

}
