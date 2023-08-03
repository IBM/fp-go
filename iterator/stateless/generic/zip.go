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

package generic

import (
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// ZipWith applies a function to pairs of elements at the same index in two iterators, collecting the results in a new iterator. If one
// input iterator is short, excess elements of the longer iterator are discarded.
func ZipWith[AS ~func() O.Option[T.Tuple2[AS, A]], BS ~func() O.Option[T.Tuple2[BS, B]], CS ~func() O.Option[T.Tuple2[CS, C]], FCT ~func(A, B) C, A, B, C any](fa AS, fb BS, f FCT) CS {
	// pre-declare to avoid cyclic reference
	var m func(T.Tuple2[O.Option[T.Tuple2[AS, A]], O.Option[T.Tuple2[BS, B]]]) O.Option[T.Tuple2[CS, C]]

	recurse := func(as AS, bs BS) CS {
		return func() O.Option[T.Tuple2[CS, C]] {
			// combine
			return F.Pipe1(
				T.MakeTuple2(as(), bs()),
				m,
			)
		}
	}

	m = F.Flow2(
		O.SequenceTuple2[T.Tuple2[AS, A], T.Tuple2[BS, B]],
		O.Map(func(t T.Tuple2[T.Tuple2[AS, A], T.Tuple2[BS, B]]) T.Tuple2[CS, C] {
			return T.MakeTuple2(recurse(t.F1.F1, t.F2.F1), f(t.F1.F2, t.F2.F2))
		}))

	// trigger the recursion
	return recurse(fa, fb)
}

// Zip takes two iterators and returns an iterators of corresponding pairs. If one input iterators is short, excess elements of the
// longer iterator are discarded
func Zip[AS ~func() O.Option[T.Tuple2[AS, A]], BS ~func() O.Option[T.Tuple2[BS, B]], CS ~func() O.Option[T.Tuple2[CS, T.Tuple2[A, B]]], A, B any](fb BS) func(AS) CS {
	return F.Bind23of3(ZipWith[AS, BS, CS, func(A, B) T.Tuple2[A, B]])(fb, T.MakeTuple2[A, B])
}
