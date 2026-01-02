package array

import "github.com/IBM/fp-go/v2/option"

type (
	// Kleisli represents a Kleisli arrow for arrays.
	// It's a function from A to []B, used for composing operations that produce arrays.
	Kleisli[A, B any] = func(A) []B

	// Operator represents a function that transforms one array into another.
	// It takes a []A and produces a []B.
	Operator[A, B any] = Kleisli[[]A, B]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]
)
