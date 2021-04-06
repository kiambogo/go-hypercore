package treeindex

import "github.com/kiambogo/go-hypercore/bitfield"

type treeIndex struct {
	bitfield bitfield.Bitfield
}

func NewTreeIndex(bitfield bitfield.Bitfield) treeIndex {
	return treeIndex{
		bitfield: bitfield,
	}
}

func (ti treeIndex) GetBit(index uint64) bool {
	return ti.bitfield.GetBit(index)
}
