package filesystem

import (
	"apps/common/message/data"
	"apps/probe/common"
	"apps/probe/function"
	"log"
	"os"
)

type FileInfo struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

// GetFile gets a File from the fileCache.
// if no File is found, it creates a new File and returns it.
func (fs *FileSystemService) getFile(path string) (*common.File, error) {

	f, ok := fs.cache.Get(path)

	if ok {
		info, err := os.Stat(path)

		if err != nil {
			return nil, err
		}
		if info.ModTime() != f.ModTime {
			return fs.LoadFileThenPutInCache(path)
		}
		return f, nil
	}

	return fs.LoadFileThenPutInCache(path)
}

func (fs *FileSystemService) LoadFileThenPutInCache(path string) (*common.File, error) {
	file, err := common.NewFile(path)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fs.cache.Put(path, file)

	return file, nil
}

func (fs *FileSystemService) ReadFile(path string) (FileInfo, error) {
	file, err := fs.getFile(path)

	if err != nil {
		log.Println(err)
		return FileInfo{}, err
	}

	content := file.GetContent()

	return FileInfo{
		Path:    path,
		Content: content,
		Size:    int64(len(content)),
	}, nil

}

type Change struct {
	Count     int      `json:"count"`
	Operation int      `json:"operation"`
	Value     []string `json:"value"`
}

func (fs *FileSystemService) EditFile(path string, changes []Change) error {
	file, err := fs.getFile(path)

	if err != nil {
		return err
	}

	lineIndex := 0

	for _, change := range changes {

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
		log.Println("unable to edit " + path)
		log.Println(err)
		log.Println("----------------")
		return err
	}

	err = file.Flush()

	return err

}

func (fs *FileSystemService) CreateFile(path string, perm int64) error {

	mode := os.FileMode(perm)
	file, err := function.CreateFile(path, mode)

	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (fs *FileSystemService) DeleteFile(path string) error {
	return function.DeleteFile(path)
}

func (fs *FileSystemService) ClearFileCache() {
	fs.cache.Range(func(key string, value *common.File) bool {
		value.Flush()
		return true
	})

	fs.cache.Clear()
}
