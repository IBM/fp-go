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

package option

import (
	F "github.com/IBM/fp-go/v2/function"
)

// TraverseArrayG transforms an array by applying a function that returns an Option to each element.
// Returns Some containing the array of results if all operations succeed, None if any fails.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	parse := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil { return None[int]() }
//	    return Some(n)
//	}
//	result := TraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"}) // Some([1, 2, 3])
//	result := TraverseArrayG[[]string, []int](parse)([]string{"1", "x", "3"}) // None
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return func(g GA) Option[GB] {
		bs := make(GB, len(g))
		for i, a := range g {
			b := f(a)
			if !b.isSome {
				return None[GB]()
			}
			bs[i] = b.value
		}
		return Some(bs)
	}
}

// TraverseArray transforms an array by applying a function that returns an Option to each element.
// Returns Some containing the array of results if all operations succeed, None if any fails.
//
// Example:
//
//	validate := func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	}
//	result := TraverseArray(validate)([]int{1, 2, 3}) // Some([2, 4, 6])
//	result := TraverseArray(validate)([]int{1, -1, 3}) // None
//
//go:inline
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return TraverseArrayG[[]A, []B](f)
}

// TraverseArrayWithIndexG transforms an array by applying an indexed function that returns an Option.
// The function receives both the index and the element.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	f := func(i int, s string) Option[string] {
//	    return Some(fmt.Sprintf("%d:%s", i, s))
//	}
//	result := TraverseArrayWithIndexG[[]string, []string](f)([]string{"a", "b"}) // Some(["0:a", "1:b"])
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) Option[B]) Kleisli[GA, GB] {
	return func(g GA) Option[GB] {
		bs := make(GB, len(g))
		for i, a := range g {
			b := f(i, a)
			if !b.isSome {
				return None[GB]()
			}
			bs[i] = b.value
		}
		return Some(bs)
	}
}

// TraverseArrayWithIndex transforms an array by applying an indexed function that returns an Option.
// The function receives both the index and the element.
//
// Example:
//
//	f := func(i int, x int) Option[int] {
//	    if x > i { return Some(x) }
//	    return None[int]()
//	}
//	result := TraverseArrayWithIndex(f)([]int{1, 2, 3}) // Some([1, 2, 3])
//
//go:inline
func TraverseArrayWithIndex[A, B any](f func(int, A) Option[B]) Kleisli[[]A, []B] {
	return TraverseArrayWithIndexG[[]A, []B](f)
}

// SequenceArrayG converts an array of Options into an Option of an array.
// Returns Some containing all values if all Options are Some, None if any is None.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	type MySlice []int
//	result := SequenceArrayG[MySlice]([]Option[int]{Some(1), Some(2)}) // Some(MySlice{1, 2})
//	result := SequenceArrayG[MySlice]([]Option[int]{Some(1), None[int]()}) // None
//
//go:inline
func SequenceArrayG[GA ~[]A, GOA ~[]Option[A], A any](ma GOA) Option[GA] {
	return TraverseArrayG[GOA, GA](F.Identity[Option[A]])(ma)
}

// SequenceArray converts an array of Options into an Option of an array.
// Returns Some containing all values if all Options are Some, None if any is None.
//
// Example:
//
//	result := SequenceArray([]Option[int]{Some(1), Some(2), Some(3)}) // Some([1, 2, 3])
//	result := SequenceArray([]Option[int]{Some(1), None[int](), Some(3)}) // None
func SequenceArray[A any](ma []Option[A]) Option[[]A] {
	return SequenceArrayG[[]A](ma)
}

// CompactArrayG filters an array of Options, keeping only the Some values and discarding None values.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	type MySlice []int
//	input := []Option[int]{Some(1), None[int](), Some(3)}
//	result := CompactArrayG[[]Option[int], MySlice](input) // MySlice{1, 3}
func CompactArrayG[A1 ~[]Option[A], A2 ~[]A, A any](fa A1) A2 {
	as := make(A2, 0, len(fa))
	for _, oa := range fa {
		if oa.isSome {
			as = append(as, oa.value)
		}
	}
	return as
}

// CompactArray filters an array of Options, keeping only the Some values and discarding None values.
//
// Example:
//
//	input := []Option[int]{Some(1), None[int](), Some(3), Some(5), None[int]()}
//	result := CompactArray(input) // [1, 3, 5]
//
//go:inline
func CompactArray[A any](fa []Option[A]) []A {
	return CompactArrayG[[]Option[A], []A](fa)
}
