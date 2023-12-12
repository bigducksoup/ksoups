package server

import (
	"config-manager/core/message"
	"config-manager/core/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Probe struct {
	Conn         *net.Conn
	Id           string
	Addr         string
	Time         time.Time
	LastPingTime time.Time
}

func CreateProbe(conn *net.Conn) (p Probe) {

	p = Probe{
		Conn:         conn,
		Id:           utils.UUID(),
		Addr:         (*conn).RemoteAddr().String(),
		Time:         time.Now(),
		LastPingTime: time.Now(),
	}
	return p
}

var Ctx Context

// Context that holds Probes metadata
type Context struct {
	Probes sync.Map
	//key string, value *chan Message.Msg
	respChanMap sync.Map
	AddrProbe   sync.Map
}

func (c *Context) GetProbe(id string) (*Probe, error) {

	value, ok := c.Probes.Load(id)
	if !ok {
		return nil, errors.New("no value matches key")
	}
	res := value.(*Probe)

	return res, nil
}

func (c *Context) GetProbeByAddr(addr string) (*Probe, error) {

	value, ok := c.AddrProbe.Load(addr)
	if !ok {
		return nil, errors.New("no value matches key")
	}
	res := value.(*Probe)

	return res, nil

}

func (c *Context) AddProbe(id string, probe *Probe) {
	c.Probes.Store(id, probe)
	addr := (*probe).Addr
	c.AddrProbe.Store(addr, probe)
	fmt.Printf("added probe id = " + id + "\n")
}

func (c *Context) RemoveProbe(id string) {

	value, ok := c.Probes.Load(id)
	if !ok {
		return
	}

	probe := value.(*Probe)

	fmt.Printf("removing probe: %s", probe.Addr)
	c.Probes.Delete(probe.Id)
	c.AddrProbe.Delete(probe.Addr)

	conn := *(probe.Conn)
	conn.Close()

}

func (c *Context) ReceiveResp(reqId string, resp message.Msg) error {
	value, ok := c.respChanMap.Load(reqId)
	if !ok {
		return errors.New("no such message id")
	}

	channel := *value.(*chan message.Msg)

	channel <- resp

	c.respChanMap.Delete(reqId)
	close(channel)

	return nil
}

// SendData  send msg to probe
func (c *Context) SendData(id string, data any, dataType message.DataType) (msg message.Msg, error error) {

	bytes, _ := json.Marshal(data)

	msg = message.Msg{
		Type:     message.REQUEST,
		Id:       utils.UUID(),
		Data:     bytes,
		ErrMark:  false,
		DataType: dataType,
	}

	p, err := c.GetProbe(id)
	if err != nil {
		return message.Msg{}, err
	}

	marshal, jsonErr := json.Marshal(msg)
	if jsonErr != nil {
		return message.Msg{}, jsonErr
	}

	encode, encodeErr := message.Encode(marshal)
	if encodeErr != nil {
		return message.Msg{}, encodeErr
	}

	conn := *(p.Conn)

	_, err = conn.Write(encode)
	return msg, err
}

func (c *Context) SendMsg(id string, msg message.Msg) error {

	probe, err := c.GetProbe(id)

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	encode, err := message.Encode(bytes)
	if err != nil {
		return err
	}

	conn := *(probe.Conn)

	_, err = conn.Write(encode)
	return err
}

func (c *Context) SendMsgExpectRes(id string, data any, dataType message.DataType) (res []byte, err error) {

	bytes, _ := json.Marshal(data)

	msg := message.Msg{
		Type:     message.REQUEST,
		Id:       utils.UUID(),
		Data:     bytes,
		ErrMark:  false,
		DataType: dataType,
	}

	resChan := make(chan message.Msg)
	c.respChanMap.Store(msg.Id, &resChan)
	sendErr := c.SendMsg(id, msg)
	if sendErr != nil {
		return []byte{}, sendErr
	}

	//等待返回结果，超时后返回错误
	select {
	case res := <-resChan:

		if res.ErrMark {
			return nil, errors.New(string(res.Data))
		}
		return res.Data, nil
	case <-time.After(10 * time.Second):
		c.respChanMap.Delete(msg.Id)
		close(resChan)
		return []byte{}, errors.New("response time out reqId =" + msg.Id)
	}

}

func init() {
	fmt.Println("server ctx initializing")
	Ctx = Context{
		Probes:      sync.Map{},
		respChanMap: sync.Map{},
	}
}
