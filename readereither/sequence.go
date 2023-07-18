package readereither

import (
	G "github.com/IBM/fp-go/readereither/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[L, E, A any](a ReaderEither[E, L, A]) ReaderEither[E, L, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderEither[E, L, A],
		ReaderEither[E, L, T.Tuple1[A]],
	](a)
}

func SequenceT2[L, E, A, B any](
	a ReaderEither[E, L, A],
	b ReaderEither[E, L, B],
) ReaderEither[E, L, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderEither[E, L, A],
		ReaderEither[E, L, B],
		ReaderEither[E, L, T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[L, E, A, B, C any](
	a ReaderEither[E, L, A],
	b ReaderEither[E, L, B],
	c ReaderEither[E, L, C],
) ReaderEither[E, L, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderEither[E, L, A],
		ReaderEither[E, L, B],
		ReaderEither[E, L, C],
		ReaderEither[E, L, T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[L, E, A, B, C, D any](
	a ReaderEither[E, L, A],
	b ReaderEither[E, L, B],
	c ReaderEither[E, L, C],
	d ReaderEither[E, L, D],
) ReaderEither[E, L, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderEither[E, L, A],
		ReaderEither[E, L, B],
		ReaderEither[E, L, C],
		ReaderEither[E, L, D],
		ReaderEither[E, L, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
