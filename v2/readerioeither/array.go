package readerioeither

import (
	"github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/monoid"
)

//go:inline
func MonadReduceArray[R, E, A, B any](as []ReaderIOEither[R, E, A], reduce func(B, A) B, initial B) ReaderIOEither[R, E, B] {
	return RA.MonadTraverseReduce(
		Of,
		Map,
		Ap,

		as,

		function.Identity[ReaderIOEither[R, E, A]],
		reduce,
		initial,
	)
}

//go:inline
func ReduceArray[R, E, A, B any](reduce func(B, A) B, initial B) Kleisli[R, E, []ReaderIOEither[R, E, A], B] {
	return RA.TraverseReduce[[]ReaderIOEither[R, E, A]](
		Of,
		Map,
		Ap,

		function.Identity[ReaderIOEither[R, E, A]],
		reduce,
		initial,
	)
}

//go:inline
func MonadReduceArrayM[R, E, A any](as []ReaderIOEither[R, E, A], m monoid.Monoid[A]) ReaderIOEither[R, E, A] {
	return MonadReduceArray(as, m.Concat, m.Empty())
}

//go:inline
func ReduceArrayM[R, E, A any](m monoid.Monoid[A]) Kleisli[R, E, []ReaderIOEither[R, E, A], A] {
	return ReduceArray[R, E](m.Concat, m.Empty())
}

//go:inline
func MonadTraverseReduceArray[R, E, A, B, C any](as []A, trfrm Kleisli[R, E, A, B], reduce func(C, B) C, initial C) ReaderIOEither[R, E, C] {
	return RA.MonadTraverseReduce(
		Of,
		Map,
		Ap,

		as,

		trfrm,
		reduce,
		initial,
	)
}

//go:inline
func TraverseReduceArray[R, E, A, B, C any](trfrm Kleisli[R, E, A, B], reduce func(C, B) C, initial C) Kleisli[R, E, []A, C] {
	return RA.TraverseReduce[[]A](
		Of,
		Map,
		Ap,

		trfrm,
		reduce,
		initial,
	)
}

//go:inline
func MonadTraverseReduceArrayM[R, E, A, B any](as []A, trfrm Kleisli[R, E, A, B], m monoid.Monoid[B]) ReaderIOEither[R, E, B] {
	return MonadTraverseReduceArray(as, trfrm, m.Concat, m.Empty())
}

//go:inline
func TraverseReduceArrayM[R, E, A, B any](trfrm Kleisli[R, E, A, B], m monoid.Monoid[B]) Kleisli[R, E, []A, B] {
	return TraverseReduceArray(trfrm, m.Concat, m.Empty())
}
