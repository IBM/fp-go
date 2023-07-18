package reader

import (
	G "github.com/ibm/fp-go/reader/generic"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[R, A any](a Reader[R, A]) Reader[R, T.Tuple1[A]] {
	return G.SequenceT1[Reader[R, A], Reader[R, T.Tuple1[A]]](a)
}

func SequenceT2[R, A, B any](a Reader[R, A], b Reader[R, B]) Reader[R, T.Tuple2[A, B]] {
	return G.SequenceT2[Reader[R, A], Reader[R, B], Reader[R, T.Tuple2[A, B]]](a, b)
}

func SequenceT3[R, A, B, C any](a Reader[R, A], b Reader[R, B], c Reader[R, C]) Reader[R, T.Tuple3[A, B, C]] {
	return G.SequenceT3[Reader[R, A], Reader[R, B], Reader[R, C], Reader[R, T.Tuple3[A, B, C]]](a, b, c)
}

func SequenceT4[R, A, B, C, D any](a Reader[R, A], b Reader[R, B], c Reader[R, C], d Reader[R, D]) Reader[R, T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[Reader[R, A], Reader[R, B], Reader[R, C], Reader[R, D], Reader[R, T.Tuple4[A, B, C, D]]](a, b, c, d)
}
