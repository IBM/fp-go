package generic

import (
	"github.com/ibm/fp-go/internal/apply"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[
	GA ~func(E) GIOA,
	GTA ~func(E) GIOTA,
	GIOA ~func() A,
	GIOTA ~func() T.Tuple1[A],
	E, A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, GIOA, GIOTA, E, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[
	GA ~func(E) GIOA,
	GB ~func(E) GIOB,
	GTAB ~func(E) GIOTAB,
	GIOA ~func() A,
	GIOB ~func() B,
	GIOTAB ~func() T.Tuple2[A, B],
	E, A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func(E) func() func(B) T.Tuple2[A, B], GIOA, func() func(B) T.Tuple2[A, B], E, A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func(E) func() func(B) T.Tuple2[A, B], GIOB, GIOTAB, func() func(B) T.Tuple2[A, B], E, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[
	GA ~func(E) GIOA,
	GB ~func(E) GIOB,
	GC ~func(E) GIOC,
	GTABC ~func(E) GIOTABC,
	GIOA ~func() A,
	GIOB ~func() B,
	GIOC ~func() C,
	GIOTABC ~func() T.Tuple3[A, B, C],
	E, A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func(E) func() func(B) func(C) T.Tuple3[A, B, C], GIOA, func() func(B) func(C) T.Tuple3[A, B, C], E, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func(E) func() func(C) T.Tuple3[A, B, C], func(E) func() func(B) func(C) T.Tuple3[A, B, C], GIOB, func() func(C) T.Tuple3[A, B, C], func() func(B) func(C) T.Tuple3[A, B, C], E, B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func(E) func() func(C) T.Tuple3[A, B, C], GIOC, GIOTABC, func() func(C) T.Tuple3[A, B, C], E, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[
	GA ~func(E) GIOA,
	GB ~func(E) GIOB,
	GC ~func(E) GIOC,
	GD ~func(E) GIOD,
	GTABCD ~func(E) GIOTABCD,
	GIOA ~func() A,
	GIOB ~func() B,
	GIOC ~func() C,
	GIOD ~func() D,
	GIOTABCD ~func() T.Tuple4[A, B, C, D],
	E, A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func(E) func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], GIOA, func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], E, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func(E) func() func(C) func(D) T.Tuple4[A, B, C, D], func(E) func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], GIOB, func() func(C) func(D) T.Tuple4[A, B, C, D], func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], E, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func(E) func() func(D) T.Tuple4[A, B, C, D], func(E) func() func(C) func(D) T.Tuple4[A, B, C, D], GIOC, func() func(D) T.Tuple4[A, B, C, D], func() func(C) func(D) T.Tuple4[A, B, C, D], E, C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func(E) func() func(D) T.Tuple4[A, B, C, D], GIOD, GIOTABCD, func() func(D) T.Tuple4[A, B, C, D], E, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
