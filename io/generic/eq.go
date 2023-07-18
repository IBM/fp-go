package generic

import (
	EQ "github.com/IBM/fp-go/eq"
	G "github.com/IBM/fp-go/internal/eq"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[GA ~func() A, A any](e EQ.Eq[A]) EQ.Eq[GA] {
	// comparator for the monad
	eq := G.Eq(
		MonadMap[GA, func() func(A) bool, A, func(A) bool],
		MonadAp[GA, func() bool, func() func(A) bool, A, bool],
		e,
	)
	// eagerly execute
	return EQ.FromEquals(func(l, r GA) bool {
		return eq(l, r)()
	})
}
