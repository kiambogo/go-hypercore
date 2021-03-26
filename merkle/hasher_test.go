package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"

	b2b "golang.org/x/crypto/blake2b"
)

func Test_BLAKE2b512_Node(t *testing.T) {
	blake2b := BLAKE2b512{}
	assert.Equal(t, &DefaultNode{}, blake2b.Node())
}

func Test_BLAKE2b512_HashLeaf(t *testing.T) {
	blake2b := BLAKE2b512{}

	node := PartialNode{
		index:  10,
		parent: 0,
		kind:   leaf,
		data:   []byte("greetings"),
	}

	leafHash := blake2b.HashLeaf(node)
	compareBytes(t, leafHash, b2b.Sum512([]byte("greetings")))
}

func Test_BLAKE2b512_HashParent(t *testing.T) {
	blake2b := BLAKE2b512{}

	leftHash := toDynamicallySizedBuffer(b2b.Sum512([]byte("hello")))
	rightHash := toDynamicallySizedBuffer(b2b.Sum512([]byte("world")))

	left := DefaultNode{
		index:  10,
		parent: 0,
		kind:   leaf,
		data:   []byte("hello"),
		hash:   leftHash,
	}
	right := DefaultNode{
		index:  10,
		parent: 0,
		kind:   leaf,
		data:   []byte("world"),
		hash:   rightHash,
	}

	parentHash := blake2b.HashParent(left, right)
	expected := b2b.Sum512(append(leftHash, rightHash...))
	compareBytes(t, parentHash, expected)
}

func toDynamicallySizedBuffer(bytes [64]byte) []byte {
	b := []byte{}
	for _, byte := range bytes {
		b = append(b, byte)
	}
	return b
}

func compareBytes(t *testing.T, expected []byte, actual [64]byte) {
	assert.Equal(t, len(expected), len(actual))
	for i, b := range expected {
		assert.Equal(t, b, actual[i])
	}
}
