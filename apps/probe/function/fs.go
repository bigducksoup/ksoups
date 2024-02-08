package function

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"syscall"
)

func CreateFile(path string, mode os.FileMode) (file *os.File, err error) {

	file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, mode)
	err = file.Chmod(mode)

	return
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func FindOwner(info fs.FileInfo) (Id uint32, name string, err error) {

	uid := info.Sys().(*syscall.Stat_t).Uid

	u, err := user.LookupId(fmt.Sprintf("%d", uid))

	if err != nil {
		name = "unknown"
		return
	}

	return uid, u.Name, nil

}

func FindGroup(info fs.FileInfo) (Id uint32, name string, err error) {

	gid := info.Sys().(*syscall.Stat_t).Gid

	g, err := user.LookupGroupId(fmt.Sprintf("%d", gid))

	if err != nil {
		log.Printf("look up file user group Id failed,gid:%d,err:%s", gid, err.Error())
		name = "unknown"
		return
	}

	return gid, g.Name, nil
}
