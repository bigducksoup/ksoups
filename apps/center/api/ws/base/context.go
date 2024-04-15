package base

//
//import (
//	"errors"
//	"github.com/gorilla/websocket"
//	"log"
//)
//
//type Context struct {
//	Clients   map[string]*Client
//	RegChan   chan *Client
//	DeRegChan chan *Client
//}
//
//func (c *Context) SendMsg(id string, messageType int, data []byte) error {
//
//	client, ok := c.Clients[id]
//
//	if !ok {
//		return errors.New("could not find Client that id = " + id)
//	}
//
//	client.send <- &wsMsg{
//		messageType: messageType,
//		data:        data,
//	}
//
//	return nil
//}
//
//// Setup 初始化ws.Context
//func (c *Context) Setup() {
//	go func() {
//		for {
//			select {
//			case client := <-c.RegChan:
//				c.Clients[client.Id] = client
//				log.Println(client.Id + " register to WS CTX")
//			case client := <-c.DeRegChan:
//				_, ok := c.Clients[client.Id]
//				if ok {
//					close(client.send)
//					client.Conn.WriteMessage(websocket.CloseInternalServerErr, []byte("something going wrong"))
//					client.Conn.Close()
//					delete(c.Clients, client.Id)
//					log.Println(client.Id + " deregister to WS CTX")
//				}
//			}
//		}
//	}()
//}
//
//func NewContext() *Context {
//	return &Context{
//		Clients:   make(map[string]*Client),
//		RegChan:   make(chan *Client),
//		DeRegChan: make(chan *Client),
//	}
//}
