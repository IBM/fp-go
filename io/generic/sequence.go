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
	"github.com/IBM/fp-go/internal/apply"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[GA ~func() A, GTA ~func() T.Tuple1[A], A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, A, T.Tuple1[A]],
		a,
	)
}

func SequenceT2[GA ~func() A, GB ~func() B, GTAB ~func() T.Tuple2[A, B], A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func() func(B) T.Tuple2[A, B], A, func(B) T.Tuple2[A, B]],
		Ap[GTAB, func() func(B) T.Tuple2[A, B], GB],

		a, b,
	)
}

func SequenceT3[GA ~func() A, GB ~func() B, GC ~func() C, GTABC ~func() T.Tuple3[A, B, C], A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func() func(B) func(C) T.Tuple3[A, B, C], A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[func() func(C) T.Tuple3[A, B, C], func() func(B) func(C) T.Tuple3[A, B, C], GB],
		Ap[GTABC, func() func(C) T.Tuple3[A, B, C], GC],

		a, b, c,
	)
}

func SequenceT4[GA ~func() A, GB ~func() B, GC ~func() C, GD ~func() D, GTABCD ~func() T.Tuple4[A, B, C, D], A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[func() func(C) func(D) T.Tuple4[A, B, C, D], func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], GB],
		Ap[func() func(D) T.Tuple4[A, B, C, D], func() func(C) func(D) T.Tuple4[A, B, C, D], GC],
		Ap[GTABCD, func() func(D) T.Tuple4[A, B, C, D], GD],

		a, b, c, d,
	)
}
