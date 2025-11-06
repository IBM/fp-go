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

package readerio

import (
	G "github.com/IBM/fp-go/v2/readerio/generic"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[R, A any](a ReaderIO[R, A]) ReaderIO[R, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderIO[R, A],
		ReaderIO[R, T.Tuple1[A]],
	](a)
}

func SequenceT2[R, A, B any](a ReaderIO[R, A], b ReaderIO[R, B]) ReaderIO[R, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[R, A, B, C any](a ReaderIO[R, A], b ReaderIO[R, B], c ReaderIO[R, C]) ReaderIO[R, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, C],
		ReaderIO[R, T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[R, A, B, C, D any](a ReaderIO[R, A], b ReaderIO[R, B], c ReaderIO[R, C], d ReaderIO[R, D]) ReaderIO[R, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, C],
		ReaderIO[R, D],
		ReaderIO[R, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
