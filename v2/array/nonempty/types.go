package nonempty

import "github.com/IBM/fp-go/v2/option"

type (

	// NonEmptyArray represents an array with at least one element
	NonEmptyArray[A any] []A

	Kleisli[A, B any] = func(A) NonEmptyArray[B]

	Operator[A, B any] = Kleisli[NonEmptyArray[A], B]

	Option[A any] = option.Option[A]
)
