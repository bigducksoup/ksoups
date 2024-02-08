package action

import (
	"apps/center/model"
	"apps/center/server"
	"apps/common/message"
	"apps/common/message/data"
	"encoding/json"
	"time"
)

type Runner struct {
}

func (s *Runner) Run(sc model.Shortcut) (string, bool) {

	oneLineShortcutRun := data.ShortcutRun{
		Type:    sc.Type,
		Timeout: time.Duration(sc.Timeout) * time.Millisecond,
		JustRun: sc.JustRun,
		Payload: sc.Payload,
	}

	bytes, err := server.Ctx.SendMsgExpectRes(sc.ProbeId, oneLineShortcutRun, message.RUN_SHORTCUT)

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
