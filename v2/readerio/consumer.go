package readerio

import "github.com/IBM/fp-go/v2/io"

//go:inline
func ChainConsumer[R, A any](c Consumer[A]) Operator[R, A, Void] {
	return ChainIOK[R](io.FromConsumer(c))
}

//go:inline
func ChainFirstConsumer[R, A any](c Consumer[A]) Operator[R, A, A] {
	return ChainFirstIOK[R](io.FromConsumer(c))
}
