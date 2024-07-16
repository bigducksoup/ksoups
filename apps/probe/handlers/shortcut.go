package handlers

import (
	"apps/common/message"
	. "apps/common/message/data"
	"apps/common/utils"
	"apps/probe/script"
	"apps/probe/service"
	"context"
	"sync"
)

func handleRunShortcut(data []byte) (any, message.DataType, error) {

	scRun, err := utils.Unmarshal[ShortcutRun](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	var result any

	script, err := service.ShortcutManage.GetScript(scRun.Id, script.ScriptType(scRun.Type))

	if err != nil {
		return nil, message.ERROR, err
	}

	stdOut, stdErr, err := service.ShortExecutionService.Exec(script, scRun.Timeout)

	if err != nil {
		return nil, message.ERROR, err
	}

	ok := true

	if len(stdErr) != 0 {
		ok = false
	}

	result = ShortcutRunResp{
		Ok:     ok,
		StdOut: string(stdOut),
		StdErr: string(stdErr),
	}

	return result, message.RUN_SHORTCUT_RESP, nil
}

func handleCreateScript(data []byte) (any, message.DataType, error) {

	c, err := utils.Unmarshal[CreateScript](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	script, err := service.ShortcutManage.CreateScript(c.Name, script.ScriptType(c.ScriptType), []byte(c.Content), c.Args)

	if err != nil {
		return nil, message.ERROR, err
	}

	resp := CreateScriptResp{
		Id:      script.Id(),
		Name:    c.Name,
		AbsPath: script.Path(),
	}

	return resp, message.CREATE_SCRIPT_RESP, nil
}

func handleRunshortcutAsync(data []byte) (any, message.DataType, error) {

	runMeta, err := utils.Unmarshal[AsyncShortCutRun](data)

	if err != nil {
		return nil, message.ERROR, err
	}

	script, err := service.ShortcutManage.GetScript(runMeta.Id, script.ScriptType(runMeta.Type))

	if err != nil {
		return nil, message.ERROR, err
	}

	// call exec async
	_, outPipe, errPipe, err := service.ShortExecutionService.ExecAsync(context.TODO(), script, script.Type())

	if err != nil {
		return nil, message.ERROR, err
	}

	// channel for stream
	streamChan := make(chan any, 10)

	w := sync.WaitGroup{}
	w.Add(2)

	// read out pipe
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := outPipe.Read(buf)
			if err != nil {
				outPipe.Close()
				w.Done()
				return
			}

			// add stream to channel
			streamChan <- AsyncShortCutRespStream{
				OutputType: STD_OUT,
				Content:    string(buf[0:n]),
			}
		}
	}()

	// read err pipe
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := errPipe.Read(buf)
			if err != nil {
				errPipe.Close()
				w.Done()
				return
			}

			streamChan <- AsyncShortCutRespStream{
				OutputType: STD_ERR,
				Content:    string(buf[0:n]),
			}
		}
	}()

	// wait all goroutine finishied reading data
	// then close channel
	go func() {
		w.Wait()
		streamChan <- AsyncShortCutRespStream{
			OutputType: END,
		}
		close(streamChan)
	}()

	// return stream channel
	return streamChan, message.SHORTCUT_STREAM, nil

}
