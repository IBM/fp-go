package readereither

import (
	G "github.com/IBM/fp-go/readereither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[R, A any](f func(R) (A, error)) func() ReaderEither[R, error, A] {
	return G.From0[ReaderEither[R, error, A]](f)
}

func From1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderEither[R, error, A] {
	return G.From1[ReaderEither[R, error, A]](f)
}

func From2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) ReaderEither[R, error, A] {
	return G.From2[ReaderEither[R, error, A]](f)
}

func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderEither[R, error, A] {
	return G.From3[ReaderEither[R, error, A]](f)
}
