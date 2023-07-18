package testing

import (
	"testing"

	ET "github.com/ibm/fp-go/either"
	EQ "github.com/ibm/fp-go/eq"
	L "github.com/ibm/fp-go/internal/monad/testing"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[E, A, B, C any](t *testing.T,
	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		ET.Eq(eqe, eqa),
		ET.Eq(eqe, eqb),
		ET.Eq(eqe, eqc),

		ET.Of[E, A],
		ET.Of[E, B],
		ET.Of[E, C],

		ET.Of[E, func(A) A],
		ET.Of[E, func(A) B],
		ET.Of[E, func(B) C],
		ET.Of[E, func(func(A) B) B],

		ET.MonadMap[E, A, A],
		ET.MonadMap[E, A, B],
		ET.MonadMap[E, A, C],
		ET.MonadMap[E, B, C],

		ET.MonadMap[E, func(B) C, func(func(A) B) func(A) C],

		ET.MonadChain[E, A, A],
		ET.MonadChain[E, A, B],
		ET.MonadChain[E, A, C],
		ET.MonadChain[E, B, C],

		ET.MonadAp[A, E, A],
		ET.MonadAp[B, E, A],
		ET.MonadAp[C, E, B],
		ET.MonadAp[C, E, A],

		ET.MonadAp[B, E, func(A) B],
		ET.MonadAp[func(A) C, E, func(A) B],

		ab,
		bc,
	)

}
