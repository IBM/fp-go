package readereither

import (
	RE "github.com/IBM/fp-go/readereither/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) ReaderEither[B]) func([]A) ReaderEither[[]B] {
	return RE.TraverseArray[ReaderEither[B], ReaderEither[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderEither[A]) ReaderEither[[]A] {
	return RE.SequenceArray[ReaderEither[A], ReaderEither[[]A]](ma)
}
