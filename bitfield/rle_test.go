package bitfield

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encode(t *testing.T) {
	t.Parallel()

	orig := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABCCCCDDEFGHHH")
	encodedData, encoded := Encode(orig)
	assert.True(t, encoded)
	assert.Equal(t, "", string(encodedData))

}
