// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	"testing"

	E "github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
	FCT "github.com/IBM/fp-go/v2/internal/functor/testing"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/stretchr/testify/assert"
)

// Apply associative composition law
//
// F.ap(F.ap(F.map(fbc, bc => ab => a => bc(ab(a))), fab), fa) <-> F.ap(fbc, F.ap(fab, fa))
//
// Deprecated: use [ApplyAssertAssociativeComposition] instead
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
	t.Helper()
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

// Apply associative composition law
//
// F.ap(F.ap(F.map(fbc, bc => ab => a => bc(ab(a))), fab), fa) <-> F.ap(fbc, F.ap(fab, fa))
func ApplyAssertAssociativeComposition[HKTA, HKTB, HKTC, HKTAB, HKTBC, HKTAC, HKTABAC, A, B, C any](t *testing.T,
	eq E.Eq[HKTC],

	fofab pointed.Pointed[func(A) B, HKTAB],
	fofbc pointed.Pointed[func(B) C, HKTBC],

	fmap functor.Functor[func(B) C, func(func(A) B) func(A) C, HKTBC, HKTABAC],

	fapab apply.Apply[A, B, HKTA, HKTB, HKTAB],
	fapbc apply.Apply[B, C, HKTB, HKTC, HKTBC],
	fapac apply.Apply[A, C, HKTA, HKTC, HKTAC],

	fapabac apply.Apply[func(A) B, func(A) C, HKTAB, HKTAC, HKTABAC],

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	t.Helper()
	return func(fa HKTA) bool {

		fab := fofab.Of(ab)
		fbc := fofbc.Of(bc)

		left := fapac.Ap(fa)(fapabac.Ap(fab)(fmap.Map(func(bc func(B) C) func(func(A) B) func(A) C {
			return func(ab func(A) B) func(A) C {
				return func(a A) C {
					return bc(ab(a))
				}
			}
		})(fbc)))

		right := fapbc.Ap(fapab.Ap(fa)(fab))(fbc)

		return assert.True(t, eq.Equals(left, right), "Apply associative composition")
	}
}

// AssertLaws asserts the apply laws `identity`, `composition` and `associative composition`
//
// Deprecated: use [ApplyAssertLaws] instead
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
	// mark as test helper
	t.Helper()
	// functor laws
	functor := FCT.AssertLaws(t, eqa, eqc, faa, fab, fac, fbc, ab, bc)
	// associative composition laws
	composition := AssertAssociativeComposition(t, eqc, fofab, fofbc, fmap, fapab, fapbc, fapac, fapabac, ab, bc)

	return func(fa HKTA) bool {
		return functor(fa) && composition(fa)
	}
}

// ApplyAssertLaws asserts the apply laws `identity`, `composition` and `associative composition`
func ApplyAssertLaws[HKTA, HKTB, HKTC, HKTAB, HKTBC, HKTAC, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqc E.Eq[HKTC],

	fofab pointed.Pointed[func(A) B, HKTAB],
	fofbc pointed.Pointed[func(B) C, HKTBC],

	faa functor.Functor[A, A, HKTA, HKTA],

	fmap functor.Functor[func(B) C, func(func(A) B) func(A) C, HKTBC, HKTABAC],

	fapab apply.Apply[A, B, HKTA, HKTB, HKTAB],
	fapbc apply.Apply[B, C, HKTB, HKTC, HKTBC],
	fapac apply.Apply[A, C, HKTA, HKTC, HKTAC],

	fapabac apply.Apply[func(A) B, func(A) C, HKTAB, HKTAC, HKTABAC],

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	// mark as test helper
	t.Helper()
	// functor laws
	functor := FCT.FunctorAssertLaws(t, eqa, eqc, faa, apply.ToFunctor(fapab), apply.ToFunctor(fapac), apply.ToFunctor(fapbc), ab, bc)
	// associative composition laws
	composition := ApplyAssertAssociativeComposition(t, eqc, fofab, fofbc, fmap, fapab, fapbc, fapac, fapabac, ab, bc)

	return func(fa HKTA) bool {
		return functor(fa) && composition(fa)
	}
}
