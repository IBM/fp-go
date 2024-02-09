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
	"github.com/IBM/fp-go/internal/applicative"
	LA "github.com/IBM/fp-go/internal/applicative/testing"
	"github.com/IBM/fp-go/internal/chain"
	LC "github.com/IBM/fp-go/internal/chain/testing"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
	"github.com/stretchr/testify/assert"
)

// Apply monad left identity law
//
// M.chain(M.of(a), f) <-> f(a)
//
// Deprecated: use [MonadAssertLeftIdentity] instead
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

// Apply monad left identity law
//
// M.chain(M.of(a), f) <-> f(a)
func MonadAssertLeftIdentity[HKTA, HKTB, HKTFAB, A, B any](t *testing.T,
	eq E.Eq[HKTB],

	fofb pointed.Pointed[B, HKTB],

	ma monad.Monad[A, B, HKTA, HKTB, HKTFAB],

	ab func(A) B,
) func(a A) bool {
	return func(a A) bool {

		f := func(a A) HKTB {
			return fofb.Of(ab(a))
		}

		left := ma.Chain(f)(ma.Of(a))
		right := f(a)

		return assert.True(t, eq.Equals(left, right), "Monad left identity")
	}
}

// Apply monad right identity law
//
// M.chain(fa, M.of) <-> fa
//
// Deprecated: use [MonadAssertRightIdentity] instead
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

// Apply monad right identity law
//
// M.chain(fa, M.of) <-> fa
func MonadAssertRightIdentity[HKTA, HKTAA, A any](t *testing.T,
	eq E.Eq[HKTA],

	ma monad.Monad[A, A, HKTA, HKTA, HKTAA],

) func(fa HKTA) bool {
	return func(fa HKTA) bool {

		left := ma.Chain(ma.Of)(fa)
		right := fa

		return assert.True(t, eq.Equals(left, right), "Monad right identity")
	}
}

// AssertLaws asserts the apply laws `identity`, `composition`, `associative composition`, 'applicative identity', 'homomorphism', 'interchange', `associativity`, `left identity`, `right identity`
//
// Deprecated: use [MonadAssertLaws] instead
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

// MonadAssertLaws asserts the apply laws `identity`, `composition`, `associative composition`, 'applicative identity', 'homomorphism', 'interchange', `associativity`, `left identity`, `right identity`
func MonadAssertLaws[HKTA, HKTB, HKTC, HKTAA, HKTAB, HKTBC, HKTAC, HKTABB, HKTABAC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqb E.Eq[HKTB],
	eqc E.Eq[HKTC],

	fofc pointed.Pointed[C, HKTC],
	fofaa pointed.Pointed[func(A) A, HKTAA],
	fofbc pointed.Pointed[func(B) C, HKTBC],
	fofabb pointed.Pointed[func(func(A) B) B, HKTABB],

	fmap functor.Functor[func(B) C, func(func(A) B) func(A) C, HKTBC, HKTABAC],

	fapabb applicative.Applicative[func(A) B, B, HKTAB, HKTB, HKTABB],
	fapabac applicative.Applicative[func(A) B, func(A) C, HKTAB, HKTAC, HKTABAC],

	maa monad.Monad[A, A, HKTA, HKTA, HKTAA],
	mab monad.Monad[A, B, HKTA, HKTB, HKTAB],
	mac monad.Monad[A, C, HKTA, HKTC, HKTAC],
	mbc monad.Monad[B, C, HKTB, HKTC, HKTBC],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {
	// derivations
	fofa := monad.ToPointed(maa)
	fofb := monad.ToPointed(mbc)
	fofab := applicative.ToPointed(fapabb)
	fapaa := monad.ToApplicative(maa)
	fapab := monad.ToApplicative(mab)
	chainab := monad.ToChainable(mab)
	chainac := monad.ToChainable(mac)
	chainbc := monad.ToChainable(mbc)
	fapbc := chain.ToApply(chainbc)
	fapac := chain.ToApply(chainac)

	faa := monad.ToFunctor(maa)

	// applicative laws
	apLaw := LA.ApplicativeAssertLaws(t, eqa, eqb, eqc, fofb, fofaa, fofbc, fofabb, faa, fmap, fapaa, fapab, fapbc, fapac, fapabb, fapabac, ab, bc)
	// chain laws
	chainLaw := LC.ChainAssertLaws(t, eqa, eqc, fofb, fofc, fofab, fofbc, faa, fmap, chainab, chainac, chainbc, applicative.ToApply(fapabac), ab, bc)
	// monad laws
	leftIdentity := MonadAssertLeftIdentity(t, eqb, fofb, mab, ab)
	rightIdentity := MonadAssertRightIdentity(t, eqa, maa)

	return func(a A) bool {
		fa := fofa.Of(a)
		return apLaw(a) && chainLaw(fa) && leftIdentity(a) && rightIdentity(fa)
	}
}
