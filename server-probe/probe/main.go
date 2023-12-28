package main

import (
	"config-manager/probe/config"
	"config-manager/probe/connect"
	fileservice "config-manager/probe/service/fs"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
)

func main() {

	path := flag.String("c", "./probe/conf.yaml", "config file path")
	flag.Parse()

	//解析配置文件
	err := config.LoadConf(*path)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go connect.InitConnect(config.Conf.CenterAddr, ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	fileservice.ClearFileCache()

	log.Println("Server exiting")

}
