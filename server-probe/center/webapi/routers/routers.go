package routers

import (
	"config-manager/center/static"
	"config-manager/center/webapi/handler"
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

	apiGroup := engine.Group("/api")
	// use auth middleware
	//apiGroup.Use(middleware.AuthMiddleWare())
	{
		fileGroup := apiGroup.Group("/file")
		{
			fileGroup.GET("/read", handler.FileRead)
			fileGroup.POST("/modify", handler.FileModify)
			fileGroup.POST("/create", handler.FileCreate)
		}

		dirGroup := apiGroup.Group("/dir")
		{
			dirGroup.GET("/read", handler.DirRead)
			dirGroup.POST("/create", handler.DirCreate)
		}

		shortcutGroup := apiGroup.Group("/shortcut")
		{
			shortcutGroup.POST("/create", handler.ShortcutCreate)
			shortcutGroup.GET("/list", handler.ListShortcuts)
			shortcutGroup.POST("/run", handler.RunShortcut)
			shortcutGroup.DELETE("/delete", handler.DeleteShortcut)
		}
		chainGroup := apiGroup.Group("/chain")
		{
			chainGroup.GET("/info", handler.ChainInfo)
			chainGroup.GET("/list", handler.ChainList)
			chainGroup.POST("/create", handler.ChainCreate)
			chainGroup.POST("/node/create", handler.NodeCreate)
			chainGroup.POST("/node/bind/shortcut", handler.BindShortcut)
			chainGroup.POST("/node/link", handler.LinkNodes)
			chainGroup.PUT("/node/set/root", handler.SetChainRoot)
			chainGroup.PUT("/exec", handler.ChainExec)
			chainGroup.GET("/exec/log", handler.ChainExecLog)

		}

		infoGroup := apiGroup.Group("/info")
		{
			infoGroup.GET("/nodes", handler.OnlineNode)
		}

		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/login", handler.Login)
			authGroup.POST("/check_login", handler.CheckLogin)
		}

	}
}
