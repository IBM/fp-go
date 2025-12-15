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

package result

import (
	"github.com/IBM/fp-go/v2/either"
)

// TraverseArrayG transforms an array by applying a function that returns an Either to each element.
// If any element produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the array of all Right values.
// The G suffix indicates support for generic slice types.
//
// Example:
//
//	parse := func(s string) either.Result[int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"})
//	// result is Right([]int{1, 2, 3})
//
//go:inline
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return either.TraverseArrayG[GA, GB](f)
}

// TraverseArray transforms an array by applying a function that returns an Either to each element.
// If any element produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the array of all Right values.
//
// Example:
//
//	parse := func(s string) either.Result[int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseArray(parse)([]string{"1", "2", "3"})
//	// result is Right([]int{1, 2, 3})
//
//go:inline
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return either.TraverseArray(f)
}

// TraverseArrayWithIndexG transforms an array by applying an indexed function that returns an Either.
// The function receives both the index and the element.
// If any element produces a Left, the entire result is that Left (short-circuits).
// The G suffix indicates support for generic slice types.
//
// Example:
//
//	validate := func(i int, s string) either.Result[string] {
//	    if S.IsNonEmpty(s) {
//	        return either.Right[error](fmt.Sprintf("%d:%s", i, s))
//	    }
//	    return either.Left[string](fmt.Errorf("empty at index %d", i))
//	}
//	result := either.TraverseArrayWithIndexG[[]string, []string](validate)([]string{"a", "b"})
//	// result is Right([]string{"0:a", "1:b"})
//
//go:inline
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) Result[B]) Kleisli[GA, GB] {
	return either.TraverseArrayWithIndexG[GA, GB](f)
}

// TraverseArrayWithIndex transforms an array by applying an indexed function that returns an Either.
// The function receives both the index and the element.
// If any element produces a Left, the entire result is that Left (short-circuits).
//
// Example:
//
//	validate := func(i int, s string) either.Result[string] {
//	    if S.IsNonEmpty(s) {
//	        return either.Right[error](fmt.Sprintf("%d:%s", i, s))
//	    }
//	    return either.Left[string](fmt.Errorf("empty at index %d", i))
//	}
//	result := either.TraverseArrayWithIndex(validate)([]string{"a", "b"})
//	// result is Right([]string{"0:a", "1:b"})
//
//go:inline
func TraverseArrayWithIndex[A, B any](f func(int, A) Result[B]) Kleisli[[]A, []B] {
	return either.TraverseArrayWithIndex(f)
}

//go:inline
func SequenceArrayG[GA ~[]A, GOA ~[]Result[A], A any](ma GOA) Result[GA] {
	return either.SequenceArrayG[GA](ma)
}

// SequenceArray converts a homogeneous sequence of Either into an Either of sequence.
// If any element is Left, returns that Left (short-circuits).
// Otherwise, returns Right containing all the Right values.
//
// Example:
//
//	eithers := []either.Result[int]{
//	    either.Right[error](1),
//	    either.Right[error](2),
//	    either.Right[error](3),
//	}
//	result := either.SequenceArray(eithers)
//	// result is Right([]int{1, 2, 3})
//
//go:inline
func SequenceArray[A any](ma []Result[A]) Result[[]A] {
	return either.SequenceArray(ma)
}

// CompactArrayG discards all Left values and keeps only the Right values.
// The G suffix indicates support for generic slice types.
//
// Example:
//
//	eithers := []either.Result[int]{
//	    either.Right[error](1),
//	    either.Left[int](errors.New("error")),
//	    either.Right[error](3),
//	}
//	result := either.CompactArrayG[[]either.Result[int], []int](eithers)
//	// result is []int{1, 3}
//
//go:inline
func CompactArrayG[A1 ~[]Result[A], A2 ~[]A, A any](fa A1) A2 {
	return either.CompactArrayG[A1, A2](fa)
}

// CompactArray discards all Left values and keeps only the Right values.
//
// Example:
//
//	eithers := []either.Result[int]{
//	    either.Right[error](1),
//	    either.Left[int](errors.New("error")),
//	    either.Right[error](3),
//	}
//	result := either.CompactArray(eithers)
//	// result is []int{1, 3}
//
//go:inline
func CompactArray[A any](fa []Result[A]) []A {
	return either.CompactArray(fa)
}
