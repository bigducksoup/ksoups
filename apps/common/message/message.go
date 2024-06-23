package message

type Type int

const (
	REQUEST        Type = 1
	RESPONSE       Type = 2
	HEARTBEAT      Type = 3
	REGISTER       Type = 4
	PROACTIVE_PUSH Type = 5
)

type Msg struct {
	Type     Type     `json:"type"`
	Id       string   `json:"id"`
	Data     []byte   `json:"data"`
	ErrMark  bool     `json:"errMark"`
	DataType DataType `json:"dataType"`
}
