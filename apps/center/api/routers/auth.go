package routers

import (
	"apps/center/api/handler"
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(e *gin.RouterGroup, middlewares ...gin.HandlerFunc) {
	// 鉴权相关
	authGroup := e.Group("/auth")
	authGroup.Use(middlewares...)
	{
		// 登录
		authGroup.POST("/login", handler.Login)

		// 检查session
		authGroup.POST("/check_login", handler.CheckLogin)

	}
}
