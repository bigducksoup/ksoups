package routers

import (
	"apps/center/api/handler"
	"github.com/gin-gonic/gin"
)

func InitShortcutRouter(e *gin.RouterGroup, middlewares ...gin.HandlerFunc) {

	// 快捷指令相关
	shortcutGroup := e.Group("/shortcut")
	shortcutGroup.Use(middlewares...)
	{

		// 创建快捷指令
		shortcutGroup.POST("/create", handler.ShortcutCreate)

		// 列出快捷指令
		shortcutGroup.GET("/list", handler.ListShortcuts)

		// 运行快捷指令
		shortcutGroup.POST("/run", handler.RunShortcut)

		// 实时快捷指令
		shortcutGroup.POST("/realtime/run", handler.RealTimeRunShortcut)

		// 删除快捷指令
		shortcutGroup.DELETE("/delete", handler.DeleteShortcut)

		// 快捷指令分组 根据probeId
		shortcutGroup.GET("/group", handler.ShortcutGroup)

		// 更新快捷指令
		shortcutGroup.POST("/update", handler.UpdateShortcut)

		// 快捷指令运行历史
		shortcutGroup.GET("/run/history", handler.ShortcutRunHistory)
	}
}
