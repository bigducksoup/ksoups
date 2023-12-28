package shortcut

import (
	"config-manager/common/utils"
	"config-manager/probe/config"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

// OneLineShortcut is  a shortcut just execute one line command
type OneLineShortcut struct {
	Id          string `gorm:"primaryKey;size:20"`
	Name        string
	Description string
	CreateTime  time.Time
	//超时时间
	Timeout time.Duration
	//仅运行，不处理结果
	JustRun bool
	Command string
	ProbeId string
}

// ScriptShortcut is a shortcut that point to a script file
type ScriptShortcut struct {
	Id          string
	Name        string
	Description string
	CreateTime  time.Time
	Timeout     time.Duration
	JustRun     bool
	Ext         string
	Path        string
	ProbeId     string
}

type ExecResult struct {
	Ok  bool
	Out string
	Err string
}

func NewOneLineShortcut(name string, cmd string, timeOut time.Duration, justRun bool, description string) OneLineShortcut {

	return OneLineShortcut{
		Id:          utils.UUID(),
		Name:        name,
		CreateTime:  time.Now(),
		Description: description,
		Timeout:     timeOut,
		JustRun:     justRun,
		Command:     cmd,
	}

}

func NewScriptShortcutWithContent(name string, ext string, content []byte) (ScriptShortcut, error) {

	//脚本目录
	scriptPath := config.Conf.ScriptsPath + "/scripts"

	//文件全路径
	path := scriptPath + "/" + name + "." + ext

	//是否存在文件
	info, err := os.Stat(path)

	if info != nil {
		return ScriptShortcut{}, errors.New("file already exists,path:" + path)
	}

	//如果不存在目录，则创建
	err = os.Mkdir(scriptPath, 0o777)

	if err != nil && !errors.Is(err, os.ErrExist) {
		return ScriptShortcut{}, err
	}

	//创建文件
	file, err := os.Create(path)

	if err != nil {
		log.Printf("unable to create file,path:%s,Err:%s \n", path, err.Error())
		return ScriptShortcut{}, err
	}

	defer file.Close()

	//写入内容
	_, err = file.Write(content)

	if err != nil {
		log.Printf("unable to write file,path:%s,Err:%s \n", path, err.Error())
		return ScriptShortcut{}, err
	}

	err = file.Chmod(0o755)

	if err != nil {
		log.Printf("unable to chmod file,path:%s,Err:%s \n", path, err.Error())
		return ScriptShortcut{}, err
	}

	return ScriptShortcut{
		Id:         utils.UUID(),
		Path:       file.Name(),
		CreateTime: time.Now(),
		Name:       name,
		Ext:        ext,
	}, nil

}

func NewScriptShortcut(name string, path string, timeOut time.Duration, justRun bool, description string) (ScriptShortcut, error) {

	//是否存在文件
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return ScriptShortcut{}, err
	}

	return ScriptShortcut{
		Id:          utils.UUID(),
		Name:        name,
		Path:        path,
		CreateTime:  time.Now(),
		Timeout:     timeOut,
		JustRun:     justRun,
		Description: description,
		Ext:         filepath.Ext(path),
	}, nil

}
