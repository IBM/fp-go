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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

// ZipWith applies a function to pairs of elements at the same index in two iterators, collecting the results in a new iterator. If one
// input iterator is short, excess elements of the longer iterator are discarded.
func ZipWith[AS ~func() O.Option[P.Pair[AS, A]], BS ~func() O.Option[P.Pair[BS, B]], CS ~func() O.Option[P.Pair[CS, C]], FCT ~func(A, B) C, A, B, C any](fa AS, fb BS, f FCT) CS {
	// pre-declare to avoid cyclic reference
	var m func(P.Pair[O.Option[P.Pair[AS, A]], O.Option[P.Pair[BS, B]]]) O.Option[P.Pair[CS, C]]

	recurse := func(as AS, bs BS) CS {
		return func() O.Option[P.Pair[CS, C]] {
			// combine
			return F.Pipe1(
				P.MakePair(as(), bs()),
				m,
			)
		}
	}

	m = F.Flow2(
		O.SequencePair[P.Pair[AS, A], P.Pair[BS, B]],
		O.Map(func(t P.Pair[P.Pair[AS, A], P.Pair[BS, B]]) P.Pair[CS, C] {
			return P.MakePair(recurse(P.Head(P.Head(t)), P.Head(P.Tail(t))), f(P.Tail(P.Head(t)), P.Tail(P.Tail(t))))
		}))

	// trigger the recursion
	return recurse(fa, fb)
}

// Zip takes two iterators and returns an iterators of corresponding pairs. If one input iterators is short, excess elements of the
// longer iterator are discarded
func Zip[AS ~func() O.Option[P.Pair[AS, A]], BS ~func() O.Option[P.Pair[BS, B]], CS ~func() O.Option[P.Pair[CS, P.Pair[A, B]]], A, B any](fb BS) func(AS) CS {
	return F.Bind23of3(ZipWith[AS, BS, CS, func(A, B) P.Pair[A, B]])(fb, P.MakePair[A, B])
}
