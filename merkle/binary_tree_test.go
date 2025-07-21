// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"crypto/sha256"
	"encoding"
	"errors"
	"hash"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMarshalBinaryFailed = errors.New("failed to marshal binary")

type binaryMarshalerFunc func() ([]byte, error)

func (f binaryMarshalerFunc) MarshalBinary() ([]byte, error) {
	return f()
}

type stringBinaryMarshaler string

func (s stringBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

func TestConstructBinaryTree(t *testing.T) {
	t.Parallel()

	globalHasher := sha256.New()

	testCases := []struct {
		Name   string
		Hasher func() hash.Hash
		Leafs  [][]encoding.BinaryMarshaler
		Assert func(t *testing.T, err error, trees []*BinaryTree)
	}{
		{
			Name: "should fail if leaf fails to MarshalBinary",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]encoding.BinaryMarshaler{
				{
					binaryMarshalerFunc(func() ([]byte, error) {
						return nil, errMarshalBinaryFailed
					}),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.ErrorIs(t, err, errMarshalBinaryFailed)
				require.Empty(t, trees)
			},
		},
		{
			Name: "should be work with non power of 2 number of leafs",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]encoding.BinaryMarshaler{
				{
					stringBinaryMarshaler("a"),
					stringBinaryMarshaler("b"),
					stringBinaryMarshaler("c"),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.Nil(t, err)
				require.Len(t, trees, 1)

				tree := trees[0]
				require.Equal(t, "6632753d6ca30fea890f37fc150eaed8d068acf596acb2251b8fafd72db977d3", tree.String())
			},
		},
		{
			Name: "should be idempotent",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]encoding.BinaryMarshaler{
				{
					stringBinaryMarshaler("a"),
					stringBinaryMarshaler("b"),
				},
				{
					stringBinaryMarshaler("a"),
					stringBinaryMarshaler("b"),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.Nil(t, err)
				require.Len(t, trees, 2)

				treeA := trees[0]
				treeB := trees[1]
				require.Equal(t, treeA.String(), treeB.String())
			},
		},
		{
			Name: "should be order dependent",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]encoding.BinaryMarshaler{
				{
					stringBinaryMarshaler("a"),
					stringBinaryMarshaler("b"),
				},
				{
					stringBinaryMarshaler("b"),
					stringBinaryMarshaler("a"),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.Nil(t, err)
				require.Len(t, trees, 2)

				treeA := trees[0]
				treeB := trees[1]

				require.NotEqual(t, treeA.String(), treeB.String())
			},
		},
		{
			Name: "should succeed if hasher is reused",
			Hasher: func() hash.Hash {
				return globalHasher
			},
			Leafs: [][]encoding.BinaryMarshaler{
				{
					stringBinaryMarshaler("a"),
					stringBinaryMarshaler("b"),
				},
				{
					stringBinaryMarshaler("a"),
					stringBinaryMarshaler("b"),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.Nil(t, err)
				require.Len(t, trees, 2)

				treeA := trees[0]
				treeB := trees[1]

				require.Equal(t, treeA.String(), treeB.String())
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			constructTrees := func() ([]*BinaryTree, error) {
				trees := make([]*BinaryTree, len(testCase.Leafs))

				for i, leafs := range testCase.Leafs {
					tree, err := ConstructBinaryTree(testCase.Hasher(), leafs...)
					if err != nil {
						return nil, err
					}

					trees[i] = tree
				}
				return trees, nil
			}

			trees, err := constructTrees()

			testCase.Assert(t, err, trees)
		})
	}
}
