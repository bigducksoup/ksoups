package main

import (
	"config-manager/probe/connect"
	fileservice "config-manager/probe/service/fileService"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {

	ctx, cancle := context.WithCancel(context.Background())

	defer cancle()

	go connect.InitConnect("127.0.0.1:9999", ctx)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	fileservice.ClearFileCache()

	log.Println("Server exiting")

}
