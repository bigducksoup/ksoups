package main

import (
	"apps/probe/config"
	"apps/probe/connect"
	"apps/probe/handlers"
	fileservice "apps/probe/service/fs"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	log.SetFlags(log.Lshortfile)
	log.SetPrefix("[Probe]")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	path := flag.String("c", "./probe/probe.yaml", "config file path")
	flag.Parse()

	//解析配置文件
	err := config.LoadConf(*path)
	if err != nil {
		panic(err)
	}

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	registerInfo, err := connect.GenerateRegisterInfo(config.Conf.PublicKey, config.Conf.Name)

	if err != nil {
		panic(err)
	}

	connect.InitProbe(connect.ProbeOptions{
		Address:           config.Conf.CenterAddr,
		Ping:              false,
		PingInterval:      10 * time.Second,
		RequestTimeOut:    10 * time.Second,
		RegisterInfo:      registerInfo,
		Context:           context,
		Reconnect:         false,
		ReconnectGapTime:  5 * time.Second,
		MaxReconnectCount: 100,
		DataHandlers:      handlers.DataHandlePolicy,
	})
	connect.ProbeInstance.StartWorking()

	// go base.InitConnect(config.Conf.CenterAddr, ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	fileservice.ClearFileCache()

	log.Println("Server exiting")

}
