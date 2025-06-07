// Copyright 2025 Z5Labs and Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
