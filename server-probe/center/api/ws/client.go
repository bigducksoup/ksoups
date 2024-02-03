package ws

import (
	"config-manager/common/utils"
	"context"
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
	send              chan *wsMsg
	sendJSON          chan any
	messageHandleFunc func(messageType int, bytes []byte, err error)
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

func (c *Client) SendJSON(obj any) {
	c.sendJSON <- obj
}

func (c *Client) setup(ctx context.Context) error {

	if c.messageHandleFunc == nil {
		return errors.New("messageHandleFunc is nil")
	}

	c.handleSend(ctx)
	c.handleReceive(ctx, c.messageHandleFunc)

	return nil
}

func (c *Client) SetMessageHandleFunc(f func(messageType int, bytes []byte, err error)) {
	c.messageHandleFunc = f
}

// handleReceive 处理接收，回调onReceiveMessage
func (c *Client) handleReceive(ctx context.Context, onReceiveMessage func(int, []byte, error)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				messageType, bytes, err := c.Conn.ReadMessage()
				onReceiveMessage(messageType, bytes, err)
			}
		}
	}()
}

// handleSend 处理发送
func (c *Client) handleSend(ctx context.Context) {
	go func() {
		for {
			select {
			case msg := <-c.send:
				if msg == nil {
					log.Println("未知原因，chan c.send 关闭")
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
			case <-ctx.Done():
				return
			}
		}
	}()
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Id:                utils.UUID(),
		Conn:              conn,
		send:              make(chan *wsMsg),
		messageHandleFunc: nil,
	}
}
