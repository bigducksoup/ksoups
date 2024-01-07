package handlers

import (
	. "config-manager/common/message/data"
	shortcutService "config-manager/probe/service/shortcut"
	"encoding/json"
)

func handleRunSC(data []byte) ([]byte, error) {

	scRun, err := readData[ShortcutRun](data)

	if err != nil {
		return nil, err
	}

	shortcutRunResp := shortcutService.ExecuteShortcut(scRun)

	bytes, err := json.Marshal(shortcutRunResp)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
