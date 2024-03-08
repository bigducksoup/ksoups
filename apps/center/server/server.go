package server

import (
	"apps/center/server/core"
)

func MakeServer(options core.CenterServerOptions) *core.CenterServer {
	s := core.CreateCenterServer(options)
	return s
}
