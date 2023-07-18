package readerioeither

import (
	G "github.com/IBM/fp-go/context/readerioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[R, A, ANY any](onCreate ReaderIOEither[R], onRelease func(R) ReaderIOEither[ANY]) func(func(R) ReaderIOEither[A]) ReaderIOEither[A] {
	// wraps the callback functions with a context check
	return G.WithResource[ReaderIOEither[A]](onCreate, onRelease)
}
