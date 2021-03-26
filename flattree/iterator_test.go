package flattree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewIterator(t *testing.T) {
	iter := NewIterator(0)
	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on new iterator")
	assert.Equal(t, uint64(0), iter.Offset(), "Offset not correct on new iterator")
	assert.Equal(t, uint64(0), iter.Index(), "Index not correct on new iterator")

	iter = NewIterator(13)
	assert.Equal(t, uint64(4), iter.Factor(), "Factor not correct on new iterator")
	assert.Equal(t, uint64(3), iter.Offset(), "Offset not correct on new iterator")
	assert.Equal(t, uint64(13), iter.Index(), "Index not correct on new iterator")

	iter = NewIterator(18)
	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on new iterator")
	assert.Equal(t, uint64(9), iter.Offset(), "Offset not correct on new iterator")
	assert.Equal(t, uint64(18), iter.Index(), "Index not correct on new iterator")
}

func Test_Iterator(t *testing.T) {
	iter := NewIterator(0)
	iter.Seek(7)

	assert.Equal(t, uint64(16), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(0), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(7), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.Prev()

	assert.Equal(t, uint64(16), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(0), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(7), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.Next()

	assert.Equal(t, uint64(16), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(1), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(23), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, false, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, true, iter.IsRight(), "IsRight not correct on iterator")

	iter.Next()

	assert.Equal(t, uint64(16), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(2), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(39), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.Prev()

	assert.Equal(t, uint64(16), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(1), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(23), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, false, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, true, iter.IsRight(), "IsRight not correct on iterator")

	iter.Sibling()

	assert.Equal(t, uint64(16), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(0), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(7), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.Parent()

	assert.Equal(t, uint64(32), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(0), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(15), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.LeftSpan()

	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(0), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(0), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.Seek(15)
	iter.RightSpan()

	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(15), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(30), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, false, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, true, iter.IsRight(), "IsRight not correct on iterator")

	iter.Seek(25)

	assert.Equal(t, uint64(4), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(6), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(25), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.LeftChild()

	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(12), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(24), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, true, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, false, iter.IsRight(), "IsRight not correct on iterator")

	iter.Seek(25)

	iter.RightChild()

	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(13), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(26), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, false, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, true, iter.IsRight(), "IsRight not correct on iterator")

	iter.Seek(14)

	assert.Equal(t, uint64(2), iter.Factor(), "Factor not correct on iterator")
	assert.Equal(t, uint64(7), iter.Offset(), "Offset not correct on iterator")
	assert.Equal(t, uint64(14), iter.Index(), "Index not correct on iterator")
	assert.Equal(t, false, iter.IsLeft(), "IsLeft not correct on iterator")
	assert.Equal(t, true, iter.IsRight(), "IsRight not correct on iterator")
}

func Test_TwoPow(t *testing.T) {
	assert.Equal(t, uint64(2), twoPow(1))
	assert.Equal(t, uint64(4), twoPow(2))
	assert.Equal(t, uint64(8), twoPow(3))
	assert.Equal(t, uint64(16), twoPow(4))
	assert.Equal(t, uint64(32), twoPow(5))
	assert.Equal(t, uint64(2147483648), twoPow(31))
	assert.Equal(t, uint64(17179869184), twoPow(34))
	assert.Equal(t, uint64(9223372036854775808), twoPow(63))
}
