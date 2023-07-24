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
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[GA ~func() O.Option[A], GTA ~func() O.Option[T.Tuple1[A]], A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[GA ~func() O.Option[A], GB ~func() O.Option[B], GTAB ~func() O.Option[T.Tuple2[A, B]], A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func() O.Option[func(B) T.Tuple2[A, B]], A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func() O.Option[func(B) T.Tuple2[A, B]], B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[GA ~func() O.Option[A], GB ~func() O.Option[B], GC ~func() O.Option[C], GTABC ~func() O.Option[T.Tuple3[A, B, C]], A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func() O.Option[func(B) func(C) T.Tuple3[A, B, C]], A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func() O.Option[func(C) T.Tuple3[A, B, C]], func() O.Option[func(B) func(C) T.Tuple3[A, B, C]], B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func() O.Option[func(C) T.Tuple3[A, B, C]], C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[GA ~func() O.Option[A], GB ~func() O.Option[B], GC ~func() O.Option[C], GD ~func() O.Option[D], GTABCD ~func() O.Option[T.Tuple4[A, B, C, D]], A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func() O.Option[func(B) func(C) func(D) T.Tuple4[A, B, C, D]], A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func() O.Option[func(C) func(D) T.Tuple4[A, B, C, D]], func() O.Option[func(B) func(C) func(D) T.Tuple4[A, B, C, D]], B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func() O.Option[func(D) T.Tuple4[A, B, C, D]], func() O.Option[func(C) func(D) T.Tuple4[A, B, C, D]], C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func() O.Option[func(D) T.Tuple4[A, B, C, D]], D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
