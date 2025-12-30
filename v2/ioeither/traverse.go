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

package ioeither

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/record"
)

// TraverseArray transforms an array
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B] {
	return array.Traverse[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		Ap[[]B, E, B],

		f,
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[E, A, B any](f func(int, A) IOEither[E, B]) Kleisli[E, []A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		Ap[[]B, E, B],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return TraverseArray(function.Identity[IOEither[E, A]])(ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		Ap[map[K]B, E, B],

		f,
	)
}

// TraverseRecordWithIndex transforms a record
func TraverseRecordWithIndex[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		Ap[map[K]B, E, B],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return TraverseRecord[K](function.Identity[IOEither[E, A]])(ma)
}

// TraverseArraySeq transforms an array
func TraverseArraySeq[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B] {
	return array.Traverse[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		ApSeq[[]B, E, B],

		f,
	)
}

// TraverseArrayWithIndexSeq transforms an array
func TraverseArrayWithIndexSeq[E, A, B any](f func(int, A) IOEither[E, B]) Kleisli[E, []A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		ApSeq[[]B, E, B],

		f,
	)
}

// SequenceArraySeq converts a homogeneous sequence of either into an either of sequence
func SequenceArraySeq[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return TraverseArraySeq(function.Identity[IOEither[E, A]])(ma)
}

// TraverseRecordSeq transforms a record
func TraverseRecordSeq[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, E, B],

		f,
	)
}

// TraverseRecordWithIndexSeq transforms a record
func TraverseRecordWithIndexSeq[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, E, B],

		f,
	)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
func SequenceRecordSeq[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return TraverseRecordSeq[K](function.Identity[IOEither[E, A]])(ma)
}

// TraverseArrayPar transforms an array
func TraverseArrayPar[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B] {
	return array.Traverse[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		ApPar[[]B, E, B],

		f,
	)
}

// TraverseArrayWithIndexPar transforms an array
func TraverseArrayWithIndexPar[E, A, B any](f func(int, A) IOEither[E, B]) Kleisli[E, []A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		ApPar[[]B, E, B],

		f,
	)
}

// SequenceArrayPar converts a homogeneous sequence of either into an either of sequence
func SequenceArrayPar[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return TraverseArrayPar(function.Identity[IOEither[E, A]])(ma)
}

// TraverseRecordPar transforms a record
func TraverseRecordPar[K comparable, E, A, B any](f Kleisli[E, A, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.Traverse[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		ApPar[map[K]B, E, B],

		f,
	)
}

// TraverseRecordWithIndexPar transforms a record
func TraverseRecordWithIndexPar[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		ApPar[map[K]B, E, B],

		f,
	)
}

// SequenceRecordPar converts a homogeneous sequence of either into an either of sequence
func SequenceRecordPar[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return TraverseRecordPar[K](function.Identity[IOEither[E, A]])(ma)
}
