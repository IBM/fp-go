package readerio

import (
	IO "github.com/ibm/fp-go/io"
	G "github.com/ibm/fp-go/readerio/generic"
)

// TraverseArray transforms an array
func TraverseArray[R, A, B any](f func(A) ReaderIO[R, B]) func([]A) ReaderIO[R, []B] {
	return G.TraverseArray[ReaderIO[R, B], ReaderIO[R, []B], IO.IO[B], IO.IO[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of Readers into a Reader of sequence
func SequenceArray[R, A any](ma []ReaderIO[R, A]) ReaderIO[R, []A] {
	return G.SequenceArray[ReaderIO[R, A], ReaderIO[R, []A]](ma)
}
