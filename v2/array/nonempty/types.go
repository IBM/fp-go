package nonempty

import "github.com/IBM/fp-go/v2/option"

type (
	// NonEmptyArray represents an array that is guaranteed to have at least one element.
	// This provides compile-time safety for operations that require non-empty collections.
	NonEmptyArray[A any] []A

	// Kleisli represents a Kleisli arrow for the NonEmptyArray monad.
	// It's a function from A to NonEmptyArray[B], used for composing operations that produce non-empty arrays.
	Kleisli[A, B any] = func(A) NonEmptyArray[B]

	// Operator represents a function that transforms one NonEmptyArray into another.
	// It takes a NonEmptyArray[A] and produces a NonEmptyArray[B].
	Operator[A, B any] = Kleisli[NonEmptyArray[A], B]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]
)
