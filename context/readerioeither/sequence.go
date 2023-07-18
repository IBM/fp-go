package readerioeither

import (
	G "github.com/IBM/fp-go/context/readerioeither/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a ReaderIOEither[A]) ReaderIOEither[T.Tuple1[A]] {
	return G.SequenceT1[ReaderIOEither[T.Tuple1[A]]](a)
}

func SequenceT2[A, B any](a ReaderIOEither[A], b ReaderIOEither[B]) ReaderIOEither[T.Tuple2[A, B]] {
	return G.SequenceT2[ReaderIOEither[T.Tuple2[A, B]]](a, b)
}

func SequenceT3[A, B, C any](a ReaderIOEither[A], b ReaderIOEither[B], c ReaderIOEither[C]) ReaderIOEither[T.Tuple3[A, B, C]] {
	return G.SequenceT3[ReaderIOEither[T.Tuple3[A, B, C]]](a, b, c)
}

func SequenceT4[A, B, C, D any](a ReaderIOEither[A], b ReaderIOEither[B], c ReaderIOEither[C], d ReaderIOEither[D]) ReaderIOEither[T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[ReaderIOEither[T.Tuple4[A, B, C, D]]](a, b, c, d)
}
