package service

import (
	scrun "apps/center/action"
	"apps/center/global"
	"apps/center/server"
	"apps/center/service/chain"
	"apps/center/service/fs"
	"apps/center/service/shortcut"
	"apps/center/service/ssh"
)

var (
	ShortcutCRUD shortcut.CRUDService
	ShortcutRUN  shortcut.RUNService
)

var (
	ChainCRUD chain.CRUDService
	ChainEXEC chain.ExecService
	ChainLOG  chain.LogService
)

var (
	SSHCRUD ssh.CRUDService
)

var (
	FS_OPERATION fs.OperationService
)

func Init() {
	ShortcutCRUD = shortcut.CRUDService{Db: global.DB}
	ShortcutRUN = shortcut.RUNService{
		Runner: scrun.Runner{},
		Db:     global.DB,
	}
	ChainLOG = chain.LogService{Db: global.DB}
	ChainCRUD = chain.CRUDService{Db: global.DB}
	ChainEXEC = chain.ExecService{
		ChainCRUD: &ChainCRUD,
		Log:       &ChainLOG,
	}

	SSHCRUD = ssh.CRUDService{Db: global.DB}
	FS_OPERATION = fs.OperationService{ServerCtx: &(server.Ctx)}

}
