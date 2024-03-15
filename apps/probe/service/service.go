package service

import "apps/probe/service/sys"

var (
	SysInfoService *sys.InfoService
)

func Init() {
	SysInfoService = &sys.InfoService{}
}
