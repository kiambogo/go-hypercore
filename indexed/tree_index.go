package indexed

import "github.com/kiambogo/go-hypercore/bitfield"

type tree struct {
	bitfield bitfield.Bitfield
}

func NewTree(bitfield bitfield.Bitfield) tree {
	return tree{
		bitfield: bitfield,
	}
}

func (t tree) GetBit(index uint64) bool {
	return t.bitfield.GetBit(index)
}
