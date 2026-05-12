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

package option

// Exists creates a predicate that tests whether an Option value is Some and its value satisfies the given predicate.
// It returns a function that takes an Option[T] and returns true only if the Option is Some and the predicate p
// returns true for the contained value.
//
// This function is useful for checking if an Option contains a value that meets certain criteria,
// commonly used in filtering operations, validation chains, or conditional logic where you need to verify
// both the presence of a value and a property of that value.
//
// The behavior is as follows:
//   - If the input is None, returns false (regardless of the predicate)
//   - If the input is Some and p returns true for the contained value, returns true
//   - If the input is Some and p returns false for the contained value, returns false
//
// Type Parameters:
//   - T: The type of the value contained in the Option
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes an Option[T] and returns true if it's Some and satisfies p
//
// Example:
//
//	import (
//	    O "github.com/IBM/fp-go/v2/option"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	// Check if Option contains a positive number
//	isPositive := N.MoreThan(0)
//	hasPositive := O.Exists(isPositive)
//
//	result1 := hasPositive(O.Some(5))
//	// result1 = true (Some with positive value)
//
//	result2 := hasPositive(O.Some(-3))
//	// result2 = false (Some with non-positive value)
//
//	result3 := hasPositive(O.None[int]())
//	// result3 = false (None value)
//
//	// Use in filtering
//	values := []O.Option[int]{
//	    O.Some(5),
//	    O.None[int](),
//	    O.Some(-3),
//	    O.Some(10),
//	}
//	hasPositiveValue := func(opt O.Option[int]) bool {
//	    return hasPositive(opt)
//	}
//	// Filter to keep only Options with positive Some values
//	// filtered would contain: [Some(5), Some(10)]
//
// See Also:
//   - Filter: Converts Some values that fail a predicate to None
//   - IsSome: Tests if Option is Some without checking the value
func Exists[T any](p Predicate[T]) Predicate[Option[T]] {
	return func(o Option[T]) bool {
		return o.isSome && p(o.value)
	}
}

// ForAll creates a predicate that tests whether an Option value is None or its Some value satisfies the given predicate.
// It returns a function that takes an Option[T] and returns true if the Option is None (regardless of content)
// or if it's Some and the predicate p returns true for the contained value.
//
// This function implements universal quantification over the Option type. In logical terms, it states:
// "for all values in the Option, the predicate holds" - which is vacuously true for None values (empty case)
// and requires the predicate to hold for Some values (non-empty case).
//
// The behavior is as follows:
//   - If the input is None, returns true (vacuous truth - predicate holds for empty case)
//   - If the input is Some and p returns true for the contained value, returns true
//   - If the input is Some and p returns false for the contained value, returns false
//
// Relationship to Haskell and Category Theory:
//
// In Haskell, this corresponds to the all function for the Maybe type when viewed as a Foldable:
//
//   all :: Foldable t => (a -> Bool) -> t a -> Bool
//   all p (Just x)  = p x
//   all p Nothing   = True
//
// From a category theory perspective, Option[T] is a sum type representing optional values.
// ForAll implements a natural transformation from predicates on T to predicates on Option[T],
// preserving the logical structure where:
//   - The None case represents the "empty" or "absent" case
//   - The Some case represents the "present" case that must satisfy the predicate
//
// This is dual to Exists, which implements existential quantification:
//   - ForAll: "all elements satisfy p" (true for empty, requires p for non-empty)
//   - Exists: "some element satisfies p" (false for empty, requires p for non-empty)
//
// The relationship follows De Morgan's laws:
//   - ForAll(p) ≡ not(Exists(not(p)))
//   - Exists(p) ≡ not(ForAll(not(p)))
//
// Type Parameters:
//   - T: The type of the value contained in the Option
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes an Option[T] and returns true if it's None or Some with p satisfied
//
// Example:
//
//	import (
//	    O "github.com/IBM/fp-go/v2/option"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	// Check if Option is None or contains a positive number
//	isPositive := N.MoreThan(0)
//	allPositive := O.ForAll(isPositive)
//
//	result1 := allPositive(O.Some(5))
//	// result1 = true (Some with positive value satisfies predicate)
//
//	result2 := allPositive(O.Some(-3))
//	// result2 = false (Some with non-positive value fails predicate)
//
//	result3 := allPositive(O.None[int]())
//	// result3 = true (None is vacuously true - no value to check)
//
//	// Use in validation: ensure all present values meet criteria
//	values := []O.Option[int]{
//	    O.Some(5),
//	    O.None[int](),        // Ignored (vacuously true)
//	    O.Some(10),
//	    O.Some(-3),           // Fails validation
//	}
//	allValid := true
//	for _, v := range values {
//	    if !allPositive(v) {
//	        allValid = false
//	        break
//	    }
//	}
//	// allValid would be false because Some(-3) fails the predicate
//
//	// Contrast with Exists:
//	hasPositive := O.Exists(isPositive)
//	// hasPositive checks if there EXISTS a Some value satisfying p
//	// allPositive checks if ALL Some values satisfy p (None is ignored)
//
// See Also:
//   - Exists: Tests if Option is Some and satisfies a predicate (existential quantification)
//   - Filter: Converts Some values that fail a predicate to None
//   - IsSome: Tests if Option is Some without checking the value
func ForAll[T any](p Predicate[T]) Predicate[Option[T]] {
	return func(o Option[T]) bool {
		return !o.isSome || p(o.value)
	}
}
