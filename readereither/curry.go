package readereither

import (
	G "github.com/ibm/fp-go/readereither/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[R, A any](f func(R) (A, error)) ReaderEither[R, error, A] {
	return G.Curry0[ReaderEither[R, error, A]](f)
}

func Curry1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderEither[R, error, A] {
	return G.Curry1[ReaderEither[R, error, A]](f)
}

func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1) func(T2) ReaderEither[R, error, A] {
	return G.Curry2[ReaderEither[R, error, A]](f)
}

func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderEither[R, error, A] {
	return G.Curry3[ReaderEither[R, error, A]](f)
}

func Uncurry1[R, T1, A any](f func(T1) ReaderEither[R, error, A]) func(R, T1) (A, error) {
	return G.Uncurry1(f)
}

func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderEither[R, error, A]) func(R, T1, T2) (A, error) {
	return G.Uncurry2(f)
}

func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderEither[R, error, A]) func(R, T1, T2, T3) (A, error) {
	return G.Uncurry3(f)
}
