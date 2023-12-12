package connect

import (
	"bufio"
	"config-manager/core/message"
	"config-manager/core/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Connection struct {
	conn      *net.Conn
	addr      string
	time      time.Time
	localAddr string
	reader    *bufio.Reader
}

func createConnection(addr string) Connection {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	connection := Connection{
		conn:      &conn,
		addr:      conn.RemoteAddr().String(),
		time:      time.Now(),
		localAddr: conn.LocalAddr().String(),
		reader:    bufio.NewReader(conn),
	}

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

func (c *Connection) ping() bool {

	ping := message.Msg{
		Type:    message.HEARTBEAT,
		Id:      utils.UUID(),
		Data:    []byte(c.localAddr),
		ErrMark: false,
	}

	err := c.SendMessage(ping)
	if err == io.EOF {
		return false
	}

	if err != nil {
		log.Fatalln(err)
	}

	return true
}

func (c *Connection) reconnect() error {
	if c.localAddr == "" {
		return errors.New("address can not be empty")
	}

	(*c.conn).Close()

	fmt.Println("reconnecting to center......")

	conn, err := net.Dial("tcp", c.localAddr)

	if err != nil {
		return err
	}

	c.conn = &conn

	c.addr = conn.RemoteAddr().String()
	c.localAddr = conn.LocalAddr().String()
	c.time = time.Now()

	reader := bufio.NewReader(conn)

	c.reader = reader

	fmt.Println("reconnected to center")

	return nil

}
