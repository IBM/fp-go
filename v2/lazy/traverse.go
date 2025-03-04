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

import "github.com/IBM/fp-go/v2/io"

func MonadTraverseArray[A, B any](tas []A, f func(A) Lazy[B]) Lazy[[]B] {
	return io.MonadTraverseArray(tas, f)
}

// TraverseArray applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArray[A, B any](f func(A) Lazy[B]) func([]A) Lazy[[]B] {
	return io.TraverseArray(f)
}

// TraverseArrayWithIndex applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArrayWithIndex[A, B any](f func(int, A) Lazy[B]) func([]A) Lazy[[]B] {
	return io.TraverseArrayWithIndex(f)
}

// SequenceArray converts an array of [IO] to an [IO] of an array
func SequenceArray[A any](tas []Lazy[A]) Lazy[[]A] {
	return io.SequenceArray(tas)
}

func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f func(A) Lazy[B]) Lazy[map[K]B] {
	return io.MonadTraverseRecord(tas, f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecord[K comparable, A, B any](f func(A) Lazy[B]) func(map[K]A) Lazy[map[K]B] {
	return io.TraverseRecord[K](f)
}

// TraverseRecord applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Lazy[B]) func(map[K]A) Lazy[map[K]B] {
	return io.TraverseRecordWithIndex[K](f)
}

// SequenceRecord converts a record of [IO] to an [IO] of a record
func SequenceRecord[K comparable, A any](tas map[K]Lazy[A]) Lazy[map[K]A] {
	return io.SequenceRecord(tas)
}
