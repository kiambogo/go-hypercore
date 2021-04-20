package indexed

import (
	"fmt"
	"testing"

	"github.com/kiambogo/go-hypercore/bitfield"
	"github.com/stretchr/testify/assert"
)

func Test_Digest(t *testing.T) {
	testCases := []struct {
		name           string
		index          uint64
		ops            func(tree tree)
		expectedDigest uint64
	}{
		{
			name:           "empty tree",
			index:          0,
			ops:            func(tree tree) {},
			expectedDigest: 0b0,
		},
		{
			name:  "full tree",
			index: 0,
			ops: func(tree tree) {
				tree.Set(0)
			},
			expectedDigest: 0b1,
		},
		{
			name:  "rooted, no sibling, no parent",
			index: 0,
			ops: func(tree tree) {
				tree.Set(1)
			},
			expectedDigest: 0b101,
		},
		{
			name:  "not rooted, has sibling",
			index: 0,
			ops: func(tree tree) {
				tree.Set(2)
			},
			expectedDigest: 0b10,
		},
		{
			name:  "full tree, 2",
			index: 0,
			ops: func(tree tree) {
				tree.Set(1)
				tree.Set(2)
			},
			expectedDigest: 0b1,
		},
		{
			name:  "rooted, sibling, no uncle, grand parents",
			index: 0,
			ops: func(tree tree) {
				tree.Set(3)
				tree.Set(2)
			},
			expectedDigest: 0b1011,
		},
		{
			name:  "not rooted, has sibling",
			index: 1,
			ops: func(tree tree) {
				tree.Set(5)
			},
			expectedDigest: 0b10,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bf := bitfield.NewBitfield(0)
			tree := NewTree(bf)
			tc.ops(tree)
			digest := tree.Digest(tc.index)

			assert.Equal(t, tc.expectedDigest, digest, tc.name)
		})
	}

}

func Test_VerifiedBy(t *testing.T) {
	t.Parallel()

	bf := bitfield.NewBitfield(0)
	tree := NewTree(bf)

	verify := func(index, node, top uint64) {
		verification := tree.VerifiedBy(index)
		assert.Equal(t, Verification{node: node, top: top}, verification, fmt.Sprintf("Index: %d, Node %d, Top %d", index, node, top))
	}

	verify(0, 0, 0)

	tree.Set(0)
	verify(0, 2, 0)

	tree.Set(2)
	verify(0, 4, 4)

	tree.Set(5)
	verify(0, 8, 8)

	tree.Set(8)
	verify(0, 10, 8)

	bf = bitfield.NewBitfield(0)
	tree = NewTree(bf)
	tree.Set(10)
	tree.Set(8)
	tree.Set(13)
	tree.Set(3)
	tree.Set(17)
	verify(10, 20, 20)

	bf = bitfield.NewBitfield(0)
	tree = NewTree(bf)
	tree.Set(7)
	tree.Set(16)
	tree.Set(18)
	tree.Set(21)
	tree.Set(25)
	tree.Set(28)
	verify(16, 30, 28)
	verify(18, 30, 28)
	verify(17, 30, 28)
}
