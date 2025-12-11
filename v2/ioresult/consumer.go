package ioresult

import (
	"github.com/IBM/fp-go/v2/ioeither"
)

//go:inline
func ChainConsumer[A any](c Consumer[A]) Operator[A, struct{}] {
	return ioeither.ChainConsumer[error](c)
}

//go:inline
func ChainFirstConsumer[A any](c Consumer[A]) Operator[A, A] {
	return ioeither.ChainFirstConsumer[error](c)
}
