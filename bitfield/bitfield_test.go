package hypercore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bitfield_PageSize(t *testing.T) {
	t.Parallel()

	bitField := NewBitfield(100)
	assert.Equal(t, 100, bitField.PageSize())

	bitField = NewBitfield(200)
	assert.Equal(t, 200, bitField.PageSize())
}

func Test_Bitfield_ByteLength(t *testing.T) {
	t.Parallel()

	bitField := NewBitfield(100)
	assert.Equal(t, uint64(0), bitField.ByteLength())
	bitField.SetBit(10, true)
	assert.Equal(t, uint64(2), bitField.ByteLength())
	bitField.SetBit(100, true)
	assert.Equal(t, uint64(13), bitField.ByteLength())
	bitField.SetBit(101, true)
	assert.Equal(t, uint64(13), bitField.ByteLength())
}

func Test_Bitfield_Len(t *testing.T) {
	t.Parallel()

	bitField := NewBitfield(100)
	assert.Equal(t, uint64(0), bitField.Len())
	bitField.SetBit(10, true)
	assert.Equal(t, uint64(16), bitField.Len())
	bitField.SetBit(100, true)
	assert.Equal(t, uint64(104), bitField.Len())
	bitField.SetBit(101, true)
	assert.Equal(t, uint64(104), bitField.Len())
}

func Test_Bitfield_IsEmpty(t *testing.T) {
	t.Parallel()

	bitField := NewBitfield(100)
	assert.True(t, bitField.IsEmpty())
	bitField.SetBit(10, true)
	assert.False(t, bitField.IsEmpty())
}

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
	t.Parallel()

	bitField := NewBitfield(0)

	bitChanged := bitField.SetBit(0, true)
	assert.True(t, bitChanged)
	updatedByte := bitField.GetByte(0)
	assert.Equal(t, byte(0x1), updatedByte)

	bitChanged = bitField.SetBit(1, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0x3), updatedByte)

	bitChanged = bitField.SetBit(2, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0x7), updatedByte)

	bitChanged = bitField.SetBit(3, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0xf), updatedByte)

	bitChanged = bitField.SetBit(4, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0x1f), updatedByte)

	bitChanged = bitField.SetBit(5, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0x3f), updatedByte)

	bitChanged = bitField.SetBit(6, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0x7f), updatedByte)

	bitChanged = bitField.SetBit(7, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0xff), updatedByte)

	bitChanged = bitField.SetBit(8, true)
	updatedByte = bitField.GetByte(0)
	assert.True(t, bitChanged)
	assert.Equal(t, byte(0xff), updatedByte)
	updatedByte = bitField.GetByte(1)
	assert.Equal(t, byte(0x1), updatedByte)
}

func Test_Bitfield_GetBit(t *testing.T) {
	t.Parallel()

	bitfield := NewBitfield(0)

	_ = bitfield.SetBit(1, true)
	assert.Equal(t, true, bitfield.GetBit(1))

	_ = bitfield.SetBit(8, true)
	assert.Equal(t, true, bitfield.GetBit(8))

	_ = bitfield.SetBit(42, true)
	assert.Equal(t, true, bitfield.GetBit(42))

	_ = bitfield.SetBit(142, true)
	assert.Equal(t, true, bitfield.GetBit(142))

	_ = bitfield.SetBit(1420, true)
	assert.Equal(t, true, bitfield.GetBit(1420))
}
