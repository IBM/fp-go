package option

import (
	Apply "github.com/ibm/fp-go/apply"
	F "github.com/ibm/fp-go/function"
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
	return Apply.SequenceT1(
		MonadMap[A, T.Tuple1[A]],

		a,
	)
}

func SequenceT2[A, B any](a Option[A], b Option[B]) Option[T.Tuple2[A, B]] {
	return Apply.SequenceT2(
		MonadMap[A, func(B) T.Tuple2[A, B]],
		MonadAp[B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[A, B, C any](a Option[A], b Option[B], c Option[C]) Option[T.Tuple3[A, B, C]] {
	return Apply.SequenceT3(
		MonadMap[A, func(B) func(C) T.Tuple3[A, B, C]],
		MonadAp[B, func(C) T.Tuple3[A, B, C]],
		MonadAp[C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[A, B, C, D any](a Option[A], b Option[B], c Option[C], d Option[D]) Option[T.Tuple4[A, B, C, D]] {
	return Apply.SequenceT4(
		MonadMap[A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		MonadAp[B, func(C) func(D) T.Tuple4[A, B, C, D]],
		MonadAp[C, func(D) T.Tuple4[A, B, C, D]],
		MonadAp[D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
