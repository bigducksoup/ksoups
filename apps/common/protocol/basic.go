package protocol

// HeaderLength how many bytes header takes
const HeaderLength int = 12

// magicNumber takes 2 bytes
const magicNumber uint16 = 0x0529

// 0 - 2
const magicNumberPos = 0

// version takes 2 bytes
const version uint16 = 0x0001

// 2 - 4
const versionPos = 2

// bodyLength start pos, int64
const bodyLengthPos = 5

type Header struct {
	magicNumber uint16
	version     uint16
	bodyLength  int64
}
