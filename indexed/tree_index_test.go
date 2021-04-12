package indexed

import (
	"testing"

	"github.com/kiambogo/go-hypercore/bitfield"
	"github.com/stretchr/testify/assert"
)

func Test_treeDigest(t *testing.T) {
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
