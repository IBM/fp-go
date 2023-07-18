package readerioeither

import (
	G "github.com/IBM/fp-go/readerioeither/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[R, E, A any](a ReaderIOEither[R, E, A]) ReaderIOEither[R, E, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, T.Tuple1[A]],
	](a)
}

func SequenceT2[R, E, A, B any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B]) ReaderIOEither[R, E, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, B],
		ReaderIOEither[R, E, T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[R, E, A, B, C any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B], c ReaderIOEither[R, E, C]) ReaderIOEither[R, E, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, B],
		ReaderIOEither[R, E, C],
		ReaderIOEither[R, E, T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[R, E, A, B, C, D any](a ReaderIOEither[R, E, A], b ReaderIOEither[R, E, B], c ReaderIOEither[R, E, C], d ReaderIOEither[R, E, D]) ReaderIOEither[R, E, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderIOEither[R, E, A],
		ReaderIOEither[R, E, B],
		ReaderIOEither[R, E, C],
		ReaderIOEither[R, E, D],
		ReaderIOEither[R, E, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
