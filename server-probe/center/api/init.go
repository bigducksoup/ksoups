package api

import (
	"config-manager/center/api/middleware"
	"config-manager/center/api/routers"
	"config-manager/center/api/ws"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func InitApiServer(port string) {

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

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
