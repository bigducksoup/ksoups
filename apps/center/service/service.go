package service

import (
	scrun "config-manager/center/action"
	"config-manager/center/global"
	"config-manager/center/service/chain"
	"config-manager/center/service/shortcut"
	"config-manager/center/service/ssh"
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

}
