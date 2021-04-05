package bitfield

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Encode(data []byte) ([]byte, bool) {
	dataLength := len(data)
	encodedData := []byte{}

	if dataLength <= 1 {
		return data, false
	}

	currentRunByte := data[0]
	var currentRunLength int64 = 0
	for i, b := range data {
		if b == currentRunByte {
			currentRunLength++
			if dataLength-i > 1 {
				continue
			}
		}

		crlBuf := make([]byte, binary.MaxVarintLen64)
		bytesWritten := binary.PutVarint(crlBuf, currentRunLength)
		crlBuf = crlBuf[:bytesWritten]
		crlBuf = append(crlBuf, currentRunByte)
		encodedData = append(encodedData, crlBuf...)
		currentRunByte = b
		currentRunLength = 1
	}

	if len(encodedData) >= len(data) {
		return data, false
	}

	return encodedData, true
}

func Decode(encoded []byte) ([]byte, error) {
	if len(encoded) == 0 {
		return []byte{}, nil
	}

	decoded := bytes.NewBuffer([]byte{})
	bufReader := bytes.NewReader(encoded)

	for bufReader.Len() > 0 {
		count, err := binary.ReadVarint(bufReader)
		if err != nil {
			return nil, err
		}
		charByte, err := bufReader.ReadByte()
		if err != nil {
			return nil, err
		}

		s := fmt.Sprintf("%d%s", count, string(charByte))
		if _, err = decoded.WriteString(s); err != nil {
			return nil, err
		}
	}

	return decoded.Bytes(), nil
}
