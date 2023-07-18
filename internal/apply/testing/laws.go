package testing

import (
	"testing"

	E "github.com/IBM/fp-go/eq"
	FCT "github.com/IBM/fp-go/internal/functor/testing"
	"github.com/stretchr/testify/assert"
)

// Apply associative composition law
//
// F.ap(F.ap(F.map(fbc, bc => ab => a => bc(ab(a))), fab), fa) <-> F.ap(fbc, F.ap(fab, fa))
func AssertAssociativeComposition[HKTA, HKTB, HKTC, HKTAB, HKTBC, HKTAC, HKTABAC, A, B, C any](t *testing.T,
	eq E.Eq[HKTC],

	fofab func(func(A) B) HKTAB,
	fofbc func(func(B) C) HKTBC,

	fmap func(HKTBC, func(func(B) C) func(func(A) B) func(A) C) HKTABAC,

	fapab func(HKTAB, HKTA) HKTB,
	fapbc func(HKTBC, HKTB) HKTC,
	fapac func(HKTAC, HKTA) HKTC,

	fapabac func(HKTABAC, HKTAB) HKTAC,

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	return func(fa HKTA) bool {

		fab := fofab(ab)
		fbc := fofbc(bc)

		left := fapac(fapabac(fmap(fbc, func(bc func(B) C) func(func(A) B) func(A) C {
			return func(ab func(A) B) func(A) C {
				return func(a A) C {
					return bc(ab(a))
				}
			}
		}), fab), fa)

		right := fapbc(fbc, fapab(fab, fa))

		return assert.True(t, eq.Equals(left, right), "Apply associative composition")
	}
}

// AssertLaws asserts the apply laws `identity`, `composition` and `associative composition`
func AssertLaws[HKTA, HKTB, HKTC, HKTAB, HKTBC, HKTAC, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqc E.Eq[HKTC],

	fofab func(func(A) B) HKTAB,
	fofbc func(func(B) C) HKTBC,

	faa func(HKTA, func(A) A) HKTA,
	fab func(HKTA, func(A) B) HKTB,
	fac func(HKTA, func(A) C) HKTC,
	fbc func(HKTB, func(B) C) HKTC,

	fmap func(HKTBC, func(func(B) C) func(func(A) B) func(A) C) HKTABAC,

	fapab func(HKTAB, HKTA) HKTB,
	fapbc func(HKTBC, HKTB) HKTC,
	fapac func(HKTAC, HKTA) HKTC,

	fapabac func(HKTABAC, HKTAB) HKTAC,

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	// functor laws
	functor := FCT.AssertLaws(t, eqa, eqc, faa, fab, fac, fbc, ab, bc)
	// associative composition laws
	composition := AssertAssociativeComposition(t, eqc, fofab, fofbc, fmap, fapab, fapbc, fapac, fapabac, ab, bc)

	return func(fa HKTA) bool {
		return functor(fa) && composition(fa)
	}
}
