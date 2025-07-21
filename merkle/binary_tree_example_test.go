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

	// Output: 39361160903c6695c6804b7157c7bd10013e9ba89b1f954243bc8e3990b08db9
	// 39361160903c6695c6804b7157c7bd10013e9ba89b1f954243bc8e3990b08db9
}
