package middleware

import (
	"config-manager/center/webapi/response"
	"config-manager/center/webapi/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

var whiteList = map[string]any{
	"/handler/auth/login": nil,
}

func AuthMiddleWare() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		path := ctx.FullPath()

		_, ok := whiteList[path]

		if ok {
			ctx.Next()
			return
		}

		sessionId := ctx.GetHeader("sid")
		if sessionId == "" {
			ctx.Abort()
			ctx.JSON(http.StatusOK, response.InvalidReqError())
			return
		}

		_, ok = session.GetSession(sessionId)

		if !ok {
			ctx.Abort()
			ctx.JSON(http.StatusOK, response.InvalidReqError())
			return
		}

		ctx.Next()
	}

}
