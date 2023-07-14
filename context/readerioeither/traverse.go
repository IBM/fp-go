package readerioeither

import (
	IOE "github.com/ibm/fp-go/ioeither"
	RE "github.com/ibm/fp-go/readerioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return RE.TraverseArray[ReaderIOEither[B], ReaderIOEither[[]B], IOE.IOEither[error, B], IOE.IOEither[error, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderIOEither[A]) ReaderIOEither[[]A] {
	return RE.SequenceArray[ReaderIOEither[A], ReaderIOEither[[]A]](ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable, A, B any](f func(A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return RE.TraverseRecord[ReaderIOEither[B], ReaderIOEither[map[K]B], IOE.IOEither[error, B], IOE.IOEither[error, map[K]B], map[K]A](f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]ReaderIOEither[A]) ReaderIOEither[map[K]A] {
	return RE.SequenceRecord[ReaderIOEither[A], ReaderIOEither[map[K]A]](ma)
}
