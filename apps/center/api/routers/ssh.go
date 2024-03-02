package routers

import (
	"apps/center/api/handler"
	"github.com/gin-gonic/gin"
)

func InitSSHRouter(e *gin.RouterGroup, middlewares ...gin.HandlerFunc) {

	sshRouters := e.Group("/ssh")

	sshRouters.Use(middlewares...)

	sshRouters.GET("/group/content", handler.GetSSHInfo)

	sshRouters.PUT("/info/save", handler.SaveSSHInfo)

	sshRouters.PUT("/group/save", handler.NewGroup)

	sshRouters.POST("/info/update", handler.UpdateSSHInfo)

	sshRouters.DELETE("/info/delete", handler.DeleteSSHInfo)

	sshRouters.POST("/group/update", handler.UpdateGroupInfo)

	sshRouters.GET("/group/tree", handler.GroupTree)

	sshRouters.DELETE("/group/delete", handler.DeleteGroup)
}
