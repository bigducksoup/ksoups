package message

type Type int

const (
	REQUEST Type = iota
	RESPONSE
	HEARTBEAT
	REGISTER
	PROACTIVE_PUSH
	SREQUEST
	MULTI_RESPONSE
	MULTIR_ESPONSE_END
)

type Msg struct {
	Type     Type     `json:"type"`
	Id       string   `json:"id"`
	Data     []byte   `json:"data"`
	ErrMark  bool     `json:"errMark"`
	DataType DataType `json:"dataType"`
}
