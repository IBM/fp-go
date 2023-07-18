package readerio

import (
	"context"

	IO "github.com/IBM/fp-go/io"
	R "github.com/IBM/fp-go/readerio/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[A any](f func(context.Context) IO.IO[A]) func() ReaderIO[A] {
	return R.From0[ReaderIO[A]](f)
}

func From1[T1, A any](f func(context.Context, T1) IO.IO[A]) func(T1) ReaderIO[A] {
	return R.From1[ReaderIO[A]](f)
}

func From2[T1, T2, A any](f func(context.Context, T1, T2) IO.IO[A]) func(T1, T2) ReaderIO[A] {
	return R.From2[ReaderIO[A]](f)
}

func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) IO.IO[A]) func(T1, T2, T3) ReaderIO[A] {
	return R.From3[ReaderIO[A]](f)
}
