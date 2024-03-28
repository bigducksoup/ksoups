package ws

import "context"

var Ctx *Context

var Pusher *MessagePusher

func Init() {

	Pusher = &MessagePusher{
		client:  make([]*Client, 0),
		regChan: make(chan *Client),
		msgChan: make(chan Msg, 5),
		context: context.TODO(),
	}

	Pusher.StartWork()

	Ctx = newContext()
	Ctx.setup()
}
