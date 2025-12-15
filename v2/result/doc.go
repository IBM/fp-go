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

// Package result provides the Result monad, a specialized Either monad with error as the left type.
//
// Result is commonly used for error handling, where:
//   - Error represents a failure case (type error)
//   - Ok represents a success case (type A)
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
// The Result type is a discriminated union that can hold either an Error value
// or an Ok value (successful result). This makes it ideal for computations that may fail.
//
// # Basic Usage
//
//	// Creating Result values
//	success := result.Ok(42)                      // Ok value
//	failure := result.Error[int](errors.New("oops")) // Error value
//
//	// Pattern matching with Fold
//	output := result.Fold(
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Success: %d", n) },
//	)(success)
//
//	// Chaining operations (short-circuits on Error)
//	doubled := result.Chain(func(n int) result.Result[int] {
//	    return result.Ok(n * 2)
//	})(success)
//
// # Monadic Operations
//
// Result implements the Monad interface, providing:
//   - Map: Transform the Ok value
//   - Chain (FlatMap): Chain computations that may fail
//   - Ap: Apply a function wrapped in Result
//
// # Error Handling
//
// Result provides utilities for working with Go's error type:
//   - TryCatchError: Convert (value, error) tuples to Result
//   - UnwrapError: Convert Result back to (value, error) tuple
//   - FromError: Create Result from error-returning functions
//
// # Subpackages
//
//   - result/http: HTTP request builders returning Result
package result

//go:generate go run .. either --count 15 --filename gen.go
