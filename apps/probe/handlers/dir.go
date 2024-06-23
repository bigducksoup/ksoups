package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	fileservice "apps/probe/service/fs"
)

func handleReadDir(data []byte) (any, message.DataType, error) {
	//将msg.Data解析为对应类型
	dr, err := readData[DirRead](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	//调用对应处理方法
	dirReadResponse, err := fileservice.DirRead(dr)

	if err != nil {
		return nil, message.ERROR, err
	}

	return dirReadResponse, message.READDIRRESP, nil
}

func handleCreateDir(data []byte) (any, message.DataType, error) {

	dc, err := readData[DirCreate](data)
	if err != nil {
		return nil, message.ERROR, err
	}

	dirCreateResponse, err := fileservice.DirCreate(dc)
	if err != nil {
		return nil, message.ERROR, err
	}

	return dirCreateResponse, message.CREATE_DIR_RESP, nil

}
