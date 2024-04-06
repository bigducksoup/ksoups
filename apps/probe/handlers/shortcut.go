package handlers

import (
	. "apps/common/message/data"
	"apps/probe/config"
	shortcutService "apps/probe/service/shortcut"
	"encoding/json"
)

func handleRunSC(data []byte) ([]byte, error) {

	scRun, err := readData[ShortcutRun](data)

	if err != nil {
		return nil, err
	}

	var result any

	//check if run in real time way
	if scRun.RealTime {
		result = shortcutService.ExecuteShortcutRealTime(scRun)
	} else {
		result = shortcutService.ExecuteShortcut(scRun)
	}

	bytes, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// 创建脚本
func handleCreateScript(data []byte) ([]byte, error) {

	c, err := readData[CreateScript](data)

	if err != nil {
		return nil, err
	}

	scriptPath, err := shortcutService.CreateScript(c.Name, config.Conf.ScriptPath, c.Content)

	if err != nil {
		return nil, err
	}

	resp := CreateScriptResp{
		Name:    c.Name,
		AbsPath: *scriptPath,
	}

	bytes, err := json.Marshal(resp)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
