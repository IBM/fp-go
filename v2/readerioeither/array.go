package readerioeither

import (
	"github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/monoid"
)

// MonadReduceArray reduces an array of ReaderIOEither values into a single ReaderIOEither
// by applying a reduction function to accumulate the success values.
//
// If any ReaderIOEither in the array fails, the entire operation fails with that error.
// The reduction is performed sequentially from left to right.
//
// Type parameters:
//   - R: The context type
//   - E: The error type
//   - A: The element type in the array
//   - B: The accumulated result type
//
// Parameters:
//   - as: Array of ReaderIOEither values to reduce
//   - reduce: Function that combines the accumulator with each element
//   - initial: Initial value for the accumulator
//
// Returns:
//
//	A ReaderIOEither containing the final accumulated value
//
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

// ReduceArray returns a function that reduces an array of ReaderIOEither values.
// This is the curried version of MonadReduceArray.
//
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

// MonadReduceArrayM reduces an array of ReaderIOEither values using a monoid.
// The monoid provides both the combination operation and the initial (empty) value.
//
//go:inline
func MonadReduceArrayM[R, E, A any](as []ReaderIOEither[R, E, A], m monoid.Monoid[A]) ReaderIOEither[R, E, A] {
	return MonadReduceArray(as, m.Concat, m.Empty())
}

// ReduceArrayM returns a function that reduces an array using a monoid.
// This is the curried version of MonadReduceArrayM.
//
//go:inline
func ReduceArrayM[R, E, A any](m monoid.Monoid[A]) Kleisli[R, E, []ReaderIOEither[R, E, A], A] {
	return ReduceArray[R, E](m.Concat, m.Empty())
}

// MonadTraverseReduceArray transforms each element of an array using a function that returns
// a ReaderIOEither, then reduces the results into a single accumulated value.
//
// This combines traverse and reduce operations: it maps over the array with an effectful
// function and simultaneously accumulates the results.
//
// Type parameters:
//   - R: The context type
//   - E: The error type
//   - A: The input element type
//   - B: The transformed element type
//   - C: The accumulated result type
//
// Parameters:
//   - as: Array of input values
//   - trfrm: Function that transforms each element into a ReaderIOEither
//   - reduce: Function that combines the accumulator with each transformed element
//   - initial: Initial value for the accumulator
//
// Returns:
//
//	A ReaderIOEither containing the final accumulated value
//
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

// TraverseReduceArray returns a function that traverses and reduces an array.
// This is the curried version of MonadTraverseReduceArray.
//
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

// MonadTraverseReduceArrayM transforms and reduces an array using a monoid.
// The monoid provides both the combination operation and the initial value.
//
//go:inline
func MonadTraverseReduceArrayM[R, E, A, B any](as []A, trfrm Kleisli[R, E, A, B], m monoid.Monoid[B]) ReaderIOEither[R, E, B] {
	return MonadTraverseReduceArray(as, trfrm, m.Concat, m.Empty())
}

// TraverseReduceArrayM returns a function that traverses and reduces using a monoid.
// This is the curried version of MonadTraverseReduceArrayM.
//
//go:inline
func TraverseReduceArrayM[R, E, A, B any](trfrm Kleisli[R, E, A, B], m monoid.Monoid[B]) Kleisli[R, E, []A, B] {
	return TraverseReduceArray(trfrm, m.Concat, m.Empty())
}
