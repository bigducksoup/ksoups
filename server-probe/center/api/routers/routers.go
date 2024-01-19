package routers

import (
	"config-manager/center/api/handler"
	"config-manager/center/api/middleware"
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

	// 静态文件
	engine.StaticFileFS("/", "dist/", statics)

	engine.StaticFileFS("/favicon.ico", "dist/favicon.ico", statics)

	engine.StaticFileFS("/logo.svg", "dist/logo.svg", statics)

	engine.StaticFS("/assets", assets)

	// api path
	apiGroup := engine.Group("/api")
	// use auth middleware
	apiGroup.Use(middleware.AuthMiddleWare())
	{
		// 文件相关
		fileGroup := apiGroup.Group("/file")
		{
			// 读文件
			fileGroup.GET("/read", handler.FileRead)

			// 修改文件
			fileGroup.POST("/modify", handler.FileModify)

			// 创建文件
			fileGroup.POST("/create", handler.FileCreate)

		}

		// 目录相关
		dirGroup := apiGroup.Group("/dir")
		{
			// 读目录
			dirGroup.GET("/read", handler.DirRead)

			// 创建目录
			dirGroup.POST("/create", handler.DirCreate)

		}

		// 快捷指令相关
		shortcutGroup := apiGroup.Group("/shortcut")
		{

			// 创建快捷指令
			shortcutGroup.POST("/create", handler.ShortcutCreate)

			// 列出快捷指令
			shortcutGroup.GET("/list", handler.ListShortcuts)

			// 运行快捷指令
			shortcutGroup.POST("/run", handler.RunShortcut)

			// 删除快捷指令
			shortcutGroup.DELETE("/delete", handler.DeleteShortcut)

			// 快捷指令分组 根据probeId
			shortcutGroup.GET("/group", handler.ShortcutGroup)
		}

		// 链式指令调度相关
		chainGroup := apiGroup.Group("/chain")
		{

			// 链式指令信息
			chainGroup.GET("/info", handler.ChainInfo)

			// 链式指令列表
			chainGroup.GET("/list", handler.ChainList)

			// 创建链式指令
			chainGroup.POST("/create", handler.ChainCreate)

			// 创建链式指令节点
			chainGroup.POST("/node/create", handler.NodeCreate)

			// 删除指令节点
			chainGroup.DELETE("/node/delete", handler.NodeDelete)

			// 绑定快捷指令到节点
			chainGroup.POST("/node/bind/shortcut", handler.BindShortcut)

			// 解绑快捷指令
			chainGroup.POST("/node/unbind/shortcut", handler.UnBindShortcut)

			// 链接两个节点
			chainGroup.POST("/node/link", handler.LinkNodes)

			// 断开两个节点
			chainGroup.POST("/node/unlink", handler.UnLinkNodes)

			// 设置链式指令根节点
			chainGroup.PUT("/node/set/root", handler.SetChainRoot)

			// 链式指令日志
			chainGroup.GET("/exec/result", handler.ChainExecResult)

			chainGroup.PUT("/exec/dispatch/new", handler.NewDispatch)

			chainGroup.PUT("/exec/single/step/dispatch", handler.DoSingleStepDispatch)

			chainGroup.POST("/exec/all/step/dispatch", handler.DoAllDispatch)

			chainGroup.GET("/exec/history", handler.ChainExecLogHistory)

		}

		// 信息相关
		infoGroup := apiGroup.Group("/info")
		{

			// 在线探针信息
			infoGroup.GET("/nodes", handler.OnlineNode)

		}

		// 鉴权相关
		authGroup := apiGroup.Group("/auth")
		{
			// 登录
			authGroup.POST("/login", handler.Login)

			// 检查session
			authGroup.POST("/check_login", handler.CheckLogin)

		}

	}
}
