package core

import (
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/protocol"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type CenterServer struct {
	Port    int
	Ctx     *Context
	options CenterServerOptions
	handler CenterServerWorker
	context context.Context
}

type CenterServerOptions struct {
	ResponseTimeOut time.Duration
	RegisterTimeOut time.Duration
	Port            int
	Context         context.Context
	HandlePolicy    map[message.Type]func(msg message.Msg, serverContext *Context) error
}

func (c *CenterServer) Start() error {

	log.Println("Center server initializing")
	//tcp listen port
	u := fmt.Sprintf(":%d", c.Port)
	listener, err := net.Listen("tcp", u)
	if err != nil {
		return err
	}
	log.Println("Center server listening port:", c.Port)

	go func() {
		for {
			select {
			case <-c.context.Done():
				listener.Close()
				c.Ctx.Reset()
				return
			default:
				conn, accErr := listener.Accept()
				if accErr != nil {
					continue
				}
				go c.handler.HandleConnection(&conn, c)
			}
		}
	}()
	log.Println("center server started")
	return nil
}

func (c *CenterServer) SetAuthenticateMethod(authenticate func(info data.RegisterInfo) error) {
	c.handler.Authenticate = authenticate
}

func (c *CenterServer) SetOnProbeRegister(onProbeRegister func(probe *Probe, info data.RegisterInfo)) {
	c.handler.onProbeRegister = onProbeRegister
}

func (c *CenterServer) SetMsgHandlePolicy(msgType message.Type, policy func(msg message.Msg, serverContext *Context) error) {
	c.handler.msgHandlePolicies[msgType] = policy
}

func (c *CenterServer) SetOnProbeOffLine(onProbeOffLine func(probe *Probe)) {
	c.handler.onProbeOffLine = onProbeOffLine
}

func CreateCenterServer(options CenterServerOptions) *CenterServer {
	serverContext := NewCenterServerContext(options.ResponseTimeOut)

	return &CenterServer{
		options: options,
		Ctx:     serverContext,
		Port:    options.Port,
		context: options.Context,
		handler: CenterServerWorker{
			context:           context.WithoutCancel(options.Context),
			msgHandlePolicies: options.HandlePolicy,
		},
	}
}

type CenterServerWorker struct {
	Authenticate      func(info data.RegisterInfo) error
	onProbeRegister   func(probe *Probe, info data.RegisterInfo)
	onProbeOffLine    func(probe *Probe)
	context           context.Context
	msgHandlePolicies map[message.Type]func(msg message.Msg, serverContext *Context) error
}

func (c *CenterServerWorker) HandleConnection(conn *net.Conn, centerServer *CenterServer) {

	var connReader io.Reader = *conn
	var connWriter io.Writer = *conn
	protocolReaderWriter := protocol.NewReaderWriter(connReader, connWriter)

	registerInfo, err := c.ShouldRegister(protocolReaderWriter, centerServer)

	if err != nil {
		log.Println(err)
		(*conn).Close()
		return
	}

	probe := CreateProbe(conn, protocolReaderWriter, registerInfo.Name)
	centerServer.Ctx.AddProbe(probe.Id, &probe)

	if c.onProbeRegister != nil {
		c.onProbeRegister(&probe, *registerInfo)
	}

	for {
		//读取消息
		payload, err := (*protocolReaderWriter).Read()
		if err == io.EOF || len(payload) == 0 {
			//连接断开
			centerServer.Ctx.RemoveProbe(probe.Id)
			c.onProbeOffLine(&probe)
			return
		}

		msg := message.Msg{}
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		err = c.HandleMsg(msg, centerServer.Ctx)

		if err != nil {
			log.Println(err)
			continue
		}
	}

}

func (c *CenterServerWorker) ShouldRegister(readerWriter *protocol.ReaderWriter, centerServer *CenterServer) (*data.RegisterInfo, error) {

	//超时定时器
	timer := time.NewTimer(centerServer.options.RegisterTimeOut)
	msgChan := make(chan message.Msg)
	errChan := make(chan error)

	//开启协程等待注册消息
	go func() {

		//decode, err := message.DecodedToBytes(reader)

		payload, err := (*readerWriter).Read()

		if err == io.EOF || len(payload) == 0 {
			errChan <- errors.New("connection lost")
			return
		}

		registerMsg := message.Msg{}

		err = json.Unmarshal(payload, &registerMsg)
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- registerMsg
	}()

	select {
	case <-timer.C:
		//超时
		return nil, errors.New("register time out,no message send")
	case err := <-errChan:
		//发生错误
		return nil, err
	case msg := <-msgChan:
		registerInfo := data.RegisterInfo{}

		err := json.Unmarshal(msg.Data, &registerInfo)
		if err != nil {
			return nil, err
		}

		_, err = centerServer.Ctx.GetProbe(registerInfo.Name)

		if err == nil {
			m := message.Msg{Id: msg.Id, ErrMark: true, Data: []byte("There is already a probe with id :" + registerInfo.Name)}
			bytes, _ := json.Marshal(m)
			(*readerWriter).Write(bytes)
			return nil, errors.New("There is already a probe with id :" + registerInfo.Name)
		}

		sendRegResp := func(ok bool, data string) {
			m := message.Msg{Id: msg.Id, ErrMark: false, Data: []byte(data)}
			bytes, _ := json.Marshal(m)

			(*readerWriter).Write(bytes)
		}

		if c.Authenticate == nil {
			sendRegResp(true, "ok")
			return &registerInfo, nil
		}

		err = c.Authenticate(registerInfo)

		if err != nil {
			sendRegResp(false, "register info invalid")
			return nil, err
		}
		sendRegResp(true, "ok")
		return &registerInfo, nil
	}

}

func (c *CenterServerWorker) HandleMsg(msg message.Msg, centerContext *Context) error {
	policy, ok := c.msgHandlePolicies[msg.Type]
	if !ok {
		return errors.New("no such policy")
	}
	return policy(msg, centerContext)
}
