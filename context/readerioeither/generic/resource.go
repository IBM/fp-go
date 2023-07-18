package generic

import (
	"context"

	E "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	RIE "github.com/ibm/fp-go/readerioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GRA ~func(context.Context) GIOA,
	GRR ~func(context.Context) GIOR,
	GRANY ~func(context.Context) GIOANY,
	GIOR ~func() E.Either[error, R],
	GIOA ~func() E.Either[error, A],
	GIOANY ~func() E.Either[error, ANY],
	R, A, ANY any](onCreate GRR, onRelease func(R) GRANY) func(func(R) GRA) GRA {
	// wraps the callback functions with a context check
	return F.Flow2(
		F.Bind2nd(F.Flow2[func(R) GRA, func(GRA) GRA, R, GRA, GRA], WithContext[GRA]),
		RIE.WithResource[GRA](WithContext(onCreate), onRelease),
	)
}
