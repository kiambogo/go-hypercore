package bitfield

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			expectedEncoded: []byte{0x6, 0x41, 0x6, 0x42, 0x8, 0x43, 0x8, 0x44, 0x2, 0x45, 0x10, 0x46, 0x4, 0x47},
		},
	}

	for _, tc := range testCases {
		encodedData, encoded := Encode(tc.input)
		assert.Equal(t, tc.shouldEncode, encoded, tc.name)
		assert.Equal(t, tc.expectedEncoded, encodedData, tc.name)
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
	}

	for _, tc := range testCases {
		decoded, err := Decode(tc.encoded)
		if tc.expectedErr != nil {
			assert.Error(t, err, tc.name)
		}
		assert.Equal(t, tc.expectedDecoded, decoded, tc.name)
	}
}
