package core

import (
	"apps/common/message"
	"apps/common/protocol"
	"apps/common/utils"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type Probe struct {
	Conn         *net.Conn
	ReaderWriter *protocol.ReaderWriter
	Id           string
	Addr         string
	Time         time.Time
	LastPingTime time.Time
}

func CreateProbe(conn *net.Conn, readerWriter *protocol.ReaderWriter, id string) (p Probe) {

	p = Probe{
		Conn:         conn,
		Id:           id,
		Addr:         (*conn).RemoteAddr().String(),
		Time:         time.Now(),
		LastPingTime: time.Now(),
		ReaderWriter: readerWriter,
	}
	return p
}

// Context that holds Probes metadata
type Context struct {
	Probes sync.Map
	//key string, value *chan Message.Msg
	respChanMap     sync.Map
	AddrProbe       sync.Map
	responseTimeOut time.Duration
}

func (c *Context) GetProbe(id string) (*Probe, error) {

	value, ok := c.Probes.Load(id)
	if !ok {
		return nil, errors.New("probe not online")
	}
	res := value.(*Probe)

	return res, nil
}

func (c *Context) GetProbeByAddr(addr string) (*Probe, error) {

	value, ok := c.AddrProbe.Load(addr)
	if !ok {
		return nil, errors.New("probe not online")
	}
	res := value.(*Probe)

	return res, nil

}

func (c *Context) AddProbe(id string, probe *Probe) {
	c.Probes.Store(id, probe)
	addr := (*probe).Addr
	c.AddrProbe.Store(addr, probe)
	log.Printf("added probe id = " + id + "\n")
}

func (c *Context) RemoveProbe(id string) {

	value, ok := c.Probes.Load(id)
	if !ok {
		return
	}

	probe := value.(*Probe)

	log.Printf("removing probe: %s", probe.Addr)
	c.Probes.Delete(probe.Id)
	c.AddrProbe.Delete(probe.Addr)

	conn := *(probe.Conn)
	conn.Close()

}

func (c *Context) ReceiveResp(reqId string, resp message.Msg, hasNext bool) error {
	respChan, ok := c.respChanMap.Load(reqId)
	if !ok {
		return errors.New("no such message id")
	}

	channel := respChan.(chan message.Msg)
	channel <- resp

	if !hasNext || resp.DataType == message.ERROR {
		c.respChanMap.Delete(reqId)
		close(channel)
	}

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

	err = (*p.ReaderWriter).Write(marshal)
	return msg, err
}

func (c *Context) SendMsg(id string, msg message.Msg) error {

	p, err := c.GetProbe(id)

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = (*p.ReaderWriter).Write(bytes)
	return err
}

// Request send msg and receive a response
// id probeId string
// data msg any
// dataType type of data message.DataType
func (c *Context) Request(id string, data any, dataType message.DataType) (res []byte, err error) {

	bytes, _ := json.Marshal(data)

	msg := message.Msg{
		Type:     message.REQUEST,
		Id:       utils.UUID(),
		Data:     bytes,
		ErrMark:  false,
		DataType: dataType,
	}

	resChan := make(chan message.Msg)
	c.respChanMap.Store(msg.Id, resChan)
	sendErr := c.SendMsg(id, msg)
	if sendErr != nil {
		return nil, sendErr
	}

	//等待返回结果，超时后返回错误
	select {
	case res := <-resChan:

		if res.DataType == message.ERROR {
			return nil, errors.New(string(res.Data))
		}
		return res.Data, nil
	case <-time.After(c.responseTimeOut):
		c.respChanMap.Delete(msg.Id)
		close(resChan)
		return []byte{}, errors.New("response time out reqId =" + msg.Id)
	}

}

// request and get multi responses
func (c *Context) RequestResponses(probeId string, data any, dataType message.DataType) (res chan []byte, errch chan error, err error) {

	bytes, _ := json.Marshal(data)

	msg := message.Msg{
		Type:     message.SREQUEST,
		Id:       utils.UUID(),
		Data:     bytes,
		ErrMark:  false,
		DataType: dataType,
	}

	resChan := make(chan message.Msg, 10)
	c.respChanMap.Store(msg.Id, resChan)

	err = c.SendMsg(probeId, msg)
	if err != nil {
		return nil, nil, err
	}

	dataChan := make(chan []byte, 10)
	ech := make(chan error, 1)

	go func() {
		for m := range resChan {

			if m.Type == message.MULTIR_ESPONSE_END {
				continue
			}

			if m.DataType == message.ERROR {
				ech <- errors.New(string(m.Data))
				close(ech)
				break
			}
			dataChan <- m.Data
		}
		close(ech)
		close(dataChan)
	}()

	return dataChan, ech, nil

}

func (c *Context) Reset() {
	c.Probes.Range(func(key, value any) bool {
		probe := value.(*Probe)
		conn := *(probe.Conn)
		conn.Close()
		return true
	})

	c.Probes = sync.Map{}
	c.respChanMap = sync.Map{}
	c.AddrProbe = sync.Map{}
}

func NewCenterServerContext(responseTimeout time.Duration) *Context {
	return &Context{
		Probes:          sync.Map{},
		respChanMap:     sync.Map{},
		AddrProbe:       sync.Map{},
		responseTimeOut: responseTimeout,
	}
}
