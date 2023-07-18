package readerioeither

import (
	"context"

	G "github.com/ibm/fp-go/context/readerioeither/generic"
	ET "github.com/ibm/fp-go/either"
	EQ "github.com/ibm/fp-go/eq"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[A any](eq EQ.Eq[ET.Either[error, A]]) func(context.Context) EQ.Eq[ReaderIOEither[A]] {
	return G.Eq[ReaderIOEither[A]](eq)
}
