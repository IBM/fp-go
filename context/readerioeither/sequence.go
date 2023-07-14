package readerioeither

import (
	RE "github.com/ibm/fp-go/readerioeither/generic"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a ReaderIOEither[A]) ReaderIOEither[T.Tuple1[A]] {
	return RE.SequenceT1[
		ReaderIOEither[A],
		ReaderIOEither[T.Tuple1[A]],
	](a)
}

func SequenceT2[A, B any](a ReaderIOEither[A], b ReaderIOEither[B]) ReaderIOEither[T.Tuple2[A, B]] {
	return RE.SequenceT2[
		ReaderIOEither[A],
		ReaderIOEither[B],
		ReaderIOEither[T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[A, B, C any](a ReaderIOEither[A], b ReaderIOEither[B], c ReaderIOEither[C]) ReaderIOEither[T.Tuple3[A, B, C]] {
	return RE.SequenceT3[
		ReaderIOEither[A],
		ReaderIOEither[B],
		ReaderIOEither[C],
		ReaderIOEither[T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[A, B, C, D any](a ReaderIOEither[A], b ReaderIOEither[B], c ReaderIOEither[C], d ReaderIOEither[D]) ReaderIOEither[T.Tuple4[A, B, C, D]] {
	return RE.SequenceT4[
		ReaderIOEither[A],
		ReaderIOEither[B],
		ReaderIOEither[C],
		ReaderIOEither[D],
		ReaderIOEither[T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
