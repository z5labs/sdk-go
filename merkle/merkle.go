// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"encoding"
	"encoding/hex"
	"errors"
	"hash"
)

// BinaryTree
type BinaryTree struct {
	hash  []byte
	left  *BinaryTree
	right *BinaryTree
}

// ConstructBinaryTree
func ConstructBinaryTree[T encoding.BinaryMarshaler](hasher hash.Hash, leafs ...T) (*BinaryTree, error) {
	nodes := make([]*BinaryTree, len(leafs))
	for i, leaf := range leafs {
		b, err := leaf.MarshalBinary()
		if err != nil {
			return nil, err
		}

		n, err := hasher.Write(b)
		if err != nil {
			return nil, err
		}
		if n != len(b) {
			return nil, errors.New("failed to write all bytes of leaf to hasher")
		}

		nodes[i] = &BinaryTree{
			hash: hasher.Sum(nil),
		}
	}

	return constructBinaryTree(hasher, nodes)
}

func constructBinaryTree(hasher hash.Hash, nodes []*BinaryTree) (*BinaryTree, error) {
	newNodes := make([]*BinaryTree, 0, len(nodes)/2)
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		right := nodes[i+1]

		concatenatedHashes := append(left.hash, right.hash...)
		n, err := hasher.Write(concatenatedHashes)
		if err != nil {
			return nil, err
		}
		if n != len(concatenatedHashes) {
			return nil, errors.New("failed to write all bytes of concatenated hash to hasher")
		}

		newNodes = append(newNodes, &BinaryTree{
			hash:  hasher.Sum(nil),
			left:  left,
			right: right,
		})
	}
	if len(newNodes) == 1 {
		return newNodes[0], nil
	}

	return constructBinaryTree(hasher, newNodes)
}

// Hash
func (t *BinaryTree) Hash() []byte {
	return t.hash
}

func (t *BinaryTree) String() string {
	return hex.EncodeToString(t.Hash())
}

// Left
func (t *BinaryTree) Left() *BinaryTree {
	return t.left
}

// Right
func (t *BinaryTree) Right() *BinaryTree {
	return t.right
}

// IsLeaf
func (t *BinaryTree) IsLeaf() bool {
	return t.left == nil && t.right == nil
}
