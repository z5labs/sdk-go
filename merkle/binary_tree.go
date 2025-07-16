// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"encoding"
	"encoding/hex"
	"fmt"
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

		err = writeToHasher(hasher, b)
		if err != nil {
			return nil, err
		}

		nodes[i] = &BinaryTree{
			hash: hasher.Sum(nil),
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

	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		err := writeToHasher(hasher, left.hash)
		if err != nil {
			return nil, err
		}

		right := nodes[i+1]
		err = writeToHasher(hasher, right.hash)
		if err != nil {
			return nil, err
		}

		newNodes[i/2] = &BinaryTree{
			hash:  hasher.Sum(nil),
			left:  left,
			right: right,
		}
	}
	if len(newNodes) == 1 {
		return newNodes[0], nil
	}

	return constructBinaryTree(hasher, newNodes)
}

// IncompleteWriteError is returned if not all bytes are successfully written
// to the provided [hash.Hash].
type IncompleteWriteError struct {
	Expected      int
	ActualWritten int
}

// Error implements the [error] interface.
func (e IncompleteWriteError) Error() string {
	return fmt.Sprintf("expected write %d bytes but only wrote: %d", e.Expected, e.ActualWritten)
}

func writeToHasher(hasher hash.Hash, b []byte) error {
	n, err := hasher.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return IncompleteWriteError{
			Expected:      len(b),
			ActualWritten: n,
		}
	}

	return nil
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
