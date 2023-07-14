package readerioeither

import (
	G "github.com/ibm/fp-go/readerioeither/generic"
)

// these functions From a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[R, A any](f func(R) func() (A, error)) ReaderIOEither[R, error, A] {
	return G.From0[ReaderIOEither[R, error, A]](f)
}

func From1[R, T1, A any](f func(R, T1) func() (A, error)) func(T1) ReaderIOEither[R, error, A] {
	return G.From1[ReaderIOEither[R, error, A]](f)
}

func From2[R, T1, T2, A any](f func(R, T1, T2) func() (A, error)) func(T1, T2) ReaderIOEither[R, error, A] {
	return G.From2[ReaderIOEither[R, error, A]](f)
}

func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) func() (A, error)) func(T1, T2, T3) ReaderIOEither[R, error, A] {
	return G.From3[ReaderIOEither[R, error, A]](f)
}
