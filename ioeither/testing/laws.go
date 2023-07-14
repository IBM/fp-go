package testing

import (
	"testing"

	ET "github.com/ibm/fp-go/either"
	EQ "github.com/ibm/fp-go/eq"
	L "github.com/ibm/fp-go/internal/monad/testing"
	IOE "github.com/ibm/fp-go/ioeither"
)

// AssertLaws asserts the apply monad laws for the `IOEither` monad
func AssertLaws[E, A, B, C any](t *testing.T,
	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		IOE.Eq(ET.Eq(eqe, eqa)),
		IOE.Eq(ET.Eq(eqe, eqb)),
		IOE.Eq(ET.Eq(eqe, eqc)),

		IOE.Of[E, A],
		IOE.Of[E, B],
		IOE.Of[E, C],

		IOE.Of[E, func(A) A],
		IOE.Of[E, func(A) B],
		IOE.Of[E, func(B) C],
		IOE.Of[E, func(func(A) B) B],

		IOE.MonadMap[E, A, A],
		IOE.MonadMap[E, A, B],
		IOE.MonadMap[E, A, C],
		IOE.MonadMap[E, B, C],

		IOE.MonadMap[E, func(B) C, func(func(A) B) func(A) C],

		IOE.MonadChain[E, A, A],
		IOE.MonadChain[E, A, B],
		IOE.MonadChain[E, A, C],
		IOE.MonadChain[E, B, C],

		IOE.MonadAp[A, E, A],
		IOE.MonadAp[B, E, A],
		IOE.MonadAp[C, E, B],
		IOE.MonadAp[C, E, A],

		IOE.MonadAp[B, E, func(A) B],
		IOE.MonadAp[func(A) C, E, func(A) B],

		ab,
		bc,
	)

}
