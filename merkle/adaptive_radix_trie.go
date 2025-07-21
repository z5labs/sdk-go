// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"bytes"
	"hash"
	"io"
)

// AdaptiveRadixTrie
type AdaptiveRadixTrie struct {
	hasher   hash.Hash
	hash     []byte
	key      []byte
	children []*AdaptiveRadixTrie
}

// NewAdaptiveRadixTrie
func NewAdaptiveRadixTrie(hasher hash.Hash) *AdaptiveRadixTrie {
	return &AdaptiveRadixTrie{
		hasher: hasher,
	}
}

// Insert
func (t *AdaptiveRadixTrie) Insert(key []byte, value io.Reader) error {
	return nil
}

func (t *AdaptiveRadixTrie) insert(key []byte, value io.Reader) error {
	if len(key) == len(t.key) {
		hash, err := hashAll(t.hasher, value)
		if err != nil {
			return err
		}

		if bytes.Equal(t.hash, hash) {
			return nil
		}

		t.hash = hash
		// TODO: recompute hashes
		return nil
	}

	for _, child := range t.children {
		if bytes.HasPrefix(key, child.key) {
			return t.Insert(key, value)
		}

		prefix, ok := commonPrefix(key, child.key)
		if !ok {
			continue
		}

		// TODO: construct new node based in common prefix
		// TODO: recompute hashes
		return nil
	}

	// TODO: construct new node based on full key
	// TODO: recompute hashes
	return nil
}

func commonPrefix(a, b []byte) ([]byte, bool) {
	for i := range a {
		if a[i] == b[i] {
			continue
		}

		if i == 0 {
			return nil, false
		}

		return a[:i], true
	}

	return a, true
}

// Hash
func (t *AdaptiveRadixTrie) Hash() []byte {
	return t.hash
}

// Key
func (t *AdaptiveRadixTrie) Key() []byte {
	return t.key
}

// Children
func (t *AdaptiveRadixTrie) Children() []*AdaptiveRadixTrie {
	return t.children
}
