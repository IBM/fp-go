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

// Package testing provides law-based testing utilities for monads.
//
// This package implements property-based tests for the three fundamental monad laws:
//   - Left Identity: Chain(f)(Of(a)) == f(a)
//   - Right Identity: Chain(Of)(m) == m
//   - Associativity: Chain(g)(Chain(f)(m)) == Chain(x => Chain(g)(f(x)))(m)
//
// Additionally, it validates that monads satisfy all prerequisite laws from:
//   - Functor (identity, composition)
//   - Apply (composition)
//   - Applicative (identity, homomorphism, interchange)
//   - Chainable (associativity)
//
// Usage:
//
//	func TestMyMonad(t *testing.T) {
//	    // Set up equality checkers
//	    eqa := eq.FromEquals[Option[int]](...)
//	    eqb := eq.FromEquals[Option[string]](...)
//	    eqc := eq.FromEquals[Option[float64]](...)
//
//	    // Set up monad instances
//	    maa := &optionMonad[int, int]{}
//	    mab := &optionMonad[int, string]{}
//	    // ... etc
//
//	    // Run the law tests
//	    lawTest := MonadAssertLaws(t, eqa, eqb, eqc, ...)
//	    assert.True(t, lawTest(42))
//	}
package testing

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/internal/applicative"
	LA "github.com/IBM/fp-go/v2/internal/applicative/testing"
	"github.com/IBM/fp-go/v2/internal/chain"
	LC "github.com/IBM/fp-go/v2/internal/chain/testing"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	"github.com/stretchr/testify/assert"
)

// AssertLeftIdentity tests the monad left identity law:
// M.chain(M.of(a), f) <-> f(a)
//
// This law ensures that lifting a value into the monad and immediately chaining
// a function over it is equivalent to just applying the function directly.
//
// Deprecated: use [MonadAssertLeftIdentity] instead
func AssertLeftIdentity[HKTA, HKTB, A, B any](t *testing.T,
	eq E.Eq[HKTB],

	fofa func(A) HKTA,
	fofb func(B) HKTB,

	fchain func(HKTA, func(A) HKTB) HKTB,

	ab func(A) B,
) func(a A) bool {
	t.Helper()

	return func(a A) bool {

		f := func(a A) HKTB {
			return fofb(ab(a))
		}

		left := fchain(fofa(a), f)
		right := f(a)

		result := eq.Equals(left, right)
		if !result {
			t.Logf("Monad left identity violated: Chain(Of(%v), f) != f(%v)", a, a)
		}
		return assert.True(t, result, "Monad left identity")
	}
}

// MonadAssertLeftIdentity tests the monad left identity law:
// M.chain(M.of(a), f) <-> f(a)
//
// This law states that wrapping a value with Of and immediately chaining with a function f
// should be equivalent to just applying f to the value directly.
//
// In other words: lifting a value into the monad and then binding over it should have
// the same effect as just applying the function.
//
// Example:
//
//	// For Option monad with value 42 and function f(x) = Some(toString(x)):
//	Chain(f)(Of(42)) == f(42)
//	// Both produce: Some("42")
//
// Type Parameters:
//   - A: Input value type
//   - B: Output value type
//   - HKTA: Higher-kinded type containing A
//   - HKTB: Higher-kinded type containing B
//   - HKTFAB: Higher-kinded type containing function A -> B
func MonadAssertLeftIdentity[HKTA, HKTB, HKTFAB, A, B any](t *testing.T,
	eq E.Eq[HKTB],

	fofb pointed.Pointed[B, HKTB],

	ma monad.Monad[A, B, HKTA, HKTB, HKTFAB],

	ab func(A) B,
) func(a A) bool {
	t.Helper()

	return func(a A) bool {

		f := func(a A) HKTB {
			return fofb.Of(ab(a))
		}

		left := ma.Chain(f)(ma.Of(a))
		right := f(a)

		result := eq.Equals(left, right)
		if !result {
			t.Errorf("Monad left identity law violated:\n"+
				"  Chain(f)(Of(a)) != f(a)\n"+
				"  where a = %v\n"+
				"  Expected: Chain(f)(Of(a)) to equal f(a)", a)
		}
		return assert.True(t, result, "Monad left identity")
	}
}

// AssertRightIdentity tests the monad right identity law:
// M.chain(fa, M.of) <-> fa
//
// This law ensures that chaining a monadic value with the Of (pure/return) function
// returns the original monadic value unchanged.
//
// Deprecated: use [MonadAssertRightIdentity] instead
func AssertRightIdentity[HKTA, A any](t *testing.T,
	eq E.Eq[HKTA],

	fofa func(A) HKTA,

	fchain func(HKTA, func(A) HKTA) HKTA,
) func(fa HKTA) bool {
	t.Helper()

	return func(fa HKTA) bool {

		left := fchain(fa, fofa)
		right := fa

		result := eq.Equals(left, right)
		if !result {
			t.Logf("Monad right identity violated: Chain(fa, Of) != fa")
		}
		return assert.True(t, result, "Monad right identity")
	}
}

// MonadAssertRightIdentity tests the monad right identity law:
// M.chain(fa, M.of) <-> fa
//
// This law states that chaining a monadic value with the Of (pure/return) function
// should return the original monadic value unchanged. This ensures that Of doesn't
// add any additional structure or effects beyond what's already present.
//
// Example:
//
//	// For Option monad with value Some(42):
//	Chain(Of)(Some(42)) == Some(42)
//	// The value remains unchanged
//
// Type Parameters:
//   - A: The value type
//   - HKTA: Higher-kinded type containing A
//   - HKTAA: Higher-kinded type containing function A -> A
func MonadAssertRightIdentity[HKTA, HKTAA, A any](t *testing.T,
	eq E.Eq[HKTA],

	ma monad.Monad[A, A, HKTA, HKTA, HKTAA],

) func(fa HKTA) bool {
	t.Helper()

	return func(fa HKTA) bool {

		left := ma.Chain(ma.Of)(fa)
		right := fa

		result := eq.Equals(left, right)
		if !result {
			t.Errorf("Monad right identity law violated:\n" +
				"  Chain(Of)(fa) != fa\n" +
				"  Expected: Chain(Of)(fa) to equal fa")
		}
		return assert.True(t, result, "Monad right identity")
	}
}

// AssertLaws tests all monad laws including prerequisite laws from Functor, Apply,
// Applicative, and Chainable.
//
// This function validates:
//   - Functor laws: identity, composition
//   - Apply law: composition
//   - Applicative laws: identity, homomorphism, interchange
//   - Chainable law: associativity
//   - Monad laws: left identity, right identity
//
// The monad laws build upon and require all the prerequisite laws to hold.
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
	t.Helper()

	// applicative laws
	applicative := LA.AssertLaws(t, eqa, eqb, eqc, fofa, fofb, fofaa, fofab, fofbc, fofabb, faa, fab, fac, fbc, fmap, fapaa, fapab, fapbc, fapac, fapabb, fapabac, ab, bc)
	// chain laws
	chain := LC.AssertLaws(t, eqa, eqc, fofb, fofc, fofab, fofbc, faa, fab, fac, fbc, fmap, chainab, chainac, chainbc, fapab, fapbc, fapac, fapabac, ab, bc)
	// monad laws
	leftIdentity := AssertLeftIdentity(t, eqb, fofa, fofb, chainab, ab)
	rightIdentity := AssertRightIdentity(t, eqa, fofa, chainaa)

	return func(a A) bool {
		fa := fofa(a)
		appOk := applicative(a)
		chainOk := chain(fa)
		leftIdOk := leftIdentity(a)
		rightIdOk := rightIdentity(fa)

		if !appOk {
			t.Logf("Applicative laws failed for input: %v", a)
		}
		if !chainOk {
			t.Logf("Chain laws failed for input: %v", fa)
		}
		if !leftIdOk {
			t.Logf("Left identity law failed for input: %v", a)
		}
		if !rightIdOk {
			t.Logf("Right identity law failed for input: %v", fa)
		}

		return appOk && chainOk && leftIdOk && rightIdOk
	}
}

// MonadAssertLaws validates all monad laws and prerequisite laws for a monad implementation.
//
// This is the primary testing function for verifying monad implementations. It checks:
//
// Monad Laws (primary):
//   - Left Identity: Chain(f)(Of(a)) == f(a)
//   - Right Identity: Chain(Of)(m) == m
//   - Associativity: Chain(g)(Chain(f)(m)) == Chain(x => Chain(g)(f(x)))(m)
//
// Prerequisite Laws (inherited from parent type classes):
//   - Functor: identity, composition
//   - Apply: composition
//   - Applicative: identity, homomorphism, interchange
//   - Chainable: associativity
//
// Usage:
//
//	func TestOptionMonad(t *testing.T) {
//	    // Create equality checkers
//	    eqa := eq.FromEquals[Option[int]](optionEq[int])
//	    eqb := eq.FromEquals[Option[string]](optionEq[string])
//	    eqc := eq.FromEquals[Option[float64]](optionEq[float64])
//
//	    // Define test functions
//	    ab := strconv.Itoa
//	    bc := func(s string) float64 { v, _ := strconv.ParseFloat(s, 64); return v }
//
//	    // Create monad instances and other required type class instances
//	    // ... (setup code)
//
//	    // Run the law tests
//	    lawTest := MonadAssertLaws(t, eqa, eqb, eqc, ...)
//	    assert.True(t, lawTest(42))
//	}
//
// Type Parameters:
//   - A, B, C: Value types for testing transformations
//   - HKTA, HKTB, HKTC: Higher-kinded types containing A, B, C
//   - HKTAA, HKTAB, HKTBC, HKTAC: Higher-kinded types containing functions
//   - HKTABB, HKTABAC: Higher-kinded types for applicative testing
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
	t.Helper()

	// Derive required type class instances from monad instances
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

	// Test prerequisite laws from parent type classes
	apLaw := LA.ApplicativeAssertLaws(t, eqa, eqb, eqc, fofb, fofaa, fofbc, fofabb, faa, fmap, fapaa, fapab, fapbc, fapac, fapabb, fapabac, ab, bc)
	chainLaw := LC.ChainAssertLaws(t, eqa, eqc, fofb, fofc, fofab, fofbc, faa, fmap, chainab, chainac, chainbc, applicative.ToApply(fapabac), ab, bc)

	// Test monad-specific laws
	leftIdentity := MonadAssertLeftIdentity(t, eqb, fofb, mab, ab)
	rightIdentity := MonadAssertRightIdentity(t, eqa, maa)

	return func(a A) bool {
		fa := fofa.Of(a)

		// Run all law tests and collect results
		apOk := apLaw(a)
		chainOk := chainLaw(fa)
		leftIdOk := leftIdentity(a)
		rightIdOk := rightIdentity(fa)

		// Log detailed failure information
		if !apOk {
			t.Errorf("Monad prerequisite failure: Applicative laws violated for input: %v", a)
		}
		if !chainOk {
			t.Errorf("Monad prerequisite failure: Chain laws violated for monadic value: %v", fa)
		}
		if !leftIdOk {
			t.Errorf("Monad law failure: Left identity violated for input: %v", a)
		}
		if !rightIdOk {
			t.Errorf("Monad law failure: Right identity violated for monadic value: %v", fa)
		}

		allOk := apOk && chainOk && leftIdOk && rightIdOk
		if allOk {
			t.Logf("✓ All monad laws satisfied for input: %v", a)
		}

		return allOk
	}
}

// MonadAssertAssociativity is a convenience function that tests only the monad associativity law
// (which is inherited from Chainable).
//
// Associativity Law:
//
//	Chain(g)(Chain(f)(m)) == Chain(x => Chain(g)(f(x)))(m)
//
// This law ensures that the order in which we nest chain operations doesn't matter,
// as long as the sequence of operations remains the same.
//
// Example:
//
//	// For Option monad:
//	f := func(x int) Option[string] { return Some(strconv.Itoa(x)) }
//	g := func(s string) Option[float64] { return Some(parseFloat(s)) }
//	m := Some(42)
//
//	// These should be equal:
//	Chain(g)(Chain(f)(m))                    // Some(42.0)
//	Chain(func(x int) { Chain(g)(f(x)) })(m) // Some(42.0)
//
// Type Parameters:
//   - A, B, C: Value types for the transformation chain
//   - HKTA, HKTB, HKTC: Higher-kinded types containing A, B, C
//   - HKTAB, HKTAC, HKTBC: Higher-kinded types containing functions
func MonadAssertAssociativity[HKTA, HKTB, HKTC, HKTAB, HKTAC, HKTBC, A, B, C any](
	t *testing.T,
	eq E.Eq[HKTC],
	fofb pointed.Pointed[B, HKTB],
	fofc pointed.Pointed[C, HKTC],
	mab monad.Monad[A, B, HKTA, HKTB, HKTAB],
	mac monad.Monad[A, C, HKTA, HKTC, HKTAC],
	mbc monad.Monad[B, C, HKTB, HKTC, HKTBC],
	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	t.Helper()

	chainab := monad.ToChainable(mab)
	chainac := monad.ToChainable(mac)
	chainbc := monad.ToChainable(mbc)

	return LC.ChainAssertAssociativity(t, eq, fofb, fofc, chainab, chainac, chainbc, ab, bc)
}

// TestMonadLaws is a helper function that runs all monad law tests with common test values.
//
// This is a convenience wrapper around MonadAssertLaws that runs the law tests with
// a set of test values and reports the results. It's useful for quick testing with
// standard inputs.
//
// Parameters:
//   - t: The testing.T instance
//   - name: A descriptive name for this test suite
//   - testValues: A slice of values of type A to test with
//   - Other parameters: Same as MonadAssertLaws
//
// Example:
//
//	func TestOptionMonadLaws(t *testing.T) {
//	    testValues := []int{0, 1, -1, 42, 100}
//	    TestMonadLaws(t, "Option[int]", testValues, eqa, eqb, eqc, ...)
//	}
func TestMonadLaws[HKTA, HKTB, HKTC, HKTAA, HKTAB, HKTBC, HKTAC, HKTABB, HKTABAC, A, B, C any](
	t *testing.T,
	name string,
	testValues []A,
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
) {
	t.Helper()

	lawTest := MonadAssertLaws(t, eqa, eqb, eqc, fofc, fofaa, fofbc, fofabb, fmap, fapabb, fapabac, maa, mab, mac, mbc, ab, bc)

	t.Run(fmt.Sprintf("MonadLaws_%s", name), func(t *testing.T) {
		for i, val := range testValues {
			t.Run(fmt.Sprintf("Value_%d", i), func(t *testing.T) {
				result := lawTest(val)
				assert.True(t, result, "Monad laws should hold for value: %v", val)
			})
		}
	})
}
