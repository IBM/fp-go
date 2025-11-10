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

// TraverseRecordG transforms a map by applying a function that returns an Either to each value.
// If any value produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the map of all Right values.
// The G suffix indicates support for generic map types.
//
// Example:
//
//	parse := func(s string) either.Result[int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseRecordG[map[string]string, map[string]int](parse)(map[string]string{"a": "1", "b": "2"})
//	// result is Right(map[string]int{"a": 1, "b": 2})
//
//go:inline
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return either.TraverseRecordG[GA, GB](f)
}

// TraverseRecord transforms a map by applying a function that returns an Either to each value.
// If any value produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the map of all Right values.
//
// Example:
//
//	parse := func(s string) either.Result[int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseRecord[string](parse)(map[string]string{"a": "1", "b": "2"})
//	// result is Right(map[string]int{"a": 1, "b": 2})
//
//go:inline
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return either.TraverseRecord[K](f)
}

// TraverseRecordWithIndexG transforms a map by applying an indexed function that returns an Either.
// The function receives both the key and the value.
// If any value produces a Left, the entire result is that Left (short-circuits).
// The G suffix indicates support for generic map types.
//
// Example:
//
//	validate := func(k string, v string) either.Result[string] {
//	    if len(v) > 0 {
//	        return either.Right[error](k + ":" + v)
//	    }
//	    return either.Left[string](fmt.Errorf("empty value for key %s", k))
//	}
//	result := either.TraverseRecordWithIndexG[map[string]string, map[string]string](validate)(map[string]string{"a": "1"})
//	// result is Right(map[string]string{"a": "a:1"})
//
//go:inline
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) Result[B]) Kleisli[GA, GB] {
	return either.TraverseRecordWithIndexG[GA, GB](f)
}

// TraverseRecordWithIndex transforms a map by applying an indexed function that returns an Either.
// The function receives both the key and the value.
// If any value produces a Left, the entire result is that Left (short-circuits).
//
// Example:
//
//	validate := func(k string, v string) either.Result[string] {
//	    if len(v) > 0 {
//	        return either.Right[error](k + ":" + v)
//	    }
//	    return either.Left[string](fmt.Errorf("empty value for key %s", k))
//	}
//	result := either.TraverseRecordWithIndex[string](validate)(map[string]string{"a": "1"})
//	// result is Right(map[string]string{"a": "a:1"})
//
//go:inline
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Result[B]) Kleisli[map[K]A, map[K]B] {
	return either.TraverseRecordWithIndex(f)
}

//go:inline
func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Result[A], K comparable, A any](ma GOA) Result[GA] {
	return either.SequenceRecordG[GA](ma)
}

// SequenceRecord converts a map of Either values into an Either of a map.
// If any value is Left, returns that Left (short-circuits).
// Otherwise, returns Right containing a map of all the Right values.
//
// Example:
//
//	eithers := map[string]either.Result[int]{
//	    "a": either.Right[error](1),
//	    "b": either.Right[error](2),
//	}
//	result := either.SequenceRecord(eithers)
//	// result is Right(map[string]int{"a": 1, "b": 2})
//
//go:inline
func SequenceRecord[K comparable, A any](ma map[K]Result[A]) Result[map[K]A] {
	return either.SequenceRecord(ma)
}

// CompactRecordG discards all Left values and keeps only the Right values.
// The G suffix indicates support for generic map types.
//
// Example:
//
//	eithers := map[string]either.Result[int]{
//	    "a": either.Right[error](1),
//	    "b": either.Left[int](errors.New("error")),
//	    "c": either.Right[error](3),
//	}
//	result := either.CompactRecordG[map[string]either.Result[int], map[string]int](eithers)
//	// result is map[string]int{"a": 1, "c": 3}
func CompactRecordG[M1 ~map[K]Result[A], M2 ~map[K]A, K comparable, A any](m M1) M2 {
	return either.CompactRecordG[M1, M2](m)
}

// CompactRecord discards all Left values and keeps only the Right values.
//
// Example:
//
//	eithers := map[string]either.Result[int]{
//	    "a": either.Right[error](1),
//	    "b": either.Left[int](errors.New("error")),
//	    "c": either.Right[error](3),
//	}
//	result := either.CompactRecord(eithers)
//	// result is map[string]int{"a": 1, "c": 3}
//
//go:inline
func CompactRecord[K comparable, A any](m map[K]Result[A]) map[K]A {
	return either.CompactRecord(m)
}
