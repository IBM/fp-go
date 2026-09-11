package readerreaderioresult

import (
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
)

// TraverseArray transforms an array by applying a function that returns a ReaderReaderIOResult
// to each element, collecting the results. The first failure short-circuits the traversal.
//
// For an empty input array the resulting array is nil, the canonical empty array.
func TraverseArray[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, []A, []B] {
	return RRIOE.TraverseArray(f)
}
