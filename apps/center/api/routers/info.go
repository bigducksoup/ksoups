package routers

import (
	"apps/center/api/handler"
	"github.com/gin-gonic/gin"
)

func InitInfoRouter(e *gin.RouterGroup, middlewares ...gin.HandlerFunc) {

	// 信息相关
	infoGroup := e.Group("/info")
	infoGroup.Use(middlewares...)
	{

		// 在线探针信息
		infoGroup.GET("/nodes", handler.OnlineNode)

	}
}
