package handlers

import (
	"config-manager/common/message"
	"encoding/json"
)

var DataHandlePolicy = map[message.DataType]func(data []byte) ([]byte, error){
	message.READDIR:       handleReadDir,
	message.READFILE:      handleReadFile,
	message.MODIFYFILE:    handleModifyFile,
	message.CREATEFILE:    handleCreateFile,
	message.RUN_SHORTCUT:  handleRunSC,
	message.CREATE_DIR:    handleCreateDir,
	message.CREATE_SCRIPT: handleCreateScript,
}

func readData[T any](data []byte) (T, error) {

	res := new(T)

	err := json.Unmarshal(data, res)
	if err != nil {
		return *res, err
	}

	return *res, nil
}
