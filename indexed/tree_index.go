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

func (t tree) Proof(index, digest uint64) (proof Proof, verified bool, err error) {
	// if the node isnt set for the index provided, return no proof
	// always return the hash of the node, even if digest isn't provided???
	// digest & 1 == has_root ??
	// IF digest == 1 and has root, then set remote tree with next (starting at index)
	// then next node will be the sibling to index
	// get all the roots from the right span of next, add to remote tree
	// ELSE
	// go to sibling. if digest is odd and sibling is set, then add sibling to remote tree
	// go to parent and repeat
	// digest = digest/2
	var roots []uint64

	if !t.Get(index) {
		return proof, false, nil
	}

	nodes := []uint64{index}

	if digest == 1 {
		return Proof{
			index:      index,
			verifiedBy: 0,
			nodes:      nodes,
		}, true, nil
	}

	next := index
	hasRoot := digest & 1
	digest >>= 2

	for digest > 0 {
		if digest == 1 && hasRoot != 0 {
			if t.Get(next) {
				_ = t.Set(next)
			}

			nextSibling := ft.Sibling(next)
			if nextSibling < next {
				next = nextSibling
			}

			roots, err = ft.FullRoots(ft.RightSpan(next) + 2)
			if err != nil {
				return
			}
			for _, root := range roots {
				if t.Get(root) {
					_ = t.Set(root)
				}
			}
			break
		}
	}

	sibling := ft.Sibling(next)
	if !isEven(digest) && t.Get(sibling) {
		t.Set(sibling)
	}
	next = ft.Parent(next)
	digest >>= 2

	for !t.Get(next) {
		sibling = ft.Sibling(next)
		if !t.Get(sibling) {
			verifiedBy := t.VerifiedBy(next)
			roots, err = ft.FullRoots(verifiedBy)
			if err != nil {
				return
			}
			for _, root := range roots {
				if root != next && !t.Get(root) {
					nodes = append(nodes, root)
				}
			}
			return Proof{
				index:      index,
				verifiedBy: verifiedBy,
				nodes:      nodes,
			}, false, err
		} else if !t.Get(sibling) {
			nodes = append(nodes, sibling)
		}
		next = ft.Parent(next)
	}

	return Proof{
		index:      index,
		verifiedBy: 0,
		nodes:      nodes,
	}, false, nil
}

// Digest will calculate the digest of the data at a particular index
// It does this by checking the uncles in the merkle tree
func (t tree) Digest(index uint64) (digest uint64) {
	if t.Get(index) {
		return 1
	}

	depthBit := 2
	nextIndex := ft.Sibling(index)
	parentIndex := ft.Parent(index)
	maxTreeIndex := max(nextIndex+2, t.bitfield.Len())

	rightNodesLeftToConsider := func(next uint64) bool { return ft.RightSpan(next) < maxTreeIndex }
	leftNodesLeftToConsider := func(parent uint64) bool { return ft.LeftSpan(parent) > 0 }

	for (rightNodesLeftToConsider(nextIndex)) || (leftNodesLeftToConsider(parentIndex)) {
		if t.Get(nextIndex) {
			digest |= uint64(depthBit)
		}
		if t.Get(parentIndex) {
			digest |= uint64(2*depthBit + 1)
			if (digest & 1) != 1 {
				digest += uint64(1)
			}
			if digest+uint64(1) == uint64(4*depthBit) {
				return 1
			}
		}
		nextIndex = ft.Sibling(parentIndex)
		parentIndex = ft.Parent(nextIndex)
		depthBit *= 2
	}

	return
}

func max(x, y uint64) uint64 {
	if x >= y {
		return x
	}
	return y
}

func isEven(n uint64) bool {
	return n%2 == 0
}
