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

// ApplicativeMonoid lifts a monoid into an applicative functor context.
//
// This function creates a monoid for applicative functor values (HKTA) given a monoid
// for the base type (A). It uses the applicative functor's operations (of, map, ap)
// to lift the monoid operations into the applicative context.
//
// This is useful for combining values that are wrapped in applicative contexts
// (like Option, Either, IO, etc.) using the underlying monoid's combination logic.
//
// Type Parameters:
//   - A: The base type with a monoid
//   - HKTA: The higher-kinded type representing the applicative functor applied to A
//   - HKTFA: The higher-kinded type representing the applicative functor applied to func(A) A
//
// Parameters:
//   - fof: The "pure" or "of" operation that lifts a value into the applicative context
//   - fmap: The map operation for the applicative functor
//   - fap: The apply operation for the applicative functor
//   - m: The monoid for the base type A
//
// Returns:
//   - A Monoid[HKTA] that combines applicative values using the base monoid
//
// Example (conceptual with Option-like type):
//
//	type Option[A any] struct { value *A }
//
//	intAddMonoid := MakeMonoid(
//	    func(a, b int) int { return a + b },
//	    0,
//	)
//
//	optMonoid := ApplicativeMonoid(
//	    some,  // func(int) Option[int]
//	    fmap,  // func(Option[int], func(int) func(int) int) Option[func(int) int]
//	    fap,   // func(Option[func(int) int], Option[int]) Option[int]
//	    intAddMonoid,
//	)
//
//	// Combine Option values using addition
//	result := optMonoid.Concat(Some(5), Some(3))  // Some(8)
//	empty := optMonoid.Empty()                     // Some(0)
func ApplicativeMonoid[A, HKTA, HKTFA any](
	fof func(A) HKTA,
	fmap func(HKTA, func(A) func(A) A) HKTFA,
	fap func(HKTFA, HKTA) HKTA,

	m Monoid[A],
) Monoid[HKTA] {

	return MakeMonoid(
		S.ApplySemigroup(fmap, fap, m).Concat,
		fof(m.Empty()),
	)
}
