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

// Package either provides the Either monad, a data structure representing a value of one of two possible types.
//
// Either is commonly used for error handling, where by convention:
//   - Left represents an error or failure case (type E)
//   - Right represents a success case (type A)
//
// # Fantasy Land Specification
//
// This implementation corresponds to the Fantasy Land Either type:
// https://github.com/fantasyland/fantasy-land#either
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Bifunctor: https://github.com/fantasyland/fantasy-land#bifunctor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//   - Foldable: https://github.com/fantasyland/fantasy-land#foldable
//
// # Core Concepts
//
// The Either type is a discriminated union that can hold either a Left value (typically an error)
// or a Right value (typically a successful result). This makes it ideal for computations that may fail.
//
// # Basic Usage
//
//	// Creating Either values
//	success := either.Right[error](42)           // Right value
//	failure := either.Left[int](errors.New("oops")) // Left value
//
//	// Pattern matching with Fold
//	result := either.Fold(
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Success: %d", n) },
//	)(success)
//
//	// Chaining operations (short-circuits on Left)
//	result := either.Chain(func(n int) either.Either[error, int] {
//	    return either.Right[error](n * 2)
//	})(success)
//
// # Monadic Operations
//
// Either implements the Monad interface, providing:
//   - Map: Transform the Right value
//   - Chain (FlatMap): Chain computations that may fail
//   - Ap: Apply a function wrapped in Either
//
// # Error Handling
//
// Either provides utilities for working with Go's error type:
//   - TryCatchError: Convert (value, error) tuples to Either
//   - UnwrapError: Convert Either back to (value, error) tuple
//   - FromError: Create Either from error-returning functions
//
// # Subpackages
//
//   - either/exec: Execute system commands returning Either
//   - either/http: HTTP request builders returning Either
//   - either/testing: Testing utilities for Either laws
package either

//go:generate go run .. either --count 15 --filename gen.go
