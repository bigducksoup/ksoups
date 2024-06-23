package protocol

import (
	"bytes"
	"encoding/binary"
)

// Decoder decodes protocol bytes to origin payload
// DecodeHeader func will return a header from headerBytes
// see Header
type Decoder interface {
	Decode(protocolBytes []byte) (payload []byte, err error)
	DecodeHeader(headerBytes []byte) (header Header, err error)
}

// default impl for Decoder
type defaultDecoder struct {
}

func (defaultDecoder *defaultDecoder) Decode(protocolBytes []byte) (payload []byte, err error) {

	headerBytes := protocolBytes[0:HeaderLength]

	header, err := defaultDecoder.DecodeHeader(headerBytes)

	if err != nil {
		return nil, err
	}

	bodyBytes := protocolBytes[HeaderLength:(HeaderLength + int(header.bodyLength))]

	body := make([]byte, header.bodyLength)

	copy(body, bodyBytes)

	return body, err
}

func (defaultDecoder *defaultDecoder) DecodeHeader(headerBytes []byte) (header Header, err error) {

	buffer := bytes.NewBuffer(headerBytes)

	err = binary.Read(buffer, binary.LittleEndian, &(header.magicNumber))

	if err != nil {
		return Header{}, err
	}
	// check magic number
	if header.magicNumber != magicNumber {
		return Header{}, ErrInvalidMagicNumber
	}

	err = binary.Read(buffer, binary.LittleEndian, &(header.version))

	if err != nil {
		return Header{}, err
	}

	// check version
	if header.version != version {
		return Header{}, ErrInvalidProtocolVersion
	}

	err = binary.Read(buffer, binary.LittleEndian, &(header.bodyLength))

	if err != nil {
		return Header{}, err
	}

	return header, nil
}
