package script

import (
	"context"

	"io"
	"os/exec"
	"slices"
)

type ScriptRunner interface {
	Run(script Script) (result []byte, err error)
	RunAsync(ctx context.Context, script Script) (inPipe io.WriteCloser, outPipe io.ReadCloser, errPipe io.ReadCloser, err error)
}

type ShellScriptRunner struct {
}


// Run script and return output
func (s *ShellScriptRunner) Run(script Script) (result []byte, err error) {

	args := slices.Insert(script.Args(), 0, script.Path())

	cmd := exec.Command("sh", args...)

	return cmd.CombinedOutput()
}


// RunAsync run Script async 
// return inPipe,outPipe and errPipe
func (s *ShellScriptRunner) RunAsync(ctx context.Context, script Script) (inPipe io.WriteCloser, outPipe io.ReadCloser, errPipe io.ReadCloser, err error) {

	args := slices.Insert(script.Args(), 0, script.Path())

	cmd := exec.CommandContext(ctx, "sh", args...)


	if err != nil {
		return nil, nil, nil, err
	}
	inPipe, err = cmd.StdinPipe()

	if err != nil {
		return nil, nil, nil, err
	}

	outPipe, err = cmd.StdoutPipe()

	if err != nil {
		return nil, nil, nil, err
	}

	errPipe, err = cmd.StderrPipe()

	if err != nil {
		return nil, nil, nil, err
	}

	err = cmd.Start()

	return
}

func NewShellScriptRunner() *ShellScriptRunner {

	return &ShellScriptRunner{}

}
