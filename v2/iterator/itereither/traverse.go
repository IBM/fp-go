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

package itereither

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
func TraverseArrayWithIndex[E, A, B any](f func(int, A) SeqEither[E, B]) Kleisli[E, []A, []B] {
	return array.TraverseWithIndex[[]A](
		Of[E, []B],
		Map[E, []B, func(B) []B],
		Ap[[]B, E, B],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, A any](ma []SeqEither[E, A]) SeqEither[E, []A] {
	return TraverseArray(function.Identity[SeqEither[E, A]])(ma)
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
func TraverseRecordWithIndex[K comparable, E, A, B any](f func(K, A) SeqEither[E, B]) Kleisli[E, map[K]A, map[K]B] {
	return record.TraverseWithIndex[map[K]A](
		Of[E, map[K]B],
		Map[E, map[K]B, func(B) map[K]B],
		Ap[map[K]B, E, B],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, E, A any](ma map[K]SeqEither[E, A]) SeqEither[E, map[K]A] {
	return TraverseRecord[K](function.Identity[SeqEither[E, A]])(ma)
}
