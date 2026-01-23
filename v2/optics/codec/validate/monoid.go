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
