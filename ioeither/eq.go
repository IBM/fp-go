package ioeither

import (
	ET "github.com/ibm/fp-go/either"
	EQ "github.com/ibm/fp-go/eq"
	G "github.com/ibm/fp-go/ioeither/generic"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[E, A any](eq EQ.Eq[ET.Either[E, A]]) EQ.Eq[IOEither[E, A]] {
	return G.Eq[IOEither[E, A]](eq)
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[E, A comparable]() EQ.Eq[IOEither[E, A]] {
	return G.FromStrictEquals[IOEither[E, A]]()
}
