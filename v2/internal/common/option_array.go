// Copyright (c) 2025 IBM Corp.
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

package common

import (
	F "github.com/IBM/fp-go/v2/function"
)

// OptionTraverseArrayG transforms an array by applying a function that returns an Option to each element.
// Returns Some containing the array of results if all operations succeed, None if any fails.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	parse := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil { return OptionNone[int]() }
//	    return OptionSome(n)
//	}
//	result := OptionTraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"}) // OptionSome([1, 2, 3])
//	result := OptionTraverseArrayG[[]string, []int](parse)([]string{"1", "x", "3"}) // None
func OptionTraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f OptionKleisli[A, B]) OptionKleisli[GA, GB] {
	return func(g GA) Option[GB] {
		bs := make(GB, len(g))
		for i, a := range g {
			b := f(a)
			if !b.s {
				return OptionNone[GB]()
			}
			bs[i] = b.a
		}
		return OptionSome(bs)
	}
}

// OptionTraverseArray transforms an array by applying a function that returns an Option to each element.
// Returns Some containing the array of results if all operations succeed, None if any fails.
//
// Example:
//
//	validate := F.Flow2(Predicate(N.MoreThan(0)), Map(N.Mul(2)))
//	result := OptionTraverseArray(validate)([]int{1, 2, 3}) // OptionSome([2, 4, 6])
//	result := OptionTraverseArray(validate)([]int{1, -1, 3}) // None
//
//go:inline
func OptionTraverseArray[A, B any](f OptionKleisli[A, B]) OptionKleisli[[]A, []B] {
	return OptionTraverseArrayG[[]A, []B](f)
}

// OptionTraverseArrayWithIndexG transforms an array by applying an indexed function that returns an Option.
// The function receives both the index and the element.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	f := func(i int, s string) Option[string] {
//	    return OptionSome(fmt.Sprintf("%d:%s", i, s))
//	}
//	result := OptionTraverseArrayWithIndexG[[]string, []string](f)([]string{"a", "b"}) // OptionSome(["0:a", "1:b"])
func OptionTraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) Option[B]) OptionKleisli[GA, GB] {
	return func(g GA) Option[GB] {
		bs := make(GB, len(g))
		for i, a := range g {
			b := f(i, a)
			if !b.s {
				return OptionNone[GB]()
			}
			bs[i] = b.a
		}
		return OptionSome(bs)
	}
}

// OptionTraverseArrayWithIndex transforms an array by applying an indexed function that returns an Option.
// The function receives both the index and the element.
//
// Example:
//
//	f := func(i int, x int) Option[int] {
//	    if x > i { return OptionSome(x) }
//	    return OptionNone[int]()
//	}
//	result := OptionTraverseArrayWithIndex(f)([]int{1, 2, 3}) // OptionSome([1, 2, 3])
//
//go:inline
func OptionTraverseArrayWithIndex[A, B any](f func(int, A) Option[B]) OptionKleisli[[]A, []B] {
	return OptionTraverseArrayWithIndexG[[]A, []B](f)
}

// OptionSequenceArrayG converts an array of Options into an Option of an array.
// Returns Some containing all values if all Options are Some, None if any is None.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	type MySlice []int
//	result := OptionSequenceArrayG[MySlice](A.From(OptionSome(1), OptionSome(2))) // OptionSome(MySlice{1, 2})
//	result := OptionSequenceArrayG[MySlice](A.From(OptionSome(1), OptionNone[int]())) // None
//
//go:inline
func OptionSequenceArrayG[GA ~[]A, GOA ~[]Option[A], A any](ma GOA) Option[GA] {
	return OptionTraverseArrayG[GOA, GA](F.Identity[Option[A]])(ma)
}

// OptionSequenceArray converts an array of Options into an Option of an array.
// Returns Some containing all values if all Options are Some, None if any is None.
//
// Example:
//
//	result := OptionSequenceArray(A.From(OptionSome(1), OptionSome(2), OptionSome(3))) // OptionSome([1, 2, 3])
//	result := OptionSequenceArray(A.From(OptionSome(1), OptionNone[int](), OptionSome(3))) // None
func OptionSequenceArray[A any](ma []Option[A]) Option[[]A] {
	return OptionSequenceArrayG[[]A](ma)
}

// OptionCompactArrayG filters an array of Options, keeping only the Some values and discarding None values.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	type MySlice []int
//	input := A.From(OptionSome(1), OptionNone[int](), OptionSome(3))
//	result := OptionCompactArrayG[[]Option[int], MySlice](input) // MySlice{1, 3}
func OptionCompactArrayG[A1 ~[]Option[A], A2 ~[]A, A any](fa A1) A2 {
	as := make(A2, 0, len(fa))
	for _, oa := range fa {
		if oa.s {
			as = append(as, oa.a)
		}
	}
	return as
}

// OptionCompactArray filters an array of Options, keeping only the Some values and discarding None values.
//
// Example:
//
//	input := A.From(OptionSome(1), OptionNone[int](), OptionSome(3), OptionSome(5), OptionNone[int]())
//	result := OptionCompactArray(input) // [1, 3, 5]
//
//go:inline
func OptionCompactArray[A any](fa []Option[A]) []A {
	return OptionCompactArrayG[[]Option[A], []A](fa)
}

// OptionTraversableArray returns a Traversable instance for arrays.
// A Traversable represents a data structure that can be traversed from left to right,
// applying an effectful function to each element and collecting the results.
//
// This function provides a way to obtain the traversal operation for arrays as a
// first-class value, which is useful when you need to pass the traversal operation
// as a parameter or compose it with other operations.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Returns:
//   - Traversable[A, B, []A, []B]: A function that takes a Kleisli arrow and returns
//     another Kleisli arrow that operates on arrays
//
// The returned Traversable has the signature:
//
//	func(OptionKleisli[A, B]) OptionKleisli[[]A, []B]
//
// which is equivalent to:
//
//	func(func(A) Option[B]) func([]A) Option[[]B]
//
// Example:
//
//	// Get the traversable instance
//	traversable := OptionTraversableArray[string, int]()
//
//	// Use it with a parsing function
//	parse := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil { return OptionNone[int]() }
//	    return OptionSome(n)
//	}
//	result := traversable(parse)([]string{"1", "2", "3"}) // OptionSome([1, 2, 3])
//
// See Also:
//   - TraverseArray: Direct traversal without obtaining the Traversable instance
//   - TraverseArrayG: Generic version supporting custom slice types
func OptionTraversableArray[A, B any]() OptionTraversable[A, B, []A, []B] {
	return OptionTraverseArrayG[[]A, []B, A, B]
}
