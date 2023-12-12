package connect

import (
	"config-manager/core/message"
	"config-manager/probe/service"
	"os"
)

func DirReadHandler(dr message.DirRead) (message.DirResponse, error) {

	dir, err := os.ReadDir(dr.Path)

	if err != nil {
		return message.DirResponse{}, err
	}

	dres := message.DirResponse{
		Parent:   dr.Path,
		FileOnly: dr.FileOnly,
		Items:    []message.DirItem{},
	}

	for _, item := range dir {

		if dr.FileOnly && item.IsDir() {
			continue
		}

		info, _ := item.Info()

		dirItem := message.DirItem{
			Name:    item.Name(),
			IsDir:   item.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime(),
		}

		dres.Items = append(dres.Items, dirItem)

	}

	return dres, nil

}

func FileReadHandler(fr message.FileRead) (message.FileReadResponse, error) {

	path := fr.Path

	content, err := service.ReadFileContent(path)

	if err != nil {
		return message.FileReadResponse{}, err
	}

	response := message.FileReadResponse{
		Path:    path,
		Content: content,
		Size:    int64(len(content)),
	}

	return response, nil

}
