// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha256"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

type bytesBinaryMarshaler []byte

func (b bytesBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return b, nil
}

func FuzzConstructBinaryTree(f *testing.F) {
	corpus := []struct {
		numOfLeafs uint
		seed1      uint64
		seed2      uint64
	}{
		{numOfLeafs: 1, seed1: 1, seed2: 1},
		{numOfLeafs: 1, seed1: 20, seed2: 30},
		{numOfLeafs: 2, seed1: 10, seed2: 37},
		{numOfLeafs: 7, seed1: 199, seed2: 400},
	}
	for _, data := range corpus {
		f.Add(data.numOfLeafs, data.seed1, data.seed2)
	}

	f.Fuzz(func(t *testing.T, numOfLeafs uint, seed1 uint64, seed2 uint64) {
		src := rand.NewPCG(seed1, seed2)
		r := rand.New(src)

		leafs := make([]bytesBinaryMarshaler, numOfLeafs)
		for i := range numOfLeafs {
			var buf bytes.Buffer
			_, err := io.CopyN(&buf, crand.Reader, r.Int64N(1024))

			require.Nil(t, err)

			leafs[i] = bytesBinaryMarshaler(buf.Bytes())
		}

		tree, err := ConstructBinaryTree(sha256.New(), leafs...)
		require.Nil(t, err)

		var foundLeaves uint
		walk(tree, func(bt *BinaryTree) {
			if bt.IsLeaf() {
				foundLeaves += 1
			}
		})

		require.Equal(t, numOfLeafs, foundLeaves)
	})
}

func walk(tree *BinaryTree, f func(*BinaryTree)) {
	if tree.IsLeaf() {
		f(tree)
		return
	}

	f(tree.Left())
	f(tree.Right())
}
