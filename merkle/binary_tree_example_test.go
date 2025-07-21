// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package merkle

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func ExampleBinaryTree() {
	treeA, err := ConstructBinaryTree(sha256.New(), strings.NewReader("a"), strings.NewReader("b"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(treeA)

	treeB, err := ConstructBinaryTree(sha256.New(), strings.NewReader("a"), strings.NewReader("b"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(treeB)

	// Output: e5a01fee14e0ed5c48714f22180f25ad8365b53f9779f79dc4a3d7e93963f94a
	// e5a01fee14e0ed5c48714f22180f25ad8365b53f9779f79dc4a3d7e93963f94a
}
