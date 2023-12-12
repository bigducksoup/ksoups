package server

import (
	"config-manager/core/message"
	"fmt"
	"time"
)

func HandleRESPONSE(msg message.Msg) {

	if msg.ErrMark == true {

	}

	err := Ctx.ReceiveResp(msg.Id, msg)

	if err != nil {
		fmt.Println("receive response failed")
		fmt.Println(err)
	}

}

func HandleHEARTBEAT(msg message.Msg) {

	addr := string(msg.Data)

	value, ok := Ctx.AddrProbe.Load(addr)

	if !ok {
		fmt.Println("no such connection")
		return
	}

	probe := value.(*Probe)

	probe.LastPingTime = time.Now()

}
