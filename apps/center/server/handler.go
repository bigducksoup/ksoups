package server

import (
	"apps/center/server/core"
	"apps/common/message"
	"errors"
)

func HandleResponse(msg message.Msg, serverContext *core.Context) error {

	if msg.Type != message.RESPONSE && msg.Type != message.MULTI_RESPONSE && msg.Type != message.MULTIR_ESPONSE_END {
		return errors.New("msg type should be RESPONSE or MULTIRESPONSE")
	}

	if msg.ErrMark {
		return errors.New("probe error occur")
	}

	hasNext := false

	if msg.Type == message.MULTI_RESPONSE {
		hasNext = true
	}

	if msg.Type == message.MULTIR_ESPONSE_END {
		hasNext = false
	}

	err := serverContext.ReceiveResp(msg.Id, msg, hasNext)

	return err
}
