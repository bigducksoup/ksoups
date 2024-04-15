package function

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"time"
)

func ExecuteCmdWithTimeout(cmd string, timeout time.Duration) (string, error) {

	//超时上下文
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	cmdSplits := strings.Split(cmd, " ")

	command := exec.CommandContext(ctx, cmdSplits[0], cmdSplits[1:]...)

	bytes, err := command.CombinedOutput()

	if err != nil {
		if ctx.Err() != nil {
			return string(bytes), nil
		}
		return string(bytes), err
	}

	return string(bytes), nil

}

func ExecuteCmdIgnoreResult(cmd string, timeout time.Duration) error {

	//超时上下文
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	cmdSplits := strings.Split(cmd, " ")

	command := exec.CommandContext(ctx, cmdSplits[0], cmdSplits[1:]...)

	err := command.Start()

	if err != nil {
		return err
	}

	go func() {
		err = command.Wait()
	}()

	return nil
}

func ExecuteCmdRealTime(ctx context.Context, cmd string, timeout time.Duration) (err error, outPipe *io.ReadCloser, errPipe *io.ReadCloser) {

	cmdSplits := strings.Split(cmd, " ")

	var command *exec.Cmd

	if timeout <= 0*time.Millisecond {
		command = exec.Command(cmdSplits[0], cmdSplits[1:]...)
	} else {
		command = exec.CommandContext(ctx, cmdSplits[0], cmdSplits[1:]...)
	}

	stdoutPipe, err := command.StdoutPipe()

	if err != nil {
		return err, nil, nil
	}

	stderrPipe, err := command.StderrPipe()

	if err != nil {
		return err, nil, nil
	}

	err = command.Start()

	if err != nil {
		return err, nil, nil
	}

	return nil, &stdoutPipe, &stderrPipe
}
