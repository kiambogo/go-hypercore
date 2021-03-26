package merkle

import (
	"golang.org/x/crypto/blake2b"
)

type NodeHasher interface {
	Node() Node
	HashLeaf(node PartialNode) []byte
	HashParent(left, right Node) []byte
}

type BLAKE2b512 struct{}

func (b2b BLAKE2b512) Node() Node {
	return &DefaultNode{}
}

func (b2b BLAKE2b512) HashLeaf(node PartialNode) []byte {
	hash := []byte{}
	for _, h := range blake2b.Sum512(node.data) {
		hash = append(hash, h)
	}
	return hash
}

func (b2b BLAKE2b512) HashParent(left, right Node) []byte {
	hash := []byte{}
	for _, h := range blake2b.Sum512(append(left.Hash(), right.Hash()...)) {
		hash = append(hash, h)
	}
	return hash
}
