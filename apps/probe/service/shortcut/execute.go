package shortcut

import (
	"apps/center/model"
	"apps/common/message"
	"apps/common/message/data"
	"apps/probe/connect"
	"apps/probe/function"
	"apps/probe/script"
	"context"
	"fmt"
	"io"
	"log"
)

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
		err := function.ExecuteCmdIgnoreResult(shortcut.Payload, shortcut.Timeout)

		if err != nil {
			execResult.Ok = false
			execResult.Err = err.Error()
			return execResult
		}

	} else {
		result, err := function.ExecuteCmdWithTimeout(shortcut.Payload, shortcut.Timeout)

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

func ExecuteShortcutRealTime(runMeta data.ShortcutRun) data.RealTimeShortcutRunResp {

	ctx, cancelFunc := context.WithTimeout(context.TODO(), runMeta.Timeout)
	defer cancelFunc()

	err, outPipe, errPipe := function.ExecuteCmdRealTime(context.WithoutCancel(ctx), runMeta.Payload, runMeta.Timeout)

	if err != nil {
		return data.RealTimeShortcutRunResp{
			Ok:    false,
			RunId: runMeta.Id,
			Err:   err.Error(),
		}
	}

	// handle pipeline output
	handlePipe := func(pipe *io.ReadCloser, processData func([]byte)) {

		buf := make([]byte, 1024)

		for {
			// 从 pipe 中读取数据
			n, err := (*pipe).Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("An error occurred: %v", err)
				}
				break
			}

			// 处理读取的数据
			processData(buf[:n])
		}
		// 关闭 pipe
		if err := (*pipe).Close(); err != nil {
			log.Printf("Failed to close: %v", err)
		}
	}

	// async handle stdout pipeline
	go handlePipe(outPipe, func(bytes []byte) {

		outPut := data.RealTimeShortcutOutPut{
			Type:    0,
			Payload: string(bytes),
			RunId:   runMeta.Id,
		}

		// TODO with id ??
		e := connect.ProbeInstance.PushToCenter(outPut, message.SHORTCUT_OUTPUT)
		if e != nil {
			(*outPipe).Close()
		}
	})

	// async handle stderr pipeline
	go handlePipe(errPipe, func(bytes []byte) {

		outPut := data.RealTimeShortcutOutPut{
			Type:    1,
			Payload: string(bytes),
			RunId:   runMeta.Id,
		}

		e := connect.ProbeInstance.PushToCenter(outPut, message.SHORTCUT_OUTPUT)

		if e != nil {
			(*errPipe).Close()
		}
	})

	return data.RealTimeShortcutRunResp{
		Ok:    true,
		RunId: runMeta.Id,
		Err:   "",
	}
}

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

func (se *ShortcutExecutionService) Exec(st script.Script, scriptType script.ScriptType) (result []byte, err error) {

	runner, ok := RunnerMap[scriptType]

	if !ok {
		return nil, fmt.Errorf("could not find runner for script type : %d", scriptType)
	}

	return runner.Run(st)

}


func (se *ShortcutExecutionService) ExecAsync(st script.Script, scriptType script.ScriptType,handleOut func(inPipe io.WriteCloser, outPipe io.ReadCloser, errPipe io.ReadCloser, err error)) {
	
	runner, ok := RunnerMap[scriptType]

	if !ok {
		handleOut(nil,nil,nil,fmt.Errorf("could not find runner for script type : %d", scriptType))
		return
	}

	handleOut(runner.RunAsync(context.TODO(),st))
}