package io

import (
	G "github.com/ibm/fp-go/io/generic"
	T "github.com/ibm/fp-go/tuple"
)

// SequenceT converts n inputs of higher kinded types into a higher kinded types of n strongly typed values, represented as a tuple

func SequenceT1[A any](a IO[A]) IO[T.Tuple1[A]] {
	return G.SequenceT1[IO[A], IO[T.Tuple1[A]]](a)
}

func SequenceT2[A, B any](a IO[A], b IO[B]) IO[T.Tuple2[A, B]] {
	return G.SequenceT2[IO[A], IO[B], IO[T.Tuple2[A, B]]](a, b)
}

func SequenceT3[A, B, C any](a IO[A], b IO[B], c IO[C]) IO[T.Tuple3[A, B, C]] {
	return G.SequenceT3[IO[A], IO[B], IO[C], IO[T.Tuple3[A, B, C]]](a, b, c)
}

func SequenceT4[A, B, C, D any](a IO[A], b IO[B], c IO[C], d IO[D]) IO[T.Tuple4[A, B, C, D]] {
	return G.SequenceT4[IO[A], IO[B], IO[C], IO[D], IO[T.Tuple4[A, B, C, D]]](a, b, c, d)
}
