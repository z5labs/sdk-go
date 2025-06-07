// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package ptr provides helpers for working references.
package ptr

// Deref returns either the dereferenced value of the given value of
// the given reference or the zero value of type T, if the reference is nil.
func Deref[T any](t *T) T {
	var zero T
	if t == nil {
		return zero
	}
	return *t
}

// Ref returns a reference to the given value.
func Ref[T any](t T) *T {
	return &t
}
