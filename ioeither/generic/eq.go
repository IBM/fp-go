package generic

import (
	ET "github.com/IBM/fp-go/either"
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/io/generic"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[GA ~func() ET.Either[E, A], E, A any](eq EQ.Eq[ET.Either[E, A]]) EQ.Eq[GA] {
	return G.Eq[GA](eq)
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[GA ~func() ET.Either[E, A], E, A comparable]() EQ.Eq[GA] {
	return Eq[GA](ET.FromStrictEquals[E, A]())
}
