package shortcut

import (
	. "config-manager/common/model/shortcut"
	. "config-manager/common/shortcut"
	"context"
	"os"
	"os/exec"
	"path/filepath"
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

func ExecuteOneLineShortcut(shortcut OneLineShortcut) ExecResult {

	//执行结果
	execResult := shortcut.ExecResult{}
	//仅运行
	if shortcut.JustRun {
		err := executeCmdIgnoreResult(shortcut.Command, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}

	} else {
		result, err := executeCmdWithTimeout(shortcut.Command, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}
		execResult.Out = result
	}

	return execResult

}

func ExecuteScriptShortcut(shortcut ScriptShortcut) ExecResult {

	path := shortcut.Path

	result := shortcut.ExecResult{}

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		result.Ok = false
		result.Err = err.Error()
		return result
	}

	if filepath.Ext(info.Name()) != ".sh" {
		result.Ok = false
		result.Err = "unacceptable file extension name"
		return result
	}

	cmd := "sh " + path
	//执行结果
	execResult := shortcut.ExecResult{
		Ok: true,
	}

	//仅运行
	if shortcut.JustRun {
		err := executeCmdIgnoreResult(cmd, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}

	} else {
		result, err := executeCmdWithTimeout(cmd, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}
		execResult.Out = result
	}

	return execResult

}
