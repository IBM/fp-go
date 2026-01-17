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

package functor

// Functor represents a type that can be mapped over, allowing transformation of values
// contained within a context while preserving the structure of that context.
//
// A Functor must satisfy the following laws:
//
// Identity:
//
//	Map(identity) == identity
//
// Composition:
//
//	Map(f ∘ g) == Map(f) ∘ Map(g)
//
// Type Parameters:
//   - A: The input value type contained in the functor
//   - B: The output value type after transformation
//   - HKTA: The higher-kinded type containing A (e.g., Option[A], Either[E, A])
//   - HKTB: The higher-kinded type containing B (e.g., Option[B], Either[E, B])
//
// Example:
//
//	// Given a functor for Option[int]
//	var f Functor[int, string, Option[int], Option[string]]
//	mapFn := f.Map(strconv.Itoa)
//	result := mapFn(Some(42)) // Returns Some("42")
type Functor[A, B, HKTA, HKTB any] interface {
	// Map transforms the value inside the functor using the provided function,
	// preserving the structure of the functor.
	//
	// Returns a function that takes a functor containing A and returns a functor containing B.
	Map(func(A) B) func(HKTA) HKTB
}

type MapType[A, B, HKTA, HKTB any] = func(func(A) B) func(HKTA) HKTB
