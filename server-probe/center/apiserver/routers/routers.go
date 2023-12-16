package routers

import (
	"config-manager/center/apiserver/api"

	"github.com/gin-gonic/gin"
)

func SetUpRouters(engine *gin.Engine) {

	engine.StaticFile("/", "./static/dist/index.html")
	engine.StaticFile("/favicon.ico", "./static/dist/favicon.ico")
	engine.Static("/assets", "./static/dist/assets")

	apiGroup := engine.Group("/api")
	{
		fileGroup := apiGroup.Group("/file")
		{
			fileGroup.GET("/read", api.FileRead)
			fileGroup.POST("/modify", api.FileModify)
		}

		dirGroup := apiGroup.Group("/dir")
		{
			dirGroup.GET("/read", api.DirRead)
		}

		infoGroup := apiGroup.Group("/info")
		{
			infoGroup.GET("/nodes", api.OnlineNode)
		}
	}
}
