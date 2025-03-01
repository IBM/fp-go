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

package readereither

import (
	RE "github.com/IBM/fp-go/v2/readereither/generic"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a ReaderEither[A]) ReaderEither[T.Tuple1[A]] {
	return RE.SequenceT1[
		ReaderEither[A],
		ReaderEither[T.Tuple1[A]],
	](a)
}

func SequenceT2[A, B any](a ReaderEither[A], b ReaderEither[B]) ReaderEither[T.Tuple2[A, B]] {
	return RE.SequenceT2[
		ReaderEither[A],
		ReaderEither[B],
		ReaderEither[T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[A, B, C any](a ReaderEither[A], b ReaderEither[B], c ReaderEither[C]) ReaderEither[T.Tuple3[A, B, C]] {
	return RE.SequenceT3[
		ReaderEither[A],
		ReaderEither[B],
		ReaderEither[C],
		ReaderEither[T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[A, B, C, D any](a ReaderEither[A], b ReaderEither[B], c ReaderEither[C], d ReaderEither[D]) ReaderEither[T.Tuple4[A, B, C, D]] {
	return RE.SequenceT4[
		ReaderEither[A],
		ReaderEither[B],
		ReaderEither[C],
		ReaderEither[D],
		ReaderEither[T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
