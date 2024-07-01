package script

import "os"

type ScriptType int

const (
	Shell ScriptType = iota
	Go
	Python
)

type Script interface {
	Id() string
	Type() ScriptType
	Name() string
	Path() string
	Args() []string
	Content() ([]byte, error)
}

type ShellScript struct {
	id         string
	scriptType ScriptType
	path       string
	name       string
	args       []string
}

func (s *ShellScript) Id() string {
	return s.id
}

func (s *ShellScript) Type() ScriptType {
	return s.scriptType
}

func (s *ShellScript) Name() string {
	return s.name
}

func (s *ShellScript) Path() string {
	return s.path
}

func (s *ShellScript) Args() []string {

	return s.args
}

func (s *ShellScript) Content() ([]byte, error) {
	return os.ReadFile(s.path)
}
