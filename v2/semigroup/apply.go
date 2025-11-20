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

package semigroup

import (
	F "github.com/IBM/fp-go/v2/function"
)

// ApplySemigroup creates a Semigroup for applicative functors (types with map and ap operations).
// Given a Semigroup[A], it lifts it to work with higher-kinded types containing A values.
//
// The resulting semigroup combines two HKTA values by:
//  1. Mapping the curried concat operation over the first value
//  2. Applying the result to the second value using the ap operation
//
// This allows semigroup operations to be performed inside applicative contexts like
// Option, Either, or Array.
//
// Type parameters:
//   - A: The base type with a semigroup operation
//   - HKTA: The higher-kinded type containing A (e.g., Option[A], Either[E, A])
//   - HKTFA: The higher-kinded type containing a function (e.g., Option[func(A) A])
//
// Parameters:
//   - fmap: Maps a function over the applicative functor
//   - fap: Applies a function in the applicative context to a value in the context
//   - s: The base semigroup for type A
//
// Example:
//
//	import (
//	    O "github.com/IBM/fp-go/v2/option"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	intSum := N.SemigroupSum[int]()
//	optionSG := semigroup.ApplySemigroup(
//	    O.Map[int, func(int) int],
//	    O.Ap[int, int],
//	    intSum,
//	)
//
//	result := optionSG.Concat(O.Some(5), O.Some(10))  // Some(15)
//	result2 := optionSG.Concat(O.None[int](), O.Some(10))  // None
func ApplySemigroup[A, HKTA, HKTFA any](
	fmap func(HKTA, func(A) func(A) A) HKTFA,
	fap func(HKTFA, HKTA) HKTA,

	s Semigroup[A],
) Semigroup[HKTA] {

	cb := F.Curry2(s.Concat)
	return MakeSemigroup(func(first HKTA, second HKTA) HKTA {
		return fap(fmap(first, cb), second)
	})
}
