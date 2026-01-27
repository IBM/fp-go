package readerreaderioresult

import (
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
)

func TraverseArray[R, A, B any](f Kleisli[R, A, B]) Kleisli[R, []A, []B] {
	return RRIOE.TraverseArray(f)
}
