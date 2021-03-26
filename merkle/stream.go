package merkle

import (
	"sync"

	"github.com/kiambogo/go-hypercore/flattree"
)

type stream struct {
	NodeHasher         // hashing implementation to use when building the merkle tree
	roots      *[]Node // the current set of root nodes in the tree
	nodes      *[]Node // the current set of all nodes in the tree
	blocks     uint64  // size of the tree (might be ok to just len(nodes))
	wg         *sync.Mutex
}

func NewStream(hasher NodeHasher, roots *[]Node, nodes *[]Node) *stream {
	if roots == nil {
		roots = new([]Node)
	}
	if nodes == nil {
		nodes = new([]Node)
	}
	return &stream{
		NodeHasher: hasher,
		roots:      roots,
		nodes:      nodes,
		blocks:     0,
		wg:         &sync.Mutex{},
	}
}

func (s stream) Roots() *[]Node {
	return s.roots
}

func (s stream) Nodes() *[]Node {
	return s.nodes
}

func (s stream) Blocks() uint64 {
	return s.blocks
}

func (s *stream) Append(data []byte) {
	// apply a mutex lock for stream thread safety
	s.wg.Lock()
	defer s.wg.Unlock()

	// construct new node with data from the method argument
	index := uint64(s.blocks * 2)
	leafPartial := PartialNode{
		index:  index,
		parent: flattree.Parent(index),
		data:   data,
		kind:   leaf,
	}
	leaf := s.Node().Build(leafPartial, s.HashLeaf(leafPartial))

	*s.roots = append(*s.roots, leaf)
	*s.nodes = append(*s.nodes, leaf)
	s.blocks++

	for len(*s.roots) > 1 {
		numRoots := len(*s.roots)
		left := (*s.roots)[numRoots-2]
		right := (*s.roots)[numRoots-1]

		if left.Parent() != right.Parent() {
			break
		}

		// construct a new parent node
		newParentPart := PartialNode{
			index:  left.Parent(),
			parent: flattree.Parent(left.Parent()),
			data:   nil,
			kind:   parent,
		}

		newParent := s.Node().Build(newParentPart, s.HashParent(left, right))

		// remove the last two elements of the roots
		*s.roots = (*s.roots)[:len(*s.roots)-2]

		// append new or rehashed parent node to roots and nodes
		*s.roots = append(*s.roots, newParent)
		*s.nodes = append(*s.nodes, newParent)
	}
}
