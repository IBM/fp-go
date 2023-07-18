package readereither

import (
	"context"

	E "github.com/IBM/fp-go/either"
)

// withContext wraps an existing ReaderEither and performs a context check for cancellation before deletating
func WithContext[A any](ma ReaderEither[A]) ReaderEither[A] {
	return func(ctx context.Context) E.Either[error, A] {
		if err := context.Cause(ctx); err != nil {
			return E.Left[A](err)
		}
		return ma(ctx)
	}
}
