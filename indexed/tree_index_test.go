package indexed

import (
	"fmt"
	"testing"

	"github.com/kiambogo/go-hypercore/bitfield"
	"github.com/stretchr/testify/assert"
)

// Small tree example
//       3
//   1       5
// 0   2   4   6

func Test_GetAndSet(t *testing.T) {
	t.Parallel()

	t.Run("get an unset index returns false", func(t *testing.T) {
		t.Parallel()
		tree := NewDefaultTree()
		assert.False(t, tree.Get(0))
	})

	t.Run("set an already set index returns false", func(t *testing.T) {
		t.Parallel()
		tree := NewDefaultTree()
		assert.True(t, tree.Set(0))
		assert.False(t, tree.Set(0))
	})

	t.Run("set iteratively updates the parent if sibling is also set", func(t *testing.T) {
		t.Parallel()
		tree := NewDefaultTree()

		assert.False(t, tree.Get(1))

		tree.Set(0)
		assert.False(t, tree.Get(1))

		tree.Set(2)
		assert.True(t, tree.Get(1))

		tree.Set(4)
		tree.Set(6)
		assert.True(t, tree.Get(5))
		assert.True(t, tree.Get(3))
	})

	t.Run("set when sibling and parent are already set", func(t *testing.T) {
		t.Parallel()
		tree := NewDefaultTree()
		tree.Set(2)
		tree.Set(1)

		assert.True(t, tree.Set(0))
	})
}

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

			tree := NewDefaultTree()
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

func Test_ProofWithoutDigest1(t *testing.T) {
	t.Parallel()

	tree := NewDefaultTree()
	proof, verified, err := tree.Proof(0, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 0, nodes: nil}, proof)
	assert.False(t, verified)

	tree.Set(0)
	proof, verified, err = tree.Proof(0, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 2, nodes: []uint64{0}}, proof)
	assert.True(t, verified)

	tree.Set(2)
	proof, verified, err = tree.Proof(0, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 4, nodes: []uint64{0, 2}}, proof)
	assert.True(t, verified)

	tree.Set(5)
	proof, verified, err = tree.Proof(0, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 8, nodes: []uint64{0, 2, 5}}, proof)
	assert.True(t, verified)

	tree.Set(8)
	proof, verified, err = tree.Proof(0, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 10, nodes: []uint64{0, 2, 5, 8}}, proof)
	assert.True(t, verified)
}

func Test_ProofWithoutDigest2(t *testing.T) {
	t.Parallel()

	tree := NewDefaultTree()
	tree.Set(10)
	tree.Set(8)
	tree.Set(13)
	tree.Set(3)
	tree.Set(17)
	proof, verified, err := tree.Proof(10, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 10, verifiedBy: 20, nodes: []uint64{10, 8, 13, 3, 17}}, proof)
	assert.True(t, verified)
}

func Test_ProofWithoutDigest3(t *testing.T) {
	t.Parallel()

	tree := NewDefaultTree()
	tree.Set(7)
	tree.Set(16)
	tree.Set(18)
	tree.Set(21)
	tree.Set(25)
	tree.Set(28)
	proof, verified, err := tree.Proof(16, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 16, verifiedBy: 30, nodes: []uint64{16, 18, 21, 7, 25, 28}}, proof)
	assert.True(t, verified)

	proof, verified, err = tree.Proof(18, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 18, verifiedBy: 30, nodes: []uint64{18, 16, 21, 7, 25, 28}}, proof)
	assert.True(t, verified)

	proof, verified, err = tree.Proof(17, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 17, verifiedBy: 30, nodes: []uint64{17, 21, 7, 25, 28}}, proof)
	assert.True(t, verified)
}

func Test_ProofWithDigest1(t *testing.T) {
	t.Parallel()

	tree := NewDefaultTree()
	proof, verified, err := tree.Proof(0, 0, NewDefaultTree())
	assert.NoError(t, err)
	assert.False(t, verified)

	tree.Set(0)
	tree.Set(2)

	proof, verified, err = tree.Proof(0, 0b1, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 0, nodes: []uint64{0}}, proof)
	assert.True(t, verified)

	tree.Set(5)

	proof, verified, err = tree.Proof(0, 0b10, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 8, nodes: []uint64{0, 5}}, proof)
	assert.True(t, verified)

	proof, verified, err = tree.Proof(0, 0b110, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 8, nodes: []uint64{0}}, proof)
	assert.True(t, verified)

	tree.Set(8)

	proof, verified, err = tree.Proof(0, 0b101, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 0, nodes: []uint64{0, 2}}, proof)
	assert.True(t, verified)

	proof, verified, err = tree.Proof(0, 0b10, NewDefaultTree())
	assert.NoError(t, err)
	assert.Equal(t, Proof{index: 0, verifiedBy: 10, nodes: []uint64{0, 5, 8}}, proof)
	assert.True(t, verified)
}
