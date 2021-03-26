package hypercore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This is probably not testing as expected.
func TestBitfield(t *testing.T) {
	bitField := NewBitfield()
	bitField.Set(0, 1)
	bitField.Set(2, 1)

	assert.Equal(t, true, bitField.Get(0), "Expected bit value to be %v", 1)
	assert.Equal(t, true, bitField.Get(2), "Expected bit value to be %v", 1)
	assert.Equal(t, false, bitField.Get(1), "Expected bit value to be %v", 0)
}
