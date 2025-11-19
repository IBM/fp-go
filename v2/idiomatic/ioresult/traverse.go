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

package ioresult

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/record"
)

// TraverseArray transforms an array
// TraverseArray transforms an array by applying an IOResult-producing function to each element.
// Uses parallel execution by default. If any element fails, the entire traversal fails.
//
//go:inline
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return array.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// TraverseArrayWithIndex transforms an array
// TraverseArrayWithIndex transforms an array with access to element indices.
// Uses parallel execution by default.
//
//go:inline
func TraverseArrayWithIndex[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
// SequenceArray converts an array of IOResults into an IOResult of an array.
// Uses parallel execution by default.
func SequenceArray[A any](ma []IOResult[A]) IOResult[[]A] {
	return TraverseArray(function.Identity[IOResult[A]])(ma)
}

// TraverseRecord transforms a record
// TraverseRecord transforms a map by applying an IOResult-producing function to each value.
// Uses parallel execution by default.
//
//go:inline
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndex transforms a record
// TraverseRecordWithIndex transforms a map with access to keys.
// Uses parallel execution by default.
//
//go:inline
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
// SequenceRecord converts a map of IOResults into an IOResult of a map.
// Uses parallel execution by default.
func SequenceRecord[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A] {
	return TraverseRecord[K](function.Identity[IOResult[A]])(ma)
}

// TraverseArraySeq transforms an array
// TraverseArraySeq transforms an array sequentially.
// Elements are processed one at a time in order.
//
//go:inline
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return array.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// TraverseArrayWithIndexSeq transforms an array
// TraverseArrayWithIndexSeq transforms an array sequentially with indices.
//
//go:inline
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// SequenceArraySeq converts a homogeneous sequence of either into an either of sequence
// SequenceArraySeq converts an array of IOResults sequentially.
func SequenceArraySeq[A any](ma []IOResult[A]) IOResult[[]A] {
	return TraverseArraySeq(function.Identity[IOResult[A]])(ma)
}

// TraverseRecordSeq transforms a record
// TraverseRecordSeq transforms a map sequentially.
//
//go:inline
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndexSeq transforms a record
// TraverseRecordWithIndexSeq transforms a map sequentially with keys.
//
//go:inline
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
// SequenceRecordSeq converts a map of IOResults sequentially.
func SequenceRecordSeq[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A] {
	return TraverseRecordSeq[K](function.Identity[IOResult[A]])(ma)
}

// TraverseArrayPar transforms an array
// TraverseArrayPar transforms an array in parallel (explicit).
// This is equivalent to TraverseArray but makes parallelism explicit.
//
//go:inline
func TraverseArrayPar[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return array.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApPar[[]B, B],

		f,
	)
}

// TraverseArrayWithIndexPar transforms an array
// TraverseArrayWithIndexPar transforms an array in parallel with indices (explicit).
//
//go:inline
func TraverseArrayWithIndexPar[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApPar[[]B, B],

		f,
	)
}

// SequenceArrayPar converts a homogeneous Paruence of either into an either of Paruence
// SequenceArrayPar converts an array of IOResults in parallel (explicit).
func SequenceArrayPar[A any](ma []IOResult[A]) IOResult[[]A] {
	return TraverseArrayPar(function.Identity[IOResult[A]])(ma)
}

// TraverseRecordPar transforms a record
// TraverseRecordPar transforms a map in parallel (explicit).
//
//go:inline
func TraverseRecordPar[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApPar[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndexPar transforms a record
// TraverseRecordWithIndexPar transforms a map in parallel with keys (explicit).
//
//go:inline
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// SequenceRecordPar converts a homogeneous Paruence of either into an either of Paruence
// SequenceRecordPar converts a map of IOResults in parallel (explicit).
func SequenceRecordPar[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A] {
	return TraverseRecordPar[K](function.Identity[IOResult[A]])(ma)
}
