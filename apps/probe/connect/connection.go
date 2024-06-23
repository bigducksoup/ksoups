package connect

import (
	"apps/common/protocol"
	"io"
	"net"
	"time"
)

type TransportProtocol string

const TCP TransportProtocol = "tcp"
const UDP TransportProtocol = "udp"

type ConnectionOptions struct {
	Address   string
	Transport TransportProtocol
}

type CenterConnection struct {
	conn         *net.Conn
	addr         string
	time         time.Time
	LocalAddr    string
	ReaderWriter *protocol.ReaderWriter
	Options      ConnectionOptions
}

func (c *CenterConnection) Close() {
	(*c.conn).Close()
}

func (c *CenterConnection) Read() ([]byte, error) {
	return (*c.ReaderWriter).Read()
}

func (c *CenterConnection) Write(payload []byte) error {
	return (*c.ReaderWriter).Write(payload)
}

func (c *CenterConnection) Reconnect() error {
	c.Close()
	conn, err := net.Dial(c.Options.Address, string(c.Options.Transport))

	if err != nil {
		return err
	}

	var reader io.Reader = conn
	var writer io.Writer = conn

	readerWriter := protocol.NewReaderWriter(reader, writer)

	c.conn = &conn
	c.ReaderWriter = readerWriter

	return err
}

func NewCenterConnection(options ConnectionOptions) (*CenterConnection, error) {

	conn, err := net.Dial(string(options.Transport), options.Address)

	if err != nil {
		return nil, err
	}

	var reader io.Reader = conn
	var writer io.Writer = conn

	readerWriter := protocol.NewReaderWriter(reader, writer)

	connection := CenterConnection{
		conn:         &conn,
		addr:         conn.RemoteAddr().String(),
		time:         time.Now(),
		LocalAddr:    conn.LocalAddr().String(),
		ReaderWriter: readerWriter,
		Options:      options,
	}

	return &connection, err
}
