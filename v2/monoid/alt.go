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

package monoid

import (
	S "github.com/IBM/fp-go/v2/semigroup"
)

// AlternativeMonoid creates a monoid for types that are both Applicative and Alternative.
//
// This combines the behavior of ApplicativeMonoid with Alternative semantics, providing
// both applicative combination and fallback/choice behavior. The resulting monoid tries
// to combine values using the applicative monoid, but falls back to alternative behavior
// when needed.
//
// This is useful for types like Option, Either, or Parser that support both applicative
// composition and alternative/fallback semantics.
//
// Type Parameters:
//   - A: The base type with a monoid
//   - HKTA: The higher-kinded type representing the functor applied to A
//   - HKTFA: The higher-kinded type representing the functor applied to func(A) A
//   - LAZYHKTA: A lazy/deferred computation of HKTA (typically func() HKTA)
//
// Parameters:
//   - fof: The "pure" operation that lifts a value into the context
//   - fmap: The map operation for the functor
//   - fap: The apply operation for the applicative
//   - falt: The alternative operation providing fallback/choice behavior
//   - m: The monoid for the base type A
//
// Returns:
//   - A Monoid[HKTA] combining applicative and alternative semantics
//
// Example (conceptual with Option-like type):
//
//	intAddMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,
//	)
//
//	optMonoid := AlternativeMonoid(
//	    some,  // pure
//	    fmap,  // map
//	    fap,   // apply
//	    falt,  // alternative (fallback)
//	    intAddMonoid,
//	)
//
//	// Combines Some values using addition
//	result1 := optMonoid.Concat(Some(5), Some(3))  // Some(8)
//	// Falls back when first is None
//	result2 := optMonoid.Concat(None(), Some(3))   // Some(3)
func AlternativeMonoid[A, HKTA, HKTFA any, LAZYHKTA ~func() HKTA](
	fof func(A) HKTA,

	fmap func(HKTA, func(A) func(A) A) HKTFA,
	fap func(HKTFA, HKTA) HKTA,

	falt func(HKTA, LAZYHKTA) HKTA,

	m Monoid[A],

) Monoid[HKTA] {

	sg := ApplicativeMonoid(fof, fmap, fap, m)

	return MakeMonoid(
		func(first, second HKTA) HKTA {
			snd := func() HKTA { return second }

			return falt(sg.Concat(first, second), func() HKTA {
				return falt(first, snd)
			})
		},
		sg.Empty(),
	)
}

// AltMonoid creates a monoid from an Alt type class (alternative/choice operation).
//
// This creates a monoid for types that support alternative/fallback semantics,
// where the Concat operation tries the first value and falls back to the second
// if the first fails or is empty. The Empty value is provided by fzero.
//
// This is commonly used with Option, Either, Parser, and similar types that
// represent computations that may fail or have multiple alternatives.
//
// Type Parameters:
//   - HKTA: The higher-kinded type (e.g., Option[A], Either[E, A])
//   - LAZYHKTA: A lazy/deferred computation of HKTA (typically func() HKTA)
//
// Parameters:
//   - fzero: A lazy computation that produces the empty/zero value
//   - falt: The alternative operation that provides fallback behavior
//
// Returns:
//   - A Monoid[HKTA] with alternative/choice semantics
//
// Example (conceptual with Option-like type):
//
//	optMonoid := AltMonoid(
//	    func() Option[int] { return None() },  // empty
//	    func(first Option[int], second func() Option[int]) Option[int] {
//	        if first.IsSome() {
//	            return first
//	        }
//	        return second()
//	    },
//	)
//
//	// First Some wins
//	result1 := optMonoid.Concat(Some(5), Some(3))  // Some(5)
//	// Falls back to second when first is None
//	result2 := optMonoid.Concat(None(), Some(3))   // Some(3)
//	// Both None returns None
//	result3 := optMonoid.Concat(None(), None())    // None()
func AltMonoid[HKTA any, LAZYHKTA ~func() HKTA](
	fzero LAZYHKTA,
	falt func(HKTA, LAZYHKTA) HKTA,

) Monoid[HKTA] {

	return MakeMonoid(
		S.AltSemigroup(falt).Concat,
		fzero(),
	)
}
