package generic

import (
	"context"

	E "github.com/IBM/fp-go/either"
	ET "github.com/IBM/fp-go/either"
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/readerioeither/generic"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[GRA ~func(context.Context) GIOA, GIOA ~func() E.Either[error, A], A any](eq EQ.Eq[ET.Either[error, A]]) func(context.Context) EQ.Eq[GRA] {
	return G.Eq[GRA](eq)
}
