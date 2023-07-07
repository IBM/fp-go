package either

import (
	Apply "github.com/ibm/fp-go/apply"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[E, A any](a Either[E, A]) Either[E, T.Tuple1[A]] {
	return Apply.SequenceT1(
		MonadMap[E, A, T.Tuple1[A]],
		a,
	)
}

func SequenceT2[E, A, B any](a Either[E, A], b Either[E, B]) Either[E, T.Tuple2[A, B]] {
	return Apply.SequenceT2(
		MonadMap[E, A, func(B) T.Tuple2[A, B]],
		MonadAp[E, B, T.Tuple2[A, B]],

		a, b,
	)
}

func SequenceT3[E, A, B, C any](a Either[E, A], b Either[E, B], c Either[E, C]) Either[E, T.Tuple3[A, B, C]] {
	return Apply.SequenceT3(
		MonadMap[E, A, func(B) func(C) T.Tuple3[A, B, C]],
		MonadAp[E, B, func(C) T.Tuple3[A, B, C]],
		MonadAp[E, C, T.Tuple3[A, B, C]],

		a, b, c,
	)
}

func SequenceT4[E, A, B, C, D any](a Either[E, A], b Either[E, B], c Either[E, C], d Either[E, D]) Either[E, T.Tuple4[A, B, C, D]] {
	return Apply.SequenceT4(
		MonadMap[E, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		MonadAp[E, B, func(C) func(D) T.Tuple4[A, B, C, D]],
		MonadAp[E, C, func(D) T.Tuple4[A, B, C, D]],
		MonadAp[E, D, T.Tuple4[A, B, C, D]],

		a, b, c, d,
	)
}
