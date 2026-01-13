package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/function"
)

// Promap is the profunctor map operation that transforms both the input and output of a ReaderResult.
// It applies f to the input context (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Adapt the context before passing it to the ReaderResult (via f)
//   - Transform the success value after the computation completes (via g)
//
// The error type is fixed as error and remains unchanged through the transformation.
//
// Type Parameters:
//   - A: The original success type produced by the ReaderResult
//   - B: The new output success type
//
// Parameters:
//   - f: Function to transform the input context, returning a new context and cancel function (contravariant)
//   - g: Function to transform the output success value from A to B (covariant)
//
// Returns:
//   - An Operator that takes a ReaderResult[A] and returns a ReaderResult[B]
//
//go:inline
func Promap[A, B any](f func(context.Context) (context.Context, context.CancelFunc), g func(A) B) Operator[A, B] {
	return function.Flow2(
		Local[A](f),
		Map(g),
	)
}

// Contramap changes the value of the local context during the execution of a ReaderResult.
// This is the contravariant functor operation that transforms the input context.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is useful for adapting a ReaderResult to work with a modified context
// by providing a function that creates a new context (and optional cancel function) from the current one.
//
// Type Parameters:
//   - A: The success type (unchanged)
//
// Parameters:
//   - f: Function to transform the context, returning a new context and cancel function
//
// Returns:
//   - A Kleisli arrow that takes a ReaderResult[A] and returns a ReaderResult[A]
//
//go:inline
func Contramap[A any](f func(context.Context) (context.Context, context.CancelFunc)) Kleisli[ReaderResult[A], A] {
	return Local[A](f)
}
