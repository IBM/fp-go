package io

import (
	G "github.com/IBM/fp-go/io/generic"
)

func MonadTraverseArray[A, B any](tas []A, f func(A) IO[B]) IO[[]B] {
	return G.MonadTraverseArray[IO[B], IO[[]B]](tas, f)
}

// TraverseArray applies a function returning an [IO] to all elements in an array and the
// transforms this into an [IO] of that array
func TraverseArray[A, B any](f func(A) IO[B]) func([]A) IO[[]B] {
	return G.TraverseArray[IO[B], IO[[]B], []A](f)
}

// SequenceArray converts an array of [IO] to an [IO] of an array
func SequenceArray[A any](tas []IO[A]) IO[[]A] {
	return G.SequenceArray[IO[A], IO[[]A]](tas)
}

func MonadTraverseRecord[K comparable, A, B any](tas map[K]A, f func(A) IO[B]) IO[map[K]B] {
	return G.MonadTraverseRecord[IO[B], IO[map[K]B]](tas, f)
}

// TraverseArray applies a function returning an [IO] to all elements in a record and the
// transforms this into an [IO] of that record
func TraverseRecord[K comparable, A, B any](f func(A) IO[B]) func(map[K]A) IO[map[K]B] {
	return G.TraverseRecord[IO[B], IO[map[K]B], map[K]A](f)
}

// SequenceRecord converts a record of [IO] to an [IO] of a record
func SequenceRecord[K comparable, A any](tas map[K]IO[A]) IO[map[K]A] {
	return G.SequenceRecord[IO[A], IO[map[K]A]](tas)
}
