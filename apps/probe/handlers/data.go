package handlers

import (
	"apps/common/message"
	"encoding/json"
)

var DataHandlePolicy = map[message.DataType]func(data []byte) (any, message.DataType, error){
	message.READDIR:       handleReadDir,
	message.READ_FILE:     handleReadFile,
	message.MODIFY_FILE:   handleModifyFile,
	message.CREATE_FILE:   handleCreateFile,
	message.RUN_SHORTCUT:  handleRunSC,
	message.CREATE_DIR:    handleCreateDir,
	message.CREATE_SCRIPT: handleCreateScript,
	message.DELETE_FILE:   handleDeleteFile,
}

func readData[T any](data []byte) (T, error) {

	res := new(T)

	err := json.Unmarshal(data, res)
	if err != nil {
		return *res, err
	}

	return *res, nil
}
