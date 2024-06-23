package protocol

import (
	"bufio"
	"io"
)

type Writer interface {
	Write(payload []byte) error
}

type defaultWriter struct {
	encoder      *Encoder
	originWriter *io.Writer
	bufWriter    *bufio.Writer
}

func (d *defaultWriter) Write(payload []byte) error {
	protocolBytes := (*d.encoder).Encode(payload)
	_, err := (*d.bufWriter).Write(protocolBytes)
	err = (*d.bufWriter).Flush()
	return err
}

func NewWriter(w io.Writer) *Writer {

	var encoder Encoder = &defaultEncoder{}

	bufWriter := bufio.NewWriter(w)

	var writer Writer = &defaultWriter{
		encoder:      &encoder,
		originWriter: &w,
		bufWriter:    bufWriter,
	}

	return &writer
}
