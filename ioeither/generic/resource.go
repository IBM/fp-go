package generic

import (
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IF "github.com/IBM/fp-go/internal/file"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GA ~func() ET.Either[E, A],
	GR ~func() ET.Either[E, R],
	GANY ~func() ET.Either[E, ANY],
	E, R, A, ANY any](onCreate GR, onRelease func(R) GANY) func(func(R) GA) GA {
	return IF.WithResource(
		MonadChain[GR, GA, E, R, A],
		MonadFold[GA, GA, E, A, ET.Either[E, A]],
		MonadFold[GANY, GA, E, ANY, ET.Either[E, A]],
		MonadMap[GANY, GA, E, ANY, A],
		Left[GA, E, A],
	)(F.Constant(onCreate), onRelease)
}
