package api

import (
	"apps/center/api/middleware"
	"apps/center/api/routers"
	"apps/center/api/ws"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func InitApiServer(port string, ctx context.Context) {

	//初始化websocket
	ws.Init()

	gin.SetMode(gin.ReleaseMode)
	e := gin.New()

	e.Use(gin.Recovery())

	e.Use(middleware.LogMiddleWare())
	//设置跨域中间件
	e.Use(middleware.Cors())

	//设置路由
	routers.SetUpRouters(e)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	go func() {
		// 服务连接
		log.Println("gin listening port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exited")
	}()

}
