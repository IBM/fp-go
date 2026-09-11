package readerreaderioeither

import (
	RA "github.com/IBM/fp-go/v2/internal/array"
)

// TraverseArray transforms an array by applying a function that returns a ReaderReaderIOEither
// to each element, collecting the results. The first Left short-circuits the traversal.
//
// For an empty input array the resulting array is nil, the canonical empty array.
func TraverseArray[R, C, E, A, B any](f Kleisli[R, C, E, A, B]) Kleisli[R, C, E, []A, []B] {
	return RA.Traverse[[]A, []B](
		Of,
		Map,
		Ap,

		f,
	)
}
