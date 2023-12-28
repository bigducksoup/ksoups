package main

import (
	"config-manager/center/apiserver"
	"config-manager/center/config"
	"config-manager/center/db"
	"config-manager/center/server"
	"flag"
)

func main() {

	path := flag.String("c", "./center/conf.yaml", "config file path")
	flag.Parse()

	//解析配置文件
	err := config.LoadConf(*path)
	if err != nil {
		panic(err)
	}

	hp := config.Conf.Api.Port
	sp := config.Conf.Center.Port

	db.InitDB()

	server.Start(sp)
	apiserver.InitApiServer(hp)

}
