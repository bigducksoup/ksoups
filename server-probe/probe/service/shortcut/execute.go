package shortcut

import (
	"config-manager/common/message/data"
	"context"
	"os/exec"
	"strings"
	"time"
)

func executeCmdWithTimeout(cmd string, timeout time.Duration) (string, error) {

	//超时上下文
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	cmdSplits := strings.Split(cmd, " ")

	command := exec.CommandContext(ctx, cmdSplits[0], cmdSplits[1:]...)

	bytes, err := command.Output()

	if err != nil {
		if ctx.Err() != nil {
			return string(bytes), nil
		}
		return "", err
	}

	return string(bytes), nil

}

func executeCmdIgnoreResult(cmd string, timeout time.Duration) error {

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

func ExecuteShortcut(shortcut data.ShortcutRun) data.ShortcutRunResp {

	//执行结果
	execResult := data.ShortcutRunResp{
		Ok: true,
	}
	//仅运行
	if shortcut.JustRun {
		err := executeCmdIgnoreResult(shortcut.Payload, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}

	} else {
		result, err := executeCmdWithTimeout(shortcut.Payload, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}
		execResult.Out = result
	}

	return execResult

}
