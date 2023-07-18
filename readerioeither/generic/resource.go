package generic

import (
	ET "github.com/IBM/fp-go/either"
	IOE "github.com/IBM/fp-go/ioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GEA ~func(L) TEA,
	GER ~func(L) TER,
	GEANY ~func(L) TEANY,

	TEA ~func() ET.Either[E, A],
	TER ~func() ET.Either[E, R],
	TEANY ~func() ET.Either[E, ANY],

	L, E, R, A, ANY any](onCreate GER, onRelease func(R) GEANY) func(func(R) GEA) GEA {

	return func(f func(R) GEA) GEA {
		return func(l L) TEA {
			// dispatch to the generic implementation
			return IOE.WithResource[TEA](
				onCreate(l),
				func(r R) TEANY {
					return onRelease(r)(l)
				},
			)(func(r R) TEA {
				return f(r)(l)
			})
		}
	}
}
