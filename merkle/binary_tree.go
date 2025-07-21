// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"bytes"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// BinaryTree is a merkle tree where each node in the tree has exactly 2 child nodes.
type BinaryTree struct {
	hash  []byte
	left  *BinaryTree
	right *BinaryTree
}

// ErrAtLeastOneLeafRequired is returned if a [BinaryTree] is zero leaf nodes are
// provided when calling [ConstructBinaryTree].
var ErrAtLeastOneLeafRequired = errors.New("must provide at least one leaf")

// ConstructBinaryTree will construct a full merkle [BinaryTree] from the given leaf nodes.
func ConstructBinaryTree[T io.Reader](hasher hash.Hash, leafs ...T) (*BinaryTree, error) {
	if len(leafs) == 0 {
		return nil, ErrAtLeastOneLeafRequired
	}

	nodes := make([]*BinaryTree, len(leafs))
	for i, leaf := range leafs {
		hash, err := hashAll(hasher, leaf)
		if err != nil {
			return nil, err
		}

		nodes[i] = &BinaryTree{
			hash: hash,
		}
	}

	return constructBinaryTree(hasher, nodes)
}

func constructBinaryTree(hasher hash.Hash, nodes []*BinaryTree) (*BinaryTree, error) {
	numOfNewNodes := len(nodes) / 2

	hasOddNumOfNodes := len(nodes)%2 != 0
	if hasOddNumOfNodes {
		numOfNewNodes += 1
	}

	newNodes := make([]*BinaryTree, numOfNewNodes)
	if hasOddNumOfNodes {
		newNodes[numOfNewNodes-1] = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
	}

	var buf bytes.Buffer
	for i := 0; i < len(nodes); i += 2 {
		buf.Reset()

		left := nodes[i]
		buf.Write(left.hash)

		right := nodes[i+1]
		buf.Write(right.hash)

		hash, err := hashAll(hasher, &buf)
		if err != nil {
			return nil, err
		}

		newNodes[i/2] = &BinaryTree{
			hash:  hash,
			left:  left,
			right: right,
		}
	}
	if len(newNodes) == 1 {
		return newNodes[0], nil
	}

	return constructBinaryTree(hasher, newNodes)
}

func hashAll(hasher hash.Hash, r io.Reader) ([]byte, error) {
	// ensure state independent hashes aka each node hash is reproducible
	// and independent of the hashing operations that came before it
	hasher.Reset()

	_, err := io.Copy(hasher, r)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

// Hash returns the raw hash value for this tree.
func (t *BinaryTree) Hash() []byte {
	return t.hash
}

// String returns a hex encoded representation of the hash.
func (t *BinaryTree) String() string {
	return hex.EncodeToString(t.Hash())
}

// Left returns the left child tree. Note, this will be nil if [IsLeaf] returns true.
func (t *BinaryTree) Left() *BinaryTree {
	return t.left
}

// Right returns the right child tree. Note, this will be nil if [IsLeaf] returns true.
func (t *BinaryTree) Right() *BinaryTree {
	return t.right
}

// IsLeaf reports whether this tree represents a leaf value.
func (t *BinaryTree) IsLeaf() bool {
	return t.left == nil && t.right == nil
}
