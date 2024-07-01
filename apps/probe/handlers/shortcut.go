package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	"apps/probe/script"
	"apps/probe/service"
	shortcutService "apps/probe/service/shortcut"
)

func handleRunSC(data []byte) (any, message.DataType, error) {

	scRun, err := readData[ShortcutRun](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	var result any

	//check if run in real time way
	if scRun.RealTime {
		result = shortcutService.ExecuteShortcutRealTime(scRun)
	} else {
		result = shortcutService.ExecuteShortcut(scRun)
	}

	return result, message.RUN_SHORTCUT_RESP, nil
}

// 创建脚本
func handleCreateScript(data []byte) (any, message.DataType, error) {

	c, err := readData[CreateScript](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	script, err := service.ShortcutManage.CreateScript(c.Name, script.Shell, []byte(c.Content), nil)

	if err != nil {
		return nil, message.ERROR, err
	}

	resp := CreateScriptResp{
		Name:    c.Name,
		AbsPath: script.Path(),
	}

	return resp, message.CREATE_SCRIPT_RESP, nil
}
