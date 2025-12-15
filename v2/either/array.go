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

package either

import (
	F "github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
)

// TraverseArrayG transforms an array by applying a function that returns an Either to each element.
// If any element produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the array of all Right values.
// The G suffix indicates support for generic slice types.
//
// Example:
//
//	parse := func(s string) either.Either[error, int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"})
//	// result is Right([]int{1, 2, 3})
//
//go:inline
func TraverseArrayG[GA ~[]A, GB ~[]B, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, GA, GB] {
	return func(ga GA) Either[E, GB] {
		bs := make(GB, len(ga))
		for i, a := range ga {
			b := f(a)
			if b.isLeft {
				return Left[GB](b.l)
			}
			bs[i] = b.r
		}
		return Of[E](bs)
	}
}

// TraverseArray transforms an array by applying a function that returns an Either to each element.
// If any element produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the array of all Right values.
//
// Example:
//
//	parse := func(s string) either.Either[error, int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseArray(parse)([]string{"1", "2", "3"})
//	// result is Right([]int{1, 2, 3})
//
//go:inline
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B] {
	return TraverseArrayG[[]A, []B](f)
}

// TraverseArrayWithIndexG transforms an array by applying an indexed function that returns an Either.
// The function receives both the index and the element.
// If any element produces a Left, the entire result is that Left (short-circuits).
// The G suffix indicates support for generic slice types.
//
// Example:
//
//	validate := func(i int, s string) either.Either[error, string] {
//	    if S.IsNonEmpty(s) {
//	        return either.Right[error](fmt.Sprintf("%d:%s", i, s))
//	    }
//	    return either.Left[string](fmt.Errorf("empty at index %d", i))
//	}
//	result := either.TraverseArrayWithIndexG[[]string, []string](validate)([]string{"a", "b"})
//	// result is Right([]string{"0:a", "1:b"})
//
//go:inline
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, E, A, B any](f func(int, A) Either[E, B]) Kleisli[E, GA, GB] {
	return func(ga GA) Either[E, GB] {
		bs := make(GB, len(ga))
		for i, a := range ga {
			b := f(i, a)
			if b.isLeft {
				return Left[GB](b.l)
			}
			bs[i] = b.r
		}
		return Of[E](bs)
	}
}

// TraverseArrayWithIndex transforms an array by applying an indexed function that returns an Either.
// The function receives both the index and the element.
// If any element produces a Left, the entire result is that Left (short-circuits).
//
// Example:
//
//	validate := func(i int, s string) either.Either[error, string] {
//	    if S.IsNonEmpty(s) {
//	        return either.Right[error](fmt.Sprintf("%d:%s", i, s))
//	    }
//	    return either.Left[string](fmt.Errorf("empty at index %d", i))
//	}
//	result := either.TraverseArrayWithIndex(validate)([]string{"a", "b"})
//	// result is Right([]string{"0:a", "1:b"})
//
//go:inline
func TraverseArrayWithIndex[E, A, B any](f func(int, A) Either[E, B]) Kleisli[E, []A, []B] {
	return TraverseArrayWithIndexG[[]A, []B](f)
}

//go:inline
func SequenceArrayG[GA ~[]A, GOA ~[]Either[E, A], E, A any](ma GOA) Either[E, GA] {
	return TraverseArrayG[GOA, GA](F.Identity[Either[E, A]])(ma)
}

// SequenceArray converts a homogeneous sequence of Either into an Either of sequence.
// If any element is Left, returns that Left (short-circuits).
// Otherwise, returns Right containing all the Right values.
//
// Example:
//
//	eithers := []either.Either[error, int]{
//	    either.Right[error](1),
//	    either.Right[error](2),
//	    either.Right[error](3),
//	}
//	result := either.SequenceArray(eithers)
//	// result is Right([]int{1, 2, 3})
//
//go:inline
func SequenceArray[E, A any](ma []Either[E, A]) Either[E, []A] {
	return SequenceArrayG[[]A](ma)
}

// CompactArrayG discards all Left values and keeps only the Right values.
// The G suffix indicates support for generic slice types.
//
// Example:
//
//	eithers := []either.Either[error, int]{
//	    either.Right[error](1),
//	    either.Left[int](errors.New("error")),
//	    either.Right[error](3),
//	}
//	result := either.CompactArrayG[[]either.Either[error, int], []int](eithers)
//	// result is []int{1, 3}
//
//go:inline
func CompactArrayG[A1 ~[]Either[E, A], A2 ~[]A, E, A any](fa A1) A2 {
	return RA.Reduce(fa, func(out A2, value Either[E, A]) A2 {
		return MonadFold(value, F.Constant1[E](out), F.Bind1st(RA.Append[A2, A], out))
	}, make(A2, 0, len(fa)))
}

// CompactArray discards all Left values and keeps only the Right values.
//
// Example:
//
//	eithers := []either.Either[error, int]{
//	    either.Right[error](1),
//	    either.Left[int](errors.New("error")),
//	    either.Right[error](3),
//	}
//	result := either.CompactArray(eithers)
//	// result is []int{1, 3}
//
//go:inline
func CompactArray[E, A any](fa []Either[E, A]) []A {
	return CompactArrayG[[]Either[E, A], []A](fa)
}
