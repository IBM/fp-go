package option

import (
	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/apply"
	T "github.com/ibm/fp-go/tuple"
)

// HKTA = HKT<A>
// HKTOA = HKT<Option<A>>
//
// Sequence converts an option of some higher kinded type into the higher kinded type of an option
func Sequence[A, HKTA, HKTOA any](
	_of func(Option[A]) HKTOA,
	_map func(HKTA, func(A) Option[A]) HKTOA,
) func(Option[HKTA]) HKTOA {
	return Fold(F.Nullary2(None[A], _of), F.Bind2nd(_map, Some[A]))
}

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a Option[A]) Option[T.Tuple1[A]] {
	return apply.SequenceT1(
		Map[A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[A, B any](a Option[A], b Option[B]) Option[T.Tuple2[A, B]] {
	return apply.SequenceT2(
		Map[A, func(B) T.Tuple2[A, B]],
		Ap[T.Tuple2[A, B], B],

		a, b,
	)
}

func SequenceT3[A, B, C any](a Option[A], b Option[B], c Option[C]) Option[T.Tuple3[A, B, C]] {
	return apply.SequenceT3(
		Map[A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[func(C) T.Tuple3[A, B, C], B],
		Ap[T.Tuple3[A, B, C], C],

		a, b, c,
	)
}

func SequenceT4[A, B, C, D any](a Option[A], b Option[B], c Option[C], d Option[D]) Option[T.Tuple4[A, B, C, D]] {
	return apply.SequenceT4(
		Map[A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[func(C) func(D) T.Tuple4[A, B, C, D], B],
		Ap[func(D) T.Tuple4[A, B, C, D], C],
		Ap[T.Tuple4[A, B, C, D], D],

		a, b, c, d,
	)
}
