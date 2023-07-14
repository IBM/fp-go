package ioeither

import (
	"context"

	ET "github.com/ibm/fp-go/either"
	IOE "github.com/ibm/fp-go/ioeither"
)

// withContext wraps an existing IOEither and performs a context check for cancellation before delegating
func WithContext[A any](ctx context.Context, ma IOE.IOEither[error, A]) IOE.IOEither[error, A] {
	return IOE.MakeIO(func() ET.Either[error, A] {
		if err := context.Cause(ctx); err != nil {
			return ET.Left[A](err)
		}
		return ma()
	})
}
