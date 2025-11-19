// Copyright (c) 2025 IBM Corp.
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

import (
	F "github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// MonadApV is the applicative validation functor that combines errors using a semigroup.
//
// Unlike the standard [MonadAp] which short-circuits on the first Left (error),
// MonadApV accumulates all errors using the provided semigroup's Concat operation.
// This is particularly useful for validation scenarios where you want to collect
// all validation errors rather than stopping at the first one.
//
// The function takes a semigroup for combining errors and returns a function that
// applies a wrapped function to a wrapped value, accumulating errors if both are Left.
//
// Behavior:
//   - If both fab and fa are Left, combines their errors using sg.Concat
//   - If only fab is Left, returns Left with fab's error
//   - If only fa is Left, returns Left with fa's error
//   - If both are Right, applies the function and returns Right with the result
//
// Type Parameters:
//   - B: The result type after applying the function
//   - E: The error type (must support the semigroup operation)
//   - A: The input type to the function
//
// Parameters:
//   - sg: A semigroup that defines how to combine two error values
//
// Returns:
//   - A function that takes a wrapped function and a wrapped value, returning
//     Either[E, B] with accumulated errors or the computed result
//
// Example:
//
//	// Define a semigroup that concatenates error messages
//	errorSemigroup := semigroup.MakeSemigroup(func(e1, e2 string) string {
//	    return e1 + "; " + e2
//	})
//
//	// Create the validation applicative
//	applyV := either.MonadApV[int](errorSemigroup)
//
//	// Both are errors - errors get combined
//	fab := either.Left[func(int) int]("error1")
//	fa := either.Left[int]("error2")
//	result := applyV(fab, fa) // Left("error1; error2")
//
//	// One error - returns that error
//	fab2 := either.Right[string](N.Mul(2))
//	fa2 := either.Left[int]("validation failed")
//	result2 := applyV(fab2, fa2) // Left("validation failed")
//
//	// Both success - applies function
//	fab3 := either.Right[string](N.Mul(2))
//	fa3 := either.Right[string](21)
//	result3 := applyV(fab3, fa3) // Right(42)
func MonadApV[B, A, E any](sg S.Semigroup[E]) func(fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
	return func(fab Either[E, func(a A) B], fa Either[E, A]) Either[E, B] {
		if fab.isLeft {
			if fa.isLeft {
				return Left[B](sg.Concat(fab.l, fa.l))
			}
			return Left[B](fab.l)
		}
		if fa.isLeft {
			return Left[B](fa.l)
		}
		return Of[E](fab.r(fa.r))
	}
}

// ApV is the curried version of [MonadApV] that combines errors using a semigroup.
//
// This function provides a more convenient API for validation scenarios by currying
// the arguments. It first takes the value to validate, then returns a function that
// takes the validation function. This allows for a more natural composition style.
//
// Like [MonadApV], this accumulates all errors using the provided semigroup instead
// of short-circuiting on the first error. This is the key difference from the
// standard [Ap] function.
//
// Type Parameters:
//   - B: The result type after applying the function
//   - E: The error type (must support the semigroup operation)
//   - A: The input type to the function
//
// Parameters:
//   - sg: A semigroup that defines how to combine two error values
//
// Returns:
//   - A function that takes a value Either[E, A] and returns an Operator that
//     applies validation functions while accumulating errors
//
// Example:
//
//	// Define a semigroup for combining validation errors
//	type ValidationError struct {
//	    Errors []string
//	}
//	errorSemigroup := semigroup.MakeSemigroup(func(e1, e2 ValidationError) ValidationError {
//	    return ValidationError{Errors: append(e1.Errors, e2.Errors...)}
//	})
//
//	// Create validators
//	validatePositive := func(x int) either.Either[ValidationError, int] {
//	    if x > 0 {
//	        return either.Right[ValidationError](x)
//	    }
//	    return either.Left[int](ValidationError{Errors: []string{"must be positive"}})
//	}
//
//	// Use ApV for validation
//	applyValidation := either.ApV[int](errorSemigroup)
//	value := either.Left[int](ValidationError{Errors: []string{"invalid input"}})
//	validator := either.Left[func(int) int](ValidationError{Errors: []string{"invalid validator"}})
//
//	result := applyValidation(value)(validator)
//	// Left(ValidationError{Errors: []string{"invalid validator", "invalid input"}})
//
//go:inline
func ApV[B, A, E any](sg S.Semigroup[E]) func(Either[E, A]) Operator[E, func(A) B, B] {
	apv := MonadApV[B, A](sg)
	return func(e Either[E, A]) Operator[E, func(A) B, B] {
		return F.Bind2nd(apv, e)
	}
}
