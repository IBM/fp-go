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

package pointed

// Pointed represents a type that can lift a pure value into a computational context.
//
// Pointed is the minimal extension of a Functor that adds the ability to create
// a context-wrapped value from a bare value. It provides the canonical way to
// construct values of a higher-kinded type.
//
// Type Parameters:
//   - A: The value type to be lifted into the context
//   - HKTA: The higher-kinded type containing A (e.g., Option[A], Either[E, A])
//
// Example:
//
//	// Given a pointed functor for Option[int]
//	var p Pointed[int, Option[int]]
//	result := p.Of(42) // Returns Some(42)
type Pointed[A, HKTA any] interface {
	// Of lifts a pure value into its higher-kinded type context.
	//
	// This operation wraps a value A in the minimal context required by the type HKTA,
	// creating a valid instance of the higher-kinded type.
	Of(A) HKTA
}

type OfType[A, HKTA any] = func(A) HKTA
