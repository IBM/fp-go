package readerioeither

import (
	G "github.com/ibm/fp-go/context/readerioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return G.TraverseArray[[]A, ReaderIOEither[[]B]](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderIOEither[A]) ReaderIOEither[[]A] {
	return G.SequenceArray[[]A, []ReaderIOEither[A], ReaderIOEither[[]A]](ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable, A, B any](f func(A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return G.TraverseRecord[K, map[K]A, ReaderIOEither[map[K]B]](f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]ReaderIOEither[A]) ReaderIOEither[map[K]A] {
	return G.SequenceRecord[K, map[K]A, map[K]ReaderIOEither[A], ReaderIOEither[map[K]A]](ma)
}
