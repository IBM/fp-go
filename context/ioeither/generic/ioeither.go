package generic

import (
	"context"

	E "github.com/IBM/fp-go/either"
	ET "github.com/IBM/fp-go/either"
	IOE "github.com/IBM/fp-go/ioeither/generic"
)

// withContext wraps an existing IOEither and performs a context check for cancellation before delegating
func WithContext[GIO ~func() E.Either[error, A], A any](ctx context.Context, ma GIO) GIO {
	return IOE.MakeIO[GIO](func() E.Either[error, A] {
		if err := context.Cause(ctx); err != nil {
			return ET.Left[A](err)
		}
		return ma()
	})
}
