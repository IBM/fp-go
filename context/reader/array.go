package reader

import (
	R "github.com/ibm/fp-go/reader/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) Reader[B]) func([]A) Reader[[]B] {
	return R.TraverseArray[Reader[B], Reader[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []Reader[A]) Reader[[]A] {
	return R.SequenceArray[Reader[A], Reader[[]A]](ma)
}
