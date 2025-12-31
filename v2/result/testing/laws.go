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

	ET "github.com/IBM/fp-go/v2/either/testing"
	EQ "github.com/IBM/fp-go/v2/eq"
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
func AssertLaws[A, B, C any](t *testing.T,
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {
	return ET.AssertLaws(t, EQ.FromStrictEquals[error](), eqa, eqb, eqc, ab, bc)
}
