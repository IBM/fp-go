package readerio

import (
	IO "github.com/ibm/fp-go/io"
	G "github.com/ibm/fp-go/readerio/generic"
)

// these functions From a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[R, A any](f func(R) IO.IO[A]) ReaderIO[R, A] {
	return G.From0[ReaderIO[R, A]](f)
}

func From1[R, T1, A any](f func(R, T1) IO.IO[A]) func(T1) ReaderIO[R, A] {
	return G.From1[ReaderIO[R, A]](f)
}

func From2[R, T1, T2, A any](f func(R, T1, T2) IO.IO[A]) func(T1, T2) ReaderIO[R, A] {
	return G.From2[ReaderIO[R, A]](f)
}

func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) IO.IO[A]) func(T1, T2, T3) ReaderIO[R, A] {
	return G.From3[ReaderIO[R, A]](f)
}
