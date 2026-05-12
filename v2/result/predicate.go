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

package result

import "github.com/IBM/fp-go/v2/either"

// Exists creates a predicate that tests whether a Result value is Ok and its value satisfies the given predicate.
// It returns a function that takes a Result[T] and returns true only if the Result is Ok and the predicate p
// returns true for the Ok value.
//
// This function is a specialized version of either.Exists for the Result type, where the error type is
// fixed to Go's built-in error interface. It's useful for checking if a Result contains a successful value
// that meets certain criteria, commonly used in filtering operations, validation chains, or conditional
// logic where you need to verify both the success state and a property of the success value.
//
// The behavior is as follows:
//   - If the input is an error (Err), returns false (regardless of the predicate)
//   - If the input is Ok and p returns true for the Ok value, returns true
//   - If the input is Ok and p returns false for the Ok value, returns false
//
// Type Parameters:
//   - T: The type of the Ok value (success type)
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes a Result[T] and returns true if it's Ok and satisfies p
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    N "github.com/IBM/fp-go/v2/number"
//	    "errors"
//	)
//
//	// Check if Result contains a positive number
//	isPositive := N.MoreThan(0)
//	hasPositive := R.Exists(isPositive)
//
//	result1 := hasPositive(R.Of(5))
//	// result1 = true (Ok with positive value)
//
//	result2 := hasPositive(R.Of(-3))
//	// result2 = false (Ok with non-positive value)
//
//	result3 := hasPositive(R.Error[int](errors.New("error")))
//	// result3 = false (error value)
//
//	// Use in filtering
//	values := []R.Result[int]{
//	    R.Of(5),
//	    R.Error[int](errors.New("error")),
//	    R.Of(-3),
//	    R.Of(10),
//	}
//	hasPositiveValue := func(r R.Result[int]) bool {
//	    return hasPositive(r)
//	}
//	// Filter to keep only Results with positive Ok values
//	// filtered would contain: [Ok(5), Ok(10)]
//
// See Also:
//   - ExistsError: Tests if Result is an error and satisfies a predicate
//   - Filter: Converts Ok values that fail a predicate to errors
//   - either.Exists: The underlying implementation for Either types
func Exists[T any](p Predicate[T]) Predicate[Result[T]] {
	return either.Exists[error](p)
}

// ExistsError creates a predicate that tests whether a Result value is an error and the error satisfies the given predicate.
// It returns a function that takes a Result[T] and returns true only if the Result is an error and the predicate p
// returns true for the error value.
//
// This function is a specialized version of either.ExistsLeft for the Result type, where the error type is
// fixed to Go's built-in error interface. It's useful for checking if a Result contains an error that meets
// certain criteria, commonly used in error filtering, error categorization, or conditional logic where you
// need to verify both the error state and a property of the error value.
//
// The behavior is as follows:
//   - If the input is Ok, returns false (regardless of the predicate)
//   - If the input is an error and p returns true for the error value, returns true
//   - If the input is an error and p returns false for the error value, returns false
//
// Type Parameters:
//   - T: The type of the Ok value (success type)
//
// Parameters:
//   - p: A predicate function that tests error values
//
// Returns:
//
//	A Predicate function that takes a Result[T] and returns true if it's an error and satisfies p
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    "errors"
//	    "strings"
//	)
//
//	// Check if Result contains a validation error
//	isValidationError := func(err error) bool {
//	    return strings.Contains(err.Error(), "validation")
//	}
//	hasValidationError := R.ExistsError[int](isValidationError)
//
//	result1 := hasValidationError(R.Error[int](errors.New("validation: invalid input")))
//	// result1 = true (error with validation message)
//
//	result2 := hasValidationError(R.Error[int](errors.New("network: connection failed")))
//	// result2 = false (error without validation message)
//
//	result3 := hasValidationError(R.Of(42))
//	// result3 = false (Ok value)
//
//	// Use in error categorization
//	results := []R.Result[int]{
//	    R.Error[int](errors.New("validation: empty field")),
//	    R.Of(100),
//	    R.Error[int](errors.New("network: timeout")),
//	    R.Error[int](errors.New("validation: invalid format")),
//	}
//	hasValidationErr := func(r R.Result[int]) bool {
//	    return hasValidationError(r)
//	}
//	// Filter to find validation errors
//	// filtered would contain: [Err("validation: empty field"), Err("validation: invalid format")]
//
// See Also:
//   - Exists: Tests if Result is Ok and satisfies a predicate
//   - IsError: Tests if Result is an error without checking the error value
//   - either.ExistsLeft: The underlying implementation for Either types
func ExistsError[T any](p Predicate[error]) Predicate[Result[T]] {
	return either.ExistsLeft[T](p)
}

// ForAll creates a predicate that tests whether a Result value is an error or its Ok value satisfies the given predicate.
// It returns a function that takes a Result[T] and returns true if the Result is an error (regardless of its value)
// or if it's Ok and the predicate p returns true for the Ok value.
//
// This function is a specialized version of either.ForAll for the Result type, where the error type is
// fixed to Go's built-in error interface. It implements universal quantification over the Result type.
// In logical terms, it states: "for all values in the Result, the predicate holds" - which is vacuously
// true for error values (empty case) and requires the predicate to hold for Ok values (non-empty case).
//
// The behavior is as follows:
//   - If the input is an error, returns true (vacuous truth - predicate holds for empty case)
//   - If the input is Ok and p returns true for the Ok value, returns true
//   - If the input is Ok and p returns false for the Ok value, returns false
//
// Relationship to Haskell and Category Theory:
//
// In Haskell, this corresponds to the all function for the Either type when viewed as a Foldable:
//
//   all :: Foldable t => (a -> Bool) -> t a -> Bool
//   all p (Right x) = p x
//   all p (Left _)  = True
//
// For Result, which is Either[error, T]:
//   all p (Ok x)  = p x
//   all p (Err _) = True
//
// From a category theory perspective, Result[T] is a coproduct (sum type) in the category of types.
// ForAll implements a natural transformation from predicates on T to predicates on Result[T],
// preserving the logical structure where:
//   - The error case represents the "empty" or "absent" case (like Nothing in Maybe/Option)
//   - The Ok case represents the "present" case that must satisfy the predicate
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
//   - T: The type of the Ok value (success type)
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes a Result[T] and returns true if it's an error or Ok with p satisfied
//
// Example:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/result"
//	    N "github.com/IBM/fp-go/v2/number"
//	    "errors"
//	)
//
//	// Check if Result is an error or contains a positive number
//	isPositive := N.MoreThan(0)
//	allPositive := R.ForAll(isPositive)
//
//	result1 := allPositive(R.Of(5))
//	// result1 = true (Ok with positive value satisfies predicate)
//
//	result2 := allPositive(R.Of(-3))
//	// result2 = false (Ok with non-positive value fails predicate)
//
//	result3 := allPositive(R.Error[int](errors.New("error")))
//	// result3 = true (error is vacuously true - no value to check)
//
//	// Use in validation: ensure all successful results meet criteria
//	values := []R.Result[int]{
//	    R.Of(5),
//	    R.Error[int](errors.New("error")),  // Ignored (vacuously true)
//	    R.Of(10),
//	    R.Of(-3),                            // Fails validation
//	}
//	allValid := true
//	for _, v := range values {
//	    if !allPositive(v) {
//	        allValid = false
//	        break
//	    }
//	}
//	// allValid = false because Of(-3) fails the predicate
//
//	// Contrast with Exists:
//	hasPositive := R.Exists(isPositive)
//	// hasPositive checks if there EXISTS an Ok value satisfying p
//	// allPositive checks if ALL Ok values satisfy p (errors are ignored)
//
// See Also:
//   - Exists: Tests if Result is Ok and satisfies a predicate (existential quantification)
//   - ExistsError: Tests if Result is an error and satisfies a predicate
//   - Filter: Converts Ok values that fail a predicate to errors
//   - either.ForAll: The underlying implementation for Either types
func ForAll[T any](p Predicate[T]) Predicate[Result[T]] {
	return either.ForAll[error](p)
}
