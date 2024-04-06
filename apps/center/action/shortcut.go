package action

import (
	"apps/center/global"
	"apps/center/model"
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"encoding/json"
	"errors"
	"time"
)

type Runner struct {
}

func (s *Runner) Run(sc model.Shortcut) (string, bool) {

	oneLineShortcutRun := data.ShortcutRun{
		Id:      utils.UUID(),
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		JustRun: sc.JustRun,
		Payload: sc.Payload,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(sc.ProbeId, oneLineShortcutRun, message.RUN_SHORTCUT)

	if err != nil {
		return err.Error(), false
	}

	resp := data.ShortcutRunResp{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return err.Error(), false
	}

	if !resp.Ok {
		return resp.StdErr, false
	}
	return resp.StdOut, true

}

// RealTimeRun run shortcut in a realtime way
// after call it, the output of shortcut run operation will be pushed to center.
// center handle it and push to ui by using websocket
func (s *Runner) RealTimeRun(sc model.Shortcut) (id string, err error) {

	shortcutRun := data.ShortcutRun{
		Id:       utils.UUID(),
		Type:     sc.Type,
		Timeout:  time.Duration(sc.Timeout) * time.Millisecond,
		JustRun:  false,
		Payload:  sc.Payload,
		RealTime: true,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(sc.ProbeId, shortcutRun, message.RUN_SHORTCUT)

	resp := data.RealTimeShortcutRunResp{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		return "", err
	}

	if !resp.Ok {
		return "", errors.New(resp.Err)
	}

	return resp.RunId, nil
}

func (s *Runner) ResultRun(sc *model.Shortcut) (*data.ShortcutRunResp, error) {
	oneLineShortcutRun := data.ShortcutRun{
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		JustRun: sc.JustRun,
		Payload: sc.Payload,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(sc.ProbeId, oneLineShortcutRun, message.RUN_SHORTCUT)

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
		JustRun: sc.JustRun,
		Payload: sc.Payload,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(sc.ProbeId, oneLineShortcutRun, message.RUN_SHORTCUT)

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
