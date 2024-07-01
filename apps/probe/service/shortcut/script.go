package shortcut

import (
	"apps/probe/config"
	"apps/probe/script"
	"errors"
	"fmt"
)

// ScriptRepoMap
// @see NewShortcutManageService
var ScriptRepoMap map[script.ScriptType]script.ScriptRepo = make(map[script.ScriptType]script.ScriptRepo)

type ShortcutManageService struct {
	shellRepo *script.ShellScriptRepo
}

func NewShortcutManageService() *ShortcutManageService {
	shellRepo := script.NewShellScriptRepo(config.Conf.ScriptPath)
	ScriptRepoMap[script.Shell] = shellRepo
	return &ShortcutManageService{
		shellRepo: shellRepo,
	}
}

// CreateScript create script by name,type and content
func (s *ShortcutManageService) CreateScript(name string, scriptType script.ScriptType, content []byte, args []string) (script.Script, error) {

	repo, ok := ScriptRepoMap[scriptType]

	if !ok {
		return nil, errors.New("do not support such script type")
	}

	id, err := repo.Store(content, name, args)

	if err != nil {
		return nil, err
	}

	st, ok := repo.Get(id)

	if !ok {
		return nil, fmt.Errorf("could not find script with id : %s",id)
	}

	return st, nil
}
