package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	"apps/probe/config"
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

	scriptPath, err := shortcutService.CreateScript(c.Name, config.Conf.ScriptPath, c.Content)

	if err != nil {
		return nil, message.ERROR, err
	}

	resp := CreateScriptResp{
		Name:    c.Name,
		AbsPath: *scriptPath,
	}

	return resp, message.CREATE_SCRIPT_RESP, nil
}
