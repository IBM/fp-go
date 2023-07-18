package generic

import (
	"context"

	CIOE "github.com/IBM/fp-go/context/ioeither/generic"
	E "github.com/IBM/fp-go/either"
	IOE "github.com/IBM/fp-go/ioeither/generic"
)

// withContext wraps an existing ReaderIOEither and performs a context check for cancellation before delegating
func WithContext[GRA ~func(context.Context) GIOA, GIOA ~func() E.Either[error, A], A any](ma GRA) GRA {
	return func(ctx context.Context) GIOA {
		if err := context.Cause(ctx); err != nil {
			return IOE.Left[GIOA](err)
		}
		return CIOE.WithContext(ctx, ma(ctx))
	}
}
