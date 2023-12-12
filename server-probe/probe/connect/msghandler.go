package connect

import (
	"config-manager/core/message"
	"encoding/json"
	"errors"
)

func HandleMessage(msg message.Msg, connection *Connection) {

	policy := handlePolicy[msg.Type]
	policy(msg, connection)

}

var handlePolicy map[message.Type]func(msg message.Msg, connection *Connection) = map[message.Type]func(msg message.Msg, connection *Connection){
	message.REQUEST: HandleReq,
}

func HandleReq(msg message.Msg, connection *Connection) {

	switch msg.DataType {

	case message.READDIR:

		//将msg.Data解析为对应类型
		dr, err := readData[message.DirRead](msg)

		//调用对应处理方法
		dirReadResponse, err := DirReadHandler(dr)

		if err != nil {
			responseErr(connection, msg.Id, err)
			return
		}

		//设置处理后的返回类型
		msg.DataType = message.DIRREADRESP
		bytes, err := json.Marshal(dirReadResponse)

		if err != nil {
			responseErr(connection, msg.Id, err)
			return
		}

		//设置返回的数据
		msg.Data = bytes

	case message.READFILE:

		fileRead, err := readData[message.FileRead](msg)

		if err != nil {
			responseErr(connection, msg.Id, err)
			return
		}

		fileReadResponse, err := FileReadHandler(fileRead)

		if err != nil {
			responseErr(connection, msg.Id, err)
			return
		}

		msg.DataType = message.DIRREADRESP
		bytes, err := json.Marshal(fileReadResponse)
		if err != nil {
			responseErr(connection, msg.Id, err)
			return
		}

		msg.Data = bytes

	default:
		responseErr(connection, msg.Id, errors.New("unknown datatype"))
		return
	}
	msg.Type = message.RESPONSE
	connection.SendMessage(msg)

}

func responseErr(c *Connection, id string, err error) {

	errResp := message.Msg{
		Type:     message.RESPONSE,
		Id:       id,
		ErrMark:  true,
		DataType: message.ERROR,
		Data:     []byte(err.Error()),
	}

	c.SendMessage(errResp)

}

func readData[T any](msg message.Msg) (T, error) {

	data := new(T)

	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		return *data, err
	}

	return *data, nil
}
