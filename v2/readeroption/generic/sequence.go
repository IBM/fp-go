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
	"github.com/IBM/fp-go/v2/internal/apply"
	O "github.com/IBM/fp-go/v2/option"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[
	GA ~func(E) O.Option[A],
	GTA ~func(E) O.Option[T.Tuple1[A]],
	E, A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, E, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[
	GA ~func(E) O.Option[A],
	GB ~func(E) O.Option[B],
	GTAB ~func(E) O.Option[T.Tuple2[A, B]],
	E, A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func(E) O.Option[func(B) T.Tuple2[A, B]], E, A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func(E) O.Option[func(B) T.Tuple2[A, B]], E, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[
	GA ~func(E) O.Option[A],
	GB ~func(E) O.Option[B],
	GC ~func(E) O.Option[C],
	GTABC ~func(E) O.Option[T.Tuple3[A, B, C]],
	E, A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func(E) O.Option[func(B) func(C) T.Tuple3[A, B, C]], E, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func(E) O.Option[func(C) T.Tuple3[A, B, C]], func(E) O.Option[func(B) func(C) T.Tuple3[A, B, C]], E, B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func(E) O.Option[func(C) T.Tuple3[A, B, C]], E, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[
	GA ~func(E) O.Option[A],
	GB ~func(E) O.Option[B],
	GC ~func(E) O.Option[C],
	GD ~func(E) O.Option[D],
	GTABCD ~func(E) O.Option[T.Tuple4[A, B, C, D]],
	E, A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func(E) O.Option[func(B) func(C) func(D) T.Tuple4[A, B, C, D]], E, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func(E) O.Option[func(C) func(D) T.Tuple4[A, B, C, D]], func(E) O.Option[func(B) func(C) func(D) T.Tuple4[A, B, C, D]], E, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func(E) O.Option[func(D) T.Tuple4[A, B, C, D]], func(E) O.Option[func(C) func(D) T.Tuple4[A, B, C, D]], E, C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func(E) O.Option[func(D) T.Tuple4[A, B, C, D]], E, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
