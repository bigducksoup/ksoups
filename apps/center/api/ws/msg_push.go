package ws

import (
	"apps/center/api/ws/base"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Msg[T string | []byte] struct {
	// 0 : binary
	// 1 : text
	Type    int `json:"type"`
	Payload T
}

type MessagePusher struct {
	Client     []*base.Client
	MsgChan    chan Msg[string]
	BinaryChan chan Msg[[]byte]
	RegChan    chan *base.Client
	UnRegChan  chan *base.Client
	Context    context.Context
}

func (s *MessagePusher) RegMe(client *base.Client) {
	s.RegChan <- client
}

func (s *MessagePusher) unRegMe(client *base.Client) {
	s.UnRegChan <- client
}

func (s *MessagePusher) SendMsg(msg Msg[string]) {
	s.MsgChan <- msg
}

func (s *MessagePusher) SendBinary(msg Msg[[]byte]) {
	s.BinaryChan <- msg
}

// StartWork after create a MessagePusher ,StartWork should be called.
// just common goroutine + for + select + channels
// 3 channels are being selected here
// s.RegChan for register Client
// s.MsgChan for send message to all Client
// s.Context.Done() for controlling
func (s *MessagePusher) StartWork() {
	go func() {
		for {
			select {
			case client := <-s.RegChan:
				s.Client = append(s.Client, client)

			case client := <-s.UnRegChan:
				for index := range s.Client {
					if s.Client[index] == client {
						s.Client = append(s.Client[0:index], s.Client[index+1:]...)
						client.Close()
						break
					}
				}
			case msg := <-s.MsgChan:
				bytes, err := json.Marshal(msg.Payload)

				if err != nil {
					log.Println(err)
					break
				}
				for _, client := range s.Client {
					client.Send(websocket.TextMessage, bytes)
				}
			case binary := <-s.BinaryChan:
				for _, client := range s.Client {
					client.Send(websocket.BinaryMessage, binary.Payload)
				}
			case <-s.Context.Done():
				close(s.RegChan)
				close(s.MsgChan)
				close(s.UnRegChan)
				return
			}
		}
	}()
}

// DoMessagePush after open a websocket connection ,register Client to message pusher,
// you can call Pusher.SendMsg to push any message to registered clients.
// I make this mainly for real time monitoring shortcut run result.
// In the future, this ws connection could be used for any proactive message pushing.
func DoMessagePush(client *base.Client, c *gin.Context) {
	Pusher.RegMe(client)

	ctx := context.WithoutCancel(Pusher.Context)

	err := client.Setup(func(messageType int, bytes []byte, err error) {
		if err != nil {
			Pusher.unRegMe(client)
		}
	})

	if err != nil {
		return
	}

	select {
	case <-ctx.Done():
		return
	}
}
