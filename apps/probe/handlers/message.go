package handlers

import (
	"apps/common/message"
	"apps/probe/connect"
	"errors"
)

var MessageHandler = map[message.Type]func(m message.Msg, p *connect.Probe){
	message.REQUEST:  HandleRequest,
	message.SREQUEST: HandleSRequest,
}

func HandleRequest(m message.Msg, p *connect.Probe) {
	// find data handler
	dataHandler, ok := DataHandlePolicy[m.DataType]

	if !ok {
		p.ReportErr(errors.New("no dataHandler could be found for this message"))
		return
	}

	response, dataType, handleErr := dataHandler(m.Data)

	if handleErr != nil {
		p.ResponseErr(handleErr, m)
		return
	}

	err := p.ResponseToCenter(m.Id, response, dataType)

	if err != nil {
		p.ReportErr(err)
	}

}

func HandleSRequest(m message.Msg, p *connect.Probe) {
	// find data handler
	dataHandler, ok := DataHandlePolicy[m.DataType]

	if !ok {
		p.ReportErr(errors.New("no dataHandler could be found for this message"))
		return
	}

	x, dataType, handleErr := dataHandler(m.Data)

	if handleErr != nil {
		p.ResponseErr(handleErr, m)
		return
	}

	// must return a chan any, and the dataType declears item in chan
	xchan, ok := x.(chan any)

	if !ok {
		p.ResponseErr(errors.New("assertion fail"), m)
		return
	}

	for resp := range xchan {
		p.SendToCenter(m.Id, resp, dataType, message.MULTI_RESPONSE, false)
	}

	p.SendToCenter(m.Id, nil, dataType, message.MULTIR_ESPONSE_END, false)
}
