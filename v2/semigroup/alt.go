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

// AltSemigroup creates a Semigroup for alternative functors (types with an alt operation).
// The alt operation provides a way to combine two values of the same higher-kinded type,
// typically representing alternative computations or choices.
//
// The function takes a curried alt operation that accepts a lazy (thunked) value and returns
// an operator on the first value (the shape of Alt in the concrete packages), and returns a
// Semigroup that eagerly evaluates both values before combining them.
//
// Type parameters:
//   - HKTA: The higher-kinded type (e.g., Option[A], Either[E, A])
//   - LAZYHKTA: A lazy/thunked version of HKTA (must be func() HKTA)
//
// Example:
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	// Alt operation for Option: returns first if Some, otherwise evaluates second
//	sg := semigroup.AltSemigroup(O.Alt[int])
//	result := sg.Concat(O.Some(1), O.Some(2))  // Returns: Some(1)
//	result2 := sg.Concat(O.None[int](), O.Some(2))  // Returns: Some(2)
func AltSemigroup[HKTA any, LAZYHKTA ~func() HKTA](
	falt func(LAZYHKTA) func(HKTA) HKTA,

) Semigroup[HKTA] {

	return MakeSemigroup(
		func(first, second HKTA) HKTA {
			return falt(func() HKTA { return second })(first)
		},
	)
}
