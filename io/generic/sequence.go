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
		Ap[GB, GTAB, func() func(B) T.Tuple2[A, B], B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[GA ~func() A, GB ~func() B, GC ~func() C, GTABC ~func() T.Tuple3[A, B, C], A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func() func(B) func(C) T.Tuple3[A, B, C], A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func() func(C) T.Tuple3[A, B, C], func() func(B) func(C) T.Tuple3[A, B, C], B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func() func(C) T.Tuple3[A, B, C], C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[GA ~func() A, GB ~func() B, GC ~func() C, GD ~func() D, GTABCD ~func() T.Tuple4[A, B, C, D], A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func() func(C) func(D) T.Tuple4[A, B, C, D], func() func(B) func(C) func(D) T.Tuple4[A, B, C, D], B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func() func(D) T.Tuple4[A, B, C, D], func() func(C) func(D) T.Tuple4[A, B, C, D], C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func() func(D) T.Tuple4[A, B, C, D], D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
