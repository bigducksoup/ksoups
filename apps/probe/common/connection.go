package common

import (
	"bufio"
	"config-manager/common/message"
	"config-manager/common/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

type Connection struct {
	conn      *net.Conn
	addr      string
	time      time.Time
	LocalAddr string
	reader    *bufio.Reader
}

func CreateConnection(addr string) Connection {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	connection := Connection{
		conn:      &conn,
		addr:      conn.RemoteAddr().String(),
		time:      time.Now(),
		LocalAddr: conn.LocalAddr().String(),
		reader:    bufio.NewReader(conn),
	}

	log.Println("connected to center : " + addr)

	return connection
}

func (c *Connection) Receive() (message.Msg, error) {

	reader := c.reader

	msg := message.Msg{}

	bytes, err := message.DecodedToBytes(reader)
	if err != nil {
		return msg, err
	}

	if len(bytes) == 0 {
		return message.Msg{}, errors.New("connection reset")
	}

	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return message.Msg{}, err
	}

	return msg, nil

}

func (c *Connection) ReceiveWithCtx(ctx context.Context) (message.Msg, error) {

	msgChan := make(chan message.Msg)
	errChan := make(chan error)

	go func() {
		msg, err := c.Receive()

		if err != nil {
			errChan <- err
			return
		}
		msgChan <- msg
	}()

	select {
	case err := <-errChan:
		return message.Msg{}, err
	case msg := <-msgChan:
		return msg, nil
	case <-ctx.Done():
		return message.Msg{}, ctx.Err()
	}

}

func (c *Connection) SendMessage(msg message.Msg) error {

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	encode, err := message.Encode(bytes)
	if err != nil {
		return err
	}

	_, err = (*c.conn).Write(encode)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) RespErr(id string, err error) {

	errResp := message.Msg{
		Type:     message.RESPONSE,
		Id:       id,
		ErrMark:  true,
		DataType: message.ERROR,
		Data:     []byte(err.Error()),
	}

	c.SendMessage(errResp)

}

func (c *Connection) Ping() bool {

	log.Println("ping center")

	ping := message.Msg{
		Type:    message.HEARTBEAT,
		Id:      utils.UUID(),
		Data:    []byte(c.LocalAddr),
		ErrMark: false,
	}

	err := c.SendMessage(ping)
	if err == io.EOF {
		return false
	}

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (c *Connection) Reconnect() error {
	if c.addr == "" {
		return errors.New("address can not be empty")
	}

	(*c.conn).Close()

	log.Println("reconnecting to center......")

	conn, err := net.Dial("tcp", c.addr)

	if err != nil {
		return err
	}

	c.conn = &conn

	c.addr = conn.RemoteAddr().String()
	c.LocalAddr = conn.LocalAddr().String()
	c.time = time.Now()

	reader := bufio.NewReader(conn)

	c.reader = reader

	log.Println("reconnected to center")

	return nil

}

func (c *Connection) Close() {

	(*(c.conn)).Close()

}
