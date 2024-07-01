package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	"apps/probe/service"
	"strconv"
)

func handleReadDir(data []byte) (any, message.DataType, error) {
	//将msg.Data解析为对应类型
	dr, err := readData[DirRead](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	info, err := service.FileSystem.ReadDir(dr.Path, dr.FileOnly)

	if err != nil {
		return nil, message.ERROR, err
	}

	resp := DirResponse{
		Parent:   info.Path,
		FileOnly: info.OnlyFile,
		Items:    make([]DirItem, 0, len(info.Items)),
	}

	for _, item := range info.Items {
		resp.Items = append(resp.Items, DirItem{
			Name:       item.Name,
			IsDir:      item.IsDir,
			IsLink:     item.IsLink,
			LinkTo:     item.LinkTo,
			Size:       item.Size,
			Permission: item.Permission,
			User:       item.User,
			UserGroup:  item.UserGroup,
			Mode:       item.Mode,
			ModTime:    item.ModTime,
		})
	}

	return resp, message.READDIRRESP, nil
}

func handleCreateDir(data []byte) (any, message.DataType, error) {

	dc, err := readData[DirCreate](data)
	if err != nil {
		return nil, message.ERROR, err
	}

	perm, err := strconv.ParseInt(dc.Permission, 8, 0)

	if err != nil {
		return nil, message.ERROR, err
	}

	err = service.FileSystem.CreateDir(dc.Path, perm)
	if err != nil {
		return nil, message.ERROR, err
	}

	return DirCreateResponse{
		Ok:         true,
		Path:       dc.Path,
		Permission: dc.Permission,
	}, message.CREATE_DIR_RESP, nil

}
