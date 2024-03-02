package routers

import (
	"apps/center/api/middleware"
	"apps/center/api/ws"
	"apps/center/static"
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

	// 静态文件
	engine.StaticFileFS("/", "dist/", statics)

	engine.StaticFileFS("/favicon.ico", "dist/favicon.ico", statics)

	engine.StaticFileFS("/logo.svg", "dist/logo.svg", statics)

	engine.StaticFileFS("/bg.jpg", "dist/bg.jpg", statics)

	engine.StaticFS("/assets", assets)

	engine.GET("/ws/:app", ws.HandleWS)

	// api path
	apiGroup := engine.Group("/api")
	apiGroup.Use(middleware.AuthMiddleWare())

	InitAuthRouter(apiGroup)
	InitFSRouter(apiGroup)
	InitShortcutRouter(apiGroup)
	InitChainRouter(apiGroup)
	InitInfoRouter(apiGroup)
	InitSSHRouter(apiGroup)

}
