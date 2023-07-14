package either

import (
	"github.com/ibm/fp-go/internal/apply"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[E, A any](a Either[E, A]) Either[E, T.Tuple1[A]] {
	return apply.SequenceT1(
		Map[E, A, T.Tuple1[A]],
		a,
	)
}

func SequenceT2[E, A, B any](a Either[E, A], b Either[E, B]) Either[E, T.Tuple2[A, B]] {
	return apply.SequenceT2(
		Map[E, A, func(B) T.Tuple2[A, B]],
		Ap[T.Tuple2[A, B], E, B],

		a, b,
	)
}

func SequenceT3[E, A, B, C any](a Either[E, A], b Either[E, B], c Either[E, C]) Either[E, T.Tuple3[A, B, C]] {
	return apply.SequenceT3(
		Map[E, A, func(B) func(C) T.Tuple3[A, B, C]],
		Ap[func(C) T.Tuple3[A, B, C], E, B],
		Ap[T.Tuple3[A, B, C], E, C],

		a, b, c,
	)
}

func SequenceT4[E, A, B, C, D any](a Either[E, A], b Either[E, B], c Either[E, C], d Either[E, D]) Either[E, T.Tuple4[A, B, C, D]] {
	return apply.SequenceT4(
		Map[E, A, func(B) func(C) func(D) T.Tuple4[A, B, C, D]],
		Ap[func(C) func(D) T.Tuple4[A, B, C, D], E, B],
		Ap[func(D) T.Tuple4[A, B, C, D], E, C],
		Ap[T.Tuple4[A, B, C, D], E, D],

		a, b, c, d,
	)
}
