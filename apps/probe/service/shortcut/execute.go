package shortcut

import (
	"apps/probe/script"
	"context"
	"fmt"
	"io"
	"time"
)

var RunnerMap map[script.ScriptType]script.ScriptRunner = make(map[script.ScriptType]script.ScriptRunner)

type ShortcutExecutionService struct {
	shellRunner *script.ShellScriptRunner
}

func NewShortcutExecutionService() *ShortcutExecutionService {

	shellRunner := script.NewShellScriptRunner()

	RunnerMap[script.Shell] = shellRunner

	return &ShortcutExecutionService{
		shellRunner: shellRunner,
	}
}

func (se *ShortcutExecutionService) Exec(st script.Script, timeout time.Duration) (stdOut []byte, stdErr []byte, err error) {

	runner, ok := RunnerMap[st.Type()]

	if !ok {
		return nil, nil, fmt.Errorf("could not find runner for script type : %d", st.Type())
	}

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	return runner.Run(ctx, st)

}

func (se *ShortcutExecutionService) ExecAsync(ctx context.Context, st script.Script, scriptType script.ScriptType) (inPipe io.WriteCloser, outPipe io.ReadCloser, errPipe io.ReadCloser, err error) {

	runner, ok := RunnerMap[scriptType]

	if !ok {
		return nil, nil, nil, fmt.Errorf("could not find runner for script type : %d", scriptType)
	}

	return runner.RunAsync(ctx, st)
}
