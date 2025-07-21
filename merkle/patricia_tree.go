// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"hash"
	"io"
)

// AdaptiveRadixTree
type AdaptiveRadixTree struct {
	hasher hash.Hash
	key    []byte
	hash   []byte
}

// NewAdaptiveRadixTree
func NewAdaptiveRadixTree(hasher hash.Hash) *AdaptiveRadixTree {
	return &AdaptiveRadixTree{
		hasher: hasher,
	}
}

// Insert
func (t *AdaptiveRadixTree) Insert(key []byte, value io.Reader) *AdaptiveRadixTree {
	return nil
}
