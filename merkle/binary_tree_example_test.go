// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"crypto/sha256"
	"encoding"
	"fmt"
)

func ExampleBinaryTree() {
	leafs := []encoding.BinaryMarshaler{
		stringBinaryMarshaler("a"),
		stringBinaryMarshaler("b"),
	}

	treeA, err := ConstructBinaryTree(sha256.New(), leafs...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(treeA)

	treeB, err := ConstructBinaryTree(sha256.New(), leafs...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(treeB)

	// Output: 2b1acd11da9daf1bbc3580547130566fce9e509bb26fb4e3ebf6ac4f2d3c02a1
	// 2b1acd11da9daf1bbc3580547130566fce9e509bb26fb4e3ebf6ac4f2d3c02a1
}
