package handlers

import (
	"config-manager/common/message"
	fileservice "config-manager/probe/service/fileService"
	"encoding/json"
)

func readData[T any](data []byte) (T, error) {

	res := new(T)

	err := json.Unmarshal(data, res)
	if err != nil {
		return *res, err
	}

	return *res, nil
}

var DataHandlePolicy map[message.DataType]func(data []byte) ([]byte, error) = map[message.DataType]func(data []byte) ([]byte, error){
	message.READDIR:    handleReadDir,
	message.READFILE:   handleReadFile,
	message.MODIFYFILE: handleModifyFile,
}

func handleReadDir(data []byte) ([]byte, error) {
	//将msg.Data解析为对应类型
	dr, err := readData[message.DirRead](data)

	//调用对应处理方法
	dirReadResponse, err := fileservice.DirRead(dr)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(dirReadResponse)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func handleReadFile(data []byte) ([]byte, error) {
	fileRead, err := readData[message.FileRead](data)

	if err != nil {

		return nil, err
	}

	fileReadResponse, err := fileservice.FileRead(fileRead)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(fileReadResponse)

	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func handleModifyFile(data []byte) ([]byte, error) {
	return nil, nil
}
