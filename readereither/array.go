package readereither

import (
	G "github.com/ibm/fp-go/readereither/generic"
)

// TraverseArray transforms an array
func TraverseArray[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B] {
	return G.TraverseArray[ReaderEither[E, L, B], ReaderEither[E, L, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, L, A any](ma []ReaderEither[E, L, A]) ReaderEither[E, L, []A] {
	return G.SequenceArray[ReaderEither[E, L, A], ReaderEither[E, L, []A]](ma)
}
