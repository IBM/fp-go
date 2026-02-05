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

// Package validate provides functional validation primitives for building composable validators.
//
// This package implements a validation framework based on functional programming principles,
// allowing you to build complex validators from simple, composable pieces. It uses the
// Reader monad pattern to thread validation context through nested structures.
//
// # Core Concepts
//
// The validate package is built around several key types:
//
//   - Validate[I, A]: A validator that transforms input I to output A with validation context
//   - Validation[A]: The result of validation, either errors or a valid value A
//   - Context: Tracks the path through nested structures for detailed error messages
//
// # Type Structure
//
// A Validate[I, A] is defined as:
//
//	Reader[I, Decode[A]]]
//
// This means:
//  1. It takes an input of type I
//  2. Returns a Reader that depends on validation Context
//  3. That Reader produces a Validation[A] (Either[Errors, A])
//
// This layered structure allows validators to:
//   - Access the input value
//   - Track validation context (path in nested structures)
//   - Accumulate multiple validation errors
//   - Compose with other validators
//
// # Validation Context
//
// The Context type tracks the path through nested data structures during validation.
// Each ContextEntry contains:
//   - Key: The field name or map key
//   - Type: The expected type name
//   - Actual: The actual value being validated
//
// This provides detailed error messages like "at user.address.zipCode: expected string, got number".
//
// # Monoid Operations
//
// The package provides ApplicativeMonoid for combining validators using monoid operations.
// This allows you to:
//   - Combine multiple validators that produce monoidal values
//   - Accumulate results from parallel validations
//   - Build complex validators from simpler ones
//
// # Example Usage
//
// Basic validation structure:
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	    "github.com/IBM/fp-go/v2/optics/codec/validation"
//	)
//
//	// A validator that checks if a string is non-empty
//	func nonEmptyString(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    if input == "" {
//	        return validation.FailureWithMessage[string](input, "string must not be empty")
//	    }
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success(input)
//	    }
//	}
//
//	// Create a Validate function
//	var validateNonEmpty validate.Validate[string, string] = func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return nonEmptyString(input)
//	}
//
// Combining validators with monoids:
//
//	import (
//	    "github.com/IBM/fp-go/v2/monoid"
//	    "github.com/IBM/fp-go/v2/string"
//	)
//
//	// Combine string validators using string concatenation monoid
//	stringMonoid := string.Monoid
//	validatorMonoid := validate.ApplicativeMonoid[string, string](stringMonoid)
//
//	// Now you can combine validators that produce strings
//	combined := validatorMonoid.Concat(validator1, validator2)
//
// # Integration with Codec
//
// This package is designed to work with the optics/codec package for building
// type-safe encoders and decoders with validation. Validators can be composed
// into codecs that handle serialization, deserialization, and validation in a
// unified way.
//
// # Error Handling
//
// Validation errors are accumulated using the Either monad's applicative instance.
// This means:
//   - Multiple validation errors can be collected in a single pass
//   - Errors include full context path for debugging
//   - Errors can be formatted for logging or user display
//
// See the validation package for error types and formatting options.
package validate

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/optics/codec/decode"
	"github.com/IBM/fp-go/v2/reader"
)

// Of creates a Validate that always succeeds with the given value.
//
// This is the "pure" or "return" operation for the Validate monad. It lifts a plain
// value into the validation context without performing any actual validation.
//
// # Type Parameters
//
//   - I: The input type (not used, but required for type consistency)
//   - A: The type of the value to wrap
//
// # Parameters
//
//   - a: The value to wrap in a successful validation
//
// # Returns
//
// A Validate[I, A] that ignores its input and always returns a successful validation
// containing the value a.
//
// # Example
//
//	// Create a validator that always succeeds with value 42
//	alwaysValid := validate.Of[string, int](42)
//	result := alwaysValid("any input")(nil)
//	// result is validation.Success(42)
//
// # Notes
//
//   - This is useful for lifting pure values into the validation context
//   - The input type I is ignored; the validator succeeds regardless of input
//   - This satisfies the monad laws: Of is the left and right identity for Chain
func Of[I, A any](a A) Validate[I, A] {
	return reader.Of[I](decode.Of[Context](a))
}

// MonadMap applies a function to the successful result of a validation.
//
// This is the functor map operation for Validate. It transforms the success value
// without affecting the validation logic or error handling. If the validation fails,
// the function is not applied and errors are preserved.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the current validation result
//   - B: The type after applying the transformation
//
// # Parameters
//
//   - fa: The validator to transform
//   - f: The transformation function to apply to successful results
//
// # Returns
//
// A new Validate[I, B] that applies f to the result if validation succeeds.
//
// # Example
//
//	// Transform a string validator to uppercase
//	validateString := func(s string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success(s)
//	    }
//	}
//
//	upperValidator := validate.MonadMap(validateString, strings.ToUpper)
//	result := upperValidator("hello")(nil)
//	// result is validation.Success("HELLO")
//
// # Notes
//
//   - Preserves validation errors unchanged
//   - Only applies the function to successful validations
//   - Satisfies the functor laws: composition and identity
func MonadMap[I, A, B any](fa Validate[I, A], f func(A) B) Validate[I, B] {
	return readert.MonadMap[
		Validate[I, A],
		Validate[I, B]](
		decode.MonadMap,
		fa,
		f,
	)
}

// Map creates an operator that transforms validation results.
//
// This is the curried version of MonadMap, returning a function that can be applied
// to validators. It's useful for creating reusable transformation pipelines.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the current validation result
//   - B: The type after applying the transformation
//
// # Parameters
//
//   - f: The transformation function to apply to successful results
//
// # Returns
//
// An Operator[I, A, B] that transforms Validate[I, A] to Validate[I, B].
//
// # Example
//
//	// Create a reusable transformation
//	toUpper := validate.Map[string, string, string](strings.ToUpper)
//
//	// Apply it to different validators
//	validator1 := toUpper(someStringValidator)
//	validator2 := toUpper(anotherStringValidator)
//
// # Notes
//
//   - This is the point-free style version of MonadMap
//   - Useful for building transformation pipelines
//   - Can be composed with other operators
func Map[I, A, B any](f func(A) B) Operator[I, A, B] {
	return readert.Map[
		Validate[I, A],
		Validate[I, B]](
		decode.Map,
		f,
	)
}

// Chain sequences two validators, where the second depends on the result of the first.
//
// This is the monadic bind operation for Validate. It allows you to create validators
// that depend on the results of previous validations, enabling complex validation logic
// that builds on earlier results.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the first validation result
//   - B: The type of the second validation result
//
// # Parameters
//
//   - f: A Kleisli arrow that takes a value of type A and returns a Validate[I, B]
//
// # Returns
//
// An Operator[I, A, B] that sequences the validations.
//
// # Example
//
//	// First validate that a string is non-empty, then validate its length
//	validateNonEmpty := func(s string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        if s == "" {
//	            return validation.FailureWithMessage[string](s, "must not be empty")(ctx)
//	        }
//	        return validation.Success(s)
//	    }
//	}
//
//	validateLength := func(s string) validate.Validate[string, int] {
//	    return func(input string) validate.Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            if len(s) < 3 {
//	                return validation.FailureWithMessage[int](len(s), "too short")(ctx)
//	            }
//	            return validation.Success(len(s))
//	        }
//	    }
//	}
//
//	// Chain them together
//	chained := validate.Chain(validateLength)(validateNonEmpty)
//
// # Notes
//
//   - If the first validation fails, the second is not executed
//   - Errors from the first validation are preserved
//   - This enables dependent validation logic
//   - Satisfies the monad laws: associativity and identity
func Chain[I, A, B any](f Kleisli[I, A, B]) Operator[I, A, B] {
	return readert.Chain[Validate[I, A]](
		decode.Chain,
		f,
	)
}

// ChainLeft sequences a computation on the failure (Left) channel of a validation.
//
// This function operates on the error path of validation, allowing you to transform,
// enrich, or recover from validation failures. It's the dual of Chain - while Chain
// operates on success values, ChainLeft operates on error values.
//
// # Key Behavior
//
// **Critical difference from standard Either operations**: This validation-specific
// implementation **aggregates errors** using the Errors monoid. When the transformation
// function returns a failure, both the original errors AND the new errors are combined,
// ensuring comprehensive error reporting.
//
//  1. **Success Pass-Through**: If validation succeeds, the handler is never called and
//     the success value passes through unchanged.
//
//  2. **Error Recovery**: The handler can recover from failures by returning a successful
//     validation, converting Left to Right.
//
//  3. **Error Aggregation**: When the handler also returns a failure, both the original
//     errors and the new errors are combined using the Errors monoid.
//
//  4. **Input Access**: The handler returns a Validate[I, A] function, giving it access
//     to the original input value I for context-aware error handling.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the validation result
//
// # Parameters
//
//   - f: A Kleisli arrow that takes Errors and returns a Validate[I, A]. This function
//     is called only when validation fails, receiving the accumulated errors.
//
// # Returns
//
// An Operator[I, A, A] that transforms validators by handling their error cases.
//
// # Example: Error Recovery
//
//	// Validator that may fail
//	validatePositive := func(n int) Reader[validation.Context, validation.Validation[int]] {
//	    return func(ctx validation.Context) validation.Validation[int] {
//	        if n > 0 {
//	            return validation.Success(n)
//	        }
//	        return validation.FailureWithMessage[int](n, "must be positive")(ctx)
//	    }
//	}
//
//	// Recover from specific errors with a default value
//	withDefault := ChainLeft(func(errs Errors) Validate[int, int] {
//	    for _, err := range errs {
//	        if err.Messsage == "must be positive" {
//	            return Of[int](0) // recover with default
//	        }
//	    }
//	    return func(input int) Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            return either.Left[int](errs)
//	        }
//	    }
//	})
//
//	validator := withDefault(validatePositive)
//	result := validator(-5)(nil)
//	// Result: Success(0) - recovered from failure
//
// # Example: Error Context Addition
//
//	// Add contextual information to errors
//	addContext := ChainLeft(func(errs Errors) Validate[string, int] {
//	    return func(input string) Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            return either.Left[int](validation.Errors{
//	                {
//	                    Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
//	                    Messsage: "failed to validate user age",
//	                },
//	            })
//	        }
//	    }
//	})
//
//	validator := addContext(someValidator)
//	// Errors will include both original error and context
//
// # Example: Input-Dependent Recovery
//
//	// Recover with different defaults based on input
//	smartDefault := ChainLeft(func(errs Errors) Validate[string, int] {
//	    return func(input string) Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            // Use input to determine appropriate default
//	            if strings.Contains(input, "http") {
//	                return validation.Of(80)
//	            }
//	            if strings.Contains(input, "https") {
//	                return validation.Of(443)
//	            }
//	            return validation.Of(8080)
//	        }
//	    }
//	})
//
// # Notes
//
//   - Errors are accumulated, not replaced - this ensures no validation failures are lost
//   - The handler has access to both the errors and the original input
//   - Success values bypass the handler completely
//   - This enables sophisticated error handling strategies including recovery, enrichment, and transformation
//   - Use OrElse as a semantic alias when emphasizing fallback/alternative logic
func ChainLeft[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return readert.Chain[Validate[I, A]](
		decode.ChainLeft,
		f,
	)
}

// MonadChainLeft sequences a computation on the failure (Left) channel of a validation.
//
// This is the direct application version of ChainLeft. It operates on the error path
// of validation, allowing you to transform, enrich, or recover from validation failures.
// It's the dual of Chain - while Chain operates on success values, MonadChainLeft
// operates on error values.
//
// # Key Behavior
//
// **Critical difference from standard Either operations**: This validation-specific
// implementation **aggregates errors** using the Errors monoid. When the transformation
// function returns a failure, both the original errors AND the new errors are combined,
// ensuring comprehensive error reporting.
//
//  1. **Success Pass-Through**: If validation succeeds, the handler is never called and
//     the success value passes through unchanged.
//
//  2. **Error Recovery**: The handler can recover from failures by returning a successful
//     validation, converting Left to Right.
//
//  3. **Error Aggregation**: When the handler also returns a failure, both the original
//     errors and the new errors are combined using the Errors monoid.
//
//  4. **Input Access**: The handler returns a Validate[I, A] function, giving it access
//     to the original input value I for context-aware error handling.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the validation result
//
// # Parameters
//
//   - fa: The Validate[I, A] to transform
//   - f: A Kleisli arrow that takes Errors and returns a Validate[I, A]. This function
//     is called only when validation fails, receiving the accumulated errors.
//
// # Returns
//
// A Validate[I, A] that handles error cases according to the provided function.
//
// # Example: Error Recovery
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	    "github.com/IBM/fp-go/v2/optics/codec/validation"
//	)
//
//	// Validator that may fail
//	validatePositive := func(n int) validate.Reader[validation.Context, validation.Validation[int]] {
//	    return func(ctx validation.Context) validation.Validation[int] {
//	        if n > 0 {
//	            return validation.Success(n)
//	        }
//	        return validation.FailureWithMessage[int](n, "must be positive")(ctx)
//	    }
//	}
//
//	// Recover from specific errors with a default value
//	withDefault := func(errs validation.Errors) validate.Validate[int, int] {
//	    for _, err := range errs {
//	        if err.Messsage == "must be positive" {
//	            return validate.Of[int](0) // recover with default
//	        }
//	    }
//	    // Propagate other errors
//	    return func(input int) validate.Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            return either.Left[int](errs)
//	        }
//	    }
//	}
//
//	validator := validate.MonadChainLeft(validatePositive, withDefault)
//	result := validator(-5)(nil)
//	// Result: Success(0) - recovered from failure
//
// # Example: Error Context Addition
//
//	// Add contextual information to errors
//	addContext := func(errs validation.Errors) validate.Validate[string, int] {
//	    return func(input string) validate.Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            // Add context error (will be aggregated with original)
//	            return either.Left[int](validation.Errors{
//	                {
//	                    Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
//	                    Messsage: "failed to validate user age",
//	                },
//	            })
//	        }
//	    }
//	}
//
//	validator := validate.MonadChainLeft(someValidator, addContext)
//	// Errors will include both original error and context
//
// # Example: Input-Dependent Recovery
//
//	// Recover with different defaults based on input
//	smartDefault := func(errs validation.Errors) validate.Validate[string, int] {
//	    return func(input string) validate.Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            // Use input to determine appropriate default
//	            if strings.Contains(input, "http:") {
//	                return validation.Success(80)
//	            }
//	            if strings.Contains(input, "https:") {
//	                return validation.Success(443)
//	            }
//	            return validation.Success(8080)
//	        }
//	    }
//	}
//
//	validator := validate.MonadChainLeft(parsePort, smartDefault)
//
// # Notes
//
//   - Errors are accumulated, not replaced - this ensures no validation failures are lost
//   - The handler has access to both the errors and the original input
//   - Success values bypass the handler completely
//   - This is the direct application version of ChainLeft
//   - This enables sophisticated error handling strategies including recovery, enrichment, and transformation
//
// # See Also
//
//   - ChainLeft: The curried, point-free version
//   - OrElse: Semantic alias for ChainLeft emphasizing fallback logic
//   - MonadAlt: Simplified alternative that ignores error details
//   - Alt: Curried version of MonadAlt
func MonadChainLeft[I, A any](fa Validate[I, A], f Kleisli[I, Errors, A]) Validate[I, A] {
	return readert.MonadChain(
		decode.MonadChainLeft,
		fa,
		f,
	)
}

// OrElse provides an alternative validation when the primary validation fails.
//
// This is a semantic alias for ChainLeft with identical behavior. The name "OrElse"
// emphasizes the intent of providing fallback or alternative validation logic, making
// code more readable when that's the primary use case.
//
// # Relationship to ChainLeft
//
// **OrElse and ChainLeft are functionally identical** - they produce exactly the same
// results for all inputs. The choice between them is purely about code readability:
//
//   - Use **OrElse** when emphasizing fallback/alternative validation logic
//   - Use **ChainLeft** when emphasizing technical error channel transformation
//
// Both maintain the critical property of **error aggregation**, ensuring all validation
// failures are preserved and reported together.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the validation result
//
// # Parameters
//
//   - f: A Kleisli arrow that takes Errors and returns a Validate[I, A]. This function
//     is called only when validation fails, receiving the accumulated errors.
//
// # Returns
//
// An Operator[I, A, A] that transforms validators by providing alternative validation.
//
// # Example: Fallback Validation
//
//	// Primary validator that may fail
//	validateFromConfig := func(key string) Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        // Try to get value from config
//	        if value, ok := config[key]; ok {
//	            return validation.Success(value)
//	        }
//	        return validation.FailureWithMessage[string](key, "not found in config")(ctx)
//	    }
//	}
//
//	// Use OrElse for semantic clarity - "try config, or else use environment"
//	withEnvFallback := OrElse(func(errs Errors) Validate[string, string] {
//	    return func(key string) Reader[validation.Context, validation.Validation[string]] {
//	        return func(ctx validation.Context) validation.Validation[string] {
//	            if value := os.Getenv(key); value != "" {
//	                return validation.Success(value)
//	            }
//	            return either.Left[string](errs) // propagate original errors
//	        }
//	    }
//	})
//
//	validator := withEnvFallback(validateFromConfig)
//	result := validator("DATABASE_URL")(nil)
//	// Tries config first, falls back to environment variable
//
// # Example: Default Value on Failure
//
//	// Provide a default value when validation fails
//	withDefault := OrElse(func(errs Errors) Validate[int, int] {
//	    return Of[int](0) // default to 0 on any failure
//	})
//
//	validator := withDefault(someValidator)
//	result := validator(input)(nil)
//	// Always succeeds, using default value if validation fails
//
// # Example: Pipeline with Multiple Fallbacks
//
//	// Build a validation pipeline with multiple fallback strategies
//	validator := F.Pipe2(
//	    validateFromDatabase,
//	    OrElse(func(errs Errors) Validate[string, Config] {
//	        // Try cache as first fallback
//	        return validateFromCache
//	    }),
//	    OrElse(func(errs Errors) Validate[string, Config] {
//	        // Use default config as final fallback
//	        return Of[string](defaultConfig)
//	    }),
//	)
//	// Tries database, then cache, then default
//
// # Notes
//
//   - Identical behavior to ChainLeft - they are aliases
//   - Errors are accumulated when transformations fail
//   - Success values pass through unchanged
//   - The handler has access to both errors and original input
//   - Choose OrElse for better readability when providing alternatives
//   - See ChainLeft documentation for detailed behavior and additional examples
func OrElse[I, A any](f Kleisli[I, Errors, A]) Operator[I, A, A] {
	return ChainLeft(f)
}

// MonadAp applies a validator containing a function to a validator containing a value.
//
// This is the applicative apply operation for Validate. It allows you to apply
// functions wrapped in validation context to values wrapped in validation context,
// accumulating errors from both if either fails.
//
// # Type Parameters
//
//   - B: The result type after applying the function
//   - I: The input type
//   - A: The type of the value to which the function is applied
//
// # Parameters
//
//   - fab: A validator that produces a function from A to B
//   - fa: A validator that produces a value of type A
//
// # Returns
//
// A Validate[I, B] that applies the function to the value if both validations succeed.
//
// # Example
//
//	// Create a validator that produces a function
//	validateFunc := validate.Of[string, func(int) int](func(x int) int { return x * 2 })
//
//	// Create a validator that produces a value
//	validateValue := validate.Of[string, int](21)
//
//	// Apply them
//	result := validate.MonadAp(validateFunc, validateValue)
//	// When run, produces validation.Success(42)
//
// # Notes
//
//   - Both validators receive the same input
//   - If either validation fails, all errors are accumulated
//   - If both succeed, the function is applied to the value
//   - This enables parallel validation with error accumulation
//   - Satisfies the applicative functor laws
func MonadAp[B, I, A any](fab Validate[I, func(A) B], fa Validate[I, A]) Validate[I, B] {
	return readert.MonadAp[
		Validate[I, A],
		Validate[I, B],
		Validate[I, func(A) B], I, A](
		decode.MonadAp[B, Context, A],
		fab,
		fa,
	)
}

// Ap creates an operator that applies a function validator to a value validator.
//
// This is the curried version of MonadAp, returning a function that can be applied
// to function validators. It's useful for creating reusable applicative patterns.
//
// # Type Parameters
//
//   - B: The result type after applying the function
//   - I: The input type
//   - A: The type of the value to which the function is applied
//
// # Parameters
//
//   - fa: A validator that produces a value of type A
//
// # Returns
//
// An Operator[I, func(A) B, B] that applies function validators to the value validator.
//
// # Example
//
//	// Create a value validator
//	validateValue := validate.Of[string, int](21)
//
//	// Create an applicative operator
//	applyTo21 := validate.Ap[int, string, int](validateValue)
//
//	// Create a function validator
//	validateDouble := validate.Of[string, func(int) int](func(x int) int { return x * 2 })
//
//	// Apply it
//	result := applyTo21(validateDouble)
//	// When run, produces validation.Success(42)
//
// # Notes
//
//   - This is the point-free style version of MonadAp
//   - Useful for building applicative pipelines
//   - Enables parallel validation with error accumulation
//   - Can be composed with other applicative operators
func Ap[B, I, A any](fa Validate[I, A]) Operator[I, func(A) B, B] {
	return readert.Ap[
		Validate[I, A],
		Validate[I, B],
		Validate[I, func(A) B], I, A](
		decode.Ap[B, Context, A],
		fa,
	)
}

// Alt provides an alternative validator when the primary validator fails.
//
// This is the curried, point-free version of MonadAlt. It creates an operator that
// transforms a validator by adding a fallback alternative. When the first validator
// fails, the second (lazily evaluated) validator is tried. If both fail, errors are
// aggregated.
//
// Alt implements the Alternative typeclass pattern, providing a way to express
// "try this, or else try that" logic in a composable way.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the validation result
//
// # Parameters
//
//   - second: A lazy Validate[I, A] that serves as the fallback. It's only evaluated
//     if the first validator fails.
//
// # Returns
//
// An Operator[I, A, A] that transforms validators by adding alternative fallback logic.
//
// # Behavior
//
//   - **First succeeds**: Returns the first result, second is never evaluated
//   - **First fails, second succeeds**: Returns the second result
//   - **Both fail**: Aggregates errors from both validators
//
// # Example: Fallback Validation
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	    "github.com/IBM/fp-go/v2/optics/codec/validation"
//	)
//
//	// Primary validator that may fail
//	validateFromConfig := func(key string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        // Try to get value from config
//	        if value, ok := config[key]; ok {
//	            return validation.Success(value)
//	        }
//	        return validation.FailureWithMessage[string](key, "not in config")(ctx)
//	    }
//	}
//
//	// Fallback to environment variable
//	validateFromEnv := func(key string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        if value := os.Getenv(key); value != "" {
//	            return validation.Success(value)
//	        }
//	        return validation.FailureWithMessage[string](key, "not in env")(ctx)
//	    }
//	}
//
//	// Use Alt to add fallback - point-free style
//	withFallback := validate.Alt(func() validate.Validate[string, string] {
//	    return validateFromEnv
//	})
//
//	validator := withFallback(validateFromConfig)
//	result := validator("DATABASE_URL")(nil)
//	// Tries config first, falls back to environment variable
//
// # Example: Pipeline with Multiple Alternatives
//
//	// Chain multiple alternatives using function composition
//	validator := F.Pipe2(
//	    validateFromDatabase,
//	    validate.Alt(func() validate.Validate[string, Config] {
//	        return validateFromCache
//	    }),
//	    validate.Alt(func() validate.Validate[string, Config] {
//	        return validate.Of[string](defaultConfig)
//	    }),
//	)
//	// Tries database, then cache, then default
//
// # Notes
//
//   - The second validator is lazily evaluated for efficiency
//   - First success short-circuits evaluation
//   - Errors are aggregated when both fail
//   - This is the point-free version of MonadAlt
//   - Useful for building validation pipelines with F.Pipe
//
// # See Also
//
//   - MonadAlt: The direct application version
//   - ChainLeft: The more general error transformation operator
//   - OrElse: Semantic alias for ChainLeft
//   - AltMonoid: For combining multiple alternatives with monoid structure
func Alt[I, A any](second Lazy[Validate[I, A]]) Operator[I, A, A] {
	return ChainLeft(function.Ignore1of1[Errors](second))
}

// MonadAlt provides an alternative validator when the primary validator fails.
//
// This is the direct application version of Alt. It takes two validators and returns
// a new validator that tries the first, and if it fails, tries the second. If both
// fail, errors from both are aggregated.
//
// MonadAlt implements the Alternative typeclass pattern, enabling "try this, or else
// try that" logic with comprehensive error reporting.
//
// # Type Parameters
//
//   - I: The input type
//   - A: The type of the validation result
//
// # Parameters
//
//   - first: The primary Validate[I, A] to try first
//   - second: A lazy Validate[I, A] that serves as the fallback. It's only evaluated
//     if the first validator fails.
//
// # Returns
//
// A Validate[I, A] that tries the first validator, falling back to the second if needed.
//
// # Behavior
//
//   - **First succeeds**: Returns the first result, second is never evaluated
//   - **First fails, second succeeds**: Returns the second result
//   - **Both fail**: Aggregates errors from both validators
//
// # Example: Configuration with Fallback
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	    "github.com/IBM/fp-go/v2/optics/codec/validation"
//	)
//
//	// Primary validator
//	validateFromConfig := func(key string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        if value, ok := config[key]; ok {
//	            return validation.Success(value)
//	        }
//	        return validation.FailureWithMessage[string](key, "not in config")(ctx)
//	    }
//	}
//
//	// Fallback validator
//	validateFromEnv := func(key string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        if value := os.Getenv(key); value != "" {
//	            return validation.Success(value)
//	        }
//	        return validation.FailureWithMessage[string](key, "not in env")(ctx)
//	    }
//	}
//
//	// Combine with MonadAlt
//	validator := validate.MonadAlt(
//	    validateFromConfig,
//	    func() validate.Validate[string, string] { return validateFromEnv },
//	)
//	result := validator("DATABASE_URL")(nil)
//	// Tries config first, falls back to environment variable
//
// # Example: Multiple Fallbacks
//
//	// Chain multiple alternatives
//	validator := validate.MonadAlt(
//	    validate.MonadAlt(
//	        validateFromDatabase,
//	        func() validate.Validate[string, Config] { return validateFromCache },
//	    ),
//	    func() validate.Validate[string, Config] { return validate.Of[string](defaultConfig) },
//	)
//	// Tries database, then cache, then default
//
// # Example: Error Aggregation
//
//	failing1 := func(input string) validate.Reader[validation.Context, validation.Validation[int]] {
//	    return func(ctx validation.Context) validation.Validation[int] {
//	        return validation.FailureWithMessage[int](input, "error 1")(ctx)
//	    }
//	}
//	failing2 := func(input string) validate.Reader[validation.Context, validation.Validation[int]] {
//	    return func(ctx validation.Context) validation.Validation[int] {
//	        return validation.FailureWithMessage[int](input, "error 2")(ctx)
//	    }
//	}
//
//	validator := validate.MonadAlt(
//	    failing1,
//	    func() validate.Validate[string, int] { return failing2 },
//	)
//	result := validator("input")(nil)
//	// result contains both "error 1" and "error 2"
//
// # Notes
//
//   - The second validator is lazily evaluated for efficiency
//   - First success short-circuits evaluation (second not called)
//   - Errors are aggregated when both fail
//   - This is equivalent to Alt but with direct application
//   - Both validators receive the same input value
//
// # See Also
//
//   - Alt: The curried, point-free version
//   - MonadChainLeft: The underlying error transformation operation
//   - OrElse: Semantic alias for ChainLeft
//   - AltMonoid: For combining multiple alternatives with monoid structure
func MonadAlt[I, A any](first Validate[I, A], second Lazy[Validate[I, A]]) Validate[I, A] {
	return MonadChainLeft(first, function.Ignore1of1[Errors](second))
}
