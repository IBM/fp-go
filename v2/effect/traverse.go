package effect

import "github.com/IBM/fp-go/v2/context/readerreaderioresult"

func TraverseArray[C, A, B any](f Kleisli[C, A, B]) Kleisli[C, []A, []B] {
	return readerreaderioresult.TraverseArray(f)
}
