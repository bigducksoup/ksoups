package function

import "os"

func CreateFile(path string, mode os.FileMode) (file *os.File, err error) {

	file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, mode)
	err = file.Chmod(mode)

	return
}
