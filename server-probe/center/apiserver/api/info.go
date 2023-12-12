package api

import (
	"config-manager/center/server"

	"github.com/gin-gonic/gin"
)

func OnlineNode(c *gin.Context) {

	nodeAddr := [](map[string]any){}

	server.Ctx.AddrProbe.Range(func(key, value any) bool {

		nodeAddr = append(nodeAddr, map[string]any{
			"addr": key,
			"id":   value.(*server.Probe).Id,
			"time": value.(*server.Probe).Time,
		})
		return true
	})

	c.JSON(200, nodeAddr)

}
