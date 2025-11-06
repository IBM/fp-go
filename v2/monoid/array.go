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

package monoid

import (
	S "github.com/IBM/fp-go/v2/semigroup"
)

// GenericConcatAll combines all elements in a generic slice using a monoid.
//
// This function works with custom slice types (types that are defined as ~[]A).
// It uses the monoid's Concat operation to combine elements and returns the
// monoid's Empty value for empty slices.
//
// Type Parameters:
//   - GA: A generic slice type constraint (~[]A)
//   - A: The element type
//
// Parameters:
//   - m: The monoid to use for combining elements
//
// Returns:
//   - A function that takes a slice and returns the combined result
//
// Example:
//
//	type IntSlice []int
//
//	addMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,
//	)
//
//	concatAll := GenericConcatAll[IntSlice](addMonoid)
//	result := concatAll(IntSlice{1, 2, 3, 4, 5})  // 15
//	empty := concatAll(IntSlice{})                 // 0
func GenericConcatAll[GA ~[]A, A any](m Monoid[A]) func(GA) A {
	return S.GenericConcatAll[GA](S.MakeSemigroup(m.Concat))(m.Empty())
}

// ConcatAll combines all elements in a slice using a monoid.
//
// This function reduces a slice to a single value by repeatedly applying the
// monoid's Concat operation. For an empty slice, it returns the monoid's Empty value.
//
// This is the standard version that works with []A slices. For custom slice types,
// use GenericConcatAll.
//
// Parameters:
//   - m: The monoid to use for combining elements
//
// Returns:
//   - A function that takes a slice and returns the combined result
//
// Example:
//
//	addMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,
//	)
//
//	concatAll := ConcatAll(addMonoid)
//	sum := concatAll([]int{1, 2, 3, 4, 5})  // 15
//	empty := concatAll([]int{})              // 0
//
//	stringMonoid := MakeMonoid(
//	    func(a, b string) string { return a + b },
//	    "",
//	)
//	concat := ConcatAll(stringMonoid)
//	result := concat([]string{"Hello", " ", "World"})  // "Hello World"
func ConcatAll[A any](m Monoid[A]) func([]A) A {
	return GenericConcatAll[[]A](m)
}

// Fold combines all elements in a slice using a monoid.
//
// This is an alias for ConcatAll, providing a more functional programming-style name.
// It performs a left fold (reduce) operation using the monoid's Concat function,
// starting with the monoid's Empty value.
//
// Parameters:
//   - m: The monoid to use for combining elements
//
// Returns:
//   - A function that takes a slice and returns the combined result
//
// Example:
//
//	mulMonoid := MakeMonoid(
//	    func(a, b int) int { return a * b },
//	    1,
//	)
//
//	fold := Fold(mulMonoid)
//	product := fold([]int{2, 3, 4})  // 24 (1 * 2 * 3 * 4)
//	empty := fold([]int{})            // 1
func Fold[A any](m Monoid[A]) func([]A) A {
	return GenericConcatAll[[]A](m)
}
