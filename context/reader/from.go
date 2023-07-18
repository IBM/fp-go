package reader

import (
	"context"

	R "github.com/IBM/fp-go/reader/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[A any](f func(context.Context) A) func() Reader[A] {
	return R.From0[Reader[A]](f)
}

func From1[T1, A any](f func(context.Context, T1) A) func(T1) Reader[A] {
	return R.From1[Reader[A]](f)
}

func From2[T1, T2, A any](f func(context.Context, T1, T2) A) func(T1, T2) Reader[A] {
	return R.From2[Reader[A]](f)
}

func From3[T1, T2, T3, A any](f func(context.Context, T1, T2, T3) A) func(T1, T2, T3) Reader[A] {
	return R.From3[Reader[A]](f)
}
