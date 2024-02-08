package shortcut

import (
	"apps/common/utils"
	"apps/probe/function"
	"fmt"
	"os"
	"strings"
)

func CreateScript(name string, prefixPath string, content string) (absolutePath *string, err error) {

	name = strings.ReplaceAll(name, " ", "-")

	// 去除prefix 末尾的 ‘/’
	u := prefixPath[len(prefixPath)-1]
	if string(u) == "/" {
		prefixPath = prefixPath[0 : len(prefixPath)-1]
	}

	// 拼接绝对路径
	path := fmt.Sprintf("%s/%s-%s.sh", prefixPath, name, utils.UUID())

	// 创建文件
	scriptFile, err := function.CreateFile(path, 0o777)
	defer scriptFile.Close()

	if err != nil {
		os.Remove(scriptFile.Name())
		return nil, err
	}

	// 写入内容
	_, err = scriptFile.WriteString(content)

	if err != nil {
		os.Remove(scriptFile.Name())
		return nil, err
	}

	r := scriptFile.Name()

	return &r, nil
}
