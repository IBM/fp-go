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

package reader

import (
	R "github.com/IBM/fp-go/reader/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a Reader[A]) Reader[T.Tuple1[A]] {
	return R.SequenceT1[
		Reader[A],
		Reader[T.Tuple1[A]],
	](a)
}

func SequenceT2[A, B any](a Reader[A], b Reader[B]) Reader[T.Tuple2[A, B]] {
	return R.SequenceT2[
		Reader[A],
		Reader[B],
		Reader[T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[A, B, C any](a Reader[A], b Reader[B], c Reader[C]) Reader[T.Tuple3[A, B, C]] {
	return R.SequenceT3[
		Reader[A],
		Reader[B],
		Reader[C],
		Reader[T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[A, B, C, D any](a Reader[A], b Reader[B], c Reader[C], d Reader[D]) Reader[T.Tuple4[A, B, C, D]] {
	return R.SequenceT4[
		Reader[A],
		Reader[B],
		Reader[C],
		Reader[D],
		Reader[T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
