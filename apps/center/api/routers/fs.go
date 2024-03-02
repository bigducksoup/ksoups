package routers

import (
	"apps/center/api/handler"
	"github.com/gin-gonic/gin"
)

func InitFSRouter(e *gin.RouterGroup, middlewares ...gin.HandlerFunc) {

	// 文件相关
	fileGroup := e.Group("/file")
	fileGroup.Use(middlewares...)
	{
		// 读文件
		fileGroup.GET("/read", handler.FileRead)

		// 修改文件
		fileGroup.POST("/modify", handler.FileModify)

		// 创建文件
		fileGroup.POST("/create", handler.FileCreate)

	}

	// 目录相关
	dirGroup := e.Group("/dir")
	dirGroup.Use(middlewares...)
	{
		// 读目录
		dirGroup.GET("/read", handler.DirRead)

		// 创建目录
		dirGroup.POST("/create", handler.DirCreate)

	}
}
