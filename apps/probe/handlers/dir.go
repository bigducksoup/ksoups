package handlers

import (
	. "config-manager/common/message/data"
	fileservice "config-manager/probe/service/fs"
	"encoding/json"
)

func handleReadDir(data []byte) ([]byte, error) {
	//将msg.Data解析为对应类型
	dr, err := readData[DirRead](data)

	if err != nil {
		return nil, err
	}

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

func handleCreateDir(data []byte) ([]byte, error) {

	dc, err := readData[DirCreate](data)
	if err != nil {
		return nil, err
	}

	dirCreateResponse, err := fileservice.DirCreate(dc)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(dirCreateResponse)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}
