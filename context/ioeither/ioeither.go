package ioeither

import (
	"context"

	G "github.com/IBM/fp-go/context/ioeither/generic"
	IOE "github.com/IBM/fp-go/ioeither"
)

// withContext wraps an existing IOEither and performs a context check for cancellation before delegating
func WithContext[A any](ctx context.Context, ma IOE.IOEither[error, A]) IOE.IOEither[error, A] {
	return G.WithContext(ctx, ma)
}
