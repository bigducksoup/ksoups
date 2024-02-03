package fileservice

import (
	"config-manager/common/message/data"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

// DirRead reads dir info by accepting a path
func DirRead(dr data.DirRead) (data.DirResponse, error) {

	dir, err := os.ReadDir(dr.Path)

	if err != nil {
		log.Println(err)
		return data.DirResponse{}, err
	}

	dres := data.DirResponse{
		Parent:   dr.Path,
		FileOnly: dr.FileOnly,
		Items:    []data.DirItem{},
	}

	for _, item := range dir {

		if dr.FileOnly && item.IsDir() {
			continue
		}

		info, err := item.Info()

		if err != nil {
			log.Println(err)
			return data.DirResponse{}, err
		}

		// user id
		uid := info.Sys().(*syscall.Stat_t).Uid

		//user
		var fileUserName string
		fileUser, err := user.LookupId(fmt.Sprintf("%d", uid))
		if err != nil {
			log.Printf("look up file userId failed,file:%s,uid:%d,err:%s", info.Name(), uid, err.Error())
			fileUserName = "unknown"
		} else {
			fileUserName = fileUser.Name
		}

		// group id
		gid := info.Sys().(*syscall.Stat_t).Gid

		//group
		var fileGroupName string
		fileGroup, err := user.LookupGroupId(fmt.Sprintf("%d", gid))
		if err != nil {
			log.Printf("look up file user group Id failed,gid:%d,err:%s", gid, err.Error())
			fileGroupName = "unknown"
		} else {
			fileGroupName = fileGroup.Name
		}

		absPath := filepath.Join(dr.Path, info.Name())

		//是否为软链接
		var isLink bool = false

		//链接指向
		var link string
		//是软链接
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err = os.Readlink(absPath)
			if err != nil {
				log.Printf("read link err, path : %s , err : %s", absPath, err.Error())
				return data.DirResponse{}, err
			}
			isLink = true
		}

		dirItem := data.DirItem{
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

		dres.Items = append(dres.Items, dirItem)

	}

	return dres, nil

}

func DirCreate(dc data.DirCreate) (data.DirCreateResponse, error) {

	perm, err := strconv.ParseInt(dc.Permission, 8, 0)
	err = os.MkdirAll(dc.Path, os.FileMode(perm))

	if err != nil {
		return data.DirCreateResponse{}, err
	}

	return data.DirCreateResponse{
		Ok:         true,
		Path:       dc.Path,
		Permission: dc.Permission,
	}, nil

}
