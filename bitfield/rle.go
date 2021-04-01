package bitfield

import (
	"encoding/binary"
	"log"
)

func Encode(data []byte) ([]byte, bool) {
	encodedData := []byte{}

	if len(data) == 0 {
		return data, false
	}

	currentRunByte := data[0]
	var currentRunLength int64 = 1
	for _, b := range data {
		if b == currentRunByte {
			currentRunLength++
			continue
		}

		crlBuf := make([]byte, binary.MaxVarintLen64)
		_ = binary.PutVarint(crlBuf, currentRunLength)
		crlBuf = append(crlBuf, currentRunByte)
		encodedData = append(encodedData, crlBuf...)
		currentRunByte = b
		currentRunLength = 1
	}

	if len(encodedData) >= len(data) {
		log.Printf("encoded: %d is gt than orig: %d", len(encodedData), len(data))
		return data, false
	}

	return encodedData, true
}

func Decode() {

}
