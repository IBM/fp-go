package reader

import (
	G "github.com/ibm/fp-go/reader/generic"
)

// TraverseArray transforms an array
func TraverseArray[R, A, B any](f func(A) Reader[R, B]) func([]A) Reader[R, []B] {
	return G.TraverseArray[Reader[R, B], Reader[R, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[R, A any](ma []Reader[R, A]) Reader[R, []A] {
	return G.SequenceArray[Reader[R, A], Reader[R, []A]](ma)
}
