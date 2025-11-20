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

package semigroup

import (
	M "github.com/IBM/fp-go/v2/magma"
)

// GenericMonadConcatAll creates a function that concatenates all elements in a generic slice
// with an initial value using the provided semigroup operation. This is the uncurried version
// that takes both the slice and initial value as parameters.
//
// The function processes elements left-to-right, combining them with the initial value.
//
// Example:
//
//	type MyInts []int
//	import N "github.com/IBM/fp-go/v2/number"
//	sum := N.SemigroupSum[int]()
//	concatAll := semigroup.GenericMonadConcatAll[MyInts](sum)
//	result := concatAll(MyInts{1, 2, 3}, 10)  // 10 + 1 + 2 + 3 = 16
func GenericMonadConcatAll[GA ~[]A, A any](s Semigroup[A]) func(GA, A) A {
	return M.GenericMonadConcatAll[GA](M.MakeMagma(s.Concat))
}

// GenericConcatAll creates a curried function that concatenates all elements in a generic slice
// with an initial value using the provided semigroup operation.
//
// The returned function first takes the initial value, then takes the slice, and finally
// returns the result of combining all elements left-to-right with the initial value.
//
// Example:
//
//	type MyInts []int
//	import N "github.com/IBM/fp-go/v2/number"
//	sum := N.SemigroupSum[int]()
//	concatAll := semigroup.GenericConcatAll[MyInts](sum)
//	result := concatAll(10)(MyInts{1, 2, 3})  // 10 + 1 + 2 + 3 = 16
func GenericConcatAll[GA ~[]A, A any](s Semigroup[A]) func(A) func(GA) A {
	return M.GenericConcatAll[GA](M.MakeMagma(s.Concat))
}

// MonadConcatAll creates a function that concatenates all elements in a slice with an initial
// value using the provided semigroup operation. This is the uncurried version that takes both
// the slice and initial value as parameters.
//
// This is a convenience wrapper around GenericMonadConcatAll for standard slices.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//	sum := N.SemigroupSum[int]()
//	concatAll := semigroup.MonadConcatAll(sum)
//	result := concatAll([]int{1, 2, 3}, 10)  // 10 + 1 + 2 + 3 = 16
func MonadConcatAll[A any](s Semigroup[A]) func([]A, A) A {
	return GenericMonadConcatAll[[]A](s)
}

// ConcatAll creates a curried function that concatenates all elements in a slice with an
// initial value using the provided semigroup operation.
//
// The returned function first takes the initial value, then takes the slice, and finally
// returns the result of combining all elements left-to-right with the initial value.
//
// This is a convenience wrapper around GenericConcatAll for standard slices.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//	sum := N.SemigroupSum[int]()
//	concatAll := semigroup.ConcatAll(sum)
//	result := concatAll(10)([]int{1, 2, 3})  // 10 + 1 + 2 + 3 = 16
func ConcatAll[A any](s Semigroup[A]) func(A) func([]A) A {
	return GenericConcatAll[[]A](s)
}
