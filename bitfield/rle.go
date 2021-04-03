package bitfield

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		bytesWritten := binary.PutVarint(crlBuf, currentRunLength)
		crlBuf = append(crlBuf, currentRunByte)
		encodedData = append(encodedData, crlBuf[:bytesWritten]...)
		currentRunByte = b
		currentRunLength = 1
	}

	if len(encodedData) >= len(data) {
		log.Printf("encoded: %d is gt than orig: %d", len(encodedData), len(data))
		return data, false
	}

	return encodedData, true
}

func Decode(encoded []byte) ([]byte, error) {
	log.Println(encoded)
	if len(encoded) == 0 {
		return []byte{}, nil
	}

	decoded := bytes.NewBuffer([]byte{})
	bufReader := bytes.NewReader(encoded)
	log.Printf("buf has %d unread bytes", bufReader.Len())

	for bufReader.Len() > 0 {
		count, err := binary.ReadVarint(bufReader)
		if err != nil {
			return nil, err
		}
		log.Printf("read varint value: %d", count)
		log.Printf("buf has %d unread bytes", bufReader.Len())
		charByte, err := bufReader.ReadByte()
		if err != nil {
			return nil, err
		}

		s := fmt.Sprintf("%d%s", count, string(charByte))
		log.Println(s)
		_, err = decoded.WriteString(s)
		if err != nil {
			return nil, err
		}
	}

	return decoded.Bytes(), nil
}
