package main

import (
	"apps/center/api"
	"apps/center/config"
	"apps/center/db"
	"apps/center/global"
	"apps/center/server"
	"apps/center/service"
	"flag"
)

func main() {

	path := flag.String("c", "./center/conf.yaml", "config file path")
	flag.Parse()

	//解析配置文件
	conf, err := config.LoadConf(*path)
	if err != nil {
		panic(err)
	}
	global.Conf = *conf

	hp := global.Conf.WebApi.Port
	sp := global.Conf.Center.Port
	db.InitDB()
	service.Init()
	server.Start(sp)

	api.InitApiServer(hp)

}
