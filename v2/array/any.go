// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package array

import (
	G "github.com/IBM/fp-go/v2/array/generic"
)

// AnyWithIndex tests if any of the elements in the array matches the predicate.
// The predicate receives both the index and the element.
// Returns true if at least one element satisfies the predicate, false otherwise.
//
// Example:
//
//	hasEvenAtEvenIndex := array.AnyWithIndex(func(i, x int) bool {
//	    return i%2 == 0 && x%2 == 0
//	})
//	result := hasEvenAtEvenIndex([]int{1, 3, 4, 5}) // true (4 is at index 2)
//
//go:inline
func AnyWithIndex[A any](pred func(int, A) bool) func([]A) bool {
	return G.AnyWithIndex[[]A](pred)
}

// Any tests if any of the elements in the array matches the predicate.
// Returns true if at least one element satisfies the predicate, false otherwise.
// Returns false for an empty array.
//
// Example:
//
//	hasEven := array.Any(func(x int) bool { return x%2 == 0 })
//	result := hasEven([]int{1, 3, 4, 5}) // true
//
//go:inline
func Any[A any](pred func(A) bool) func([]A) bool {
	return G.Any[[]A](pred)
}
