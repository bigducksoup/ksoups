package protocol

import (
	"bufio"
	"io"
)

type Reader interface {
	ReadHeader() (Header, error)
	Read() (payload []byte, err error)
}

type defaultReader struct {
	originReader *io.Reader
	bufReader    *bufio.Reader
	decoder      *Decoder
}

func (d *defaultReader) ReadHeader() (Header, error) {

	headerBytes := make([]byte, HeaderLength)

	n, err := d.bufReader.Read(headerBytes)

	if err != nil {
		return Header{}, err
	}

	if n != HeaderLength {
		return Header{}, ErrBytesLessThanHeaderLength
	}

	header, err := (*d.decoder).DecodeHeader(headerBytes)

	if err != nil {
		return Header{}, err
	}

	return header, nil
}

func (d *defaultReader) Read() (payload []byte, err error) {

	header, err := d.ReadHeader()

	if err != nil {
		return nil, err
	}

	limitReader := io.LimitReader(d.bufReader, header.bodyLength)

	payload, err = io.ReadAll(limitReader)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewReader(r io.Reader) *Reader {

	bufReader := bufio.NewReader(r)

	var decoder Decoder = &defaultDecoder{}

	var reader Reader = &defaultReader{
		originReader: &r,
		bufReader:    bufReader,
		decoder:      &decoder,
	}

	return &reader
}
