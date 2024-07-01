package connect

import (
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"sync"
	"time"
)

type ProbeOptions struct {
	Address           string
	Ping              bool
	PingInterval      time.Duration
	RequestTimeOut    time.Duration
	Context           context.Context
	Reconnect         bool
	MaxReconnectCount int
	ReconnectGapTime  time.Duration
	RegisterInfo      data.RegisterInfo
	DataHandlers      map[message.DataType]func(data []byte) (any, message.DataType, error)
	BeforeStart       func(*Probe)
	Encoder           func(v any) ([]byte, error)
	Decoder           func(bytes []byte, v any) error
}

type Probe struct {
	toCenterConnection *CenterConnection
	Context            context.Context
	respChanMap        sync.Map
	LocalAddress       string
	ProbeOptions       ProbeOptions
	DataHandlers       map[message.DataType]func(data []byte) (any, message.DataType, error)
	BeforStart         func(*Probe)
	Encoder            func(v any) ([]byte, error)
	Decoder            func(bytes []byte, v any) error
}

func (p *Probe) StartWorking() {

	if p.BeforStart != nil {
		p.BeforStart(p)
	}

	p.RegisterToCenter()

	go func() {
		for {
			// loop read from connection
			payload, err := p.toCenterConnection.Read()
			if err != nil {
				// connection reset
				if err == io.EOF {
					// TODO reconnect
					log.Println("Lost connection to center")
					break
				}

				p.ReportErr(err)
				continue
			}

			var msg message.Msg

			err = p.Decoder(payload, &msg)

			if err != nil {
				p.ReportErr(err)
				continue
			}

			if msg.ErrMark {
				// handle msg error
				continue
			}

			switch msg.Type {
			case message.RESPONSE:
				// TODO handle response
				p.respChanMap.Store(msg.Id, msg)
			case message.REQUEST:
				// find data handler
				dataHandler, ok := p.DataHandlers[msg.DataType]

				if !ok {
					p.ReportErr(errors.New("no dataHandler could be found for this message"))
					continue
				}

				// async use handler
				go func() {
					response, dataType, handleErr := dataHandler(msg.Data)

					if handleErr != nil {
						p.ResponseErr(handleErr, msg)
						return
					}

					err = p.ResponseToCenter(msg.Id, response, dataType)

					if err != nil {
						p.ReportErr(err)
						return
					}

				}()
			case message.HEARTBEAT:
				// handle heartbeat
			case message.PROACTIVE_PUSH:
				// handle push
			}

		}
	}()

}

func (p *Probe) RegisterToCenter() error {
	return p.SendToCenter(utils.UUID(), p.ProbeOptions.RegisterInfo, message.DEFAULT, message.REGISTER, false)
}

func (p *Probe) SendToCenter(id string, v any, dataType message.DataType, messageType message.Type, errMark bool) error {

	var bytes []byte

	switch v := v.(type) {
	case string:
		bytes = []byte(v)
	default:
		encodedData, err := p.Encoder(v)

		if err != nil {
			return err
		}
		bytes = encodedData
	}

	msg := message.Msg{
		Type:     messageType,
		Id:       id,
		Data:     bytes,
		DataType: dataType,
		ErrMark:  errMark,
	}

	msgBytes, err := p.Encoder(msg)

	if err != nil {
		return err
	}

	err = p.toCenterConnection.Write(msgBytes)

	return err

}

func (p *Probe) PushToCenter(v any, dataType message.DataType) error {

	return p.SendToCenter(utils.UUID(), v, dataType, message.PROACTIVE_PUSH, false)

}

func (p *Probe) ResponseToCenter(id string, body any, dataType message.DataType) error {
	return p.SendToCenter(id, body, dataType, message.RESPONSE, false)
}

// Request send a request to center, center must return response.
// this function will return timeout err if center did not  response in time.
func (p *Probe) Request(body any, dataType message.DataType, receiver any) error {
	bytes, err := p.Encoder(body)

	if err != nil {
		return err
	}

	msg := message.Msg{
		Type:     message.REQUEST,
		Id:       utils.UUID(),
		Data:     bytes,
		DataType: dataType,
		ErrMark:  false,
	}

	msgBytes, err := p.Encoder(msg)

	if err != nil {
		return err
	}

	p.respChanMap.Store(msg.Id, make(chan message.Msg))
	err = p.toCenterConnection.Write(msgBytes)

	if err != nil {
		return err
	}

	respChan, ok := p.respChanMap.Load(msg.Id)

	if !ok {
		return errors.New("can not find response channel")
	}

	timer := time.NewTimer(p.ProbeOptions.RequestTimeOut)

	select {
	case <-timer.C:
		p.respChanMap.Delete(msg.Id)
		return errors.New("request time out")
	case <-p.Context.Done():
		return p.Context.Err()
	case respMsg := <-respChan.(chan message.Msg):
		if respMsg.ErrMark {
			return errors.New(string(respMsg.Data))
		}

		err = p.Decoder(respMsg.Data, receiver)
		return err
	}

}

func (p *Probe) ResponseErr(err error, originMessage message.Msg) error {
	log.Printf("Response error to center, error : %s", err.Error())

	return p.SendToCenter(originMessage.Id, err.Error(), message.ERROR, message.RESPONSE, true)
}

// ReportErr report error message to center
func (p *Probe) ReportErr(err error) error {

	body := []byte(err.Error())

	log.Printf("Report error to center, error : %s", err.Error())

	msg := message.Msg{
		Type:     message.REQUEST,
		Id:       utils.UUID(),
		Data:     body,
		DataType: message.ERROR,
		ErrMark:  false,
	}

	msgBytes, err := p.Encoder(msg)

	if err != nil {
		return err
	}

	return p.toCenterConnection.Write(msgBytes)
}

var ProbeInstance *Probe
var once = sync.Once{}

func InitProbe(options ProbeOptions) {
	once.Do(func() {

		log.Println("Initializing Probe")

		log.Println("Opening connection to center")
		toCenterConnection, err := NewCenterConnection(ConnectionOptions{
			Address:   options.Address,
			Transport: TCP,
		})

		if err != nil {
			panic(err)
		}

		log.Println("successfully opened connection to center")

		decoder := options.Decoder
		if decoder == nil {
			decoder = json.Unmarshal
		}

		encoder := options.Encoder
		if encoder == nil {
			encoder = json.Marshal
		}

		options.RequestTimeOut = utils.UnexpectThenDefault(options.RequestTimeOut, 0, 10*time.Second)
		options.MaxReconnectCount = utils.UnexpectThenDefault(options.MaxReconnectCount, 0, 0)
		options.ReconnectGapTime = utils.UnexpectThenDefault(options.ReconnectGapTime, 0, 5*time.Second)
		options.PingInterval = utils.UnexpectThenDefault(options.PingInterval, 0, 5*time.Second)

		ProbeInstance = &Probe{
			toCenterConnection: toCenterConnection,
			Context:            options.Context,
			respChanMap:        sync.Map{},
			LocalAddress:       toCenterConnection.LocalAddr,
			ProbeOptions:       options,
			Encoder:            encoder,
			Decoder:            decoder,
			DataHandlers:       options.DataHandlers,
		}

	})
}
