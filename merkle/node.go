package merkle

import "fmt"

type nodeKind int

const (
	leaf nodeKind = iota
	parent
)

type PartialNode struct {
	index  uint64
	parent uint64
	kind   nodeKind
	data   []byte
}

type Node interface {
	Index() uint64
	Parent() uint64
	Kind() nodeKind
	Hash() []byte
	Build(part PartialNode, hash []byte) Node
}

type DefaultNode struct {
	index  uint64
	parent uint64
	kind   nodeKind
	data   []byte
	hash   []byte
}

func (dn DefaultNode) Index() uint64 {
	return dn.index
}
func (dn DefaultNode) Parent() uint64 {
	return dn.parent
}
func (dn DefaultNode) Data() []byte {
	return dn.data
}
func (dn DefaultNode) Kind() nodeKind {
	return dn.kind
}
func (dn DefaultNode) Hash() []byte {
	return dn.hash
}
func (dn DefaultNode) Build(part PartialNode, hash []byte) Node {
	return DefaultNode{
		index:  part.index,
		parent: part.parent,
		kind:   part.kind,
		data:   part.data,
		hash:   hash,
	}
}

func (dn DefaultNode) String() string {
	return fmt.Sprintf("{Index: %d, Data: %s}", dn.index, dn.data)
}
