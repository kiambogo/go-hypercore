package bitfield

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Varint(t *testing.T) {
	testCases := []int64{
		0,
		1,
		10,
		100,
		999999999,
	}

	for _, tc := range testCases {
		crlBuf := make([]byte, binary.MaxVarintLen64)
		bytesWritten := binary.PutVarint(crlBuf, tc)
		crlBuf = crlBuf[:bytesWritten]

		bufReader := bytes.NewReader(crlBuf)
		count, err := binary.ReadVarint(bufReader)
		assert.NoError(t, err, tc)
		assert.Equal(t, count, tc)
	}
}

func Test_Encode(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name            string
		input           []byte
		shouldEncode    bool
		expectedEncoded []byte
	}{
		{
			name:            "empty input",
			input:           []byte{},
			shouldEncode:    false,
			expectedEncoded: []byte{},
		},
		{
			name:            "smaller decoded than encoded",
			input:           []byte("a"),
			shouldEncode:    false,
			expectedEncoded: []byte("a"),
		},
		{
			name:            "smaller decoded than encoded, 2",
			input:           []byte("abcdefghijklmnopqrstuv"),
			shouldEncode:    false,
			expectedEncoded: []byte("abcdefghijklmnopqrstuv"),
		},
		{
			name:            "encoded, 1",
			input:           []byte("aaa"),
			shouldEncode:    true,
			expectedEncoded: []byte{0x6, 0x61},
		},
		{
			name:            "encoded, 2",
			input:           []byte("aaaaaaaa"),
			shouldEncode:    true,
			expectedEncoded: []byte{0x10, 0x61},
		},
		{
			name:            "encoded, 3",
			input:           []byte("AAABBBCCCCDDDDEFFFFFFFFGGH"),
			shouldEncode:    true,
			expectedEncoded: []byte{0x6, 0x41, 0x6, 0x42, 0x8, 0x43, 0x8, 0x44, 0x2, 0x45, 0x10, 0x46, 0x4, 0x47, 0x2, 0x48},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			encodedData, encoded := Encode(tc.input)
			assert.Equal(t, tc.shouldEncode, encoded, tc.name)
			assert.Equal(t, tc.expectedEncoded, encodedData, tc.name)
		})
	}
}

func Test_Decode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		encoded         []byte
		expectedErr     error
		expectedDecoded []byte
	}{
		{
			name:            "empty input",
			encoded:         []byte{},
			expectedErr:     nil,
			expectedDecoded: []byte{},
		},
		{
			name:            "decode, 1",
			encoded:         []byte("\x02\x41"),
			expectedErr:     nil,
			expectedDecoded: []byte("A"),
		},
		{
			name:            "decode, 2",
			encoded:         []byte("\x14\x41"),
			expectedErr:     nil,
			expectedDecoded: []byte("AAAAAAAAAA"),
		},
		{
			name:            "decode, 3",
			encoded:         []byte("\x02\x41\x14\x42"),
			expectedErr:     nil,
			expectedDecoded: []byte("ABBBBBBBBBB"),
		},
		{
			name:            "invalid, error 1",
			encoded:         []byte("\x42"),
			expectedErr:     io.EOF,
			expectedDecoded: nil,
		},
		{
			name:            "invalid, error 2",
			encoded:         []byte("\x02\x41\x42"),
			expectedErr:     io.EOF,
			expectedDecoded: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			decoded, err := Decode(tc.encoded)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.name)
				assert.Equal(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err, tc.name)
				assert.Equal(t, string(tc.expectedDecoded), string(decoded), tc.name)
			}
		})
	}
}

func Test_EncodeAndDecode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       []byte
		expectedErr error
	}{
		{
			name:        "empty input",
			input:       []byte{},
			expectedErr: nil,
		},
		{
			name:        "valid, 1",
			input:       []byte("aaaaaa"),
			expectedErr: nil,
		},
		{
			name:        "valid, 2",
			input:       []byte("aaabcccccccccddd"),
			expectedErr: nil,
		},
		{
			name:        "valid, 3",
			input:       []byte("aaabcccccccccdddeeeeeeeeeefghi"),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			encoded, _ := Encode(tc.input)

			decoded, err := Decode(encoded)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.name)
				assert.Equal(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err, tc.name)
				assert.Equal(t, string(tc.input), string(decoded), tc.name)
			}
		})
	}
}
