// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var errReadFailed = errors.New("failed to read")

type readFunc func([]byte) (int, error)

func (f readFunc) Read(b []byte) (int, error) {
	return f(b)
}

func TestConstructBinaryTree(t *testing.T) {
	t.Parallel()

	globalHasher := sha256.New()

	testCases := []struct {
		Name   string
		Hasher func() hash.Hash
		Leafs  [][]io.Reader
		Assert func(t *testing.T, err error, trees []*BinaryTree)
	}{
		{
			Name: "should fail if leaf fails to MarshalBinary",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]io.Reader{
				{
					readFunc(func(b []byte) (int, error) {
						return 0, errReadFailed
					}),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.ErrorIs(t, err, errReadFailed)
				require.Empty(t, trees)
			},
		},
		{
			Name: "should be work with non power of 2 number of leafs",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]io.Reader{
				{
					strings.NewReader("a"),
					strings.NewReader("b"),
					strings.NewReader("c"),
				},
			},
			Assert: func(t *testing.T, err error, trees []*BinaryTree) {
				require.Nil(t, err)
				require.Len(t, trees, 1)

				tree := trees[0]
				require.Equal(t, "7075152d03a5cd92104887b476862778ec0c87be5c2fa1c0a90f87c49fad6eff", tree.String())
			},
		},
		{
			Name: "should be idempotent",
			Hasher: func() hash.Hash {
				return sha256.New()
			},
			Leafs: [][]io.Reader{
				{
					strings.NewReader("a"),
					strings.NewReader("b"),
				},
				{
					strings.NewReader("a"),
					strings.NewReader("b"),
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
			Leafs: [][]io.Reader{
				{
					strings.NewReader("a"),
					strings.NewReader("b"),
				},
				{
					strings.NewReader("b"),
					strings.NewReader("a"),
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
			Leafs: [][]io.Reader{
				{
					strings.NewReader("a"),
					strings.NewReader("b"),
				},
				{
					strings.NewReader("a"),
					strings.NewReader("b"),
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
