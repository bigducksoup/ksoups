package fs

import (
	"apps/center/server/core"
	"apps/common/message"
	"apps/common/message/data"
	"encoding/json"
)

type OperationService struct {
	ServerCtx *core.Context
}

func (o *OperationService) DeleteFile(probeId string, path string) error {

	fd := data.FileDelete{
		Path: path,
	}

	bytes, err := o.ServerCtx.SendMsgExpectRes(probeId, fd, message.DELETEFILE)

	if err != nil {
		return err
	}

	var res data.FileDeleteResponse

	err = json.Unmarshal(bytes, &res)

	return err
}
