package ws

import (
	"apps/center/api/ws/base"
	"context"
)

var Pusher *MessagePusher
var Ctx *base.Context

func Init() {

	Pusher = &MessagePusher{
		Client:  make([]*base.Client, 0),
		MsgChan: make(chan Msg, 5),
		RegChan: make(chan *base.Client, 5),
		Context: context.TODO(),
	}

	Pusher.StartWork()

	Ctx = base.NewContext()
	Ctx.Setup()

}
