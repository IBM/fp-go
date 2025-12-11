package readerio

import "github.com/IBM/fp-go/v2/io"

//go:inline
func ChainConsumer[A any](c Consumer[A]) Operator[A, struct{}] {
	return ChainIOK(io.FromConsumerK(c))
}

//go:inline
func ChainFirstConsumer[A any](c Consumer[A]) Operator[A, A] {
	return ChainFirstIOK(io.FromConsumerK(c))
}
