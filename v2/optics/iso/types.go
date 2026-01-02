package iso

import (
	"github.com/IBM/fp-go/v2/array/nonempty"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
)

type (
	// Number represents a numeric type constraint.
	Number = number.Number

	// Pair represents a tuple of two values of types A and B.
	Pair[A, B any] = pair.Pair[A, B]

	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// NonEmptyArray represents an array that is guaranteed to have at least one element.
	NonEmptyArray[A any] = nonempty.NonEmptyArray[A]
)
