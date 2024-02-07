// Copyright (c) 2023 IBM Corp.
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

	E "github.com/IBM/fp-go/eq"
	LA "github.com/IBM/fp-go/internal/applicative/testing"
	LC "github.com/IBM/fp-go/internal/chain/testing"
	"github.com/stretchr/testify/assert"
)

// Apply monad left identity law
//
// M.chain(M.of(a), f) <-> f(a)
func AssertLeftIdentity[HKTA, HKTB, A, B any](t *testing.T,
	eq E.Eq[HKTB],

	fofa func(A) HKTA,
	fofb func(B) HKTB,

	fchain func(HKTA, func(A) HKTB) HKTB,

	ab func(A) B,
) func(a A) bool {
	return func(a A) bool {

		f := func(a A) HKTB {
			return fofb(ab(a))
		}

		left := fchain(fofa(a), f)
		right := f(a)

		return assert.True(t, eq.Equals(left, right), "Monad left identity")
	}
}

// Apply monad right identity law
//
// M.chain(fa, M.of) <-> fa
func AssertRightIdentity[HKTA, A any](t *testing.T,
	eq E.Eq[HKTA],

	fofa func(A) HKTA,

	fchain func(HKTA, func(A) HKTA) HKTA,
) func(fa HKTA) bool {
	return func(fa HKTA) bool {

		left := fchain(fa, fofa)
		right := fa

		return assert.True(t, eq.Equals(left, right), "Monad right identity")
	}
}

// AssertLaws asserts the apply laws `identity`, `composition`, `associative composition`, 'applicative identity', 'homomorphism', 'interchange', `associativity`, `left identity`, `right identity`
func AssertLaws[HKTA, HKTB, HKTC, HKTAA, HKTAB, HKTBC, HKTAC, HKTABB, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqb E.Eq[HKTB],
	eqc E.Eq[HKTC],

	fofa func(A) HKTA,
	fofb func(B) HKTB,
	fofc func(C) HKTC,

	fofaa func(func(A) A) HKTAA,
	fofab func(func(A) B) HKTAB,
	fofbc func(func(B) C) HKTBC,
	fofabb func(func(func(A) B) B) HKTABB,

	faa func(HKTA, func(A) A) HKTA,
	fab func(HKTA, func(A) B) HKTB,
	fac func(HKTA, func(A) C) HKTC,
	fbc func(HKTB, func(B) C) HKTC,

	fmap func(HKTBC, func(func(B) C) func(func(A) B) func(A) C) HKTABAC,

	chainaa func(HKTA, func(A) HKTA) HKTA,
	chainab func(HKTA, func(A) HKTB) HKTB,
	chainac func(HKTA, func(A) HKTC) HKTC,
	chainbc func(HKTB, func(B) HKTC) HKTC,

	fapaa func(HKTAA, HKTA) HKTA,
	fapab func(HKTAB, HKTA) HKTB,
	fapbc func(HKTBC, HKTB) HKTC,
	fapac func(HKTAC, HKTA) HKTC,

	fapabb func(HKTABB, HKTAB) HKTB,
	fapabac func(HKTABAC, HKTAB) HKTAC,

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {
	// applicative laws
	applicative := LA.AssertLaws(t, eqa, eqb, eqc, fofa, fofb, fofaa, fofab, fofbc, fofabb, faa, fab, fac, fbc, fmap, fapaa, fapab, fapbc, fapac, fapabb, fapabac, ab, bc)
	// chain laws
	chain := LC.AssertLaws(t, eqa, eqc, fofb, fofc, fofab, fofbc, faa, fab, fac, fbc, fmap, chainab, chainac, chainbc, fapab, fapbc, fapac, fapabac, ab, bc)
	// monad laws
	leftIdentity := AssertLeftIdentity(t, eqb, fofa, fofb, chainab, ab)
	rightIdentity := AssertRightIdentity(t, eqa, fofa, chainaa)

	return func(a A) bool {
		fa := fofa(a)
		return applicative(a) && chain(fa) && leftIdentity(a) && rightIdentity(fa)
	}
}
