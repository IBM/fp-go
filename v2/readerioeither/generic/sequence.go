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
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/internal/apply"
	T "github.com/IBM/fp-go/v2/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[
	GA ~func(E) GIOA,
	GTA ~func(E) GIOTA,
	GIOA ~func() either.Either[L, A],
	GIOTA ~func() either.Either[L, T.Tuple1[A]],
	E, L, A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, GIOA, GIOTA, E, L, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[
	GA ~func(E) GIOA,
	GB ~func(E) GIOB,
	GTAB ~func(E) GIOTAB,
	GIOA ~func() either.Either[L, A],
	GIOB ~func() either.Either[L, B],
	GIOTAB ~func() either.Either[L, T.Tuple2[A, B]],
	E, L, A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func(E) func() either.Either[L, func(B) T.Tuple2[A, B]], GIOA, func() either.Either[L, func(B) T.Tuple2[A, B]], E, L, A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func(E) func() either.Either[L, func(B) T.Tuple2[A, B]], GIOB, GIOTAB, func() either.Either[L, func(B) T.Tuple2[A, B]], E, L, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[
	GA ~func(E) GIOA,
	GB ~func(E) GIOB,
	GC ~func(E) GIOC,
	GTABC ~func(E) GIOTABC,
	GIOA ~func() either.Either[L, A],
	GIOB ~func() either.Either[L, B],
	GIOC ~func() either.Either[L, C],
	GIOTABC ~func() either.Either[L, T.Tuple3[A, B, C]],
	E, L, A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func(E) func() either.Either[L, func(B) func(C) T.Tuple3[A, B, C]], GIOA, func() either.Either[L, func(B) func(C) T.Tuple3[A, B, C]], E, L, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func(E) func() either.Either[L, func(C) T.Tuple3[A, B, C]], func(E) func() either.Either[L, func(B) func(C) T.Tuple3[A, B, C]], GIOB, func() either.Either[L, func(C) T.Tuple3[A, B, C]], func() either.Either[L, func(B) func(C) T.Tuple3[A, B, C]], E, L, B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func(E) func() either.Either[L, func(C) T.Tuple3[A, B, C]], GIOC, GIOTABC, func() either.Either[L, func(C) T.Tuple3[A, B, C]], E, L, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[
	GA ~func(E) GIOA,
	GB ~func(E) GIOB,
	GC ~func(E) GIOC,
	GD ~func(E) GIOD,
	GTABCD ~func(E) GIOTABCD,
	GIOA ~func() either.Either[L, A],
	GIOB ~func() either.Either[L, B],
	GIOC ~func() either.Either[L, C],
	GIOD ~func() either.Either[L, D],
	GIOTABCD ~func() either.Either[L, T.Tuple4[A, B, C, D]],
	E, L, A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func(E) func() either.Either[L, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], GIOA, func() either.Either[L, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], E, L, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func(E) func() either.Either[L, func(C) func(D) T.Tuple4[A, B, C, D]], func(E) func() either.Either[L, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], GIOB, func() either.Either[L, func(C) func(D) T.Tuple4[A, B, C, D]], func() either.Either[L, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], E, L, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func(E) func() either.Either[L, func(D) T.Tuple4[A, B, C, D]], func(E) func() either.Either[L, func(C) func(D) T.Tuple4[A, B, C, D]], GIOC, func() either.Either[L, func(D) T.Tuple4[A, B, C, D]], func() either.Either[L, func(C) func(D) T.Tuple4[A, B, C, D]], E, L, C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func(E) func() either.Either[L, func(D) T.Tuple4[A, B, C, D]], GIOD, GIOTABCD, func() either.Either[L, func(D) T.Tuple4[A, B, C, D]], E, L, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
