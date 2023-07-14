package readerio

import (
	IO "github.com/ibm/fp-go/io"
	R "github.com/ibm/fp-go/readerio/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) ReaderIO[B]) func([]A) ReaderIO[[]B] {
	return R.TraverseArray[ReaderIO[B], ReaderIO[[]B], IO.IO[B], IO.IO[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderIO[A]) ReaderIO[[]A] {
	return R.SequenceArray[ReaderIO[A], ReaderIO[[]A]](ma)
}
