package generic

import (
	EQ "github.com/ibm/fp-go/eq"
	G "github.com/ibm/fp-go/internal/eq"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[GEA ~func(R) GIOA, GIOA ~func() A, R, A any](e EQ.Eq[A]) func(r R) EQ.Eq[GEA] {
	// comparator for the monad
	eq := G.Eq(
		MonadMap[GEA, func(R) func() func(A) bool, GIOA, func() func(A) bool, R, A, func(A) bool],
		MonadAp[GEA, func(R) func() bool, func(R) func() func(A) bool],
		e,
	)
	// eagerly execute
	return func(ctx R) EQ.Eq[GEA] {
		return EQ.FromEquals(func(l, r GEA) bool {
			return eq(l, r)(ctx)()
		})
	}
}
