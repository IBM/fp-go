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
	"github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Monoid returns a Monoid instance for arrays.
// The Monoid combines arrays through concatenation, with an empty array as the identity element.
//
// Example:
//
//	m := array.Monoid[int]()
//	result := m.Concat([]int{1, 2}, []int{3, 4}) // [1, 2, 3, 4]
//	empty := m.Empty() // []
//
//go:inline
func Monoid[T any]() M.Monoid[[]T] {
	return G.Monoid[[]T]()
}

// Semigroup returns a Semigroup instance for arrays.
// The Semigroup combines arrays through concatenation.
//
// Example:
//
//	s := array.Semigroup[int]()
//	result := s.Concat([]int{1, 2}, []int{3, 4}) // [1, 2, 3, 4]
//
//go:inline
func Semigroup[T any]() S.Semigroup[[]T] {
	return G.Semigroup[[]T]()
}

func addLen[A any](count int, data []A) int {
	return count + len(data)
}

// ArrayConcatAll efficiently concatenates multiple arrays into a single array.
// This function pre-allocates the exact amount of memory needed and performs
// a single copy operation for each input array, making it more efficient than
// repeated concatenations.
//
// Example:
//
//	result := array.ArrayConcatAll(
//	    []int{1, 2},
//	    []int{3, 4},
//	    []int{5, 6},
//	) // [1, 2, 3, 4, 5, 6]
func ArrayConcatAll[A any](data ...[]A) []A {
	// get the full size
	count := array.Reduce(data, addLen[A], 0)
	buf := make([]A, count)
	// copy
	array.Reduce(data, func(idx int, seg []A) int {
		return idx + copy(buf[idx:], seg)
	}, 0)
	// returns the final array
	return buf
}
