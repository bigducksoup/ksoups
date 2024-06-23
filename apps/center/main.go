package main

import (
	"apps/center/api"
	"apps/center/config"
	"apps/center/db"
	"apps/center/global"
	"apps/center/server"
	"apps/center/service"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	path := flag.String("c", "./center/center.yaml", "config file path")
	flag.Parse()

	//解析配置文件
	conf, err := config.LoadConf(*path)
	if err != nil {
		panic(err)
	}
	global.Conf = *conf

	apiPort := global.Conf.WebApi.Port
	rpcPort := global.Conf.Center.Port

	db.InitDB()

	ctx, cancel := context.WithCancel(context.Background())

	cs, csCancel := context.WithCancel(ctx)
	centerServer := server.InitCenterServer(rpcPort, cs)
	global.CenterServer = centerServer

	defer csCancel()

	// init service
	service.Init()
	apiContext, apiCancel := context.WithCancel(ctx)
	defer apiCancel()

	// init api
	api.InitApiServer(apiPort, apiContext)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down everything ...")
	cancel()
	time.Sleep(2 * time.Second)

}
