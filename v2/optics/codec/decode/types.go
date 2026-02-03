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

package decode

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	Errors = validation.Errors

	// Validation represents the result of a validation operation that may contain
	// validation errors or a successfully validated value of type A.
	// This is an alias for validation.Validation[A], which is Either[Errors, A].
	//
	// In the decode context:
	//   - Left(Errors): Decoding failed with one or more validation errors
	//   - Right(A): Successfully decoded value of type A
	//
	// Example:
	//
	//	// Success case
	//	valid := validation.Success(42)  // Right(42)
	//
	//	// Failure case
	//	invalid := validation.Failures[int](validation.Errors{
	//	    &validation.ValidationError{Messsage: "invalid format"},
	//	})  // Left([...])
	Validation[A any] = validation.Validation[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	// This is an alias for reader.Reader[R, A], which is func(R) A.
	//
	// In the decode context, Reader is used to access the input data being decoded.
	// The environment R is typically the raw input (e.g., JSON, string, bytes) that
	// needs to be decoded into a structured type A.
	//
	// Example:
	//
	//	// A reader that extracts a field from a map
	//	getField := func(data map[string]any) string {
	//	    return data["name"].(string)
	//	}  // Reader[map[string]any, string]
	Reader[R, A any] = reader.Reader[R, A]

	// Decode is a function that decodes input I to type A with validation.
	// It combines the Reader pattern (for accessing input) with Validation (for error handling).
	//
	// Type: func(I) Validation[A]
	//
	// A Decode function:
	//  1. Takes raw input of type I (e.g., JSON, string, bytes)
	//  2. Attempts to decode/parse it into type A
	//  3. Returns a Validation[A] with either:
	//     - Success(A): Successfully decoded value
	//     - Failures(Errors): Validation errors describing what went wrong
	//
	// This type is the foundation of the decode package, enabling composable,
	// type-safe decoding with comprehensive error reporting.
	//
	// Example:
	//
	//	// Decode a string to an integer
	//	decodeInt := func(input string) Validation[int] {
	//	    n, err := strconv.Atoi(input)
	//	    if err != nil {
	//	        return validation.Failures[int](validation.Errors{
	//	            &validation.ValidationError{
	//	                Value:    input,
	//	                Messsage: "not a valid integer",
	//	                Cause:    err,
	//	            },
	//	        })
	//	    }
	//	    return validation.Success(n)
	//	}  // Decode[string, int]
	//
	//	result := decodeInt("42")  // Success(42)
	//	result := decodeInt("abc") // Failures([...])
	Decode[I, A any] = Reader[I, Validation[A]]

	// Kleisli represents a function from A to a decoded B given input type I.
	// It's a Reader that takes an input A and produces a Decode[I, B] function.
	// This enables composition of decoding operations in a functional style.
	//
	// Type: func(A) Decode[I, B]
	//       which expands to: func(A) func(I) Validation[B]
	//
	// Kleisli arrows are the fundamental building blocks for composing decoders.
	// They allow you to chain decoding operations where each step can:
	//  1. Depend on the result of the previous step (the A parameter)
	//  2. Access the original input (the I parameter via Decode)
	//  3. Fail with validation errors (via Validation[B])
	//
	// This is particularly useful for:
	//  - Conditional decoding based on previously decoded values
	//  - Multi-stage decoding pipelines
	//  - Dependent field validation
	//
	// Example:
	//
	//	// Decode a user, then decode their age based on their type
	//	decodeAge := func(userType string) Decode[map[string]any, int] {
	//	    return func(data map[string]any) Validation[int] {
	//	        if userType == "admin" {
	//	            // Admins must be 18+
	//	            age := data["age"].(int)
	//	            if age < 18 {
	//	                return validation.Failures[int](/* error */)
	//	            }
	//	            return validation.Success(age)
	//	        }
	//	        // Regular users can be any age
	//	        return validation.Success(data["age"].(int))
	//	    }
	//	}  // Kleisli[map[string]any, string, int]
	//
	//	// Use with Chain to compose decoders
	//	decoder := F.Pipe2(
	//	    decodeUserType,           // Decode[map[string]any, string]
	//	    Chain(decodeAge),         // Chains with Kleisli
	//	    Map(func(age int) User {  // Transform to final type
	//	        return User{Age: age}
	//	    }),
	//	)
	Kleisli[I, A, B any] = Reader[A, Decode[I, B]]

	// Operator represents a decoding transformation that takes a decoded A and produces a decoded B.
	// It's a specialized Kleisli arrow for composing decode operations where the input is already decoded.
	// This allows chaining multiple decode transformations together.
	//
	// Type: func(Decode[I, A]) Decode[I, B]
	//
	// Operators are higher-order functions that transform one decoder into another.
	// They are the result of partially applying functions like Map, Chain, and Ap,
	// making them ideal for use in composition pipelines with F.Pipe.
	//
	// Key characteristics:
	//  - Takes a Decode[I, A] as input
	//  - Returns a Decode[I, B] as output
	//  - Preserves the input type I (the raw data being decoded)
	//  - Transforms the output type from A to B
	//
	// Common operators:
	//  - Map(f): Transforms successful decode results
	//  - Chain(f): Sequences dependent decode operations
	//  - Ap(fa): Applies function decoders to value decoders
	//
	// Example:
	//
	//	// Create reusable operators
	//	toString := Map(func(n int) string {
	//	    return strconv.Itoa(n)
	//	})  // Operator[string, int, string]
	//
	//	validatePositive := Chain(func(n int) Decode[string, int] {
	//	    return func(input string) Validation[int] {
	//	        if n <= 0 {
	//	            return validation.Failures[int](/* error */)
	//	        }
	//	        return validation.Success(n)
	//	    }
	//	})  // Operator[string, int, int]
	//
	//	// Compose operators in a pipeline
	//	decoder := F.Pipe2(
	//	    decodeInt,          // Decode[string, int]
	//	    validatePositive,   // Operator[string, int, int]
	//	    toString,           // Operator[string, int, string]
	//	)  // Decode[string, string]
	//
	//	result := decoder("42")  // Success("42")
	//	result := decoder("-5")  // Failures([...])
	Operator[I, A, B any] = Kleisli[I, Decode[I, A], B]

	// Endomorphism represents a function from a type to itself: func(A) A.
	// This is an alias for endomorphism.Endomorphism[A].
	//
	// In the decode context, endomorphisms are used with LetL to transform
	// decoded values using pure functions that don't change the type.
	//
	// Endomorphisms are useful for:
	//  - Normalizing data (e.g., trimming strings, rounding numbers)
	//  - Applying business rules (e.g., clamping values to ranges)
	//  - Data sanitization (e.g., removing special characters)
	//
	// Example:
	//
	//	// Normalize a string by trimming and lowercasing
	//	normalize := func(s string) string {
	//	    return strings.ToLower(strings.TrimSpace(s))
	//	}  // Endomorphism[string]
	//
	//	// Clamp an integer to a range
	//	clamp := func(n int) int {
	//	    if n < 0 { return 0 }
	//	    if n > 100 { return 100 }
	//	    return n
	//	}  // Endomorphism[int]
	//
	//	// Use with LetL to transform decoded values
	//	decoder := F.Pipe1(
	//	    decodeString,
	//	    LetL(nameLens, normalize),
	//	)
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
