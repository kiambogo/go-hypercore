package hypercore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bitfield_SetByte(t *testing.T) {
	t.Parallel()

	bitField := NewBitfield(0)
	changed := bitField.SetByte(0, 1)

	assert.True(t, changed)
	b := bitField.GetByte(0)
	assert.Equal(t, uint8(1), b&1)
	assert.Equal(t, uint8(0), b&8)

	changed = bitField.SetByte(0, 1)
	assert.False(t, changed)

	changed = bitField.SetByte(6, 0)
	assert.False(t, changed)
	changed = bitField.SetByte(6, 1)
	b = bitField.GetByte(6)
	assert.Equal(t, uint8(1), b&1)
	assert.Equal(t, uint8(0), b&8)
}

func Test_Bitfield_SetByteTwice(t *testing.T) {
	t.Parallel()

	bitField := NewBitfield(0)
	changed := bitField.SetByte(0, 1)
	assert.True(t, changed)
	changed = bitField.SetByte(0, 1)
	assert.False(t, changed)
}

func Test_Bitfield_SetBit(t *testing.T) {
	bitField := NewBitfield(0)
	changed := bitField.SetBit(0, true)

	assert.True(t, changed)
	updatedByte := bitField.GetByte(0)
	assert.Equal(t, uint8(1), updatedByte&1)
	assert.Equal(t, uint8(0), updatedByte&8)
}
