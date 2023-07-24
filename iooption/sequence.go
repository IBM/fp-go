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

package iooption

import (
	G "github.com/IBM/fp-go/iooption/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a IOOption[A]) IOOption[T.Tuple1[A]] {
	return G.SequenceT1[
		IOOption[A],
		IOOption[T.Tuple1[A]],
	](a)
}

func SequenceT2[A, B any](
	a IOOption[A],
	b IOOption[B],
) IOOption[T.Tuple2[A, B]] {
	return G.SequenceT2[
		IOOption[A],
		IOOption[B],
		IOOption[T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[A, B, C any](
	a IOOption[A],
	b IOOption[B],
	c IOOption[C],
) IOOption[T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		IOOption[A],
		IOOption[B],
		IOOption[C],
		IOOption[T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[A, B, C, D any](
	a IOOption[A],
	b IOOption[B],
	c IOOption[C],
	d IOOption[D],
) IOOption[T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		IOOption[A],
		IOOption[B],
		IOOption[C],
		IOOption[D],
		IOOption[T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
