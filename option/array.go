package option

import (
	F "github.com/ibm/fp-go/function"
	RA "github.com/ibm/fp-go/internal/array"
)

// TraverseArray transforms an array
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f func(A) Option[B]) func(GA) Option[GB] {
	return RA.Traverse[GA](
		Of[GB],
		MonadMap[GB, func(B) GB],
		MonadAp[B, GB],

		f,
	)
}

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) Option[B]) func([]A) Option[[]B] {
	return TraverseArrayG[[]A, []B](f)
}

func SequenceArrayG[GA ~[]A, GOA ~[]Option[A], A any](ma GOA) Option[GA] {
	return TraverseArrayG[GOA, GA](F.Identity[Option[A]])(ma)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []Option[A]) Option[[]A] {
	return SequenceArrayG[[]A](ma)
}
