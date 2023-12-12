package middleware

import (
	"config-manager/center/apiserver/global"
	"github.com/gin-gonic/gin"
)

func LogMiddleWare() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		global.Logger.Println(ctx.Request.RemoteAddr)
		ctx.Next()
	}

}
