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

package lazy

import "github.com/IBM/fp-go/v2/io"

// MonadTraverseArray applies a function returning a lazy computation to all elements
// in an array and transforms this into a lazy computation of that array.
//
// This is the monadic version of TraverseArray, taking the array as the first parameter.
//
// Example:
//
//	numbers := []int{1, 2, 3}
//	result := lazy.MonadTraverseArray(numbers, func(x int) lazy.Lazy[int] {
//	    return lazy.Of(x * 2)
//	})()
//	// result is []int{2, 4, 6}
func MonadTraverseArray[A, B any](tas []A, f Kleisli[A, B]) Lazy[[]B] {
	return io.MonadTraverseArray(tas, f)
}

// TraverseArray applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return io.TraverseArray(f)
}

// TraverseArrayWithIndex applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArrayWithIndex[A, B any](f func(int, A) Lazy[B]) Kleisli[[]A, []B] {
	return io.TraverseArrayWithIndex(f)
}

// SequenceArray converts an array of [IO] to an [IO] of an array
func SequenceArray[A any](tas []Lazy[A]) Lazy[[]A] {
	return io.SequenceArray(tas)
}

// MonadTraverseRecord applies a function returning a lazy computation to all values
// in a record (map) and transforms this into a lazy computation of that record.
//
// This is the monadic version of TraverseRecord, taking the record as the first parameter.
//
// Example:
//
//	record := map[string]int{"a": 1, "b": 2}
//	result := lazy.MonadTraverseRecord(record, func(x int) lazy.Lazy[int] {
//	    return lazy.Of(x * 2)
//	})()
//	// result is map[string]int{"a": 2, "b": 4}
func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f Kleisli[A, B]) Lazy[map[K]B] {
	return io.MonadTraverseRecord(tas, f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return io.TraverseRecord[K](f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Lazy[B]) Kleisli[map[K]A, map[K]B] {
	return io.TraverseRecordWithIndex(f)
}

// SequenceRecord converts a record of [IO] to an [IO] of a record
func SequenceRecord[K comparable, A any](tas map[K]Lazy[A]) Lazy[map[K]A] {
	return io.SequenceRecord(tas)
}
