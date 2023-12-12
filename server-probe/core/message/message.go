// socket_stick/proto/proto.go
package message

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

// Encode 将消息编码
func Encode(message []byte) ([]byte, error) {
	var length = int64(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	toBytes, err := DecodedToBytes(reader)

	if err != nil {
		return "", err
	}

	return string(toBytes), nil

}

func DecodedToBytesa(reader *bufio.Reader) ([]byte, error) {
	lengthByte, _ := reader.Peek(8) // 读取前4个字节的数据
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int64
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return []byte{}, err
	}
	// Buffered返回缓冲中现有的可读取的字节数。
	if int64(reader.Buffered()) < length+8 {
		return []byte{}, err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(8+length))
	_, err = reader.Read(pack)
	if err != nil {
		return []byte{}, err
	}
	return pack[8:], nil
}

func DecodedToBytes(reader *bufio.Reader) ([]byte, error) {
	// 读取前8个字节的数据
	lengthByte := make([]byte, 8)
	_, err := reader.Read(lengthByte)
	if err != nil {
		return []byte{}, err
	}
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int64
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return []byte{}, err
	}

	// 创建一个有限的读取器
	limitReader := io.LimitReader(reader, length)

	// 读取真正的消息数据
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

type Type int

const (
	REQUEST   Type = 1
	RESPONSE  Type = 2
	HEARTBEAT Type = 3
)

type DataType string

const (
	READDIR      DataType = "DIRREAD"
	DIRREADRESP  DataType = "DIRREADRESP"
	DEFAULT      DataType = "DEFAULT"
	READFILE     DataType = "READFILE"
	READFILERESP DataType = "READFILERESP"
	ERROR        DataType = "ERROR"
)

type Msg struct {
	Type     Type     `json:"type"`
	Id       string   `json:"id"`
	Data     []byte   `json:"data"`
	ErrMark  bool     `json:"errMark"`
	DataType DataType `json:"dataType"`
}
