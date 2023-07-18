package readerioeither

import (
	G "github.com/ibm/fp-go/context/readerioeither/generic"
)

// WithContext wraps an existing ReaderIOEither and performs a context check for cancellation before delegating
func WithContext[A any](ma ReaderIOEither[A]) ReaderIOEither[A] {
	return G.WithContext(ma)
}
