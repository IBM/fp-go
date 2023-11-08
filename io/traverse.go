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
	G "github.com/IBM/fp-go/io/generic"
)

func MonadTraverseArray[A, B any](tas []A, f func(A) IO[B]) IO[[]B] {
	return G.MonadTraverseArray[IO[B], IO[[]B]](tas, f)
}

// TraverseArray applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArray[A, B any](f func(A) IO[B]) func([]A) IO[[]B] {
	return G.TraverseArray[IO[B], IO[[]B], []A](f)
}

// TraverseArrayWithIndex applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArrayWithIndex[A, B any](f func(int, A) IO[B]) func([]A) IO[[]B] {
	return G.TraverseArrayWithIndex[IO[B], IO[[]B], []A](f)
}

// SequenceArray converts an array of [IO] to an [IO] of an array
func SequenceArray[A any](tas []IO[A]) IO[[]A] {
	return G.SequenceArray[IO[A], IO[[]A]](tas)
}

func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f func(A) IO[B]) IO[map[K]B] {
	return G.MonadTraverseRecord[IO[map[K]B]](tas, f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecord[K comparable, A, B any](f func(A) IO[B]) func(map[K]A) IO[map[K]B] {
	return G.TraverseRecord[IO[map[K]B], map[K]A, IO[B]](f)
}

// TraverseRecordWithIndex applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IO[B]) func(map[K]A) IO[map[K]B] {
	return G.TraverseRecordWithIndex[IO[B], IO[map[K]B], map[K]A](f)
}

// SequenceRecord converts a record of [IO] to an [IO] of a record
func SequenceRecord[K comparable, A any](tas map[K]IO[A]) IO[map[K]A] {
	return G.SequenceRecord[IO[A], IO[map[K]A]](tas)
}
