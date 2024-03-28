package ws

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Msg struct {
	Type    int `json:"type"`
	Payload any
}

type MessagePusher struct {
	client  []*Client
	msgChan chan Msg
	regChan chan *Client
	context context.Context
}

func (s *MessagePusher) RegMe(client *Client) {
	s.regChan <- client
}

func (s *MessagePusher) SendMsg(msg Msg) {
	s.msgChan <- msg
}

// StartWork after create a MessagePusher ,StartWork should be called.
// just common goroutine + for + select + channels
// 3 channels are being selected here
// s.regChan for register client
// s.msgChan for send message to all client
// s.context.Done() for controlling
func (s *MessagePusher) StartWork() {
	go func() {
		for {
			select {
			case client := <-s.regChan:
				s.client = append(s.client, client)
			case msg := <-s.msgChan:
				bytes, err := json.Marshal(msg)

				if err != nil {
					log.Println(err)
					break
				}
				for _, client := range s.client {
					client.Send(websocket.TextMessage, bytes)
				}
			case <-s.context.Done():
				close(s.regChan)
				close(s.msgChan)
				return
			}
		}
	}()
}

// DoMessagePush after open a websocket connection ,register client to message pusher,
// you can call Pusher.SendMsg to push any message to registered clients.
// I make this mainly for real time monitoring shortcut run result.
// In the future, this ws connection could be used for any proactive message pushing.
func DoMessagePush(client *Client, c *gin.Context) {
	Pusher.RegMe(client)

	ctx, cancelFunc := context.WithCancel(Pusher.context)

	client.SetMessageHandleFunc(func(messageType int, bytes []byte, err error) {
		if err != nil {
			cancelFunc()
			Ctx.deRegChan <- client
		}
	})
	err := client.setup(ctx)

	if err != nil {
		return
	}

	select {
	case <-ctx.Done():
		return
	}
}
