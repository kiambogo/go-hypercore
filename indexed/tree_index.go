package indexed

import (
	"github.com/kiambogo/go-hypercore/bitfield"
	. "github.com/kiambogo/go-hypercore/flattree"
)

type tree struct {
	bitfield bitfield.Bitfield
}

func NewTree(bitfield bitfield.Bitfield) tree {
	return tree{
		bitfield: bitfield,
	}
}

func (t tree) Get(index uint64) bool {
	return t.bitfield.GetBit(index)
}

func (t tree) Set(index uint64) bool {
	// update the element in the tree at index
	if !t.bitfield.SetBit(int(index), true) {
		return false
	}

	// iteratively update the tree, setting the parent of index to true if the sibling is also set
	for t.bitfield.GetBit(Sibling(index)) {
		index = Parent(index)
		if !t.bitfield.SetBit(int(index), true) {
			break
		}
	}

	return true
}

func (t tree) Proof() {
}

func (t tree) Digest() {
}
