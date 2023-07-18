package readerio

import (
	G "github.com/ibm/fp-go/readerio/generic"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[R, A any](a ReaderIO[R, A]) ReaderIO[R, T.Tuple1[A]] {
	return G.SequenceT1[
		ReaderIO[R, A],
		ReaderIO[R, T.Tuple1[A]],
	](a)
}

func SequenceT2[R, A, B any](a ReaderIO[R, A], b ReaderIO[R, B]) ReaderIO[R, T.Tuple2[A, B]] {
	return G.SequenceT2[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[R, A, B, C any](a ReaderIO[R, A], b ReaderIO[R, B], c ReaderIO[R, C]) ReaderIO[R, T.Tuple3[A, B, C]] {
	return G.SequenceT3[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, C],
		ReaderIO[R, T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[R, A, B, C, D any](a ReaderIO[R, A], b ReaderIO[R, B], c ReaderIO[R, C], d ReaderIO[R, D]) ReaderIO[R, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[
		ReaderIO[R, A],
		ReaderIO[R, B],
		ReaderIO[R, C],
		ReaderIO[R, D],
		ReaderIO[R, T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
