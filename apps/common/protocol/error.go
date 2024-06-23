package protocol

import "errors"

var ErrInvalidMagicNumber = errors.New("magic number is invalid")
var ErrInvalidProtocolVersion = errors.New("version is abandoned")
var ErrBytesLessThanHeaderLength = errors.New("reading bytes from reader,bytes count less than header length")
