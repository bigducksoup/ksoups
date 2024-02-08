package handlers

import (
	"apps/common/message"
	"apps/probe/common"
	"errors"
	"log"
)

func HandleMessage(msg message.Msg, connection *common.Connection) {

	policy, ok := msgHandlePolicy[msg.Type]

	if !ok {
		connection.RespErr(msg.Id, errors.New("unknown message type , could not find handle policy"))
		return
	}

	policy(msg, connection)

}

var msgHandlePolicy = map[message.Type]func(msg message.Msg, connection *common.Connection){
	message.REQUEST: handleReq,
}

func handleReq(msg message.Msg, connection *common.Connection) {

	policy, ok := DataHandlePolicy[msg.DataType]

	if !ok {
		connection.RespErr(msg.Id, errors.New("unknown data type , could not find handle policy"))
		return
	}

	result, err := policy(msg.Data)

	if err != nil {
		connection.RespErr(msg.Id, err)
		return
	}

	msg.Type = message.RESPONSE
	msg.Data = result
	err = connection.SendMessage(msg)
	if err != nil {
		log.Println(err)
	}

}
