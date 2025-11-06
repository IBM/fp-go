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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	L "github.com/IBM/fp-go/v2/internal/apply/testing"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/stretchr/testify/assert"
)

// Chain associativity law
//
// F.chain(F.chain(fa, afb), bfc) <-> F.chain(fa, a => F.chain(afb(a), bfc))
//
// Deprecated: use [ChainAssertAssociativity] instead
func AssertAssociativity[HKTA, HKTB, HKTC, A, B, C any](t *testing.T,
	eq E.Eq[HKTC],

	fofb func(B) HKTB,
	fofc func(C) HKTC,

	chainab func(HKTA, func(A) HKTB) HKTB,
	chainac func(HKTA, func(A) HKTC) HKTC,
	chainbc func(HKTB, func(B) HKTC) HKTC,

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	return func(fa HKTA) bool {

		afb := F.Flow2(ab, fofb)
		bfc := F.Flow2(bc, fofc)

		left := chainbc(chainab(fa, afb), bfc)

		right := chainac(fa, func(a A) HKTC {
			return chainbc(afb(a), bfc)
		})

		return assert.True(t, eq.Equals(left, right), "Chain associativity")
	}
}

// Chain associativity law
//
// F.chain(F.chain(fa, afb), bfc) <-> F.chain(fa, a => F.chain(afb(a), bfc))
func ChainAssertAssociativity[HKTA, HKTB, HKTC, HKTAB, HKTAC, HKTBC, A, B, C any](t *testing.T,
	eq E.Eq[HKTC],

	fofb pointed.Pointed[B, HKTB],
	fofc pointed.Pointed[C, HKTC],

	chainab chain.Chainable[A, B, HKTA, HKTB, HKTAB],
	chainac chain.Chainable[A, C, HKTA, HKTC, HKTAC],
	chainbc chain.Chainable[B, C, HKTB, HKTC, HKTBC],

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	return func(fa HKTA) bool {

		afb := F.Flow2(ab, fofb.Of)
		bfc := F.Flow2(bc, fofc.Of)

		left := chainbc.Chain(bfc)(chainab.Chain(afb)(fa))

		right := chainac.Chain(func(a A) HKTC {
			return chainbc.Chain(bfc)(afb(a))
		})(fa)

		return assert.True(t, eq.Equals(left, right), "Chain associativity")
	}
}

// AssertLaws asserts the apply laws `identity`, `composition`, `associative composition` and `associativity`
//
// Deprecated: use [ChainAssertLaws] instead
func AssertLaws[HKTA, HKTB, HKTC, HKTAB, HKTBC, HKTAC, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqc E.Eq[HKTC],

	fofb func(B) HKTB,
	fofc func(C) HKTC,

	fofab func(func(A) B) HKTAB,
	fofbc func(func(B) C) HKTBC,

	faa func(HKTA, func(A) A) HKTA,
	fab func(HKTA, func(A) B) HKTB,
	fac func(HKTA, func(A) C) HKTC,
	fbc func(HKTB, func(B) C) HKTC,

	fmap func(HKTBC, func(func(B) C) func(func(A) B) func(A) C) HKTABAC,

	chainab func(HKTA, func(A) HKTB) HKTB,
	chainac func(HKTA, func(A) HKTC) HKTC,
	chainbc func(HKTB, func(B) HKTC) HKTC,

	fapab func(HKTAB, HKTA) HKTB,
	fapbc func(HKTBC, HKTB) HKTC,
	fapac func(HKTAC, HKTA) HKTC,

	fapabac func(HKTABAC, HKTAB) HKTAC,

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	// apply laws
	apply := L.AssertLaws(t, eqa, eqc, fofab, fofbc, faa, fab, fac, fbc, fmap, fapab, fapbc, fapac, fapabac, ab, bc)
	// chain laws
	associativity := AssertAssociativity(t, eqc, fofb, fofc, chainab, chainac, chainbc, ab, bc)

	return func(fa HKTA) bool {
		return apply(fa) && associativity(fa)
	}
}

// ChainAssertLaws asserts the apply laws `identity`, `composition`, `associative composition` and `associativity`
func ChainAssertLaws[HKTA, HKTB, HKTC, HKTAB, HKTBC, HKTAC, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqc E.Eq[HKTC],

	fofb pointed.Pointed[B, HKTB],
	fofc pointed.Pointed[C, HKTC],

	fofab pointed.Pointed[func(A) B, HKTAB],
	fofbc pointed.Pointed[func(B) C, HKTBC],

	faa functor.Functor[A, A, HKTA, HKTA],

	fmap functor.Functor[func(B) C, func(func(A) B) func(A) C, HKTBC, HKTABAC],

	chainab chain.Chainable[A, B, HKTA, HKTB, HKTAB],
	chainac chain.Chainable[A, C, HKTA, HKTC, HKTAC],
	chainbc chain.Chainable[B, C, HKTB, HKTC, HKTBC],

	fapabac apply.Apply[func(A) B, func(A) C, HKTAB, HKTAC, HKTABAC],

	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	// apply laws
	apply := L.ApplyAssertLaws(t, eqa, eqc, fofab, fofbc, faa, fmap, chain.ToApply(chainab), chain.ToApply(chainbc), chain.ToApply(chainac), fapabac, ab, bc)
	// chain laws
	associativity := ChainAssertAssociativity(t, eqc, fofb, fofc, chainab, chainac, chainbc, ab, bc)

	return func(fa HKTA) bool {
		return apply(fa) && associativity(fa)
	}
}
