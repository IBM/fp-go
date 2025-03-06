// Copyright (c) 2023 IBM Corp.
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

package readerioeither

import (
	G "github.com/IBM/fp-go/v2/context/readerioeither/generic"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/record"
)

// TraverseArray uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArray[A, B any](f func(A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return array.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// TraverseArrayWithIndex uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArrayWithIndex[A, B any](f func(int, A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return array.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderIOEither[A]) ReaderIOEither[[]A] {
	return TraverseArray(function.Identity[ReaderIOEither[A]])(ma)
}

// TraverseRecord uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecord[K comparable, A, B any](f func(A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return record.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndex uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]ReaderIOEither[A]) ReaderIOEither[map[K]A] {
	return TraverseRecord[K](function.Identity[ReaderIOEither[A]])(ma)
}

// TraverseArraySeq uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArraySeq[A, B any](f func(A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return array.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// TraverseArrayWithIndexSeq uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return G.TraverseArrayWithIndexSeq[[]A, ReaderIOEither[[]B]](f)
}

// SequenceArraySeq converts a homogeneous sequence of either into an either of sequence
func SequenceArraySeq[A any](ma []ReaderIOEither[A]) ReaderIOEither[[]A] {
	return G.SequenceArraySeq[[]A, []ReaderIOEither[A], ReaderIOEither[[]A]](ma)
}

// TraverseRecordSeq uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecordSeq[K comparable, A, B any](f func(A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return G.TraverseRecordSeq[K, map[K]A, ReaderIOEither[map[K]B]](f)
}

// TraverseRecordWithIndexSeq uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return G.TraverseRecordWithIndexSeq[K, map[K]A, ReaderIOEither[map[K]B]](f)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
func SequenceRecordSeq[K comparable, A any](ma map[K]ReaderIOEither[A]) ReaderIOEither[map[K]A] {
	return G.SequenceRecordSeq[K, map[K]A, map[K]ReaderIOEither[A], ReaderIOEither[map[K]A]](ma)
}

// TraverseArrayPar uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArrayPar[A, B any](f func(A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return G.TraverseArrayPar[[]A, ReaderIOEither[[]B]](f)
}

// TraverseArrayWithIndexPar uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArrayWithIndexPar[A, B any](f func(int, A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return G.TraverseArrayWithIndexPar[[]A, ReaderIOEither[[]B]](f)
}

// SequenceArrayPar converts a homogeneous sequence of either into an either of sequence
func SequenceArrayPar[A any](ma []ReaderIOEither[A]) ReaderIOEither[[]A] {
	return G.SequenceArrayPar[[]A, []ReaderIOEither[A], ReaderIOEither[[]A]](ma)
}

// TraverseRecordPar uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecordPar[K comparable, A, B any](f func(A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return G.TraverseRecordPar[K, map[K]A, ReaderIOEither[map[K]B]](f)
}

// TraverseRecordWithIndexPar uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return G.TraverseRecordWithIndexPar[K, map[K]A, ReaderIOEither[map[K]B]](f)
}

// SequenceRecordPar converts a homogeneous sequence of either into an either of sequence
func SequenceRecordPar[K comparable, A any](ma map[K]ReaderIOEither[A]) ReaderIOEither[map[K]A] {
	return G.SequenceRecordPar[K, map[K]A, map[K]ReaderIOEither[A], ReaderIOEither[map[K]A]](ma)
}
