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
	"github.com/IBM/fp-go/v2/option"
)

// FindFirst finds the first element which satisfies a predicate function.
// Returns Some(element) if found, None if no element matches.
//
// Example:
//
//	findGreaterThan3 := array.FindFirst(func(x int) bool { return x > 3 })
//	result := findGreaterThan3([]int{1, 2, 4, 5}) // Some(4)
//	result2 := findGreaterThan3([]int{1, 2, 3}) // None
//
//go:inline
func FindFirst[A any](pred func(A) bool) option.Kleisli[[]A, A] {
	return G.FindFirst[[]A](pred)
}

// FindFirstWithIndex finds the first element which satisfies a predicate function that also receives the index.
// Returns Some(element) if found, None if no element matches.
//
// Example:
//
//	findEvenAtEvenIndex := array.FindFirstWithIndex(func(i, x int) bool {
//	    return i%2 == 0 && x%2 == 0
//	})
//	result := findEvenAtEvenIndex([]int{1, 3, 4, 5}) // Some(4)
//
//go:inline
func FindFirstWithIndex[A any](pred func(int, A) bool) option.Kleisli[[]A, A] {
	return G.FindFirstWithIndex[[]A](pred)
}

// FindFirstMap finds the first element for which the selector function returns Some.
// This combines finding and mapping in a single operation.
//
// Example:
//
//	import "strconv"
//
//	parseFirst := array.FindFirstMap(func(s string) option.Option[int] {
//	    if n, err := strconv.Atoi(s); err == nil {
//	        return option.Some(n)
//	    }
//	    return option.None[int]()
//	})
//	result := parseFirst([]string{"a", "42", "b"}) // Some(42)
//
//go:inline
func FindFirstMap[A, B any](sel option.Kleisli[A, B]) option.Kleisli[[]A, B] {
	return G.FindFirstMap[[]A](sel)
}

// FindFirstMapWithIndex finds the first element for which the selector function returns Some.
// The selector receives both the index and the element.
//
//go:inline
func FindFirstMapWithIndex[A, B any](sel func(int, A) Option[B]) option.Kleisli[[]A, B] {
	return G.FindFirstMapWithIndex[[]A](sel)
}

// FindLast finds the last element which satisfies a predicate function.
// Returns Some(element) if found, None if no element matches.
//
// Example:
//
//	findGreaterThan3 := array.FindLast(func(x int) bool { return x > 3 })
//	result := findGreaterThan3([]int{1, 4, 2, 5}) // Some(5)
//
//go:inline
func FindLast[A any](pred func(A) bool) option.Kleisli[[]A, A] {
	return G.FindLast[[]A](pred)
}

// FindLastWithIndex finds the last element which satisfies a predicate function that also receives the index.
// Returns Some(element) if found, None if no element matches.
//
//go:inline
func FindLastWithIndex[A any](pred func(int, A) bool) option.Kleisli[[]A, A] {
	return G.FindLastWithIndex[[]A](pred)
}

// FindLastMap finds the last element for which the selector function returns Some.
// This combines finding and mapping in a single operation, searching from the end.
//
//go:inline
func FindLastMap[A, B any](sel option.Kleisli[A, B]) option.Kleisli[[]A, B] {
	return G.FindLastMap[[]A](sel)
}

// FindLastMapWithIndex finds the last element for which the selector function returns Some.
// The selector receives both the index and the element, searching from the end.
//
//go:inline
func FindLastMapWithIndex[A, B any](sel func(int, A) Option[B]) option.Kleisli[[]A, B] {
	return G.FindLastMapWithIndex[[]A](sel)
}
