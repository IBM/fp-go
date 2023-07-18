package ioeither

import (
	G "github.com/IBM/fp-go/ioeither/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[E, A any](a IOEither[E, A]) IOEither[E, T.Tuple1[A]] {
	return G.SequenceT1[
		IOEither[E, A],
		IOEither[E, T.Tuple1[A]],
	](a)
}

func SequenceT2[E, A, B any](
	a IOEither[E, A],
	b IOEither[E, B],
) IOEither[E, T.Tuple2[A, B]] {
	return G.SequenceT2[
		IOEither[E, A],
		IOEither[E, B],
		IOEither[E, T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[E, A, B, C any](
	a IOEither[E, A],
	b IOEither[E, B],
	c IOEither[E, C],
) IOEither[E, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		IOEither[E, A],
		IOEither[E, B],
		IOEither[E, C],
		IOEither[E, T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[E, A, B, C, D any](
	a IOEither[E, A],
	b IOEither[E, B],
	c IOEither[E, C],
	d IOEither[E, D],
) IOEither[E, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		IOEither[E, A],
		IOEither[E, B],
		IOEither[E, C],
		IOEither[E, D],
		IOEither[E, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
