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

package validation

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Result represents a computation that may succeed with a value of type A or fail with an error.
	// This is an alias for result.Result[A], which is Either[error, A].
	//
	// Used for converting validation results to standard Go error handling patterns.
	Result[A any] = result.Result[A]

	// Either represents a value that can be one of two types: Left (error) or Right (success).
	// This is an alias for either.Either[E, A], a disjoint union type.
	//
	// In the validation context:
	//   - Left[E]: Contains error information of type E
	//   - Right[A]: Contains a successfully validated value of type A
	Either[E, A any] = either.Either[E, A]

	// ContextEntry represents a single entry in the validation context path.
	// It tracks the location and type information during nested validation,
	// enabling precise error reporting with full path information.
	//
	// Fields:
	//   - Key: The key or field name (e.g., "email", "address", "items[0]")
	//   - Type: The expected type name (e.g., "string", "int", "User")
	//   - Actual: The actual value being validated (for error reporting)
	//
	// Example:
	//
	//	entry := ContextEntry{
	//	    Key:    "user.email",
	//	    Type:   "string",
	//	    Actual: 12345,
	//	}
	ContextEntry struct {
		Key    string // The key or field name (for objects/maps)
		Type   string // The expected type name
		Actual any    // The actual value being validated
	}

	// Context is a stack of ContextEntry values representing the path through
	// nested structures during validation. Used to provide detailed error messages
	// that show exactly where in a nested structure a validation failure occurred.
	//
	// The context builds up as validation descends into nested structures:
	//   - [] - root level
	//   - [{Key: "user"}] - inside user object
	//   - [{Key: "user"}, {Key: "address"}] - inside user.address
	//   - [{Key: "user"}, {Key: "address"}, {Key: "zipCode"}] - at user.address.zipCode
	//
	// Example:
	//
	//	ctx := Context{
	//	    {Key: "user", Type: "User"},
	//	    {Key: "address", Type: "Address"},
	//	    {Key: "zipCode", Type: "string"},
	//	}
	//	// Represents path: user.address.zipCode
	Context = []ContextEntry

	// ValidationError represents a single validation failure with full context information.
	// It implements the error interface and provides detailed information about what failed,
	// where it failed, and why it failed.
	//
	// Fields:
	//   - Value: The actual value that failed validation
	//   - Context: The path to the value in nested structures (e.g., user.address.zipCode)
	//   - Messsage: Human-readable error description
	//   - Cause: Optional underlying error that caused the validation failure
	//
	// The ValidationError type implements:
	//   - error interface: For standard Go error handling
	//   - fmt.Formatter: For custom formatting with %v, %+v
	//   - slog.LogValuer: For structured logging with slog
	//
	// Example:
	//
	//	err := &ValidationError{
	//	    Value:    "not-an-email",
	//	    Context:  []ContextEntry{{Key: "user"}, {Key: "email"}},
	//	    Messsage: "invalid email format",
	//	    Cause:    nil,
	//	}
	//	fmt.Printf("%v", err)   // at user.email: invalid email format
	//	fmt.Printf("%+v", err)  // at user.email: invalid email format
	//	                        //   value: "not-an-email"
	ValidationError struct {
		Value    any     // The value that failed validation
		Context  Context // The path to the value in nested structures
		Messsage string  // Human-readable error message
		Cause    error   // Optional underlying error cause
	}

	// Errors is a collection of validation errors.
	// This type is used to accumulate multiple validation failures,
	// allowing all errors to be reported at once rather than failing fast.
	//
	// Example:
	//
	//	errors := Errors{
	//	    &ValidationError{Value: "", Messsage: "name is required"},
	//	    &ValidationError{Value: "invalid", Messsage: "invalid email"},
	//	    &ValidationError{Value: -1, Messsage: "age must be positive"},
	//	}
	Errors = []*ValidationError

	// validationErrors wraps a collection of validation errors with an optional root cause.
	// It provides structured error information for validation failures and implements
	// the error interface for integration with standard Go error handling.
	//
	// This type is internal and created via MakeValidationErrors.
	// It implements:
	//   - error interface: For standard Go error handling
	//   - fmt.Formatter: For custom formatting with %v, %+v
	//   - slog.LogValuer: For structured logging with slog
	//
	// Fields:
	//   - errors: The collection of individual validation errors
	//   - cause: Optional root cause error
	validationErrors struct {
		errors Errors
		cause  error
	}

	// Validation represents the result of a validation operation.
	// It's an Either type where:
	//   - Left(Errors): Validation failed with one or more errors
	//   - Right(A): Successfully validated value of type A
	//
	// This type supports applicative operations, allowing multiple validations
	// to be combined while accumulating all errors rather than failing fast.
	//
	// Example:
	//
	//	// Success case
	//	valid := Success(42)  // Right(42)
	//
	//	// Failure case
	//	invalid := Failures[int](Errors{
	//	    &ValidationError{Messsage: "must be positive"},
	//	})  // Left([...])
	//
	//	// Combining validations (accumulates all errors)
	//	result := Ap(Ap(Of(func(x int) func(y int) int {
	//	    return func(y int) int { return x + y }
	//	}))(validateX))(validateY)
	Validation[A any] = Either[Errors, A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	// This is an alias for reader.Reader[R, A], which is func(R) A.
	//
	// In the validation context, Reader is used for context-dependent validation operations
	// where the validation logic needs access to the current validation context path.
	//
	// Example:
	//
	//	validateWithContext := func(ctx Context) Validation[int] {
	//	    // Use ctx to provide detailed error messages
	//	    return Success(42)
	//	}
	Reader[R, A any] = reader.Reader[R, A]

	// Kleisli represents a function from A to a validated B.
	// It's a Reader that takes an input A and produces a Validation[B].
	// This is the fundamental building block for composable validation operations.
	//
	// Type: func(A) Validation[B]
	//
	// Kleisli arrows can be composed using Chain/Bind operations to build
	// complex validation pipelines from simple validation functions.
	//
	// Example:
	//
	//	validatePositive := func(x int) Validation[int] {
	//	    if x > 0 {
	//	        return Success(x)
	//	    }
	//	    return Failures[int](/* error */)
	//	}
	//
	//	validateEven := func(x int) Validation[int] {
	//	    if x%2 == 0 {
	//	        return Success(x)
	//	    }
	//	    return Failures[int](/* error */)
	//	}
	//
	//	// Compose validations
	//	validatePositiveEven := Chain(validateEven)(Success(42))
	Kleisli[A, B any] = Reader[A, Validation[B]]

	// Operator represents a validation transformation that takes a validated A and produces a validated B.
	// It's a specialized Kleisli arrow for composing validation operations where the input
	// is already a Validation[A].
	//
	// Type: func(Validation[A]) Validation[B]
	//
	// Operators are used to transform and compose validation results, enabling
	// functional composition of validation pipelines.
	//
	// Example:
	//
	//	// Transform a validated int to a validated string
	//	intToString := Map(func(x int) string {
	//	    return strconv.Itoa(x)
	//	})  // Operator[int, string]
	//
	//	result := intToString(Success(42))  // Success("42")
	Operator[A, B any] = Kleisli[Validation[A], B]

	// Monoid represents an algebraic structure with an associative binary operation and an identity element.
	// This is an alias for monoid.Monoid[A].
	//
	// In the validation context, monoids are used to combine validation results:
	//   - ApplicativeMonoid: Combines successful validations using the monoid operation
	//   - AlternativeMonoid: Provides fallback behavior for failed validations
	//
	// Example:
	//
	//	import N "github.com/IBM/fp-go/v2/number"
	//
	//	intAdd := N.MonoidSum[int]()
	//	m := ApplicativeMonoid(intAdd)
	//	result := m.Concat(Success(5), Success(3))  // Success(8)
	Monoid[A any] = monoid.Monoid[A]

	// Endomorphism represents a function from a type to itself: func(A) A.
	// This is an alias for endomorphism.Endomorphism[A].
	//
	// In the validation context, endomorphisms are used with LetL to transform
	// values within a validation context using pure functions.
	//
	// Example:
	//
	//	double := func(x int) int { return x * 2 }  // Endomorphism[int]
	//	result := LetL(lens, double)(Success(21))   // Success(42)
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Lazy[A any] = lazy.Lazy[A]
)
