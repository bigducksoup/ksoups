package routers

import (
	"config-manager/center/apiserver/api"

	"github.com/gin-gonic/gin"
)

func SetUpRouters(engine *gin.Engine) {

	fileGroup := engine.Group("/file")
	{
		fileGroup.GET("/read", api.FileRead)
	}

	dirGroup := engine.Group("/dir")
	{
		dirGroup.GET("/read", api.DirRead)
	}

	infoGroup := engine.Group("/info")
	{
		infoGroup.GET("/nodes", api.OnlineNode)
	}

}
