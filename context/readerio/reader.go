package readerio

import (
	"context"

	R "github.com/IBM/fp-go/readerio/generic"
)

func MonadMap[A, B any](fa ReaderIO[A], f func(A) B) ReaderIO[B] {
	return R.MonadMap[ReaderIO[A], ReaderIO[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(ReaderIO[A]) ReaderIO[B] {
	return R.Map[ReaderIO[A], ReaderIO[B]](f)
}

func MonadChain[A, B any](ma ReaderIO[A], f func(A) ReaderIO[B]) ReaderIO[B] {
	return R.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderIO[B]) func(ReaderIO[A]) ReaderIO[B] {
	return R.Chain[ReaderIO[A]](f)
}

func Of[A any](a A) ReaderIO[A] {
	return R.Of[ReaderIO[A]](a)
}

func MonadAp[A, B any](fab ReaderIO[func(A) B], fa ReaderIO[A]) ReaderIO[B] {
	return R.MonadAp[ReaderIO[A], ReaderIO[B]](fab, fa)
}

func Ap[A, B any](fa ReaderIO[A]) func(ReaderIO[func(A) B]) ReaderIO[B] {
	return R.Ap[ReaderIO[A], ReaderIO[B], ReaderIO[func(A) B]](fa)
}

func Ask() ReaderIO[context.Context] {
	return R.Ask[ReaderIO[context.Context]]()
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() ReaderIO[A]) ReaderIO[A] {
	return R.Defer[ReaderIO[A]](gen)
}
