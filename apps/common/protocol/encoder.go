package protocol

import (
	"bytes"
	"encoding/binary"
)

// Encoder encodes payload to protocol bytes
// it also provides header making function
type Encoder interface {
	Encode(payload []byte) (protocolBytes []byte)
	MakeHeader(payloadLength int64) (header []byte)
}

// default impl for Encoder
type defaultEncoder struct {
}

func (defaultEncoder *defaultEncoder) Encode(payload []byte) (protocolBytes []byte) {

	header := defaultEncoder.MakeHeader(int64(len(payload)))
	protocolBytes = append(header, payload...)
	return protocolBytes
}

func (defaultEncoder *defaultEncoder) MakeHeader(payloadLength int64) (header []byte) {

	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, magicNumber)
	if err != nil {
		panic(err)
	}

	err = binary.Write(buffer, binary.LittleEndian, version)
	if err != nil {
		panic(err)
	}

	err = binary.Write(buffer, binary.LittleEndian, payloadLength)

	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
