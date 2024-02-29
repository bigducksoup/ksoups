package service

import (
	scrun "apps/center/action"
	"apps/center/global"
	"apps/center/server"
	"apps/center/service/chain"
	v2 "apps/center/service/chain/v2"
	"apps/center/service/fs"
	"apps/center/service/shortcut"
	"apps/center/service/ssh"
)

var (
	ShortcutCRUD shortcut.CRUDService
	ShortcutRUN  shortcut.RUNService
)

var (
	ChainCRUD   chain.CRUDService
	ChainEXECV2 v2.ExecService
	ChainLOG    chain.LogService
	ChainINFO   v2.InfoService
	Graph       v2.GraphService
)

var (
	SSHCRUD ssh.CRUDService
)

var (
	FS_OPERATION fs.OperationService
)

func Init() {

	Graph = v2.GraphService{
		Db: global.DB,
	}

	ShortcutCRUD = shortcut.CRUDService{Db: global.DB, GraphService: &Graph}
	ShortcutRUN = shortcut.RUNService{
		Runner: scrun.Runner{},
		Db:     global.DB,
	}
	ChainLOG = chain.LogService{Db: global.DB}
	ChainCRUD = chain.CRUDService{Db: global.DB}
	ChainEXECV2 = v2.ExecService{
		CRUDService: &ChainCRUD,
	}

	ChainINFO = v2.InfoService{
		Db: global.DB,
	}

	SSHCRUD = ssh.CRUDService{Db: global.DB}
	FS_OPERATION = fs.OperationService{ServerCtx: &(server.Ctx)}

}
