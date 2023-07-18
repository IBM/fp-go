package generic

import (
	ET "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/internal/apply"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[
	GA ~func(E) ET.Either[L, A],
	GTA ~func(E) ET.Either[L, T.Tuple1[A]],
	L, E, A any](a GA) GTA {
	return apply.SequenceT1(
		Map[GA, GTA, L, E, A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[
	GA ~func(E) ET.Either[L, A],
	GB ~func(E) ET.Either[L, B],
	GTAB ~func(E) ET.Either[L, T.Tuple2[A, B]],
	L, E, A, B any](a GA, b GB) GTAB {
	return apply.SequenceT2(
		Map[GA, func(E) ET.Either[L, func(B) T.Tuple2[A, B]], L, E, A, func(B) T.Tuple2[A, B]],
		Ap[GB, GTAB, func(E) ET.Either[L, func(B) T.Tuple2[A, B]], L, E, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[
	GA ~func(E) ET.Either[L, A],
	GB ~func(E) ET.Either[L, B],
	GC ~func(E) ET.Either[L, C],
	GTABC ~func(E) ET.Either[L, T.Tuple3[A, B, C]],
	L, E, A, B, C any](a GA, b GB, c GC) GTABC {
	return apply.SequenceT3(
		Map[GA, func(E) ET.Either[L, func(B) func(C) T.Tuple3[A, B, C]], L, E, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[GB, func(E) ET.Either[L, func(C) T.Tuple3[A, B, C]], func(E) ET.Either[L, func(B) func(C) T.Tuple3[A, B, C]], L, E, B, func(C) T.Tuple3[A, B, C]],
		Ap[GC, GTABC, func(E) ET.Either[L, func(C) T.Tuple3[A, B, C]], L, E, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[
	GA ~func(E) ET.Either[L, A],
	GB ~func(E) ET.Either[L, B],
	GC ~func(E) ET.Either[L, C],
	GD ~func(E) ET.Either[L, D],
	GTABCD ~func(E) ET.Either[L, T.Tuple4[A, B, C, D]],
	L, E, A, B, C, D any](a GA, b GB, c GC, d GD) GTABCD {
	return apply.SequenceT4(
		Map[GA, func(E) ET.Either[L, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], L, E, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GB, func(E) ET.Either[L, func(C) func(D) T.Tuple4[A, B, C, D]], func(E) ET.Either[L, func(B) func(C) func(D) T.Tuple4[A, B, C, D]], L, E, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[GC, func(E) ET.Either[L, func(D) T.Tuple4[A, B, C, D]], func(E) ET.Either[L, func(C) func(D) T.Tuple4[A, B, C, D]], L, E, C, func(D) T.Tuple4[A, B, C, D]],
		Ap[GD, GTABCD, func(E) ET.Either[L, func(D) T.Tuple4[A, B, C, D]], L, E, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
