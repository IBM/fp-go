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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	INTI "github.com/IBM/fp-go/v2/internal/iter"
	"github.com/IBM/fp-go/v2/internal/record"
)

// TraverseArray transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]].
// This uses the default applicative behavior (parallel or sequential based on useParallel flag).
//
// Parameters:
//   - f: Function that transforms each element into a ReaderIOResult
//
// Returns a function that transforms an array into a ReaderIOResult of an array.
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return array.Traverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		F.Flow2(f, WithContext),
	)
}

func TraverseIter[A, B any](f Kleisli[A, B]) Kleisli[Seq[A], Seq[B]] {
	return INTI.Traverse[Seq[A]](
		Map[B],

		Of[Seq[B]],
		Map[Seq[B]],
		Ap[Seq[B]],

		F.Flow2(f, WithContext),
	)
}

// TraverseArrayWithIndex transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]].
// The transformation function receives both the index and the element.
//
// Parameters:
//   - f: Function that transforms each element with its index into a ReaderIOResult
//
// Returns a function that transforms an array into a ReaderIOResult of an array.
func TraverseArrayWithIndex[A, B any](f func(int, A) ReaderIOResult[B]) Kleisli[[]A, []B] {
	return array.TraverseWithIndex(
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of ReaderIOResult into a ReaderIOResult of sequence.
// This is equivalent to TraverseArray with the identity function.
//
// Parameters:
//   - ma: Array of ReaderIOResult values
//
// Returns a ReaderIOResult containing an array of values.
func SequenceArray[A any](ma []ReaderIOResult[A]) ReaderIOResult[[]A] {
	return TraverseArray(function.Identity[ReaderIOResult[A]])(ma)
}

// TraverseRecord transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]].
//
// Parameters:
//   - f: Function that transforms each value into a ReaderIOResult
//
// Returns a function that transforms a map into a ReaderIOResult of a map.
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		F.Flow2(f, WithContext),
	)
}

// TraverseRecordWithIndex transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]].
// The transformation function receives both the key and the value.
//
// Parameters:
//   - f: Function that transforms each key-value pair into a ReaderIOResult
//
// Returns a function that transforms a map into a ReaderIOResult of a map.
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) ReaderIOResult[B]) Kleisli[map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// SequenceRecord converts a homogeneous map of ReaderIOResult into a ReaderIOResult of map.
//
// Parameters:
//   - ma: Map of ReaderIOResult values
//
// Returns a ReaderIOResult containing a map of values.
func SequenceRecord[K comparable, A any](ma map[K]ReaderIOResult[A]) ReaderIOResult[map[K]A] {
	return TraverseRecord[K](function.Identity[ReaderIOResult[A]])(ma)
}

// MonadTraverseArraySeq transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]].
// This explicitly uses sequential execution.
//
// Parameters:
//   - as: The array to traverse
//   - f: Function that transforms each element into a ReaderIOResult
//
// Returns a ReaderIOResult containing an array of transformed values.
func MonadTraverseArraySeq[A, B any](as []A, f Kleisli[A, B]) ReaderIOResult[[]B] {
	return array.MonadTraverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],
		as,
		F.Flow2(f, WithContext),
	)
}

// TraverseArraySeq transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]].
// This is the curried version of [MonadTraverseArraySeq] with sequential execution.
//
// Parameters:
//   - f: Function that transforms each element into a ReaderIOResult
//
// Returns a function that transforms an array into a ReaderIOResult of an array.
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return array.Traverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],
		F.Flow2(f, WithContext),
	)
}

// TraverseArrayWithIndexSeq uses transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]]
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) ReaderIOResult[B]) Kleisli[[]A, []B] {
	return array.TraverseWithIndex(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],
		f,
	)
}

// SequenceArraySeq converts a homogeneous sequence of ReaderIOResult into a ReaderIOResult of sequence.
// This explicitly uses sequential execution.
//
// Parameters:
//   - ma: Array of ReaderIOResult values
//
// Returns a ReaderIOResult containing an array of values.
func SequenceArraySeq[A any](ma []ReaderIOResult[A]) ReaderIOResult[[]A] {
	return MonadTraverseArraySeq(ma, function.Identity[ReaderIOResult[A]])
}

// MonadTraverseRecordSeq uses transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]]
func MonadTraverseRecordSeq[K comparable, A, B any](as map[K]A, f Kleisli[A, B]) ReaderIOResult[map[K]B] {
	return record.MonadTraverse(
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],
		as,
		F.Flow2(f, WithContext),
	)
}

// TraverseRecordSeq uses transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]]
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		F.Flow2(f, WithContext),
	)
}

// TraverseRecordWithIndexSeq uses transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]]
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) ReaderIOResult[B]) Kleisli[map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
func SequenceRecordSeq[K comparable, A any](ma map[K]ReaderIOResult[A]) ReaderIOResult[map[K]A] {
	return MonadTraverseRecordSeq(ma, function.Identity[ReaderIOResult[A]])
}

// MonadTraverseArrayPar transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]].
// This explicitly uses parallel execution.
//
// Parameters:
//   - as: The array to traverse
//   - f: Function that transforms each element into a ReaderIOResult
//
// Returns a ReaderIOResult containing an array of transformed values.
func MonadTraverseArrayPar[A, B any](as []A, f Kleisli[A, B]) ReaderIOResult[[]B] {
	return array.MonadTraverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApPar[[]B, B],
		as,
		F.Flow2(f, WithContext),
	)
}

// TraverseArrayPar transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]].
// This is the curried version of [MonadTraverseArrayPar] with parallel execution.
//
// Parameters:
//   - f: Function that transforms each element into a ReaderIOResult
//
// Returns a function that transforms an array into a ReaderIOResult of an array.
func TraverseArrayPar[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return array.Traverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApPar[[]B, B],
		F.Flow2(f, WithContext),
	)
}

// TraverseArrayWithIndexPar uses transforms an array [[]A] into [[]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[[]B]]
func TraverseArrayWithIndexPar[A, B any](f func(int, A) ReaderIOResult[B]) Kleisli[[]A, []B] {
	return array.TraverseWithIndex(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApPar[[]B, B],
		f,
	)
}

// SequenceArrayPar converts a homogeneous sequence of ReaderIOResult into a ReaderIOResult of sequence.
// This explicitly uses parallel execution.
//
// Parameters:
//   - ma: Array of ReaderIOResult values
//
// Returns a ReaderIOResult containing an array of values.
func SequenceArrayPar[A any](ma []ReaderIOResult[A]) ReaderIOResult[[]A] {
	return MonadTraverseArrayPar(ma, function.Identity[ReaderIOResult[A]])
}

// TraverseRecordPar uses transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]]
func TraverseRecordPar[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApPar[map[K]B, B],

		F.Flow2(f, WithContext),
	)
}

// TraverseRecordWithIndexPar uses transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]]
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) ReaderIOResult[B]) Kleisli[map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApPar[map[K]B, B],

		f,
	)
}

// MonadTraverseRecordPar uses transforms a record [map[K]A] into [map[K]ReaderIOResult[B]] and then resolves that into a [ReaderIOResult[map[K]B]]
func MonadTraverseRecordPar[K comparable, A, B any](as map[K]A, f Kleisli[A, B]) ReaderIOResult[map[K]B] {
	return record.MonadTraverse(
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApPar[map[K]B, B],
		as,
		F.Flow2(f, WithContext),
	)
}

// SequenceRecordPar converts a homogeneous map of ReaderIOResult into a ReaderIOResult of map.
// This explicitly uses parallel execution.
//
// Parameters:
//   - ma: Map of ReaderIOResult values
//
// Returns a ReaderIOResult containing a map of values.
func SequenceRecordPar[K comparable, A any](ma map[K]ReaderIOResult[A]) ReaderIOResult[map[K]A] {
	return MonadTraverseRecordPar(ma, function.Identity[ReaderIOResult[A]])
}

// TraversableArray returns a Traversable instance for arrays.
// This provides a higher-order function that can be used to traverse arrays
// with any transformation function.
//
// Type Parameters:
//   - A: The element type of the input array
//   - B: The element type of the output array
//
// Returns:
//   - Traversable[A, B, []A, []B]: A function that takes a Kleisli arrow and returns
//     a function that transforms arrays
//
// Example:
//
//	parse := result.Eitherize1(strconv.Atoi)
//	traversable := TraversableArray[string, int]()
//	transform := traversable(parse)
//	result := transform([]string{"1", "2", "3"})(ctx)()
func TraversableArray[A, B any]() Traversable[A, B, []A, []B] {
	return TraverseArray[A, B]
}

// TraversableRecord returns a Traversable instance for maps.
// This provides a higher-order function that can be used to traverse maps
// with any transformation function.
//
// Type Parameters:
//   - K: The key type of the map (must be comparable)
//   - A: The value type of the input map
//   - B: The value type of the output map
//
// Returns:
//   - Traversable[A, B, map[K]A, map[K]B]: A function that takes a Kleisli arrow and
//     returns a function that transforms maps
//
// Example:
//
//	parse := result.Eitherize1(strconv.Atoi)
//	traversable := TraversableRecord[string, string, int]()
//	transform := traversable(parse)
//	input := map[string]string{"a": "1", "b": "2"}
//	result := transform(input)(ctx)()
func TraversableRecord[K comparable, A, B any]() Traversable[A, B, map[K]A, map[K]B] {
	return TraverseRecord[K, A, B]
}

// TraversableIter returns a Traversable instance for iterators.
// This provides a higher-order function that can be used to traverse iterators
// with any transformation function.
//
// Type Parameters:
//   - A: The element type of the input iterator
//   - B: The element type of the output iterator
//
// Returns:
//   - Traversable[A, B, Seq[A], Seq[B]]: A function that takes a Kleisli arrow and
//     returns a function that transforms iterators
//
// Example:
//
//	parse := result.Eitherize1(strconv.Atoi)
//	traversable := TraversableIter[string, int]()
//	transform := traversable(parse)
//	input := slices.Values([]string{"1", "2", "3"})
//	result := transform(input)(ctx)()
func TraversableIter[A, B any]() Traversable[A, B, Seq[A], Seq[B]] {
	return TraverseIter[A, B]
}
