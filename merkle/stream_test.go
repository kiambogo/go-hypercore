package merkle

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var blake2bHasher = BLAKE2b512{}

func Test_NewStream_Empty(t *testing.T) {
	t.Parallel()

	stream := NewStream(blake2bHasher, nil, nil)

	assert.Equal(t, stream.Roots(), new([]Node))
}

func Test_NewStream_SetRootsAndNodes(t *testing.T) {
	t.Parallel()

	roots := &[]Node{
		DefaultNode{
			index:  0,
			parent: 1,
			kind:   leaf,
			data:   []byte("hello, i'm node a"),
			hash:   []byte{},
		},
		DefaultNode{
			index:  2,
			parent: 1,
			kind:   leaf,
			data:   []byte("hello, i'm node b"),
			hash:   []byte{},
		},
	}

	stream := NewStream(blake2bHasher, roots, roots)
	assert.Equal(t, stream.Roots(), roots)
	assert.Equal(t, stream.Nodes(), roots)
}

func Test_NewStream_Append(t *testing.T) {
	t.Parallel()

	stream := NewStream(blake2bHasher, nil, nil)

	stream.Append([]byte("hello, world!"))
	checkNodeCounts(t, 1, 0, stream)

	stream.Append([]byte("foo"))
	checkNodeCounts(t, 2, 1, stream)

	stream.Append([]byte("bar"))
	checkNodeCounts(t, 3, 1, stream)

	stream.Append([]byte("baz"))
	checkNodeCounts(t, 4, 3, stream)
}

func checkNodeCounts(t *testing.T, expectedLeafs, expectedParents int, stream *stream) {
	var leafNodes, parentNodes = 0, 0
	for _, n := range *stream.nodes {
		if n.Kind() == leaf {
			leafNodes++
		} else {
			parentNodes++
		}
	}
	assert.Equal(t, expectedLeafs, leafNodes)
	assert.Equal(t, expectedParents, parentNodes)
	return
}

func Benchmark_Blake2BStream100(b *testing.B) {
	benchmarkBlake2BStream(100, b)
}

func Benchmark_Blake2BStream1000(b *testing.B) {
	benchmarkBlake2BStream(1000, b)
}

func Benchmark_Blake2BStream10000(b *testing.B) {
	benchmarkBlake2BStream(10000, b)
}

func benchmarkBlake2BStream(i int, b *testing.B) {
	stream := NewStream(blake2bHasher, nil, nil)

	for n := 0; n < b.N; n++ {
		for j := 0; j < i; j++ {
			stream.Append([]byte(fmt.Sprint(j)))
		}
	}
}
