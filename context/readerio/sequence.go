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

package readerio

import (
	R "github.com/IBM/fp-go/readerio/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a ReaderIO[A]) ReaderIO[T.Tuple1[A]] {
	return R.SequenceT1[
		ReaderIO[A],
		ReaderIO[T.Tuple1[A]],
	](a)
}

func SequenceT2[A, B any](a ReaderIO[A], b ReaderIO[B]) ReaderIO[T.Tuple2[A, B]] {
	return R.SequenceT2[
		ReaderIO[A],
		ReaderIO[B],
		ReaderIO[T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[A, B, C any](a ReaderIO[A], b ReaderIO[B], c ReaderIO[C]) ReaderIO[T.Tuple3[A, B, C]] {
	return R.SequenceT3[
		ReaderIO[A],
		ReaderIO[B],
		ReaderIO[C],
		ReaderIO[T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[A, B, C, D any](a ReaderIO[A], b ReaderIO[B], c ReaderIO[C], d ReaderIO[D]) ReaderIO[T.Tuple4[A, B, C, D]] {
	return R.SequenceT4[
		ReaderIO[A],
		ReaderIO[B],
		ReaderIO[C],
		ReaderIO[D],
		ReaderIO[T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
