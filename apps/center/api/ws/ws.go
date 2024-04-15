package ws

import (
	"apps/center/api/ws/base"
	"context"
)

var Pusher *MessagePusher

func Init() {

	Pusher = &MessagePusher{
		Client:     make([]*base.Client, 0),
		MsgChan:    make(chan Msg[string], 5),
		BinaryChan: make(chan Msg[[]byte], 5),
		RegChan:    make(chan *base.Client, 5),
		Context:    context.TODO(),
	}

	Pusher.StartWork()

}
