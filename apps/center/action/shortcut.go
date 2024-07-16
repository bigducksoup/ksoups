package action

import (
	"apps/center/global"
	"apps/center/model"
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"encoding/json"
	"time"
)

type Runner struct {
}

func (s *Runner) Run(sc model.Shortcut) (stdout string, stderr string, err error) {

	runMeta := data.ShortcutRun{
		Id:      sc.Id,
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		Payload: sc.Payload,
		Args:    sc.Args,
	}

	bytes, err := global.CenterServer.Ctx.Request(sc.ProbeId, runMeta, message.RUN_SHORTCUT)

	if err != nil {
		return "", "", err
	}

	resp := data.ShortcutRunResp{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return "", "", err
	}

	return resp.StdOut, resp.StdErr, nil

}

// RealTimeRun run shortcut in a realtime way
// after call it, the output of shortcut run operation will be pushed to center.
// center handle it and push to ui by using websocket
func (s *Runner) RealTimeRun(sc model.Shortcut) (streams chan data.AsyncShortCutRespStream, err error) {

	shortcutRun := data.AsyncShortCutRun{
		Id:      sc.Id,
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		Payload: sc.Payload,
	}

	bytesChan, errch, err := global.CenterServer.Ctx.RequestResponses(sc.ProbeId, shortcutRun, message.RUN_SHORTCUT_ASYNC)

	if err != nil {
		return nil, err
	}

	streams = make(chan data.AsyncShortCutRespStream, 5)

	writeErr := func(e error) {
		streams <- data.AsyncShortCutRespStream{
			OutputType: data.INTERNAL_ERR,
			Content:    e.Error(),
		}
	}

	go func() {
		for {
			select {
			case err, ok := <-errch:
				if !ok {
					break
				}
				writeErr(err)
				return
			case bytes, ok := <-bytesChan:

				if !ok {
					close(streams)
					return
				}

				stream, err := utils.Unmarshal[data.AsyncShortCutRespStream](bytes)
				if err != nil {
					writeErr(err)
					return
				}
				streams <- stream
			}
		}

	}()

	return streams, nil

}

func (s *Runner) ResultRun(sc *model.Shortcut) (*data.ShortcutRunResp, error) {
	oneLineShortcutRun := data.ShortcutRun{
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		Payload: sc.Payload,
	}

	bytes, err := global.CenterServer.Ctx.Request(sc.ProbeId, oneLineShortcutRun, message.RUN_SHORTCUT)

	if err != nil {
		return nil, err
	}

	resp := data.ShortcutRunResp{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type ShortcutRunner interface {
	Run(shortcut model.Shortcut) (RunResult, error)
}

type NormalShortcutRunner struct {
}

func NewNormalShortcutRunner() *NormalShortcutRunner {
	return &NormalShortcutRunner{}
}

func (n *NormalShortcutRunner) Run(sc model.Shortcut) (RunResult, error) {
	var res RunResult

	oneLineShortcutRun := data.ShortcutRun{
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		Payload: sc.Payload,
	}

	bytes, err := global.CenterServer.Ctx.Request(sc.ProbeId, oneLineShortcutRun, message.RUN_SHORTCUT)

	if err != nil {
		return res, err
	}

	resp := data.ShortcutRunResp{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return res, err
	}

	res.Ok = resp.Ok
	res.StdOut = resp.StdOut
	res.StdErr = resp.StdErr

	return res, nil
}
