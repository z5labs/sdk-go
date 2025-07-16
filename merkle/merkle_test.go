// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"hash"
	"testing"

	"github.com/stretchr/testify/require"
)

type stringBinaryMarshaler string

func (s stringBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

func TestConstructBinaryTree(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name   string
		Hasher hash.Hash
		Leafs  []encoding.BinaryMarshaler
		Assert func(t *testing.T, tree *BinaryTree, err error)
	}{
		{
			Name:   "should be idempotent",
			Hasher: sha256.New(),
			Leafs: []encoding.BinaryMarshaler{
				stringBinaryMarshaler("a"),
				stringBinaryMarshaler("b"),
			},
			Assert: func(t *testing.T, tree *BinaryTree, err error) {
				require.Nil(t, err)

				s := hex.EncodeToString(tree.Hash())

				require.Equal(t, "2b1acd11da9daf1bbc3580547130566fce9e509bb26fb4e3ebf6ac4f2d3c02a1", s)

				left := tree.Left()
				require.NotNil(t, left)
				require.True(t, left.IsLeaf())

				right := tree.Right()
				require.NotNil(t, right)
				require.True(t, right.IsLeaf())
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			tree, err := ConstructBinaryTree(testCase.Hasher, testCase.Leafs...)

			testCase.Assert(t, tree, err)
		})
	}
}
