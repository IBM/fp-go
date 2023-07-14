package readerioeither

import (
	"context"

	RE "github.com/ibm/fp-go/readerioeither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Eitherize0[A any](f func(context.Context) (A, error)) ReaderIOEither[A] {
	return RE.Eitherize0[ReaderIOEither[A]](f)
}

func Eitherize1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderIOEither[A] {
	return RE.Eitherize1[ReaderIOEither[A]](f)
}

func Eitherize2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1, T2) ReaderIOEither[A] {
	return RE.Eitherize2[ReaderIOEither[A]](f)
}

func Eitherize3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderIOEither[A] {
	return RE.Eitherize3[ReaderIOEither[A]](f)
}

func Eitherize4[T1, T2, T3, T4, A any](f func(context.Context, T1, T2, T3, T4) (A, error)) func(T1, T2, T3, T4) ReaderIOEither[A] {
	return RE.Eitherize4[ReaderIOEither[A]](f)
}
