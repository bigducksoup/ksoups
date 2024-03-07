package server

import (
	"apps/center/server/ServerContext"
	"apps/common/message"
	"log"
	"time"
)

func HandleRESPONSE(msg message.Msg) {

	if msg.ErrMark == true {

	}

	err := ServerContext.Ctx.ReceiveResp(msg.Id, msg)

	if err != nil {
		log.Println("receive response failed")
		log.Println(err)
	}

}

func HandleHEARTBEAT(msg message.Msg) {

	addr := string(msg.Data)

	value, ok := ServerContext.Ctx.AddrProbe.Load(addr)

	if !ok {
		log.Println("no such connection")
		return
	}

	probe := value.(*ServerContext.Probe)

	probe.LastPingTime = time.Now()

}
