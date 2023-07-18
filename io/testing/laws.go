package testing

import (
	"testing"

	EQ "github.com/IBM/fp-go/eq"
	L "github.com/IBM/fp-go/internal/monad/testing"
	"github.com/IBM/fp-go/io"
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
		io.Eq(eqa),
		io.Eq(eqb),
		io.Eq(eqc),

		io.Of[A],
		io.Of[B],
		io.Of[C],

		io.Of[func(A) A],
		io.Of[func(A) B],
		io.Of[func(B) C],
		io.Of[func(func(A) B) B],

		io.MonadMap[A, A],
		io.MonadMap[A, B],
		io.MonadMap[A, C],
		io.MonadMap[B, C],

		io.MonadMap[func(B) C, func(func(A) B) func(A) C],

		io.MonadChain[A, A],
		io.MonadChain[A, B],
		io.MonadChain[A, C],
		io.MonadChain[B, C],

		io.MonadAp[A, A],
		io.MonadAp[B, A],
		io.MonadAp[C, B],
		io.MonadAp[C, A],

		io.MonadAp[B, func(A) B],
		io.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
