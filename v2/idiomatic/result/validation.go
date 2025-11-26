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

package result

import (
	S "github.com/IBM/fp-go/v2/semigroup"
)

// ApV applies a function wrapped in a Result to a value wrapped in a Result,
// accumulating errors using a semigroup instead of short-circuiting.
//
// This function is designed for validation scenarios where you want to collect
// all validation errors rather than stopping at the first error. It differs
// from the standard [Ap] function in that it combines errors from both the
// function and the value using the provided semigroup operation.
//
// The function works as follows:
//   - If both the value and the function have errors, it combines them using the semigroup
//   - If only one has an error, it returns that error
//   - If neither has an error, it applies the function to the value
//
// Type Parameters:
//   - B: The result type after applying the function
//   - A: The input type to the function
//
// Parameters:
//   - sg: A semigroup that defines how to combine two error values. The semigroup's
//     Concat operation determines how errors are accumulated (e.g., concatenating
//     error messages, merging error lists, etc.)
//
// Returns:
//   - A curried function that takes a value (A, error), then takes a function
//     (func(A) B, error), and returns the result (B, error) with accumulated errors
//
// Behavior:
//   - Right + Right: Applies the function to the value and returns Right(result)
//   - Right + Left: Returns the Left error from the function
//   - Left + Right: Returns the Left error from the value
//   - Left + Left: Returns Left(sg.Concat(function_error, value_error))
//
// Example:
//
//	import (
//	    "errors"
//	    "fmt"
//	    "strings"
//	    S "github.com/IBM/fp-go/v2/semigroup"
//	    "github.com/IBM/fp-go/v2/idiomatic/result"
//	)
//
//	// Create a semigroup that combines errors by concatenating their messages
//	errorSemigroup := S.MakeSemigroup(func(e1, e2 error) error {
//	    return fmt.Errorf("%v; %v", e1, e2)
//	})
//
//	// ApV with both function and value having errors
//	double := N.Mul(2)
//	apv := result.ApV[int, int](errorSemigroup)
//
//	value := result.Left[int](errors.New("invalid value"))
//	fn := result.Left[func(int) int](errors.New("invalid function"))
//
//	result := apv(value)(fn)
//	// Left(error: "invalid function; invalid value")
//
//	// ApV with successful application
//	goodValue, _ := result.Right(5)
//	goodFn, _ := result.Right(double)
//	result2 := apv(goodValue)(goodFn)
//	// Right(10)
func ApV[B, A any](sg S.Semigroup[error]) func(A, error) Operator[func(A) B, B] {
	return func(a A, aerr error) Operator[func(A) B, B] {
		return func(fab func(A) B, faberr error) (B, error) {
			// Both have errors: combine them using the semigroup
			if aerr != nil {
				if faberr != nil {
					return Left[B](sg.Concat(faberr, aerr))
				}
				// Only value has error
				return Left[B](aerr)
			}
			// Only function has error
			if faberr != nil {
				return Left[B](faberr)
			}
			// Both are successful: apply function to value
			return Of(fab(a))
		}
	}
}
