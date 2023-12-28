package fileservice

import (
	"config-manager/common/message/data"
	"config-manager/probe/common"
	"fmt"
	"log"
	"os"
	"strconv"
)

// FileRead reads file content
func FileRead(fr data.FileRead) (data.FileReadResponse, error) {

	path := fr.Path

	file, err := GetFile(path)

	if err != nil {
		log.Println(err)
		return data.FileReadResponse{}, err
	}

	response := data.FileReadResponse{
		Path:    path,
		Content: file.GetContent(),
		Size:    int64(len(file.GetContent())),
	}

	return response, nil

}

func FileModify(fm data.FileModify) (data.FileModifyResponse, error) {

	file, err := GetFile(fm.Path)

	result := data.FileModifyResponse{}

	if err != nil {
		log.Println(err)
		return result, err
	}

	lineIndex := 0

	for _, change := range fm.Changes {

		if change.Operation == data.ADDED {
			err = file.InsertLines(lineIndex, change.Value)
			if err != nil {
				break
			}
			lineIndex += change.Count
			continue
		}

		if change.Operation == data.REMOVED {
			err = file.RemoveLines(lineIndex, lineIndex+change.Count)
			if err != nil {
				log.Println(err)
				break
			}
			continue
		}

		if change.Operation == data.REMAIN {
			lineIndex += change.Count
		}

	}

	if err != nil {

		log.Println("unable to edit " + fm.Path)
		log.Println(err)
		log.Println("----------------")
		result.OK = false
		return result, err
	}

	err = file.Flush()
	if err != nil {
		log.Println(err)
		return data.FileModifyResponse{}, err
	}

	result.OK = true
	return result, nil
}

func FileCreate(fc data.FileCreate) (data.FileCreateResponse, error) {

	perm, err := strconv.ParseInt(fc.Permission, 8, 0)

	if err != nil {
		return data.FileCreateResponse{}, err
	}

	mode := os.FileMode(perm)

	fmt.Println(mode)

	file, err := os.OpenFile(fc.Path, os.O_CREATE, mode)

	if err != nil {
		return data.FileCreateResponse{}, err
	}

	err = file.Chmod(mode)
	if err != nil {
		return data.FileCreateResponse{}, err
	}

	defer file.Close()

	return data.FileCreateResponse{
		Ok:         true,
		Path:       fc.Path,
		Permission: fc.Permission,
	}, nil

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
