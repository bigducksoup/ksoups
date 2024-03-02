package routers

import (
	"apps/center/api/handler"
	"github.com/gin-gonic/gin"
)

func InitChainRouter(e *gin.RouterGroup, middlewares ...gin.HandlerFunc) {
	// 链式指令调度相关
	chainGroup := e.Group("/chain")
	chainGroup.Use(middlewares...)
	{

		v2 := chainGroup.Group("/v2")
		{
			v2.PUT("/machine/new", handler.NewExecMachine)
			v2.PUT("/machine/exec/all", handler.MachineExecAll)
			v2.PUT("/machine/exec/one", handler.MachineExecOne)

			v2.GET("/exec/list", handler.ExecList)
			v2.GET("/exec/detail", handler.ExecDetail)

		}

		chainGroup.POST("/load", handler.ChainLoadFromAllData)

		// 链式指令信息
		chainGroup.GET("/info", handler.ChainInfo)

		// 链式指令列表
		chainGroup.GET("/list", handler.ChainList)

		// 创建链式指令
		chainGroup.POST("/create", handler.ChainCreate)

		// 删除链式指令
		chainGroup.DELETE("/delete", handler.ChainDelete)

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

		chainGroup.GET("/exec/history", handler.ChainExecLogHistory)

	}
}
