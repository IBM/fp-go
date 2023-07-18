package readerioeither

import (
	G "github.com/IBM/fp-go/readerioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[L, E, R, A any](onCreate ReaderIOEither[L, E, R], onRelease func(R) ReaderIOEither[L, E, any]) func(func(R) ReaderIOEither[L, E, A]) ReaderIOEither[L, E, A] {
	return G.WithResource[ReaderIOEither[L, E, A]](onCreate, onRelease)
}
