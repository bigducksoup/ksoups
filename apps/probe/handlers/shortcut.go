package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	"apps/probe/script"
	"apps/probe/service"
)

func handleRunShortcut(data []byte) (any, message.DataType, error) {

	scRun, err := readData[ShortcutRun](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	var result any

	script, err := service.ShortcutManage.GetScript(scRun.Id, script.ScriptType(scRun.Type))

	if err != nil {
		return nil, message.ERROR, err
	}

	stdOut, stdErr, err := service.ShortExecutionService.Exec(script)

	if err != nil {
		return nil, message.ERROR, err
	}

	result = ShortcutRunResp{
		Ok:     true,
		Err:    string(stdErr),
		StdOut: string(stdOut),
		StdErr: string(stdErr),
	}

	return result, message.RUN_SHORTCUT_RESP, nil
}

// 创建脚本
func handleCreateScript(data []byte) (any, message.DataType, error) {

	c, err := readData[CreateScript](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	script, err := service.ShortcutManage.CreateScript(c.Name, script.ScriptType(c.ScriptType), []byte(c.Content), c.Args)

	if err != nil {
		return nil, message.ERROR, err
	}

	resp := CreateScriptResp{
		Id:      script.Id(),
		Name:    c.Name,
		AbsPath: script.Path(),
	}

	return resp, message.CREATE_SCRIPT_RESP, nil
}
