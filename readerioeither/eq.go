package readerioeither

import (
	ET "github.com/ibm/fp-go/either"
	EQ "github.com/ibm/fp-go/eq"
	G "github.com/ibm/fp-go/readerioeither/generic"
)

// Eq implements the equals predicate for values contained in the IOEither monad
func Eq[R, E, A any](eq EQ.Eq[ET.Either[E, A]]) func(R) EQ.Eq[ReaderIOEither[R, E, A]] {
	return G.Eq[ReaderIOEither[R, E, A]](eq)
}

// FromStrictEquals constructs an `Eq` from the canonical comparison function
func FromStrictEquals[R, E, A comparable]() func(R) EQ.Eq[ReaderIOEither[R, E, A]] {
	return G.FromStrictEquals[ReaderIOEither[R, E, A]]()
}
