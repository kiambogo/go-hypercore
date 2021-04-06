package bitfield

import (
	"bytes"
	"encoding/binary"
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
		byteMatch := b == currentRunByte
		atLastByte := i == dataLength-1

		// continued byte match, but end of the encoded data
		if byteMatch && atLastByte {
			currentRunLength++
			encodedData = appendByteCount(encodedData, currentRunLength, currentRunByte)
			break
		}

		// continued byte match, still more encoded data to iterate through
		if byteMatch {
			currentRunLength++
			continue
		}

		// end of the encoded data where the last byte is different than the previous byte
		if atLastByte {
			encodedData = appendByteCount(encodedData, currentRunLength, currentRunByte)
			currentRunByte = b
			currentRunLength = 1
			encodedData = appendByteCount(encodedData, currentRunLength, currentRunByte)
			break
		}

		// different byte found with more encoded data to process
		encodedData = appendByteCount(encodedData, currentRunLength, currentRunByte)
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

		for n := int64(0); n < count; n++ {
			if err := decoded.WriteByte(charByte); err != nil {
				return nil, err
			}
		}
	}

	return decoded.Bytes(), nil
}

func appendByteCount(slice []byte, count int64, elem byte) []byte {
	crlBuf := make([]byte, binary.MaxVarintLen64)
	bytesWritten := binary.PutVarint(crlBuf, count)
	crlBuf = crlBuf[:bytesWritten]
	crlBuf = append(crlBuf, elem)
	slice = append(slice, crlBuf...)

	return slice
}
