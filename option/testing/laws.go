package testing

import (
	"testing"

	EQ "github.com/ibm/fp-go/eq"
	L "github.com/ibm/fp-go/internal/monad/testing"
	O "github.com/ibm/fp-go/option"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[A, B, C any](t *testing.T,
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		O.Eq(eqa),
		O.Eq(eqb),
		O.Eq(eqc),

		O.Of[A],
		O.Of[B],
		O.Of[C],

		O.Of[func(A) A],
		O.Of[func(A) B],
		O.Of[func(B) C],
		O.Of[func(func(A) B) B],

		O.MonadMap[A, A],
		O.MonadMap[A, B],
		O.MonadMap[A, C],
		O.MonadMap[B, C],

		O.MonadMap[func(B) C, func(func(A) B) func(A) C],

		O.MonadChain[A, A],
		O.MonadChain[A, B],
		O.MonadChain[A, C],
		O.MonadChain[B, C],

		O.MonadAp[A, A],
		O.MonadAp[B, A],
		O.MonadAp[C, B],
		O.MonadAp[C, A],

		O.MonadAp[B, func(A) B],
		O.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
