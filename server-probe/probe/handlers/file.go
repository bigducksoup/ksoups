package handlers

import (
	. "config-manager/common/message/data"
	fileservice "config-manager/probe/service/fs"
	"encoding/json"
)

func handleReadFile(data []byte) ([]byte, error) {
	fileRead, err := readData[FileRead](data)

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

	mf, err := readData[FileModify](data)

	if err != nil {
		return nil, err
	}

	result, err := fileservice.FileModify(mf)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func handleCreateFile(data []byte) ([]byte, error) {

	fc, err := readData[FileCreate](data)

	if err != nil {
		return nil, err
	}

	fileCreateResponse, err := fileservice.FileCreate(fc)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(fileCreateResponse)

	if err != nil {
		return nil, err
	}

	return bytes, nil

}
