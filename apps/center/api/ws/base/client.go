package base

import (
	"apps/common/utils"
	"errors"
	"github.com/gorilla/websocket"
	"log"
)

type wsMsg struct {
	messageType int
	data        []byte
}

type Client struct {
	Id   string
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	send      chan *wsMsg
	sendJSON  chan any
	closeChan chan struct{}
}

func (c *Client) Send(messageType int, data []byte) {

	msg := &wsMsg{
		messageType: messageType,
		data:        data,
	}
	if msg == nil {
		log.Println("nil message ??????")
	}
	c.send <- msg
}

func (c *Client) SendJSON(obj any) error {
	c.sendJSON <- obj
	return nil
}

func (c *Client) Setup(onReceiveMessage func(int, []byte, error)) error {

	if onReceiveMessage == nil {
		return errors.New("messageHandleFunc is nil")
	}

	go func() {
		for {
			select {
			case msg := <-c.send:
				if msg == nil {
					log.Println("chan c.send 关闭")
					return
				}
				err := c.Conn.WriteMessage(msg.messageType, msg.data)
				if err != nil {
					return
				}
			case obj := <-c.sendJSON:
				err := c.Conn.WriteJSON(obj)
				if err != nil {
					return
				}
			case <-c.closeChan:
				return
			}
		}
	}()

	go func() {
		for {
			messageType, bytes, err := c.Conn.ReadMessage()
			onReceiveMessage(messageType, bytes, err)
			if err != nil {
				return
			}
		}
	}()

	return nil
}

func (c *Client) Close() {
	c.closeChan <- struct{}{}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Id:        utils.UUID(),
		Conn:      conn,
		send:      make(chan *wsMsg, 10),
		sendJSON:  make(chan any, 10),
		closeChan: make(chan struct{}, 1),
	}
}
