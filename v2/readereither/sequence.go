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

package readereither

import (
	G "github.com/IBM/fp-go/v2/readereither/generic"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[L, E, A any](a ReaderEither[E, L, A]) ReaderEither[E, L, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderEither[E, L, A],
		ReaderEither[E, L, T.Tuple1[A]],
	](a)
}

func SequenceT2[L, E, A, B any](
	a ReaderEither[E, L, A],
	b ReaderEither[E, L, B],
) ReaderEither[E, L, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderEither[E, L, A],
		ReaderEither[E, L, B],
		ReaderEither[E, L, T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[L, E, A, B, C any](
	a ReaderEither[E, L, A],
	b ReaderEither[E, L, B],
	c ReaderEither[E, L, C],
) ReaderEither[E, L, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderEither[E, L, A],
		ReaderEither[E, L, B],
		ReaderEither[E, L, C],
		ReaderEither[E, L, T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[L, E, A, B, C, D any](
	a ReaderEither[E, L, A],
	b ReaderEither[E, L, B],
	c ReaderEither[E, L, C],
	d ReaderEither[E, L, D],
) ReaderEither[E, L, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderEither[E, L, A],
		ReaderEither[E, L, B],
		ReaderEither[E, L, C],
		ReaderEither[E, L, D],
		ReaderEither[E, L, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
