package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	fileservice "apps/probe/service/fs"
)

func handleReadFile(data []byte) (any, message.DataType, error) {
	fileRead, err := readData[FileRead](data)

	if err != nil {

		return nil, message.ERROR, err
	}

	fileReadResponse, err := fileservice.FileRead(fileRead)

	if err != nil {
		return nil, message.ERROR, err
	}

	return fileReadResponse, message.READ_FILE_RESP, nil

}

func handleModifyFile(data []byte) (any, message.DataType, error) {

	mf, err := readData[FileModify](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	result, err := fileservice.FileModify(mf)

	if err != nil {
		return nil, message.ERROR, err
	}

	return result, message.MODIFY_FILE_RESP, nil
}

func handleCreateFile(data []byte) (any, message.DataType, error) {

	fc, err := readData[FileCreate](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	fileCreateResponse, err := fileservice.FileCreate(fc)

	if err != nil {
		return nil, message.ERROR, err
	}

	return fileCreateResponse, message.CREATE_FILE_RESP, nil

}

func handleDeleteFile(data []byte) (any, message.DataType, error) {

	fd, err := readData[FileDelete](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	fileDeteteResponse, err := fileservice.FileDelete(fd)

	if err != nil {
		return nil, message.ERROR, err
	}

	return fileDeteteResponse, message.DELETE_FILE_RESP, nil
}
