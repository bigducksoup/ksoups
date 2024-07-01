package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	"apps/probe/service"
	filesystem "apps/probe/service/filesystem"
	"strconv"
)

func handleReadFile(data []byte) (any, message.DataType, error) {
	fileRead, err := readData[FileRead](data)

	if err != nil {

		return nil, message.ERROR, err
	}

	fileInfo, err := service.FileSystem.ReadFile(fileRead.Path)

	if err != nil {
		return nil, message.ERROR, err
	}

	resp := FileReadResponse{
		Path:    fileInfo.Path,
		Content: fileInfo.Content,
		Size:    fileInfo.Size,
	}
	return resp, message.READ_FILE_RESP, nil

}

func handleModifyFile(data []byte) (any, message.DataType, error) {

	mf, err := readData[FileModify](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	var changes []filesystem.Change

	for _, c := range mf.Changes {
		changes = append(changes, filesystem.Change{
			Count:     c.Count,
			Operation: c.Operation,
			Value:     c.Value,
		})
	}

	err = service.FileSystem.EditFile(mf.Path, changes)

	if err != nil {
		return nil, message.ERROR, err
	}

	return FileModifyResponse{
		Path: mf.Path,
		OK:   true,
	}, message.MODIFY_FILE_RESP, nil
}

func handleCreateFile(data []byte) (any, message.DataType, error) {

	fc, err := readData[FileCreate](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	perm, err := strconv.ParseInt(fc.Permission, 8, 0)

	if err != nil {
		return nil, message.ERROR, err
	}

	err = service.FileSystem.CreateFile(fc.Path, perm)

	if err != nil {
		return nil, message.ERROR, err
	}

	return FileCreateResponse{
		Ok:         true,
		Path:       fc.Path,
		Permission: fc.Permission,
	}, message.CREATE_FILE_RESP, nil

}

func handleDeleteFile(data []byte) (any, message.DataType, error) {

	fd, err := readData[FileDelete](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	err = service.FileSystem.DeleteFile(fd.Path)

	if err != nil {
		return nil, message.ERROR, err
	}

	return FileDeleteResponse{
		Ok:   true,
		Path: fd.Path,
	}, message.DELETE_FILE_RESP, nil
}
