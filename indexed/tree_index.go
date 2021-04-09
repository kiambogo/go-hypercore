package indexed

import (
	"github.com/kiambogo/go-hypercore/bitfield"
	ft "github.com/kiambogo/go-hypercore/flattree"
)

type tree struct {
	bitfield *bitfield.Bitfield
}

func NewTree(bitfield *bitfield.Bitfield) tree {
	return tree{
		bitfield: bitfield,
	}
}

func (t tree) Get(index uint64) bool {
	return t.bitfield.GetBit(index)
}

func (t *tree) Set(index uint64) bool {
	// update the element in the tree at index
	if !t.bitfield.SetBit(int(index), true) {
		return false
	}

	// iteratively update the tree, setting the parent of index to true if the sibling is also set
	for t.bitfield.GetBit(ft.Sibling(index)) {
		index = ft.Parent(index)
		if !t.bitfield.SetBit(int(index), true) {
			break
		}
	}

	return true
}

func (t tree) Proof(index uint64) {
}

// Digest will calculate the digest of the data at a particular index
// It does this by checking the uncles in the merkle tree
func (t tree) Digest(index uint64) (digest uint64) {
  if t.Get(index) {
    return 1
  }

  bit := 2
  next := ft.Sibling(index)
  max := max(next+2, t.bitfield.Len())
  parent := ft.Parent(next)

  for (ft.RightSpan(next) < max) || (ft.LeftSpan(parent) > 0) {
    if t.Get(next) {
      digest |= uint64(bit)
    }
    if t.Get(parent) {
      digest |= uint64(2*bit+1)
      if (digest&1) != 1 {
        digest += uint64(1)
      }
      if (digest+uint64(1) == uint64(4*bit)) {
        return 1
      }
    }
    next = ft.Sibling(parent)
    parent = ft.Parent(next)
    bit *= 2
  }

  return
}

func max(x,y uint64) uint64 {
  if x >= y {
    return x
  }
  return y
}
