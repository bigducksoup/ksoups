package core

import (
	"apps/common/message"
)

func HandleRESPONSE(msg message.Msg, serverContext *Context) error {

	if msg.ErrMark == true {

	}

	err := serverContext.ReceiveResp(msg.Id, msg)

	if err != nil {
		return err
	}
	return nil
}

// TODO HandleHEARTBEAT is not implemented
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
