package readerioresult

import (
	"github.com/IBM/fp-go/v2/readerioeither"
)

//go:inline
func ChainConsumer[R, A any](c Consumer[A]) Operator[R, A, struct{}] {
	return readerioeither.ChainConsumer[R, error](c)
}

//go:inline
func ChainFirstConsumer[R, A any](c Consumer[A]) Operator[R, A, A] {
	return readerioeither.ChainFirstConsumer[R, error](c)
}
