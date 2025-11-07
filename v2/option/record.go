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
	RR "github.com/IBM/fp-go/v2/internal/record"
)

// TraverseRecordG transforms a record (map) by applying a function that returns an Option to each value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
// This is the generic version that works with custom map types.
//
// Example:
//
//	validate := func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecordG[map[string]int, map[string]int](validate)(input) // Some(map[a:2 b:4])
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return RR.Traverse[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseRecord transforms a record (map) by applying a function that returns an Option to each value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
//
// Example:
//
//	validate := func(x int) Option[string] {
//	    if x > 0 { return Some(fmt.Sprintf("%d", x)) }
//	    return None[string]()
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecord(validate)(input) // Some(map[a:"1" b:"2"])
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

// TraverseRecordWithIndexG transforms a record by applying a function that receives both key and value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
// This is the generic version that works with custom map types.
//
// Example:
//
//	f := func(k string, v int) Option[string] {
//	    return Some(fmt.Sprintf("%s:%d", k, v))
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecordWithIndexG[map[string]int, map[string]string](f)(input) // Some(map[a:"a:1" b:"b:2"])
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) Option[B]) Kleisli[GA, GB] {
	return RR.TraverseWithIndex[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseRecordWithIndex transforms a record by applying a function that receives both key and value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
//
// Example:
//
//	f := func(k string, v int) Option[int] {
//	    if v > 0 { return Some(v) }
//	    return None[int]()
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecordWithIndex(f)(input) // Some(map[a:1 b:2])
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Option[B]) Kleisli[map[K]A, map[K]B] {
	return TraverseRecordWithIndexG[map[K]A, map[K]B](f)
}

// SequenceRecordG converts a map of Options into an Option of a map.
// Returns Some containing all key-value pairs if all Options are Some, None if any is None.
// This is the generic version that works with custom map types.
//
// Example:
//
//	type MyMap map[string]int
//	input := map[string]Option[int]{"a": Some(1), "b": Some(2)}
//	result := SequenceRecordG[MyMap](input) // Some(MyMap{"a": 1, "b": 2})
func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Option[A], K comparable, A any](ma GOA) Option[GA] {
	return TraverseRecordG[GOA, GA](F.Identity[Option[A]])(ma)
}

// SequenceRecord converts a map of Options into an Option of a map.
// Returns Some containing all key-value pairs if all Options are Some, None if any is None.
//
// Example:
//
//	input := map[string]Option[int]{"a": Some(1), "b": Some(2), "c": Some(3)}
//	result := SequenceRecord(input) // Some(map[a:1 b:2 c:3])
//	input := map[string]Option[int]{"a": Some(1), "b": None[int]()}
//	result := SequenceRecord(input) // None
func SequenceRecord[K comparable, A any](ma map[K]Option[A]) Option[map[K]A] {
	return SequenceRecordG[map[K]A](ma)
}

func upsertAtReadWrite[M ~map[K]V, K comparable, V any](r M, k K, v V) M {
	r[k] = v
	return r
}

// CompactRecordG filters a map of Options, keeping only the Some values and discarding None values.
// This is the generic version that works with custom map types.
//
// Example:
//
//	type MyMap map[string]int
//	input := map[string]Option[int]{"a": Some(1), "b": None[int](), "c": Some(3)}
//	result := CompactRecordG[map[string]Option[int], MyMap](input) // MyMap{"a": 1, "c": 3}
func CompactRecordG[M1 ~map[K]Option[A], M2 ~map[K]A, K comparable, A any](m M1) M2 {
	bnd := F.Bind12of3(upsertAtReadWrite[M2])
	return RR.ReduceWithIndex(m, func(key K, m M2, value Option[A]) M2 {
		return MonadFold(value, F.Constant(m), bnd(m, key))
	}, make(M2))
}

// CompactRecord filters a map of Options, keeping only the Some values and discarding None values.
//
// Example:
//
//	input := map[string]Option[int]{"a": Some(1), "b": None[int](), "c": Some(3)}
//	result := CompactRecord(input) // map[a:1 c:3]
func CompactRecord[K comparable, A any](m map[K]Option[A]) map[K]A {
	return CompactRecordG[map[K]Option[A], map[K]A](m)
}
