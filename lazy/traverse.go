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

package lazy

import (
	G "github.com/IBM/fp-go/io/generic"
)

func MonadTraverseArray[A, B any](tas []A, f func(A) Lazy[B]) Lazy[[]B] {
	return G.MonadTraverseArray[Lazy[B], Lazy[[]B]](tas, f)
}

// TraverseArray applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArray[A, B any](f func(A) Lazy[B]) func([]A) Lazy[[]B] {
	return G.TraverseArray[Lazy[B], Lazy[[]B], []A](f)
}

// TraverseArrayWithIndex applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArrayWithIndex[A, B any](f func(int, A) Lazy[B]) func([]A) Lazy[[]B] {
	return G.TraverseArrayWithIndex[Lazy[B], Lazy[[]B], []A](f)
}

// SequenceArray converts an array of [IO] to an [IO] of an array
func SequenceArray[A any](tas []Lazy[A]) Lazy[[]A] {
	return G.SequenceArray[Lazy[A], Lazy[[]A]](tas)
}

func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f func(A) Lazy[B]) Lazy[map[K]B] {
	return G.MonadTraverseRecord[Lazy[B], Lazy[map[K]B]](tas, f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecord[K comparable, A, B any](f func(A) Lazy[B]) func(map[K]A) Lazy[map[K]B] {
	return G.TraverseRecord[Lazy[B], Lazy[map[K]B], map[K]A](f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Lazy[B]) func(map[K]A) Lazy[map[K]B] {
	return G.TraverseRecordWithIndex[Lazy[B], Lazy[map[K]B], map[K]A](f)
}

// SequenceRecord converts a record of [IO] to an [IO] of a record
func SequenceRecord[K comparable, A any](tas map[K]Lazy[A]) Lazy[map[K]A] {
	return G.SequenceRecord[Lazy[A], Lazy[map[K]A]](tas)
}
