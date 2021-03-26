package flattree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Index(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		depth, offset, expected uint64
	}{
		{0, 0, 0},
		{0, 1, 2},
		{0, 2, 4},
		{0, 3, 6},
		{0, 4, 8},
		{1, 0, 1},
		{1, 1, 5},
		{1, 2, 9},
		{1, 3, 13},
		{1, 4, 17},
		{2, 0, 3},
		{2, 1, 11},
		{2, 2, 19},
		{2, 3, 27},
		{2, 4, 35},
		{3, 0, 7},
		{3, 1, 23},
		{3, 2, 39},
		{3, 3, 55},
		{3, 4, 71},
	}

	for _, tc := range testCases {
		i := Index(tc.depth, tc.offset)
		assert.Equal(t, tc.expected, i, "Index of %d with offset %d expected value of %d; got %d", tc.depth, tc.offset, tc.expected, i)
	}
}

func Test_Depth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		depth, expected uint64
	}{
		{0, 0},
		{1, 1},
		{2, 0},
		{3, 2},
		{4, 0},
		{5, 1},
		{6, 0},
		{7, 3},
		{8, 0},
	}

	for _, tc := range testCases {
		d := Depth(tc.depth)
		assert.Equal(t, tc.expected, d, "Depth of %d with expected value of %d; got %d", tc.depth, tc.expected, d)
	}
}

func Test_Offset(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		offset, expected uint64
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 0},
		{4, 2},
		{5, 1},
		{6, 3},
		{7, 0},
		{8, 4},
		{9, 2},
		{10, 5},
		{11, 1},
		{12, 6},
	}

	for _, tc := range testCases {
		o := Offset(tc.offset)
		assert.Equal(t, tc.expected, o, "Offset of %d with expected value of %d; got %d", tc.offset, tc.expected, o)
	}
}

func Test_Parent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected uint64
	}{
		{0, 1},
		{2, 1},
		{4, 5},
		{6, 5},
		{8, 9},
		{10, 9},
		{12, 13},
		{14, 13},
		{1, 3},
		{5, 3},
		{9, 11},
		{13, 11},
		{3, 7},
		{11, 7},
	}

	for _, tc := range testCases {
		o := Parent(tc.node)
		assert.Equal(t, tc.expected, o, "Parent of %d with expected value of %d; got %d", tc.node, tc.expected, o)
	}
}

func Test_Sibling(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected uint64
	}{
		{0, 2},
		{2, 0},
		{4, 6},
		{6, 4},
		{8, 10},
		{10, 8},
		{12, 14},
		{14, 12},
		{1, 5},
		{5, 1},
		{9, 13},
		{13, 9},
		{3, 11},
		{11, 3},
	}

	for _, tc := range testCases {
		o := Sibling(tc.node)
		assert.Equal(t, tc.expected, o, "Sibling of %d with expected value of %d; got %d", tc.node, tc.expected, o)
	}
}

func Test_Uncle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected uint64
	}{
		{0, 5},
		{2, 5},
		{4, 1},
		{6, 1},
		{8, 13},
		{10, 13},
		{12, 9},
		{14, 9},
		{1, 11},
		{5, 11},
		{9, 3},
		{13, 3},
	}

	for _, tc := range testCases {
		o := Uncle(tc.node)
		assert.Equal(t, tc.expected, o, "Uncle of %d with expected value of %d; got %d", tc.node, tc.expected, o)
	}
}

func Test_Children(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected1, expected2 uint64
		expectedExists             bool
	}{
		{0, 0, 0, false},
		{2, 0, 0, false},
		{4, 0, 0, false},
		{6, 0, 0, false},
		{8, 0, 0, false},
		{10, 0, 0, false},
		{12, 0, 0, false},
		{14, 0, 0, false},
		{1, 0, 2, true},
		{5, 4, 6, true},
		{9, 8, 10, true},
		{13, 12, 14, true},
	}

	for _, tc := range testCases {
		c1, c2, exists := Children(tc.node)
		assert.Equal(t, tc.expected1, c1, "Children of %d with expected left child of %d; got %d", tc.node, tc.expected1, c1)
		assert.Equal(t, tc.expected2, c2, "Children of %d with expected right child of %d; got %d", tc.node, tc.expected2, c2)
		assert.Equal(t, tc.expectedExists, exists, "Children of %d with expected exists of %d; got %d", tc.node, tc.expectedExists, exists)
	}
}

func Test_LeftChild(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected uint64
		expectedExists bool
	}{
		{0, 0, false},
		{2, 0, false},
		{4, 0, false},
		{6, 0, false},
		{8, 0, false},
		{10, 0, false},
		{12, 0, false},
		{14, 0, false},
		{1, 0, true},
		{5, 4, true},
		{9, 8, true},
		{13, 12, true},
	}

	for _, tc := range testCases {
		child, exists := LeftChild(tc.node)
		assert.Equal(t, tc.expected, child, "Left child of %d with expected value of %d; got %d", tc.node, tc.expected, child)
		assert.Equal(t, tc.expectedExists, exists, "Left child of %d with expected exists of %d; got %d", tc.node, tc.expectedExists, exists)
	}
}

func Test_RightChild(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected uint64
		expectedExists bool
	}{
		{0, 0, false},
		{2, 0, false},
		{4, 0, false},
		{6, 0, false},
		{8, 0, false},
		{10, 0, false},
		{12, 0, false},
		{14, 0, false},
		{1, 2, true},
		{5, 6, true},
		{9, 10, true},
		{13, 14, true},
	}

	for _, tc := range testCases {
		child, exists := RightChild(tc.node)
		assert.Equal(t, tc.expected, child, "Right child of %d with expected value of %d; got %d", tc.node, tc.expected, child)
		assert.Equal(t, tc.expectedExists, exists, "Right child of %d with expected exists of %d; got %d", tc.node, tc.expectedExists, exists)
	}
}

func Test_Spans(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expectedLeft, expectedRight uint64
	}{
		{0, 0, 0},
		{1, 0, 2},
		{2, 2, 2},
		{3, 0, 6},
		{4, 4, 4},
		{5, 4, 6},
		{6, 6, 6},
		{7, 0, 14},
		{8, 8, 8},
		{23, 16, 30},
		{27, 24, 30},
	}

	for _, tc := range testCases {
		left, right := Spans(tc.node)
		assert.Equal(t, tc.expectedLeft, left, "Spans of node %d with expected left value of %d; got %d", tc.node, tc.expectedLeft, left)
		assert.Equal(t, tc.expectedRight, right, "Spans of node %d with expected right value of %d; got %d", tc.node, tc.expectedRight, right)
	}
}

func Test_Count(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node, expected uint64
	}{
		{0, 1},
		{1, 3},
		{2, 1},
		{3, 7},
		{4, 1},
		{5, 3},
		{6, 1},
		{7, 15},
		{8, 1},
		{23, 15},
		{27, 7},
	}

	for _, tc := range testCases {
		n := Count(tc.node)
		assert.Equal(t, tc.expected, n, "Count of node %d with expected value of %d; got %d", tc.node, tc.expected, n)
	}
}

func Test_FullRoots(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node        uint64
		expected    []uint64
		expectedErr error
	}{
		{0, []uint64{}, nil},
		{2, []uint64{0}, nil},
		{8, []uint64{3}, nil},
		{16, []uint64{7}, nil},
		{18, []uint64{7, 16}, nil},
		{20, []uint64{7, 17}, nil},
		{32, []uint64{15}, nil},
		{36, []uint64{15, 33}, nil},
		{3, []uint64{}, errors.New("odd index passed to FullRoots()")},
	}

	for _, tc := range testCases {
		roots, err := FullRoots(tc.node)
		assert.Equal(t, tc.expected, roots, "FullRoots of node %d with expected values of %d; got %d", tc.node, tc.expected, roots)
		if tc.expectedErr == nil {
			assert.NoError(t, err, "FullRoots of node %d with unexpected error of %d", tc.node, err)
		} else {
			assert.EqualError(t, err, tc.expectedErr.Error(), "FullRoots of node %d with expected error of %d, but no error received", tc.node, err)
		}
	}
}
