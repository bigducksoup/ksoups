package apiserver

import (
	"config-manager/center/apiserver/global"
	"config-manager/center/apiserver/middleware"
	"config-manager/center/apiserver/routers"

	"github.com/gin-gonic/gin"
)

func InitApiServer() {

	e := gin.Default()

	e.Use(middleware.LogMiddleWare())
	e.Use(middleware.Cors())

	routers.SetUpRouters(e)

	err := e.Run()
	if err != nil {
		global.Logger.Println(err)
		return
	}

}
