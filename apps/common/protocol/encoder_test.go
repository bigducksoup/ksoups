package protocol

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestDefaultWorkers(t *testing.T) {

	// encoder and decoder test
	encoder := defaultEncoder{}
	decoder := defaultDecoder{}

	header := encoder.MakeHeader(199)

	decodeHeader, _ := decoder.DecodeHeader(header)

	t.Log(decodeHeader)

	// reader test

	protocolBytes := encoder.Encode([]byte("this is a protocol body!"))

	payload, _ := decoder.Decode(protocolBytes)

	t.Log(string(payload))

	buffer := bytes.NewBuffer(protocolBytes)

	reader := NewReader(buffer)

	payloadFromReader, _ := (*reader).Read()

	t.Log(string(payloadFromReader))

}

func TestDefaultWriter_Write(t *testing.T) {

	buffer := new(bytes.Buffer)
	writer := NewWriter(buffer)

	err := (*writer).Write([]byte("hello world"))

	if err != nil {
		t.Fatal(err)
	}

	reader := NewReader(buffer)

	payload, err := (*reader).Read()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(payload))

}

func fn1() {
	go func() {
		timer := time.NewTicker(1 * time.Second)

		for _ = range timer.C {
			fmt.Println("hello")
		}
	}()
}

func TestNothing(t *testing.T) {

	fn1()

	timer := time.NewTimer(10 * time.Second)

	for _ = range timer.C {

	}

}
