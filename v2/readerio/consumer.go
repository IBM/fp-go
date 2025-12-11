package readerio

import "github.com/IBM/fp-go/v2/io"

//go:inline
func ChainConsumer[R, A any](c Consumer[A]) Operator[R, A, struct{}] {
	return ChainIOK[R](io.FromConsumerK(c))
}

//go:inline
func ChainFirstConsumer[R, A any](c Consumer[A]) Operator[R, A, A] {
	return ChainFirstIOK[R](io.FromConsumerK(c))
}
