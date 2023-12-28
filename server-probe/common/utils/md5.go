package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(bytes []byte) string {
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:])
}
