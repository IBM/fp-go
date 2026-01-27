package readerreaderioeither

import (
	RA "github.com/IBM/fp-go/v2/internal/array"
)

func TraverseArray[R, C, E, A, B any](f Kleisli[R, C, E, A, B]) Kleisli[R, C, E, []A, []B] {
	return RA.Traverse[[]A, []B](
		Of,
		Map,
		Ap,

		f,
	)
}
