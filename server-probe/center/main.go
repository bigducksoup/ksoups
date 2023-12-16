package main

import (
	"config-manager/center/apiserver"
	"config-manager/center/server"
	"flag"
	"fmt"
)

func main() {

	//解析配置文件
	conf := flag.String("c", "config.yaml", "config file")
	flag.Parse()
	fmt.Println(*conf)

	server.Start()
	apiserver.InitApiServer()

}
