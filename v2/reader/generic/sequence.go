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
	"github.com/IBM/fp-go/v2/internal/apply"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[GA ~func(R) A, GTA ~func(R) T.Tuple1[A], R, A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, R, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[GA ~func(R) A, GB ~func(R) B, GTAB ~func(R) T.Tuple2[A, B], R, A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func(R) func(B) T.Tuple2[A, B], R, A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func(R) func(B) T.Tuple2[A, B], R, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[GA ~func(R) A, GB ~func(R) B, GC ~func(R) C, GTABC ~func(R) T.Tuple3[A, B, C], R, A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func(R) func(B) func(C) T.Tuple3[A, B, C], R, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func(R) func(C) T.Tuple3[A, B, C], func(R) func(B) func(C) T.Tuple3[A, B, C], R, B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func(R) func(C) T.Tuple3[A, B, C], R, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[GA ~func(R) A, GB ~func(R) B, GC ~func(R) C, GD ~func(R) D, GTABCD ~func(R) T.Tuple4[A, B, C, D], R, A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func(R) func(B) func(C) func(D) T.Tuple4[A, B, C, D], R, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func(R) func(C) func(D) T.Tuple4[A, B, C, D], func(R) func(B) func(C) func(D) T.Tuple4[A, B, C, D], R, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func(R) func(D) T.Tuple4[A, B, C, D], func(R) func(C) func(D) T.Tuple4[A, B, C, D], R, C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func(R) func(D) T.Tuple4[A, B, C, D], R, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
