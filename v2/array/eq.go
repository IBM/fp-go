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
	"slices"

	E "github.com/IBM/fp-go/v2/eq"
)

// Eq creates an equality checker for arrays given an equality checker for elements.
// Two arrays are considered equal if they have the same length and all corresponding
// elements are equal according to the provided Eq instance.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/eq"
//
//	intArrayEq := array.Eq(eq.FromStrictEquals[int]())
//	result := intArrayEq.Equals([]int{1, 2, 3}, []int{1, 2, 3}) // true
//	result2 := intArrayEq.Equals([]int{1, 2, 3}, []int{1, 2, 4}) // false
func Eq[T any](e E.Eq[T]) E.Eq[[]T] {
	eq := e.Equals
	return E.FromEquals(func(left, right []T) bool {
		return slices.EqualFunc(left, right, eq)
	})
}

//go:inline
func StrictEquals[T comparable]() E.Eq[[]T] {
	return E.FromEquals(slices.Equal[[]T])
}
