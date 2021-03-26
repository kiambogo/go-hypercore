package flattree

import (
	"errors"
	"math/bits"
)

// Index returns the flat-tree of the node at the provided depth and offset
func Index(depth, offset uint64) uint64 {
	return (offset << (depth + 1)) | ((1 << depth) - 1)
}

// Depth returns the depth of a given node
func Depth(n uint64) uint64 {
	return uint64(bits.TrailingZeros64(^n))
}

// Offset returns the offset of a given node
// The offset is the distance from the left edge of the tree
func Offset(n uint64) uint64 {
	if isEven(n) {
		return n / 2
	}
	return n >> (Depth(n) + 1)
}

// Parent returns the parent node of the provided node
func Parent(n uint64) uint64 {
	return Index(Depth(n)+1, Offset(n)/2)
}

// Sibling returns the sibling of the provided node
func Sibling(n uint64) uint64 {
	return Index(Depth(n), Offset(n)^1)
}

// Uncle returns the parent's sibling of the provided node
func Uncle(n uint64) uint64 {
	return Index(Depth(n)+1, Offset(Parent(n))^1)
}

// Children returns the children of the provided node, if it exists
// Returns the children and a bool indicating if they exist
func Children(n uint64) (left uint64, right uint64, exists bool) {
	if isEven(n) {
		return 0, 0, false
	}

	depth := Depth(n)
	offset := Offset(n) * 2
	left = Index(depth-1, offset)
	right = Index(depth-1, offset+1)

	return left, right, true
}

// LeftChild returns the left child of the provided node, if it exists
// Returns the left child and a bool indicating if it exists
func LeftChild(n uint64) (uint64, bool) {
	if isEven(n) {
		return 0, false
	}

	return Index(Depth(n)-1, Offset(n)*2), true
}

// RightChild returns the left child of the provided node, if it exists
// Returns the right child and a bool indicating if it exists
func RightChild(n uint64) (uint64, bool) {
	if isEven(n) {
		return 0, false
	}

	return Index(Depth(n)-1, (Offset(n)*2)+1), true
}

// Spans returns the left and right most nodes in the tree which the provided node spans
func Spans(n uint64) (left uint64, right uint64) {
	if isEven(n) {
		return n, n
	}
	depth := Depth(n)
	offset := Offset(n)
	left = offset * (2 << depth)
	right = (offset+1)*(2<<depth) - 2
	return
}

// LeftSpan returns the left most node in the tree which the provided node spans
func LeftSpan(n uint64) uint64 {
	if isEven(n) {
		return n
	}
	return Offset(n) * (2 << Depth(n))
}

// RightSpan returns the right most node in the tree which the provided node spans
func RightSpan(n uint64) uint64 {
	if isEven(n) {
		return n
	}
	return (Offset(n)+1)*(2<<Depth(n)) - 2
}

// Count returns the number of nodes under the given node, including the provided node itself
func Count(n uint64) uint64 {
	return (2 << Depth(n)) - 1
}

// FullRoots returns a list of all roots less than the provided index
// A root is a subtrees where all nodes have either 2 or 0 children
func FullRoots(index uint64) (roots []uint64, err error) {
	roots = []uint64{}
	if !isEven(index) {
		err = errors.New("odd index passed to FullRoots()")
		return
	}

	index /= 2
	offset := uint64(0)
	factor := uint64(1)

	for {
		if index == 0 {
			return
		}
		for uint64(factor*2) <= index {
			factor *= 2
		}
		root := offset + factor - 1
		roots = append(roots, root)
		offset += 2 * factor
		index -= factor
		factor = 1
	}
}

func isEven(n uint64) bool {
	return n%2 == 0
}
