package server

import (
	"apps/center/api/ws"
	"apps/center/server/core"
	"apps/common/message"
	"apps/common/message/data"
	"encoding/json"
)

func HandleRESPONSE(msg message.Msg, serverContext *core.Context) error {

	if msg.ErrMark == true {

	}

	err := serverContext.ReceiveResp(msg.Id, msg)

	if err != nil {
		return err
	}
	return nil
}

// HandleProActivePush just as its name
// TODO look this
func HandleProActivePush(msg message.Msg, serverContext *core.Context) error {

	if msg.DataType == message.SHORTCUT_OUTPUT {
		outPut := data.RealTimeShortcutOutPut{}
		err := json.Unmarshal(msg.Data, &outPut)

		if err != nil {
			return err
		}

		// TODO record out put and update info in sqlite

		ws.Pusher.SendMsg(ws.Msg[string]{
			Type:    0,
			Payload: string(msg.Data),
		})

		if err != nil {
			return err
		}

	}
	return nil
}

//func HandleHEARTBEAT(msg message.Msg) {
//
//	addr := string(msg.Data)
//
//	value, ok := Ctx.AddrProbe.Load(addr)
//
//	if !ok {
//		log.Println("no such connection")
//		return
//	}
//
//	probe := value.(*Probe)
//
//	probe.LastPingTime = time.Now()
//
//}
