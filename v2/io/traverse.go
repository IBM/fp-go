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

package io

import (
	F "github.com/IBM/fp-go/v2/function"
	INTA "github.com/IBM/fp-go/v2/internal/array"
	INTR "github.com/IBM/fp-go/v2/internal/record"
)

func MonadTraverseArray[A, B any](tas []A, f func(A) IO[B]) IO[[]B] {
	return INTA.MonadTraverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		tas,
		f,
	)
}

// TraverseArray applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArray[A, B any](f func(A) IO[B]) func([]A) IO[[]B] {
	return INTA.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// TraverseArrayWithIndex applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArrayWithIndex[A, B any](f func(int, A) IO[B]) func([]A) IO[[]B] {
	return INTA.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		Ap[[]B, B],

		f,
	)
}

// SequenceArray converts an array of [IO] to an [IO] of an array
func SequenceArray[A any](tas []IO[A]) IO[[]A] {
	return MonadTraverseArray(tas, F.Identity[IO[A]])
}

func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f func(A) IO[B]) IO[map[K]B] {
	return INTR.MonadTraverse(
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		tas,
		f,
	)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecord[K comparable, A, B any](f func(A) IO[B]) func(map[K]A) IO[map[K]B] {
	return INTR.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndex applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IO[B]) func(map[K]A) IO[map[K]B] {
	return INTR.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		Ap[map[K]B, B],

		f,
	)
}

// SequenceRecord converts a record of [IO] to an [IO] of a record
func SequenceRecord[K comparable, A any](tas map[K]IO[A]) IO[map[K]A] {
	return MonadTraverseRecord(tas, F.Identity[IO[A]])
}

func MonadTraverseArraySeq[A, B any](tas []A, f func(A) IO[B]) IO[[]B] {
	return INTA.MonadTraverse(
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		tas,
		f,
	)
}

// TraverseArraySeq applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArraySeq[A, B any](f func(A) IO[B]) func([]A) IO[[]B] {
	return INTA.Traverse[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// TraverseArrayWithIndexSeq applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IO[B]) func([]A) IO[[]B] {
	return INTA.TraverseWithIndex[[]A](
		Of[[]B],
		Map[[]B, func(B) []B],
		ApSeq[[]B, B],

		f,
	)
}

// SequenceArraySeq converts an array of [IO] to an [IO] of an array
func SequenceArraySeq[A any](tas []IO[A]) IO[[]A] {
	return MonadTraverseArraySeq(tas, F.Identity[IO[A]])
}

func MonadTraverseRecordSeq[K comparable, A, B any](tas map[K]A, f func(A) IO[B]) IO[map[K]B] {
	return INTR.MonadTraverse(
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		tas,
		f,
	)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordSeq[K comparable, A, B any](f func(A) IO[B]) func(map[K]A) IO[map[K]B] {
	return INTR.Traverse[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// TraverseRecordWithIndexSeq applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordWithIndeSeq[K comparable, A, B any](f func(K, A) IO[B]) func(map[K]A) IO[map[K]B] {
	return INTR.TraverseWithIndex[map[K]A](
		Of[map[K]B],
		Map[map[K]B, func(B) map[K]B],
		ApSeq[map[K]B, B],

		f,
	)
}

// SequenceRecordSeq converts a record of [IO] to an [IO] of a record
func SequenceRecordSeq[K comparable, A any](tas map[K]IO[A]) IO[map[K]A] {
	return MonadTraverseRecordSeq(tas, F.Identity[IO[A]])
}
