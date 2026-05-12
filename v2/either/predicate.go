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

package either

// Exists creates a predicate that tests whether an Either value is Right and its value satisfies the given predicate.
// It returns a function that takes an Either[E, T] and returns true only if the Either is Right and the predicate p
// returns true for the Right value.
//
// This function is useful for checking if an Either contains a successful value that meets certain criteria,
// commonly used in filtering operations, validation chains, or conditional logic where you need to verify
// both the success state and a property of the success value.
//
// The behavior is as follows:
//   - If the input is Left, returns false (regardless of the predicate)
//   - If the input is Right and p returns true for the Right value, returns true
//   - If the input is Right and p returns false for the Right value, returns false
//
// Type Parameters:
//   - E: The type of the Left value (error type)
//   - T: The type of the Right value (success type)
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes an Either[E, T] and returns true if it's Right and satisfies p
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	// Check if Either contains a positive number
//	isPositive := N.MoreThan(0)
//	hasPositive := E.Exists[string](isPositive)
//
//	result1 := hasPositive(E.Right[string](5))
//	// result1 = true (Right with positive value)
//
//	result2 := hasPositive(E.Right[string](-3))
//	// result2 = false (Right with non-positive value)
//
//	result3 := hasPositive(E.Left[int]("error"))
//	// result3 = false (Left value)
//
//	// Use in filtering
//	values := []E.Either[string, int]{
//	    E.Right[string](5),
//	    E.Left[int]("error"),
//	    E.Right[string](-3),
//	    E.Right[string](10),
//	}
//	hasPositiveValue := func(e E.Either[string, int]) bool {
//	    return hasPositive(e)
//	}
//	// Filter to keep only Eithers with positive Right values
//	// filtered would contain: [Right(5), Right(10)]
//
// See Also:
//   - ExistsLeft: Tests if Either is Left and satisfies a predicate
//   - Filter: Converts Right values that fail a predicate to Left
func Exists[E, T any](p Predicate[T]) Predicate[Either[E, T]] {
	return func(e Either[E, T]) bool {
		return !e.isLeft && p(e.r)
	}
}

// ExistsLeft creates a predicate that tests whether an Either value is Left and its value satisfies the given predicate.
// It returns a function that takes an Either[E, T] and returns true only if the Either is Left and the predicate p
// returns true for the Left value.
//
// This function is useful for checking if an Either contains an error value that meets certain criteria,
// commonly used in error filtering, error categorization, or conditional logic where you need to verify
// both the error state and a property of the error value.
//
// The behavior is as follows:
//   - If the input is Right, returns false (regardless of the predicate)
//   - If the input is Left and p returns true for the Left value, returns true
//   - If the input is Left and p returns false for the Left value, returns false
//
// Type Parameters:
//   - T: The type of the Right value (success type)
//   - E: The type of the Left value (error type)
//
// Parameters:
//   - p: A predicate function that tests values of type E
//
// Returns:
//
//	A Predicate function that takes an Either[E, T] and returns true if it's Left and satisfies p
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    "strings"
//	)
//
//	// Check if Either contains a validation error
//	isValidationError := func(s string) bool {
//	    return strings.HasPrefix(s, "validation:")
//	}
//	hasValidationError := E.ExistsLeft[int](isValidationError)
//
//	result1 := hasValidationError(E.Left[int]("validation: invalid input"))
//	// result1 = true (Left with validation error)
//
//	result2 := hasValidationError(E.Left[int]("network: connection failed"))
//	// result2 = false (Left with non-validation error)
//
//	result3 := hasValidationError(E.Right[string](42))
//	// result3 = false (Right value)
//
//	// Use in error categorization
//	results := []E.Either[string, int]{
//	    E.Left[int]("validation: empty field"),
//	    E.Right[string](100),
//	    E.Left[int]("network: timeout"),
//	    E.Left[int]("validation: invalid format"),
//	}
//	hasValidationErr := func(e E.Either[string, int]) bool {
//	    return hasValidationError(e)
//	}
//	// Filter to find validation errors
//	// filtered would contain: [Left("validation: empty field"), Left("validation: invalid format")]
//
// See Also:
//   - Exists: Tests if Either is Right and satisfies a predicate
//   - IsLeft: Tests if Either is Left without checking the value
func ExistsLeft[T, E any](p Predicate[E]) Predicate[Either[E, T]] {
	return func(e Either[E, T]) bool {
		return e.isLeft && p(e.l)
	}
}

// ForAll creates a predicate that tests whether an Either value is Left or its Right value satisfies the given predicate.
// It returns a function that takes an Either[E, T] and returns true if the Either is Left (regardless of its value)
// or if it's Right and the predicate p returns true for the Right value.
//
// This function implements universal quantification over the Either type. In logical terms, it states:
// "for all values in the Either, the predicate holds" - which is vacuously true for Left values (empty case)
// and requires the predicate to hold for Right values (non-empty case).
//
// The behavior is as follows:
//   - If the input is Left, returns true (vacuous truth - predicate holds for empty case)
//   - If the input is Right and p returns true for the Right value, returns true
//   - If the input is Right and p returns false for the Right value, returns false
//
// Relationship to Haskell and Category Theory:
//
// In Haskell, this corresponds to the all function for the Either type when viewed as a Foldable:
//
//   all :: Foldable t => (a -> Bool) -> t a -> Bool
//   all p (Right x) = p x
//   all p (Left _)  = True
//
// From a category theory perspective, Either[E, T] is a coproduct (sum type) in the category of types.
// ForAll implements a natural transformation from predicates on T to predicates on Either[E, T],
// preserving the logical structure where:
//   - The Left case represents the "empty" or "absent" case (like Nothing in Maybe/Option)
//   - The Right case represents the "present" case that must satisfy the predicate
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
//   - E: The type of the Left value (error type)
//   - T: The type of the Right value (success type)
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes an Either[E, T] and returns true if it's Left or Right with p satisfied
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	// Check if Either is Left or contains a positive number
//	isPositive := N.MoreThan(0)
//	allPositive := E.ForAll[string](isPositive)
//
//	result1 := allPositive(E.Right[string](5))
//	// result1 = true (Right with positive value satisfies predicate)
//
//	result2 := allPositive(E.Right[string](-3))
//	// result2 = false (Right with non-positive value fails predicate)
//
//	result3 := allPositive(E.Left[int]("error"))
//	// result3 = true (Left is vacuously true - no value to check)
//
//	// Use in validation: ensure all successful results meet criteria
//	values := []E.Either[string, int]{
//	    E.Right[string](5),
//	    E.Left[int]("error"),      // Ignored (vacuously true)
//	    E.Right[string](10),
//	    E.Right[string](-3),       // Fails validation
//	}
//	allValid := func(e E.Either[string, int]) bool {
//	    return allPositive(e)
//	}
//	// Check if all non-error values are positive
//	// result would be false because Right(-3) fails the predicate
//
//	// Contrast with Exists:
//	hasPositive := E.Exists[string](isPositive)
//	// hasPositive checks if there EXISTS a Right value satisfying p
//	// allPositive checks if ALL Right values satisfy p (Left is ignored)
//
// See Also:
//   - Exists: Tests if Either is Right and satisfies a predicate (existential quantification)
//   - ExistsLeft: Tests if Either is Left and satisfies a predicate
//   - Filter: Converts Right values that fail a predicate to Left
func ForAll[E, T any](p Predicate[T]) Predicate[Either[E, T]] {
	return func(e Either[E, T]) bool {
		return e.isLeft || p(e.r)
	}
}
