package readerioeither

import (
	"context"

	CIOE "github.com/ibm/fp-go/context/ioeither"
	IOE "github.com/ibm/fp-go/ioeither"
)

// withContext wraps an existing ReaderIOEither and performs a context check for cancellation before delegating
func WithContext[A any](ma ReaderIOEither[A]) ReaderIOEither[A] {
	return func(ctx context.Context) IOE.IOEither[error, A] {
		if err := context.Cause(ctx); err != nil {
			return IOE.Left[A](err)
		}
		return CIOE.WithContext(ctx, ma(ctx))
	}
}
