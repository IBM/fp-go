package ioeither

import (
	G "github.com/IBM/fp-go/ioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[E, A, B any](f func(A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArray[IOEither[E, B], IOEither[E, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return G.SequenceArray[IOEither[E, A], IOEither[E, []A]](ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable, E, A, B any](f func(A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecord[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return G.SequenceRecord[IOEither[E, A], IOEither[E, map[K]A]](ma)
}
