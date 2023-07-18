package readerio

import (
	G "github.com/ibm/fp-go/readerio/generic"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[R, A, B any](first ReaderIO[R, A], second ReaderIO[R, B]) ReaderIO[R, A] {
	return G.MonadApFirst[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) A]](first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[R, A, B any](second ReaderIO[R, B]) func(ReaderIO[R, A]) ReaderIO[R, A] {
	return G.ApFirst[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) A]](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[R, A, B any](first ReaderIO[R, A], second ReaderIO[R, B]) ReaderIO[R, B] {
	return G.MonadApSecond[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) B]](first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[R, A, B any](second ReaderIO[R, B]) func(ReaderIO[R, A]) ReaderIO[R, B] {
	return G.ApSecond[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) B]](second)
}
