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
	"github.com/IBM/fp-go/v2/monoid"
)

// ApplicativeMonoid creates a Monoid instance for Validate[I, A] given a Monoid[A].
//
// This function lifts a monoid operation on values of type A to work with validators
// that produce values of type A. It uses the applicative functor structure of the
// nested Reader types to combine validators while preserving their validation context.
//
// The resulting monoid allows you to:
//   - Combine multiple validators that produce monoidal values
//   - Run validators in parallel and merge their results using the monoid operation
//   - Build complex validators compositionally from simpler ones
//
// # Type Parameters
//
//   - I: The input type that validators accept
//   - A: The output type that validators produce (must have a Monoid instance)
//
// # Parameters
//
//   - m: A Monoid[A] that defines how to combine values of type A
//
// # Returns
//
// A Monoid[Validate[I, A]] that can combine validators using the applicative structure.
//
// # How It Works
//
// The function composes three layers of applicative monoids:
//  1. The innermost layer uses validation.ApplicativeMonoid(m) to combine Validation[A] values
//  2. The middle layer wraps this in reader.ApplicativeMonoid for the Context dependency
//  3. The outer layer wraps everything in reader.ApplicativeMonoid for the input I dependency
//
// This creates a monoid that:
//   - Takes the same input I for both validators
//   - Threads the same Context through both validators
//   - Combines successful results using the monoid operation on A
//   - Accumulates validation errors from both validators if either fails
//
// # Example
//
// Combining string validators using string concatenation:
//
//	import (
//	    "github.com/IBM/fp-go/v2/monoid"
//	    "github.com/IBM/fp-go/v2/string"
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	    "github.com/IBM/fp-go/v2/optics/codec/validation"
//	)
//
//	// Create a monoid for string validators
//	stringMonoid := string.Monoid
//	validatorMonoid := validate.ApplicativeMonoid[string, string](stringMonoid)
//
//	// Define two validators that extract different parts
//	validator1 := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success("Hello ")
//	    }
//	}
//
//	validator2 := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success("World")
//	    }
//	}
//
//	// Combine them - results will be concatenated
//	combined := validatorMonoid.Concat(validator1, validator2)
//	// When run, produces validation.Success("Hello World")
//
// Combining numeric validators using addition:
//
//	import (
//	    "github.com/IBM/fp-go/v2/number"
//	)
//
//	// Create a monoid for int validators using addition
//	intMonoid := number.MonoidSum[int]()
//	validatorMonoid := validate.ApplicativeMonoid[string, int](intMonoid)
//
//	// Validators that extract and validate different numeric fields
//	// Results will be summed together
//
// # Notes
//
//   - Both validators receive the same input value I
//   - If either validator fails, all errors are accumulated
//   - If both succeed, their results are combined using the monoid operation
//   - The empty element of the monoid serves as the identity for the Concat operation
//   - This follows the applicative functor laws for combining effectful computations
//
// # See Also
//
//   - validation.ApplicativeMonoid: The underlying monoid for validation results
//   - reader.ApplicativeMonoid: The monoid for reader computations
//   - Monoid[A]: The monoid instance for the result type
func ApplicativeMonoid[I, A any](m Monoid[A]) Monoid[Validate[I, A]] {
	return monoid.ApplicativeMonoid[A, Validate[I, A]](
		Of,
		MonadMap,
		MonadAp,
		m,
	)
}

// AlternativeMonoid creates a Monoid instance for Validate[I, A] that combines both
// applicative and alternative semantics.
//
// This function creates a monoid that:
//  1. When both validators succeed: Combines their results using the provided monoid operation
//  2. When one validator fails: Uses the successful validator's result (alternative behavior)
//  3. When both validators fail: Aggregates all errors from both validators
//
// This is a hybrid approach that combines:
//   - ApplicativeMonoid: Combines successful results using the monoid operation
//   - AltMonoid: Provides fallback behavior when validators fail
//
// # Type Parameters
//
//   - I: The input type that validators accept
//   - A: The output type that validators produce (must have a Monoid instance)
//
// # Parameters
//
//   - m: A Monoid[A] that defines how to combine values of type A
//
// # Returns
//
// A Monoid[Validate[I, A]] that combines validators using both applicative and alternative semantics.
//
// # Behavior Details
//
// The AlternativeMonoid differs from ApplicativeMonoid in how it handles mixed success/failure:
//
//   - **Both succeed**: Results are combined using the monoid operation (like ApplicativeMonoid)
//   - **First succeeds, second fails**: Returns the first result (alternative fallback)
//   - **First fails, second succeeds**: Returns the second result (alternative fallback)
//   - **Both fail**: Aggregates errors from both validators
//
// # Example: String Concatenation with Fallback
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	    "github.com/IBM/fp-go/v2/optics/codec/validation"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	m := validate.AlternativeMonoid[string, string](S.Monoid)
//
//	// Both succeed - results are concatenated
//	validator1 := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success("Hello")
//	    }
//	}
//	validator2 := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success(" World")
//	    }
//	}
//	combined := m.Concat(validator1, validator2)
//	result := combined("input")(nil)
//	// result is validation.Success("Hello World")
//
// # Example: Fallback Behavior
//
//	// First fails, second succeeds - uses second result
//	failing := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.FailureWithMessage[string](input, "first failed")(ctx)
//	    }
//	}
//	succeeding := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success("fallback")
//	    }
//	}
//	combined := m.Concat(failing, succeeding)
//	result := combined("input")(nil)
//	// result is validation.Success("fallback")
//
// # Example: Error Aggregation
//
//	// Both fail - errors are aggregated
//	failing1 := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.FailureWithMessage[string](input, "error 1")(ctx)
//	    }
//	}
//	failing2 := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.FailureWithMessage[string](input, "error 2")(ctx)
//	    }
//	}
//	combined := m.Concat(failing1, failing2)
//	result := combined("input")(nil)
//	// result contains both "error 1" and "error 2"
//
// # Comparison with Other Monoids
//
//   - **ApplicativeMonoid**: Always combines results when both succeed, fails if either fails
//   - **AlternativeMonoid**: Combines results when both succeed, provides fallback when one fails
//   - **AltMonoid**: Always uses first success, never combines results
//
// # Use Cases
//
//   - Validation with fallback strategies and result combination
//   - Building validators that accumulate results but provide alternatives
//   - Configuration loading with multiple sources and merging
//   - Data aggregation with error recovery
//
// # Notes
//
//   - Both validators receive the same input value I
//   - The empty element of the monoid serves as the identity for the Concat operation
//   - Error aggregation ensures no validation failures are lost
//   - This follows both applicative and alternative functor laws
//
// # See Also
//
//   - ApplicativeMonoid: For pure applicative combination without fallback
//   - AltMonoid: For pure alternative behavior without result combination
//   - MonadAlt: The underlying alternative operation
func AlternativeMonoid[I, A any](m Monoid[A]) Monoid[Validate[I, A]] {
	return monoid.AlternativeMonoid(
		Of[I, A],
		MonadMap[I, A, func(A) A],
		MonadAp[A, I, A],
		MonadAlt[I, A],
		m,
	)
}

// AltMonoid creates a Monoid instance for Validate[I, A] using alternative semantics
// with a provided zero/default validator.
//
// This function creates a monoid where:
//  1. The first successful validator wins (no result combination)
//  2. If the first fails, the second is tried as a fallback
//  3. If both fail, errors are aggregated
//  4. The provided zero validator serves as the identity element
//
// Unlike AlternativeMonoid, AltMonoid does NOT combine successful results - it always
// returns the first success. This makes it ideal for fallback chains and default values.
//
// # Type Parameters
//
//   - I: The input type that validators accept
//   - A: The output type that validators produce
//
// # Parameters
//
//   - zero: A lazy Validate[I, A] that serves as the identity element. This is typically
//     a validator that always succeeds with a default value, but can also be a failing
//     validator if no default is appropriate.
//
// # Returns
//
// A Monoid[Validate[I, A]] that combines validators using alternative semantics where
// the first success wins.
//
// # Behavior Details
//
// The AltMonoid implements a "first success wins" strategy:
//
//   - **First succeeds**: Returns the first result, second is never evaluated
//   - **First fails, second succeeds**: Returns the second result
//   - **Both fail**: Aggregates errors from both validators
//   - **Concat with Empty**: The zero validator is used as fallback
//
// # Example: Default Value Fallback
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/codec/validate"
//	)
//
//	// Create a monoid with a default value of 0
//	m := validate.AltMonoid(func() validate.Validate[string, int] {
//	    return validate.Of[string, int](0)
//	})
//
//	// First validator succeeds - returns 42, second is not evaluated
//	validator1 := validate.Of[string, int](42)
//	validator2 := validate.Of[string, int](100)
//	combined := m.Concat(validator1, validator2)
//	result := combined("input")(nil)
//	// result is validation.Success(42)
//
// # Example: Fallback Chain
//
//	// Try primary, then fallback, then default
//	m := validate.AltMonoid(func() validate.Validate[string, string] {
//	    return validate.Of[string, string]("default")
//	})
//
//	primary := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.FailureWithMessage[string](input, "primary failed")(ctx)
//	    }
//	}
//	secondary := func(input string) validate.Reader[validation.Context, validation.Validation[string]] {
//	    return func(ctx validation.Context) validation.Validation[string] {
//	        return validation.Success("secondary value")
//	    }
//	}
//
//	// Chain: try primary, then secondary, then default
//	combined := m.Concat(m.Concat(primary, secondary), m.Empty())
//	result := combined("input")(nil)
//	// result is validation.Success("secondary value")
//
// # Example: Error Aggregation
//
//	// Both fail - errors are aggregated
//	m := validate.AltMonoid(func() validate.Validate[string, int] {
//	    return func(input string) validate.Reader[validation.Context, validation.Validation[int]] {
//	        return func(ctx validation.Context) validation.Validation[int] {
//	            return validation.FailureWithMessage[int](input, "no default")(ctx)
//	        }
//	    }
//	})
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
//	combined := m.Concat(failing1, failing2)
//	result := combined("input")(nil)
//	// result contains both "error 1" and "error 2"
//
// # Comparison with Other Monoids
//
//   - **ApplicativeMonoid**: Combines results when both succeed using monoid operation
//   - **AlternativeMonoid**: Combines results when both succeed, provides fallback when one fails
//   - **AltMonoid**: First success wins, never combines results (pure alternative)
//
// # Use Cases
//
//   - Configuration loading with fallback sources (try file, then env, then default)
//   - Validation with default values
//   - Parser combinators with alternative branches
//   - Error recovery with multiple strategies
//
// # Notes
//
//   - The zero validator is lazily evaluated, only when needed
//   - First success short-circuits evaluation (second validator not called)
//   - Error aggregation ensures all validation failures are reported
//   - This follows the alternative functor laws
//
// # See Also
//
//   - AlternativeMonoid: For combining results when both succeed
//   - ApplicativeMonoid: For pure applicative combination
//   - MonadAlt: The underlying alternative operation
//   - Alt: The curried version for pipeline composition
func AltMonoid[I, A any](zero Lazy[Validate[I, A]]) Monoid[Validate[I, A]] {
	return monoid.AltMonoid(
		zero,
		MonadAlt[I, A],
	)
}
