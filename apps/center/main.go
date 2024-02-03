package main

import (
	"config-manager/center/api"
	"config-manager/center/config"
	"config-manager/center/db"
	"config-manager/center/global"
	"config-manager/center/server"
	"config-manager/center/service"
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
