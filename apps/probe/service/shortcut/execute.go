package shortcut

import (
	"config-manager/center/model"
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
		return string(bytes), err
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

	if shortcut.Type == model.SCRIPT {
		shortcut.Payload = "sh " + shortcut.Payload
	}

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
			execResult.StdErr = result
			return execResult
		}
		execResult.StdOut = result
	}

	return execResult

}
