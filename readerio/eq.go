package readerio

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/readerio/generic"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[R, A any](e EQ.Eq[A]) func(r R) EQ.Eq[ReaderIO[R, A]] {
	return G.Eq[ReaderIO[R, A]](e)
}
