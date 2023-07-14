package generic

import (
	ET "github.com/ibm/fp-go/either"
	"github.com/ibm/fp-go/internal/apply"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[GA ~func() ET.Either[E, A], GTA ~func() ET.Either[E, T.Tuple1[A]], E, A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, E, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GTAB ~func() ET.Either[E, T.Tuple2[A, B]], E, A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func() ET.Either[E, func(B) T.Tuple2[A, B]], E, A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func() ET.Either[E, func(B) T.Tuple2[A, B]], E, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GC ~func() ET.Either[E, C], GTABC ~func() ET.Either[E, T.Tuple3[A, B, C]], E, A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func() ET.Either[E, func(B) func(C) T.Tuple3[A, B, C]], E, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func() ET.Either[E, func(C) T.Tuple3[A, B, C]], func() ET.Either[E, func(B) func(C) T.Tuple3[A, B, C]], E, B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func() ET.Either[E, func(C) T.Tuple3[A, B, C]], E, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GC ~func() ET.Either[E, C], GD ~func() ET.Either[E, D], GTABCD ~func() ET.Either[E, T.Tuple4[A, B, C, D]], E, A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func() ET.Either[E, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], E, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func() ET.Either[E, func(C) func(D) T.Tuple4[A, B, C, D]], func() ET.Either[E, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], E, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func() ET.Either[E, func(D) T.Tuple4[A, B, C, D]], func() ET.Either[E, func(C) func(D) T.Tuple4[A, B, C, D]], E, C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func() ET.Either[E, func(D) T.Tuple4[A, B, C, D]], E, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
