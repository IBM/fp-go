package decode

import (
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
)

type (

	// Validation represents the result of a validation operation that may contain
	// validation errors or a successfully validated value of type A.
	Validation[A any] = validation.Validation[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// Decode is a function that decodes input I to type A with validation.
	// It returns a Validation result directly.
	Decode[I, A any] = Reader[I, Validation[A]]

	// Kleisli represents a function from A to a decoded B given input type I.
	// It's a Reader that takes an input A and produces a Decode[I, B] function.
	// This enables composition of decoding operations in a functional style.
	Kleisli[I, A, B any] = Reader[A, Decode[I, B]]

	// Operator represents a decoding transformation that takes a decoded A and produces a decoded B.
	// It's a specialized Kleisli arrow for composing decode operations where the input is already decoded.
	// This allows chaining multiple decode transformations together.
	Operator[I, A, B any] = Kleisli[I, Decode[I, A], B]
)
