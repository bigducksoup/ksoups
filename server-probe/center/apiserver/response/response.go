package response

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func Success[T any](data T) Resp[T] {
	return Resp[T]{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func Fail(err error) Resp[any] {
	return Resp[any]{
		Code: 500,
		Msg:  err.Error(),
		Data: nil,
	}
}

func StringFail(message string) Resp[any] {

	return Resp[any]{
		Code: 500,
		Msg:  message,
		Data: nil,
	}
}

var paramErrResp = Resp[any]{
	Code: 400,
	Msg:  "参数错误",
	Data: nil,
}

func ParamsError() Resp[any] {
	return paramErrResp
}

var InvalidReqResp = Resp[any]{
	Code: 403,
	Msg:  "非法请求",
	Data: nil,
}

func InvalidReqError() Resp[any] {
	return InvalidReqResp
}
