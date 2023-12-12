package utils

import (
	"crypto/rand"
	"fmt"
	"log"
)

func UUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // set version to 4
	b[8] = (b[8] & 0x3f) | 0x80 // set variant to 10
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
