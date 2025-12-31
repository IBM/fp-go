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

// Package testing provides utilities for testing Either monad laws.
// This is useful for verifying that custom Either implementations satisfy the monad laws.
package testing

import (
	"testing"

	ET "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
)

// AssertLaws asserts that the Either monad satisfies the monad laws.
// This includes testing:
//   - Identity laws (left and right identity)
//   - Associativity law
//   - Functor laws
//   - Applicative laws
//
// Parameters:
//   - t: Testing context
//   - eqe, eqa, eqb, eqc: Equality predicates for the types
//   - ab: Function from A to B for testing
//   - bc: Function from B to C for testing
//
// Returns a function that takes a value of type A and returns true if all laws hold.
//
// Example:
//
//	func TestEitherLaws(t *testing.T) {
//	    eqInt := eq.FromStrictEquals[int]()
//	    eqString := eq.FromStrictEquals[string]()
//	    eqError := eq.FromStrictEquals[error]()
//
//	    ab := strconv.Itoa
//	    bc := S.IsNonEmpty
//
//	    testing.AssertLaws(t, eqError, eqInt, eqString, eq.FromStrictEquals[bool](), ab, bc)(42)
//	}
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
