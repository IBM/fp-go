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

// Package result provides an idiomatic Go approach to error handling using the (value, error) tuple pattern.
//
// This package represents the Result/Either monad idiomatically in Go, leveraging the standard
// (T, error) return pattern that Go developers are familiar with. By convention:
//   - (value, nil) represents a success case (Right)
//   - (zero, error) represents a failure case (Left)
//
// # Core Concepts
//
// The Result pattern is a functional approach to error handling that makes error flows explicit
// and composable. Instead of checking errors manually at each step, you can chain operations
// that automatically short-circuit on the first error.
//
// # Basic Usage
//
//	// Creating Result values
//	success := result.Right[error](42)           // (42, nil)
//	failure := result.Left[int](errors.New("oops")) // (0, error)
//
//	// Pattern matching with Fold
//	output := result.Fold(
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Success: %d", n) },
//	)(success)
//
//	// Chaining operations (short-circuits on Left/error)
//	output := result.Chain(func(n int) (int, error) {
//	    return result.Right[error](n * 2)
//	})(success)
//
// # Monadic Operations
//
// Result implements the Monad interface, providing:
//   - Map: Transform the Right value
//   - Chain (FlatMap): Chain computations that may fail
//   - Ap: Apply a function wrapped in Result
//
// # Error Handling
//
// Result provides utilities for working with Go's error type:
//   - FromError: Create Result from error-returning functions
//   - FromPredicate: Create Result based on a predicate
//   - ToError: Extract the error from a Result
//
// # Subpackages
//
//   - result/exec: Execute system commands returning Result
//   - result/http: HTTP request builders returning Result
package result

//go:generate go run .. either --count 15 --filename gen.go
