package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
)

type Context struct {
	clients   map[string]*Client
	regChan   chan *Client
	deRegChan chan *Client
}

func (c *Context) SendMsg(id string, messageType int, data []byte) error {

	client, ok := c.clients[id]

	if !ok {
		return errors.New("could not find client that id = " + id)
	}

	client.send <- &wsMsg{
		messageType: messageType,
		data:        data,
	}

	return nil
}

// setup 初始化ws.Context
func (c *Context) setup() {
	go func() {
		for {
			select {
			case client := <-c.regChan:
				c.clients[client.Id] = client
				log.Println(client.Id + " register to WS CTX")
			case client := <-c.deRegChan:
				_, ok := c.clients[client.Id]
				if ok {
					close(client.send)
					client.Conn.WriteMessage(websocket.CloseInternalServerErr, []byte("something going wrong"))
					client.Conn.Close()
					delete(c.clients, client.Id)
					log.Println(client.Id + " deregister to WS CTX")
				}
			}
		}
	}()
}

func newContext() *Context {
	return &Context{
		clients:   make(map[string]*Client),
		regChan:   make(chan *Client),
		deRegChan: make(chan *Client),
	}
}
