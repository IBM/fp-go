package reader

import (
	R "github.com/IBM/fp-go/reader/generic"
	T "github.com/IBM/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a Reader[A]) Reader[T.Tuple1[A]] {
	return R.SequenceT1[
		Reader[A],
		Reader[T.Tuple1[A]],
	](a)
}

func SequenceT2[A, B any](a Reader[A], b Reader[B]) Reader[T.Tuple2[A, B]] {
	return R.SequenceT2[
		Reader[A],
		Reader[B],
		Reader[T.Tuple2[A, B]],
	](a, b)
}

func SequenceT3[A, B, C any](a Reader[A], b Reader[B], c Reader[C]) Reader[T.Tuple3[A, B, C]] {
	return R.SequenceT3[
		Reader[A],
		Reader[B],
		Reader[C],
		Reader[T.Tuple3[A, B, C]],
	](a, b, c)
}

func SequenceT4[A, B, C, D any](a Reader[A], b Reader[B], c Reader[C], d Reader[D]) Reader[T.Tuple4[A, B, C, D]] {
	return R.SequenceT4[
		Reader[A],
		Reader[B],
		Reader[C],
		Reader[D],
		Reader[T.Tuple4[A, B, C, D]],
	](a, b, c, d)
}
