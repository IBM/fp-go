package either

import (
	RI "github.com/IBM/fp-go/v2/internal/iter"
)

func TraverseIterG[GA ~func(yield func(A) bool), GB ~func(yield func(B) bool), E, A, B any](f Kleisli[E, A, B]) Kleisli[E, GA, GB] {
	return RI.Traverse[GA](
		Map[E, B],
		Of[E, GB],
		Map[E, GB],
		Ap[GB, E],
		f,
	)
}

func TraversableIter[E, A, B any]() Traversable[E, A, B, Seq[A], Seq[B]] {
	return TraverseIterG[Seq[A], Seq[B], E, A, B]
}
