package script

import (
	"apps/common/utils"
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ScriptRepo interface {
	Store(scriptContent []byte, name string, args []string) (id string, err error)
	GetPath(id string) (path string, ok bool)
	Get(id string) (script Script, ok bool)
	Delete(id string) error
	DeleteByPath(path string) error
}

type ShellScriptRepo struct {
	storePath string
}

// Store save script on scripts dictionary
// the name of script and args will be written at the second line of script file
func (s *ShellScriptRepo) Store(scriptContent []byte, name string, args []string) (id string, err error) {

	// create id and set file name
	id = utils.UUID()
	fileName := fmt.Sprintf("%s.sh", id)

	// releative path
	fullPath := filepath.Join(s.storePath, fileName)

	// create file
	file, err := os.Create(fullPath)

	if err != nil {
		return "", err
	}

	// chmod
	err = file.Chmod(0o777)

	if err != nil {
		return "", err
	}

	// format args
	argsInFile := fmt.Sprintf("#%s %s\n", name, strings.Join(args, " "))

	_, err = file.Write([]byte("#!/bin/bash\n"))

	if err != nil {
		return "", err
	}

	_, err = file.Write([]byte(argsInFile))

	if err != nil {
		return "", err
	}

	_, err = file.Write(scriptContent)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (s *ShellScriptRepo) GetPath(id string) (path string, ok bool) {
	script, ok := s.Get(id)

	if !ok {
		return "", false
	}

	return script.Path(), true
}

// Get get script by id
// actually locate file by path
func (s *ShellScriptRepo) Get(id string) (script Script, ok bool) {

	fileName := fmt.Sprintf("%s.sh", id)

	fullPath := filepath.Join(s.storePath, fileName)

	file, err := os.Open(fullPath)

	if err != nil {
		return nil, false
	}

	reader := bufio.NewReader(file)

	reader.ReadString('\n')

	line, err := reader.ReadString('\n')

	if err != nil {
		return nil, false
	}

	line = line[:len(line)-1]

	arr := strings.Split(line, " ")

	name := arr[0]
	args := arr[1:]

	if len(args) == 1 && args[0] == "" {
		args = nil
	}

	absPath, err := filepath.Abs(fullPath)

	if err != nil {
		return nil, false
	}

	return &ShellScript{
		scriptType: Shell,
		id:         id,
		name:       name[1:],
		path:       absPath,
		args:       args,
	}, true

}

func (s *ShellScriptRepo) Delete(id string) error {
	script, ok := s.Get(id)

	if !ok {
		return errors.New("no such script with id = " + id)
	}

	return s.DeleteByPath(script.Path())
}

func (s *ShellScriptRepo) DeleteByPath(path string) error {
	return os.Remove(path)
}

func NewShellScriptRepo(storePath string) *ShellScriptRepo {

	_, err := os.Stat(storePath)

	if err != nil {
		mkDirErr := os.Mkdir(storePath, os.ModeDir|os.ModePerm)
		if mkDirErr != nil {
			panic(mkDirErr)
		}
	}

	return &ShellScriptRepo{
		storePath: storePath,
	}
}
