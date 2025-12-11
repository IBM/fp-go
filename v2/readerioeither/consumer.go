package readerioeither

import "github.com/IBM/fp-go/v2/io"

//go:inline
func ChainConsumer[R, E, A any](c Consumer[A]) Operator[R, E, A, struct{}] {
	return ChainIOK[R, E](io.FromConsumerK(c))
}

//go:inline
func ChainFirstConsumer[R, E, A any](c Consumer[A]) Operator[R, E, A, A] {
	return ChainFirstIOK[R, E](io.FromConsumerK(c))
}
