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

// Package readerreaderioresult provides a functional programming abstraction that combines
// four powerful concepts: Reader, Reader, IO, and Result (Either[error, A]) monads in a nested structure.
// This is a specialized version of readerreaderioeither where the error type is fixed to `error` and
// the inner context is fixed to `context.Context`.
//
// # Type Definition
//
// ReaderReaderIOResult[R, A] is defined as:
//
//	type ReaderReaderIOResult[R, A] = ReaderReaderIOEither[R, context.Context, error, A]
//
// Which expands to:
//
//	func(R) func(context.Context) func() Either[error, A]
//
// This represents a computation that:
//   - Takes an outer environment/context of type R
//   - Returns a function that takes a context.Context
//   - Returns an IO operation (a thunk/function with no parameters)
//   - Produces an Either[error, A] (Result[A]) when executed
//
// # Type Parameter Ordering Convention
//
// This package follows a consistent convention for ordering type parameters in function signatures.
// The general rule is: R -> C -> E -> T (outer context, inner context, error, type), where:
//   - R: The outer Reader context/environment type
//   - C: The inner Reader context/environment type (for the ReaderIOEither)
//   - E: The Either error type
//   - T: The value type(s) (A, B, etc.)
//
// However, when some type parameters can be automatically inferred by the Go compiler from
// function arguments, the convention is modified to minimize explicit type annotations:
//
// Rule: Undetectable types come first, followed by detectable types, while preserving
// the relative order within each group (R -> C -> E -> T).
//
// Examples:
//
//  1. All types detectable from first argument:
//     MonadMap[R, C, E, A, B](fa ReaderReaderIOEither[R, C, E, A], f func(A) B)
//     - R, C, E, A are detectable from fa
//     - B is detectable from f
//     - Order: R, C, E, A, B (standard order, all detectable)
//
//  2. Some types undetectable:
//     FromReader[C, E, R, A](ma Reader[R, A]) ReaderReaderIOEither[R, C, E, A]
//     - R, A are detectable from ma
//     - C, E are undetectable (not in any argument)
//     - Order: C, E, R, A (C, E first as undetectable, then R, A in standard order)
//
//  3. Multiple undetectable types:
//     Local[C, E, A, R1, R2](f func(R2) R1) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A]
//     - C, E, A are undetectable
//     - R1, R2 are detectable from f
//     - Order: C, E, A, R1, R2 (undetectable first, then detectable)
//
//  4. Functions returning Kleisli arrows:
//     ChainReaderOptionK[R, C, A, B, E](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, B]
//     - Canonical order would be R, C, E, A, B
//     - E is detectable from onNone parameter
//     - R, C, A, B are not detectable (they're in the Kleisli argument type)
//     - Order: R, C, A, B, E (undetectable R, C, A, B first, then detectable E)
//
// This convention allows for more ergonomic function calls:
//
//	// Without convention - need to specify all types:
//	result := FromReader[OuterCtx, InnerCtx, error, User](readerFunc)
//
//	// With convention - only specify undetectable types:
//	result := FromReader[InnerCtx, error](readerFunc)  // R and A inferred from readerFunc
//
// The reasoning behind this approach is to reduce the number of explicit type parameters
// that developers need to specify when calling functions, improving code readability and
// reducing verbosity while maintaining type safety.
//
// Additional examples demonstrating the convention:
//
//  5. FromReaderOption[R, C, A, E](onNone Lazy[E]) Kleisli[R, C, E, ReaderOption[R, A], A]
//     - Canonical order would be R, C, E, A
//     - E is detectable from onNone parameter
//     - R, C, A are not detectable (they're in the return type's Kleisli)
//     - Order: R, C, A, E (undetectable R, C, A first, then detectable E)
//
//  6. MapLeft[R, C, A, E1, E2](f func(E1) E2) func(ReaderReaderIOEither[R, C, E1, A]) ReaderReaderIOEither[R, C, E2, A]
//     - Canonical order would be R, C, E1, E2, A
//     - E1, E2 are detectable from f parameter
//     - R, C, A are not detectable (they're in the return type)
//     - Order: R, C, A, E1, E2 (undetectable R, C, A first, then detectable E1, E2)
//
// Additional special cases:
//
//   - Ap[B, R, C, E, A]: B is undetectable (in function return type), so B comes first
//   - ChainOptionK[R, C, A, B, E]: R, C, A, B are undetectable, E is detectable from onNone
//   - FromReaderIO[C, E, R, A]: C, E are undetectable, R, A are detectable from ReaderIO[R, A]
//
// All functions in this package follow this convention consistently.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - Reader monad: https://github.com/fantasyland/fantasy-land
//   - Reader monad (nested): https://github.com/fantasyland/fantasy-land
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
// # ReaderReaderIOEither
//
// ReaderReaderIOEither[R, C, E, A] represents a computation that:
//   - Depends on an outer context/environment of type R (outer Reader)
//   - Returns a computation that depends on an inner context/environment of type C (inner Reader)
//   - Performs side effects (IO)
//   - Can fail with an error of type E or succeed with a value of type A (Either)
//
// This is particularly useful for:
//   - Multi-level dependency injection patterns
//   - Layered architectures with different context requirements at each layer
//   - Composing operations that need access to multiple levels of configuration or context
//   - Building reusable components that can be configured at different stages
//
// # Core Operations
//
// Construction:
//   - Of/Right: Create a successful computation
//   - Left: Create a failed computation
//   - FromEither: Lift an Either into ReaderReaderIOEither
//   - FromIO: Lift an IO into ReaderReaderIOEither
//   - FromReader: Lift a Reader into ReaderReaderIOEither
//   - FromReaderIO: Lift a ReaderIO into ReaderReaderIOEither
//   - FromIOEither: Lift an IOEither into ReaderReaderIOEither
//   - FromReaderEither: Lift a ReaderEither into ReaderReaderIOEither
//   - FromReaderIOEither: Lift a ReaderIOEither into ReaderReaderIOEither
//   - FromReaderOption: Lift a ReaderOption into ReaderReaderIOEither
//
// Transformation:
//   - Map: Transform the success value
//   - MapLeft: Transform the error value
//   - Chain/Bind: Sequence dependent computations
//   - Flatten: Flatten nested ReaderReaderIOEither
//
// Combination:
//   - Ap: Apply a function in a context to a value in a context
//   - ApSeq: Sequential application
//   - ApPar: Parallel application
//
// Error Handling:
//   - Alt: Choose the first successful computation
//
// Context Access:
//   - Ask: Get the current outer context
//   - Asks: Get a value derived from the outer context
//   - Local: Run a computation with a modified outer context
//   - Read: Execute with a specific outer context
//
// Kleisli Composition:
//   - ChainEitherK: Chain with Either-returning functions
//   - ChainReaderK: Chain with Reader-returning functions
//   - ChainReaderIOK: Chain with ReaderIO-returning functions
//   - ChainReaderEitherK: Chain with ReaderEither-returning functions
//   - ChainReaderOptionK: Chain with ReaderOption-returning functions
//   - ChainIOEitherK: Chain with IOEither-returning functions
//   - ChainIOK: Chain with IO-returning functions
//   - ChainOptionK: Chain with Option-returning functions
//
// First/Tap Operations (execute for side effects, return original value):
//   - ChainFirst/Tap: Execute a computation but return the original value
//   - ChainFirstEitherK/TapEitherK: Tap with Either-returning functions
//   - ChainFirstReaderK/TapReaderK: Tap with Reader-returning functions
//   - ChainFirstReaderIOK/TapReaderIOK: Tap with ReaderIO-returning functions
//   - ChainFirstReaderEitherK/TapReaderEitherK: Tap with ReaderEither-returning functions
//   - ChainFirstReaderOptionK/TapReaderOptionK: Tap with ReaderOption-returning functions
//   - ChainFirstIOK/TapIOK: Tap with IO-returning functions
//
// # Example Usage
//
//	type AppConfig struct {
//	    DatabaseURL string
//	    LogLevel    string
//	}
//
//	// A computation that depends on AppConfig and context.Context
//	func fetchUser(id int) ReaderReaderIOResult[AppConfig, User] {
//	    return func(cfg AppConfig) readerioresult.ReaderIOResult[context.Context, User] {
//	        // Use cfg.DatabaseURL and cfg.LogLevel
//	        return func(ctx context.Context) ioresult.IOResult[User] {
//	            // Use ctx for cancellation/timeout
//	            return func() result.Result[User] {
//	                // Perform the actual IO operation
//	                // Return result.Of(user) or result.Error[User](err)
//	            }
//	        }
//	    }
//	}
//
//	// Compose operations
//	result := function.Pipe2(
//	    fetchUser(123),
//	    Map[AppConfig](func(u User) string { return u.Name }),
//	    Chain[AppConfig](func(name string) ReaderReaderIOResult[AppConfig, string] {
//	        return Of[AppConfig]("Hello, " + name)
//	    }),
//	)
//
//	// Execute with config and context
//	appConfig := AppConfig{DatabaseURL: "postgres://...", LogLevel: "info"}
//	ctx := t.Context()
//	outcome := result(appConfig)(ctx)() // Returns result.Result[string]
//
// # Use Cases
//
// This monad is particularly useful for:
//   - Applications with layered configuration (app config + request context)
//   - HTTP handlers that need both application config and request context
//   - Database operations with connection pool config and query context
//   - Retry logic with policy configuration and execution context
//   - Resource management with bracket pattern across multiple contexts
//
// # Relationship to Other Packages
//
//   - readerreaderioeither: The generic version with configurable error and context types
//   - readerioresult: Single reader with context.Context and error
//   - readerresult: Single reader with error (no IO)
//   - context/readerioresult: Alias for readerioresult with context.Context
package readerreaderioresult
