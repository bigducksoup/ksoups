package middleware

import (
	"github.com/gin-gonic/gin"
)

func LogMiddleWare() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		ctx.Next()
	}

}
