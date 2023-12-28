package api

import (
	"config-manager/center/global"
	"config-manager/common/model/node"

	"github.com/gin-gonic/gin"
)

func OnlineNode(c *gin.Context) {

	var nodes []node.Node
	global.DB.Find(&nodes)
	c.JSON(200, nodes)

}
