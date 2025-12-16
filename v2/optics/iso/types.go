package iso

import (
	"github.com/IBM/fp-go/v2/array/nonempty"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
)

type (
	Number               = number.Number
	Pair[A, B any]       = pair.Pair[A, B]
	Either[E, A any]     = either.Either[E, A]
	NonEmptyArray[A any] = nonempty.NonEmptyArray[A]
)
