package testing

import (
	"testing"

	E "github.com/ibm/fp-go/eq"
	F "github.com/ibm/fp-go/function"
	L "github.com/ibm/fp-go/internal/apply/testing"
	"github.com/stretchr/testify/assert"
)

// Applicative identity law
//
// A.ap(A.of(a => a), fa) <-> fa
func AssertIdentity[HKTA, HKTAA, A any](t *testing.T,
	eq E.Eq[HKTA],

	fof func(func(A) A) HKTAA,

	fap func(HKTAA, HKTA) HKTA,
) func(fa HKTA) bool {
	return func(fa HKTA) bool {

		left := fap(fof(F.Identity[A]), fa)
		right := fa

		return assert.True(t, eq.Equals(left, right), "Applicative identity")
	}
}

// Applicative homomorphism law
//
// A.ap(A.of(ab), A.of(a)) <-> A.of(ab(a))
func AssertHomomorphism[HKTA, HKTB, HKTAB, A, B any](t *testing.T,
	eq E.Eq[HKTB],

	fofa func(A) HKTA,
	fofb func(B) HKTB,
	fofab func(func(A) B) HKTAB,

	fap func(HKTAB, HKTA) HKTB,

	ab func(A) B,
) func(a A) bool {
	return func(a A) bool {

		left := fap(fofab(ab), fofa(a))
		right := fofb(ab(a))

		return assert.True(t, eq.Equals(left, right), "Applicative homomorphism")
	}
}

// Applicative interchange law
//
// A.ap(fab, A.of(a)) <-> A.ap(A.of(ab => ab(a)), fab)
func AssertInterchange[HKTA, HKTB, HKTAB, HKTABB, A, B any](t *testing.T,
	eq E.Eq[HKTB],

	fofa func(A) HKTA,
	fofb func(B) HKTB,
	fofab func(func(A) B) HKTAB,
	fofabb func(func(func(A) B) B) HKTABB,

	fapab func(HKTAB, HKTA) HKTB,
	fapabb func(HKTABB, HKTAB) HKTB,

	ab func(A) B,
) func(a A) bool {
	return func(a A) bool {

		fab := fofab(ab)

		left := fapab(fab, fofa(a))
		right := fapabb(fofabb(func(ab func(A) B) B {
			return ab(a)
		}), fab)

		return assert.True(t, eq.Equals(left, right), "Applicative homomorphism")
	}
}

// AssertLaws asserts the apply laws `identity`, `composition`, `associative composition`, 'applicative identity', 'homomorphism', 'interchange'
func AssertLaws[HKTA, HKTB, HKTC, HKTAA, HKTAB, HKTBC, HKTAC, HKTABB, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqb E.Eq[HKTB],
	eqc E.Eq[HKTC],

	fofa func(A) HKTA,
	fofb func(B) HKTB,

	fofaa func(func(A) A) HKTAA,
	fofab func(func(A) B) HKTAB,
	fofbc func(func(B) C) HKTBC,
	fofabb func(func(func(A) B) B) HKTABB,

	faa func(HKTA, func(A) A) HKTA,
	fab func(HKTA, func(A) B) HKTB,
	fac func(HKTA, func(A) C) HKTC,
	fbc func(HKTB, func(B) C) HKTC,

	fmap func(HKTBC, func(func(B) C) func(func(A) B) func(A) C) HKTABAC,

	fapaa func(HKTAA, HKTA) HKTA,
	fapab func(HKTAB, HKTA) HKTB,
	fapbc func(HKTBC, HKTB) HKTC,
	fapac func(HKTAC, HKTA) HKTC,

	fapabb func(HKTABB, HKTAB) HKTB,
	fapabac func(HKTABAC, HKTAB) HKTAC,

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {
	// apply laws
	apply := L.AssertLaws(t, eqa, eqc, fofab, fofbc, faa, fab, fac, fbc, fmap, fapab, fapbc, fapac, fapabac, ab, bc)
	// applicative laws
	identity := AssertIdentity(t, eqa, fofaa, fapaa)
	homomorphism := AssertHomomorphism(t, eqb, fofa, fofb, fofab, fapab, ab)
	interchange := AssertInterchange(t, eqb, fofa, fofb, fofab, fofabb, fapab, fapabb, ab)

	return func(a A) bool {
		fa := fofa(a)
		return apply(fa) && identity(fa) && homomorphism(a) && interchange(a)
	}
}
