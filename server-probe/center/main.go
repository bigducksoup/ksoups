package main

import (
	"config-manager/center/apiserver"
	_ "config-manager/center/server"
)

func main() {

	apiserver.InitApiServer()

}
