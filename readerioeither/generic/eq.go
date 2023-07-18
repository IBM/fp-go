package generic

import (
	ET "github.com/ibm/fp-go/either"
	EQ "github.com/ibm/fp-go/eq"
	G "github.com/ibm/fp-go/readerio/generic"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](eq EQ.Eq[ET.Either[E, A]]) func(R) EQ.Eq[GEA] {
	return G.Eq[GEA](eq)
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R any, E, A comparable]() func(R) EQ.Eq[GEA] {
	return Eq[GEA](ET.FromStrictEquals[E, A]())
}
