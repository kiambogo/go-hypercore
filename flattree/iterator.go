package flattree

type iterator struct {
	index  uint64 // keeps track of the current index of the iterator
	offset uint64 // keeps track of the current offset of the iterator
	factor uint64 // keeps track of the factor of the iterator (2^depth)
}

// NewIterator will construct a new iterator at the designated position
func NewIterator(index uint64) *iterator {
	i := &iterator{}

	i.Seek(index)

	return i
}

// Factor will return the current factor of the iterator
func (i iterator) Index() uint64 {
	return i.index
}

// Factor will return the current factor of the iterator
func (i iterator) Offset() uint64 {
	return i.offset
}

// Factor will return the current factor of the iterator
func (i iterator) Factor() uint64 {
	return i.factor
}

// Seek will position the iterator at the designated index
func (i *iterator) Seek(index uint64) {
	i.index = index
	if isEven(index) {
		i.offset = index / 2
		i.factor = 2
	} else {
		i.offset = Offset(index)
		i.factor = twoPow(Depth(index) + 1)
	}
}

// IsLeft checks if the iterator is currently at a left node
func (i iterator) IsLeft() bool {
	return isEven(i.offset)
}

// IsRight checks if the iterator is currently at a right node
func (i iterator) IsRight() bool {
	return !isEven(i.offset)
}

// Prev moves the iterator to the previous item of the current node, returning its value
func (i *iterator) Prev() uint64 {
	if i.offset == 0 {
		return i.index
	}
	i.offset -= 1
	i.index -= i.factor

	return i.index
}

// Next moves the iterator to the next item of the current node, returning its value
func (i *iterator) Next() uint64 {
	i.offset += 1
	i.index += i.factor

	return i.index
}

// Sibling moves the iterator to the sibling of the current node, returning its value
func (i *iterator) Sibling() uint64 {
	if i.IsLeft() {
		return i.Next()
	}
	return i.Prev()
}

// Parent moves the iterator to the parent of the current node, returning its value
func (i *iterator) Parent() uint64 {
	if isEven(i.offset) {
		i.index += i.factor / 2
		i.offset /= 2
	} else {
		i.index -= i.factor / 2
		i.offset = (i.offset - 1) / 2
	}

	i.factor *= 2
	return i.index
}

// LeftSpan moves the iterator to the left span current node, returning its value
func (i *iterator) LeftSpan() uint64 {
	i.index = i.index + 1 - i.factor/2
	i.offset = i.index / 2
	i.factor = 2

	return i.index
}

// RightSpan moves the iterator to the right span current node, returning its value
func (i *iterator) RightSpan() uint64 {
	i.index = i.index + i.factor/2 - 1
	i.offset = i.index / 2
	i.factor = 2

	return i.index
}

// LeftChild moves the iterator to the left child of the current node, returning its value
func (i *iterator) LeftChild() uint64 {
	if i.factor == 2 {
		return i.index
	}
	i.factor /= 2
	i.index = i.index - i.factor/2
	i.offset *= 2

	return i.index
}

// RightChild moves the iterator to the left child of the current node, returning its value
func (i *iterator) RightChild() uint64 {
	if i.factor == 2 {
		return i.index
	}
	i.factor /= 2
	i.index = i.index + i.factor/2
	i.offset = 2*i.offset + 1

	return i.index
}

// twoPow returns the value of 2 raised to an exponent n, argument to the method
func twoPow(n uint64) uint64 {
	return 1 << n
}
