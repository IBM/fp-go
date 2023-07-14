package readerioeither

import (
	"context"

	RE "github.com/ibm/fp-go/readerioeither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[A any](f func(context.Context) func() (A, error)) ReaderIOEither[A] {
	return RE.From0[ReaderIOEither[A]](f)
}

func From1[T1, A any](f func(context.Context, T1) func() (A, error)) func(T1) ReaderIOEither[A] {
	return RE.From1[ReaderIOEither[A]](f)
}

func From2[T1, T2, A any](f func(context.Context, T1, T2) func() (A, error)) func(T1, T2) ReaderIOEither[A] {
	return RE.From2[ReaderIOEither[A]](f)
}

func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) func() (A, error)) func(T1, T2, T3) ReaderIOEither[A] {
	return RE.From3[ReaderIOEither[A]](f)
}

func From4[T1, T2, T3, T4, A any](f func(context.Context, T1, T2, T3, T4) func() (A, error)) func(T1, T2, T3, T4) ReaderIOEither[A] {
	return RE.From4[ReaderIOEither[A]](f)
}
