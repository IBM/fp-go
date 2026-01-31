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
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/codec/decode"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
)

type (

	// Monoid represents an algebraic structure with an associative binary operation
	// and an identity element. Used for combining values of type A.
	//
	// A Monoid[A] must satisfy:
	//   - Associativity: Concat(Concat(a, b), c) == Concat(a, Concat(b, c))
	//   - Identity: Concat(Empty(), a) == a == Concat(a, Empty())
	//
	// Common examples:
	//   - Numbers with addition (identity: 0)
	//   - Numbers with multiplication (identity: 1)
	//   - Strings with concatenation (identity: "")
	//   - Lists with concatenation (identity: [])
	Monoid[A any] = monoid.Monoid[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	//
	// Reader[R, A] is a function type: func(R) A
	//
	// The Reader pattern is used to:
	//   - Thread configuration or context through computations
	//   - Implement dependency injection in a functional way
	//   - Defer computation until the environment is available
	//   - Compose computations that share the same environment
	//
	// Example:
	//   type Config struct { Port int }
	//   getPort := func(cfg Config) int { return cfg.Port }
	//   // getPort is a Reader[Config, int]
	Reader[R, A any] = reader.Reader[R, A]

	// Validation represents the result of a validation operation that may contain
	// validation errors or a successfully validated value of type A.
	//
	// Validation[A] is an Either[Errors, A], where:
	//   - Left(errors): Validation failed with one or more errors
	//   - Right(value): Validation succeeded with value of type A
	//
	// The Validation type supports:
	//   - Error accumulation: Multiple validation errors can be collected
	//   - Applicative composition: Parallel validations with error aggregation
	//   - Monadic composition: Sequential validations with short-circuiting
	//
	// Example:
	//   success := validation.Success(42)           // Right(42)
	//   failure := validation.Failure[int](errors)  // Left(errors)
	Validation[A any] = validation.Validation[A]

	// Context provides contextual information for validation operations,
	// tracking the path through nested data structures.
	//
	// Context is a slice of ContextEntry values, where each entry represents
	// a level in the nested structure being validated. This enables detailed
	// error messages that show exactly where validation failed.
	//
	// Example context path for nested validation:
	//   Context{
	//     {Key: "user", Type: "User"},
	//     {Key: "address", Type: "Address"},
	//     {Key: "zipCode", Type: "string"},
	//   }
	//   // Represents: user.address.zipCode
	//
	// The context is used to generate error messages like:
	//   "at user.address.zipCode: expected string, got number"
	Context = validation.Context

	// Decode represents a decoding operation that transforms input I into output A
	// within a validation context.
	//
	// Type structure:
	//   Decode[I, A] = Reader[Context, Validation[A]]
	//
	// This means:
	//   1. Takes a validation Context (path through nested structures)
	//   2. Returns a Validation[A] (Either[Errors, A])
	//
	// Decode is used as the foundation for validation operations, providing:
	//   - Context-aware error reporting with detailed paths
	//   - Error accumulation across multiple validations
	//   - Composable validation logic
	//
	// The Decode type is typically not used directly but through the Validate type,
	// which adds an additional Reader layer for accessing the input value.
	//
	// Example:
	//   decoder := func(ctx Context) Validation[int] {
	//     // Perform validation and return result
	//     return validation.Success(42)
	//   }
	//   // decoder is a Decode[any, int]
	Decode[I, A any] = decode.Decode[I, A]

	// Validate is a function that validates input I to produce type A with full context tracking.
	//
	// Type structure:
	//   Validate[I, A] = Reader[I, Decode[Context, A]]
	//
	// This means:
	//  1. Takes an input of type I
	//  2. Returns a Reader that depends on validation Context
	//  3. That Reader produces a Validation[A] (Either[Errors, A])
	//
	// The layered structure enables:
	//   - Access to the input value being validated
	//   - Context tracking through nested structures
	//   - Error accumulation with detailed paths
	//   - Composition with other validators
	//
	// Example usage:
	//   validatePositive := func(n int) Reader[Context, Validation[int]] {
	//     return func(ctx Context) Validation[int] {
	//       if n > 0 {
	//         return validation.Success(n)
	//       }
	//       return validation.FailureWithMessage[int](n, "must be positive")(ctx)
	//     }
	//   }
	//   // validatePositive is a Validate[int, int]
	//
	// The Validate type forms:
	//   - A Functor: Can map over successful results
	//   - An Applicative: Can combine validators in parallel
	//   - A Monad: Can chain dependent validations
	Validate[I, A any] = Reader[I, Decode[Context, A]]

	// Errors is a collection of validation errors that occurred during validation.
	//
	// Each error in the collection contains:
	//   - The value that failed validation
	//   - The context path where the error occurred
	//   - A human-readable error message
	//   - An optional underlying cause error
	//
	// Errors can be accumulated from multiple validation failures, allowing
	// all problems to be reported at once rather than failing fast.
	Errors = validation.Errors

	// Kleisli represents a Kleisli arrow for the Validate monad.
	//
	// A Kleisli arrow is a function from A to a monadic value Validate[I, B].
	// It's used for composing computations that produce monadic results.
	//
	// Type: Kleisli[I, A, B] = func(A) Validate[I, B]
	//
	// Kleisli arrows can be composed using the Chain function, enabling
	// sequential validation where later validators depend on earlier results.
	//
	// Example:
	//   parseString := func(s string) Validate[string, int] {
	//     // Parse string to int with validation
	//   }
	//   checkPositive := func(n int) Validate[string, int] {
	//     // Validate that int is positive
	//   }
	//   // Both are Kleisli arrows that can be composed
	Kleisli[I, A, B any] = Reader[A, Validate[I, B]]

	// Operator represents a transformation operator for validators.
	//
	// An Operator transforms a Validate[I, A] into a Validate[I, B].
	// It's a specialized Kleisli arrow where the input is itself a validator.
	//
	// Type: Operator[I, A, B] = func(Validate[I, A]) Validate[I, B]
	//
	// Operators are used to:
	//   - Transform validation results (Map)
	//   - Chain dependent validations (Chain)
	//   - Apply function validators to value validators (Ap)
	//
	// Example:
	//   toUpper := Map[string, string, string](strings.ToUpper)
	//   // toUpper is an Operator[string, string, string]
	//   // It can be applied to any string validator to uppercase the result
	Operator[I, A, B any] = Kleisli[I, Validate[I, A], B]

	// Endomorphism represents a function from a type to itself.
	//
	// Type: Endomorphism[A] = func(A) A
	//
	// An endomorphism is a morphism (structure-preserving map) where the source
	// and target are the same type. In simpler terms, it's a function that takes
	// a value of type A and returns a value of the same type A.
	//
	// Endomorphisms are useful for:
	//   - Transformations that preserve type (e.g., string normalization)
	//   - Composable updates and modifications
	//   - Building pipelines of same-type transformations
	//   - Implementing the Monoid pattern (composition as binary operation)
	//
	// Endomorphisms form a Monoid under function composition:
	//   - Identity: func(a A) A { return a }
	//   - Concat: func(f, g Endomorphism[A]) Endomorphism[A] {
	//       return func(a A) A { return f(g(a)) }
	//     }
	//
	// Example:
	//   trim := strings.TrimSpace      // Endomorphism[string]
	//   lower := strings.ToLower       // Endomorphism[string]
	//   normalize := compose(trim, lower)  // Endomorphism[string]
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
