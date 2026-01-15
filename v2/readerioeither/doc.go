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

// Package readerioeither provides a functional programming abstraction that combines
// three powerful concepts: Reader, IO, and Either monads.
//
// # Type Parameter Ordering Convention
//
// This package follows a consistent convention for ordering type parameters in function signatures.
// The general rule is: R -> E -> T (context, error, type), where:
//   - R: The Reader context/environment type
//   - E: The Either error type
//   - T: The value type(s) (A, B, etc.)
//
// However, when some type parameters can be automatically inferred by the Go compiler from
// function arguments, the convention is modified to minimize explicit type annotations:
//
// Rule: Undetectable types come first, followed by detectable types, while preserving
// the relative order within each group (R -> E -> T).
//
// Examples:
//
//  1. All types detectable from first argument:
//     MonadMap[R, E, A, B](fa ReaderIOEither[R, E, A], f func(A) B)
//     - R, E, A are detectable from fa
//     - B is detectable from f
//     - Order: R, E, A, B (standard order, all detectable)
//
//  2. Some types undetectable:
//     FromReader[E, R, A](ma Reader[R, A]) ReaderIOEither[R, E, A]
//     - R, A are detectable from ma
//     - E is undetectable (not in any argument)
//     - Order: E, R, A (E first as undetectable, then R, A in standard order)
//
//  3. Multiple undetectable types:
//     Local[E, A, R1, R2](f func(R2) R1) func(ReaderIOEither[R1, E, A]) ReaderIOEither[R2, E, A]
//     - E, A are undetectable
//     - R1, R2 are detectable from f
//
//  4. Functions returning Kleisli arrows:
//     ChainReaderOptionK[R, A, B, E](onNone func() E) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, B]
//     - Canonical order would be R, E, A, B
//     - E is detectable from onNone parameter
//     - R, A, B are not detectable (they're in the Kleisli argument type)
//     - Order: R, A, B, E (undetectable R, A, B first, then detectable E)
//
// This convention allows for more ergonomic function calls:
//
//	// Without convention - need to specify all types:
//	result := FromReader[context.Context, error, User](readerFunc)
//
//	// With convention - only specify undetectable type:
//	result := FromReader[error](readerFunc)  // R and A inferred from readerFunc
//
// The reasoning behind this approach is to reduce the number of explicit type parameters
// that developers need to specify when calling functions, improving code readability and
// reducing verbosity while maintaining type safety.
//
// Additional examples demonstrating the convention:
//
//  5. FromReaderOption[R, A, E](onNone func() E) Kleisli[R, E, ReaderOption[R, A], A]
//     - Canonical order would be R, E, A
//     - E is detectable from onNone parameter
//     - R, A are not detectable (they're in the return type's Kleisli)
//     - Order: R, A, E (undetectable R, A first, then detectable E)
//
//  6. MapLeft[R, A, E1, E2](f func(E1) E2) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A]
//     - Canonical order would be R, E1, E2, A
//     - E1, E2 are detectable from f parameter
//     - R, A are not detectable (they're in the return type)
//     - Order: R, A, E1, E2 (undetectable R, A first, then detectable E1, E2)
//
// Additional special cases:
//
//   - Ap[B, R, E, A]: B is undetectable (in function return type), so B comes first
//   - OrLeft[A, E1, R, E2]: A is undetectable, comes first before detectable E1, R, E2
//   - ReadIO[E, A, R]: E and A are undetectable, R is detectable from IO[R]
//   - ChainFirstLeft[A, R, EA, EB, B]: A is undetectable, comes first before detectable R, EA, EB, B
//   - TapLeft[A, R, EB, EA, B]: Similar to ChainFirstLeft, A is undetectable and comes first
//
// All functions in this package follow this convention consistently.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - Reader monad: https://github.com/fantasyland/fantasy-land
//   - IO monad: https://github.com/fantasyland/fantasy-land
//   - Either monad: https://github.com/fantasyland/fantasy-land#either
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Bifunctor: https://github.com/fantasyland/fantasy-land#bifunctor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//
// # ReaderIOEither
//
// ReaderIOEither[R, E, A] represents a computation that:
//   - Depends on some context/environment of type R (Reader)
//   - Performs side effects (IO)
//   - Can fail with an error of type E or succeed with a value of type A (Either)
//
// This is particularly useful for:
//   - Dependency injection patterns
//   - Error handling in effectful computations
//   - Composing operations that need access to shared configuration or context
//
// # Core Operations
//
// Construction:
//   - Of/Right: Create a successful computation
//   - Left/ThrowError: Create a failed computation
//   - FromEither: Lift an Either into ReaderIOEither
//   - FromIO: Lift an IO into ReaderIOEither
//   - FromReader: Lift a Reader into ReaderIOEither
//   - FromIOEither: Lift an IOEither into ReaderIOEither
//   - TryCatch: Wrap error-returning functions
//
// Transformation:
//   - Map: Transform the success value
//   - MapLeft: Transform the error value
//   - BiMap: Transform both success and error values
//   - Chain/Bind: Sequence dependent computations
//   - Flatten: Flatten nested ReaderIOEither
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//   - SequenceArray: Convert array of ReaderIOEither to ReaderIOEither of array
//   - TraverseArray: Map and sequence in one operation
//
// Error Handling:
//   - Fold: Handle both success and error cases
//   - GetOrElse: Provide a default value on error
//   - OrElse: Try an alternative computation on error
//   - Alt: Choose the first successful computation
//
// Context Access:
//   - Ask: Get the current context
//   - Asks: Get a value derived from the context
//   - Local: Run a computation with a modified context
//
// Resource Management:
//   - Bracket: Ensure resource cleanup
//   - WithResource: Manage resource lifecycle
//
// # Example Usage
//
//	type Config struct {
//	    BaseURL string
//	    Timeout time.Duration
//	}
//
//	// A computation that depends on Config, performs IO, and can fail
//	func fetchUser(id int) readerioeither.ReaderIOEither[Config, error, User] {
//	    return func(cfg Config) ioeither.IOEither[error, User] {
//	        return func() either.Either[error, User] {
//	            // Use cfg.BaseURL and cfg.Timeout to fetch user
//	            // Return either.Right(user) or either.Left(err)
//	        }
//	    }
//	}
//
//	// Compose operations
//	result := function.Pipe2(
//	    fetchUser(123),
//	    readerioeither.Map[Config, error](func(u User) string { return u.Name }),
//	    readerioeither.Chain[Config, error](func(name string) readerioeither.ReaderIOEither[Config, error, string] {
//	        return readerioeither.Of[Config, error]("Hello, " + name)
//	    }),
//	)
//
//	// Execute with config
//	config := Config{BaseURL: "https://api.example.com", Timeout: 30 * time.Second}
//	outcome := result(config)() // Returns either.Either[error, string]
package readerioeither

//go:generate go run .. readerioeither --count 10 --filename gen.go
