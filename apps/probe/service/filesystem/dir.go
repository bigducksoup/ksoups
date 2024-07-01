package filesystem

import (
	"apps/probe/function"
	"log"
	"os"
	"path/filepath"
	"time"
)

type DirInfo struct {
	Path     string
	OnlyFile bool
	Items    []DirItem
}

type DirItem struct {
	Name       string    `json:"name"`
	IsDir      bool      `json:"isDir"`
	IsLink     bool      `json:"isLink"`
	LinkTo     string    `json:"linkTo"`
	Size       int64     `json:"size"`
	Permission string    `json:"permission"`
	User       string    `json:"user"`
	UserGroup  string    `json:"usergroup"`
	Mode       string    `json:"mode"`
	ModTime    time.Time `json:"modTime"`
}

func (fs *FileSystemService) ReadDir(path string, onlyFile bool) (info DirInfo, err error) {

	dir, err := os.ReadDir(path)

	if err != nil {
		log.Println(err)
		return DirInfo{}, err
	}

	dirInfo := DirInfo{
		Path:     path,
		OnlyFile: onlyFile,
		Items:    []DirItem{},
	}

	for _, item := range dir {
		if onlyFile && item.IsDir() {
			continue
		}

		info, err := item.Info()

		if err != nil {
			return DirInfo{}, err
		}

		uid, fileUserName, err := function.FindOwner(info)

		if err != nil {
			log.Printf("look up file userId failed,file:%s,uid:%d,err:%s", info.Name(), uid, err.Error())
		}

		gid, fileGroupName, err := function.FindGroup(info)

		if err != nil {
			log.Printf("look up file user group Id failed,gid:%d,err:%s", gid, err.Error())
		}

		absPath := filepath.Join(path, info.Name())

		//是否为软链接
		var isLink bool = false

		//链接指向
		var link string
		//是软链接
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err = os.Readlink(absPath)
			if err != nil {
				log.Printf("read link err, path : %s , err : %s", absPath, err.Error())
				return dirInfo, err
			}
			isLink = true
		}

		dirItem := DirItem{
			Name:       item.Name(),
			IsDir:      item.IsDir(),
			IsLink:     isLink,
			LinkTo:     link,
			Size:       info.Size(),
			Permission: info.Mode().Perm().String(),
			User:       fileUserName,
			UserGroup:  fileGroupName,
			Mode:       info.Mode().String(),
			ModTime:    info.ModTime(),
		}

		dirInfo.Items = append(dirInfo.Items, dirItem)

	}

	return dirInfo, nil
}

func (fs *FileSystemService) CreateDir(path string, perm int64) error {
	return os.MkdirAll(path, os.FileMode(perm))
}
