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

package validate

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerresult"
)

// FromReaderResult converts a ReaderResult into a Validate.
//
// This function bridges the gap between simple error-based validation (ReaderResult)
// and the more sophisticated validation framework that supports error accumulation
// and detailed context tracking.
//
// # Type Parameters
//
//   - I: The input type that the validator will receive
//   - A: The output type that the validator will produce on success
//
// # Parameters
//
//   - r: A ReaderResult[I, A] which is readerresult.ReaderResult[I, A]
//     This represents a computation that:
//     1. Takes an input of type I
//     2. Returns Either[error, A] (success with A or failure with error)
//
// # Returns
//
//   - Validate[I, A]: A validator that:
//     1. Takes an input of type I
//     2. Takes a validation Context (path through nested structures)
//     3. Returns Validation[A] (Either[Errors, A])
//
// # Behavior
//
// The conversion follows this logic:
//
//  1. Success case: If the ReaderResult succeeds with value A:
//     - Wraps the value in validation.Success[A]
//     - Returns a validator that always succeeds with that value
//
//  2. Failure case: If the ReaderResult fails with an error:
//     - Creates a validation.ValidationError with:
//     - The input value that caused the failure
//     - The current validation context (path information)
//     - A generic message "unable to decode"
//     - The original error as the cause
//     - Returns a validator that fails with this detailed error
//
// # Error Handling
//
// The function enhances simple error handling by:
//   - Converting a single error into a structured validation.ValidationError
//   - Preserving the original error as the cause (accessible via Unwrap())
//   - Adding context information about where the error occurred
//   - Making the error compatible with the validation framework's error accumulation
//
// # Example Usage
//
// Basic conversion:
//
//	// A simple ReaderResult that parses an integer
//	parseIntRR := result.Eitherize1(strconv.Atoi)
//
//	// Convert to Validate
//	validateInt := FromReaderResult[string, int](parseIntRR)
//
//	// Use the validator
//	result := validateInt("42")(nil)  // Success(42)
//	result := validateInt("abc")(nil) // Failure with ValidationError
//
// Integration with validation pipeline:
//
//	// Combine with other validators
//	validatePositiveInt := F.Pipe1(
//	    FromReaderResult[string, int](parseIntRR),
//	    Chain(func(n int) Validate[string, int] {
//	        if n > 0 {
//	            return Of[string](n)
//	        }
//	        return func(input string) Reader[Context, Validation[int]] {
//	            return validation.FailureWithMessage[int](n, "must be positive")
//	        }
//	    }),
//	)
//
// # Implementation Details
//
// The function uses a functional composition approach:
//
//  1. readerresult.Map: Transforms successful results
//     - Wraps the success value in validation.Success
//     - Lifts it into a Reader context with reader.Of
//
//  2. readerresult.GetOrElse: Handles failures
//     - Uses reader.Asks to access the validation context
//     - Creates a validation.ValidationError with validation.FailureWithError
//     - Uses reader.Local to adapt the context type
//
// # See Also
//
//   - Validate: The target validation type
//   - ReaderResult: The source type
//   - validation.Success: Creates successful validations
//   - validation.FailureWithError: Creates validation failures with cause
//   - Context: Validation context for error reporting
func FromReaderResult[I, A any](r ReaderResult[I, A]) Validate[I, A] {
	return F.Pipe2(
		r,
		readerresult.Map[I](F.Flow2(
			validation.Success[A],
			reader.Of[Context],
		)),
		readerresult.GetOrElse(F.Pipe1(
			reader.Asks(F.Flip(F.Bind2nd(validation.FailureWithError[A], "unable to decode"))),
			reader.Map[error](reader.Local[Decode[Context, A]](F.ToAny[I])),
		)),
	)
}
