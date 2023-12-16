package fileservice

import (
	"config-manager/common/message"
	"config-manager/probe/common"
	"log"
	"os"
)

func DirRead(dr message.DirRead) (message.DirResponse, error) {

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

func FileRead(fr message.FileRead) (message.FileReadResponse, error) {

	path := fr.Path

	file, err := GetFile(path)

	if err != nil {
		return message.FileReadResponse{}, err
	}

	response := message.FileReadResponse{
		Path:    path,
		Content: file.GetContent(),
		Size:    int64(len(file.GetContent())),
	}

	return response, nil

}

var fileCache = common.NewLRUCache[*common.File](10)

func ClearFileCache() {

	log.Println("clear file cache, flush all files")

	fileCache.Range(func(key string, value *common.File) bool {
		value.Flush()
		return true
	})

	fileCache.Clear()

}

// GetFile gets a File from the fileCache.
// if no File is found, it creates a new File and returns it.
func GetFile(path string) (*common.File, error) {

	v, ok := fileCache.Get(path)

	if ok {
		return v, nil
	}

	file, err := common.NewFile(path)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fileCache.Put(path, file)
	log.Println(fileCache.Size())

	return file, nil
}
