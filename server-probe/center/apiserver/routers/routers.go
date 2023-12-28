package routers

import (
	"config-manager/center/apiserver/api"
	"config-manager/center/static"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRouters(engine *gin.Engine) {

	// var DistFS embed.FS
	statics := http.FS(static.DistFS)

	asset, err := fs.Sub(static.DistFS, "dist/assets")

	if err != nil {
		panic(err)
	}

	assets := http.FS(asset)

	engine.StaticFileFS("/", "dist/", statics)
	engine.StaticFileFS("/favicon.ico", "dist/favicon.ico", statics)
	engine.StaticFileFS("/logo.svg", "dist/logo.svg", statics)
	engine.StaticFS("/assets", assets)

	// engine.StaticFile("/", "static/dist/index.html")
	// engine.StaticFile("/favicon.ico", "static/dist/favicon.ico")
	// engine.Static("/assets", "static/dist/assets")

	apiGroup := engine.Group("/api")
	// use auth middleware
	//apiGroup.Use(middleware.AuthMiddleWare())
	{
		fileGroup := apiGroup.Group("/file")
		{
			fileGroup.GET("/read", api.FileRead)
			fileGroup.POST("/modify", api.FileModify)
			fileGroup.POST("/create", api.FileCreate)
		}

		dirGroup := apiGroup.Group("/dir")
		{
			dirGroup.GET("/read", api.DirRead)
		}

		infoGroup := apiGroup.Group("/info")
		{
			infoGroup.GET("/nodes", api.OnlineNode)
		}

		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/login", api.Login)
			authGroup.POST("/check_login", api.CheckLogin)
		}

	}
}
