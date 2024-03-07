package middleware

import (
	"apps/center/api/response"
	"apps/center/api/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

var whiteList = map[string]any{
	"/api/auth/login": nil,
}

func AuthenticationMiddleWare() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		if ctx.RemoteIP() == "127.0.0.1" {
			ctx.Next()
			return
		}

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

func AuthorizationRoleMiddleWare(...string) gin.HandlerFunc {
	// TODO  logic
	return func(context *gin.Context) {

	}
}

func AuthorizationPermMiddleWare(...string) gin.HandlerFunc {
	// TODO  logic
	return func(context *gin.Context) {

	}
}
