package readerioeither

import (
	IOE "github.com/IBM/fp-go/ioeither"
	G "github.com/IBM/fp-go/readerioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func([]A) ReaderIOEither[R, E, []B] {
	return G.TraverseArray[ReaderIOEither[R, E, B], ReaderIOEither[R, E, []B], IOE.IOEither[E, B], IOE.IOEither[E, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of Readers into a Reader of a sequence
func SequenceArray[R, E, A any](ma []ReaderIOEither[R, E, A]) ReaderIOEither[R, E, []A] {
	return G.SequenceArray[ReaderIOEither[R, E, A], ReaderIOEither[R, E, []A]](ma)
}

// TraverseRecord transforms an array
func TraverseRecord[R any, K comparable, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(map[K]A) ReaderIOEither[R, E, map[K]B] {
	return G.TraverseRecord[ReaderIOEither[R, E, B], ReaderIOEither[R, E, map[K]B], IOE.IOEither[E, B], IOE.IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecord converts a homogeneous sequence of Readers into a Reader of a sequence
func SequenceRecord[R any, K comparable, E, A any](ma map[K]ReaderIOEither[R, E, A]) ReaderIOEither[R, E, map[K]A] {
	return G.SequenceRecord[ReaderIOEither[R, E, A], ReaderIOEither[R, E, map[K]A]](ma)
}
