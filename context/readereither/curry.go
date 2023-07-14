package readereither

import (
	"context"

	RE "github.com/ibm/fp-go/readereither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[A any](f func(context.Context) (A, error)) ReaderEither[A] {
	return RE.Curry0[ReaderEither[A]](f)
}

func Curry1[T1, A any](f func(context.Context, T1) (A, error)) func(T1) ReaderEither[A] {
	return RE.Curry1[ReaderEither[A]](f)
}

func Curry2[T1, T2, A any](f func(context.Context, T1, T2) (A, error)) func(T1) func(T2) ReaderEither[A] {
	return RE.Curry2[ReaderEither[A]](f)
}

func Curry3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderEither[A] {
	return RE.Curry3[ReaderEither[A]](f)
}

func Uncurry1[T1, A any](f func(T1) ReaderEither[A]) func(context.Context, T1) (A, error) {
	return RE.Uncurry1(f)
}

func Uncurry2[T1, T2, A any](f func(T1) func(T2) ReaderEither[A]) func(context.Context, T1, T2) (A, error) {
	return RE.Uncurry2(f)
}

func Uncurry3[T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderEither[A]) func(context.Context, T1, T2, T3) (A, error) {
	return RE.Uncurry3(f)
}
