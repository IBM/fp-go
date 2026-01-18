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

// Package readerresult implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error.
//
// # Side Effects and Context
//
// IMPORTANT: In contrast to the functional readerresult package (readerresult.ReaderResult[R, A]),
// this context/readerresult package has side effects by design because it depends on context.Context,
// which is inherently effectful:
//   - context.Context can be cancelled (ctx.Done() channel)
//   - context.Context has deadlines and timeouts (ctx.Deadline())
//   - context.Context carries request-scoped values (ctx.Value())
//   - context.Context propagates cancellation signals across goroutines
//
// This means that ReaderResult[A] = func(context.Context) (A, error) represents an EFFECTFUL computation,
// not a pure function. The computation's behavior can change based on the context's state (cancelled,
// timed out, etc.), making it fundamentally different from a pure Reader monad.
//
// Comparison of packages:
//   - readerresult.ReaderResult[R, A] = func(R) Result[A] - PURE (R can be any type, no side effects)
//   - idiomatic/readerresult.ReaderResult[R, A] = func(R) (A, error) - EFFECTFUL (also uses context.Context)
//   - context/readerresult.ReaderResult[A] = func(context.Context) (A, error) - EFFECTFUL (uses context.Context)
//
// Use this package (context/readerresult) when you need:
//   - Cancellation support for long-running operations
//   - Timeout/deadline handling
//   - Request-scoped values (tracing IDs, user context, etc.)
//   - Integration with Go's standard context-aware APIs
//   - Idiomatic Go error handling with (value, error) tuples
//
// Use the functional readerresult package when you need:
//   - Pure dependency injection without side effects
//   - Testable computations with simple state/config objects
//   - Functional composition without context propagation
//   - Generic environment types (not limited to context.Context)
//
// # Pure vs Effectful Functions
//
// This package distinguishes between pure (side-effect free) and effectful (side-effectful) functions:
//
// EFFECTFUL FUNCTIONS (depend on context.Context):
//   - ReaderResult[A]: func(context.Context) (A, error) - Effectful computation that needs context
//   - These functions are effectful because context.Context is effectful (can be cancelled, has deadlines, carries values)
//   - Use for: operations that need cancellation, timeouts, context values, or any context-dependent behavior
//   - Examples: database queries, HTTP requests, operations that respect cancellation
//
// PURE FUNCTIONS (side-effect free):
//   - func(State) (Value, error) - Pure computation that only depends on state, not context
//   - func(State) Value - Pure transformation without errors
//   - These functions are pure because they only read from their input state and don't depend on external context
//   - Use for: parsing, validation, calculations, data transformations that don't need context
//   - Examples: JSON parsing, input validation, mathematical computations
//
// The package provides different bind operations for each:
//   - Bind: For effectful ReaderResult computations (State -> ReaderResult[Value])
//   - BindResultK: For pure functions with errors (State -> (Value, error))
//   - Let: For pure functions without errors (State -> Value)
//   - BindReaderK: For context-dependent pure functions (State -> Reader[Context, Value])
//   - BindEitherK: For pure Result/Either values (State -> Result[Value])
package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Option represents an optional value that may or may not be present.
	// This is an alias for option.Option[A].
	//
	// Type Parameters:
	//   - A: The type of the value that may be present
	//
	// Example:
	//
	//	opt := option.Some(42)           // Option[int] with value
	//	none := option.None[int]()       // Option[int] without value
	Option[A any] = option.Option[A]

	// Either represents a value that can be either a Left (error) or Right (success).
	// This is specialized to use error as the Left type.
	// This is an alias for either.Either[error, A].
	//
	// Type Parameters:
	//   - A: The type of the Right (success) value
	//
	// Example:
	//
	//	success := either.Right[error, int](42)           // Right(42)
	//	failure := either.Left[int](errors.New("failed")) // Left(error)
	Either[A any] = either.Either[error, A]

	// Result represents a computation that can either succeed with a value or fail with an error.
	// This is an alias for result.Result[A], which is equivalent to Either[error, A].
	//
	// Type Parameters:
	//   - A: The type of the success value
	//
	// Example:
	//
	//	success := result.Of[error](42)                // Right(42)
	//	failure := result.Error[int](errors.New("failed")) // Left(error)
	Result[A any] = result.Result[A]

	// Reader represents a computation that depends on an environment R to produce a value A.
	// This is an alias for reader.Reader[R, A].
	//
	// Type Parameters:
	//   - R: The type of the environment/context
	//   - A: The type of the produced value
	//
	// Example:
	//
	//	type Config struct { Port int }
	//	getPort := func(cfg Config) int { return cfg.Port }
	//	// getPort is a Reader[Config, int]
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderResult is a specialization of the Reader monad for the typical Go scenario.
	// It represents an effectful computation that:
	//   - Depends on context.Context (for cancellation, deadlines, values)
	//   - Can fail with an error
	//   - Produces a value of type A on success
	//
	// IMPORTANT: This is an EFFECTFUL type because context.Context is effectful.
	// The computation's behavior can change based on context state (cancelled, timed out, etc.).
	//
	// Type Parameters:
	//   - A: The type of the success value
	//
	// Signature:
	//
	//	type ReaderResult[A any] = func(context.Context) Result[A]
	//
	// Example:
	//
	//	getUserByID := func(ctx context.Context) result.Result[User] {
	//	    if ctx.Err() != nil {
	//	        return result.Error[User](ctx.Err())
	//	    }
	//	    // Fetch user from database
	//	    return result.Of(User{ID: 123, Name: "Alice"})
	//	}
	//	// getUserByID is a ReaderResult[User]
	ReaderResult[A any] = readereither.ReaderEither[context.Context, error, A]

	// Kleisli represents a function that takes a value of type A and returns a ReaderResult[B].
	// This is the fundamental building block for composing ReaderResult computations.
	//
	// Type Parameters:
	//   - A: The input type
	//   - B: The output type (wrapped in ReaderResult)
	//
	// Signature:
	//
	//	type Kleisli[A, B any] = func(A) ReaderResult[B]
	//
	// Example:
	//
	//	getUserByID := func(id int) readerresult.ReaderResult[User] {
	//	    return func(ctx context.Context) result.Result[User] {
	//	        // Fetch user from database
	//	        return result.Of(User{ID: id, Name: "Alice"})
	//	    }
	//	}
	//	// getUserByID is a Kleisli[int, User]
	Kleisli[A, B any] = reader.Reader[A, ReaderResult[B]]

	// Operator represents a function that transforms one ReaderResult into another.
	// This is a specialized Kleisli where the input is itself a ReaderResult.
	//
	// Type Parameters:
	//   - A: The input ReaderResult's success type
	//   - B: The output ReaderResult's success type
	//
	// Signature:
	//
	//	type Operator[A, B any] = func(ReaderResult[A]) ReaderResult[B]
	//
	// Example:
	//
	//	mapToString := readerresult.Map(func(x int) string {
	//	    return fmt.Sprintf("value: %d", x)
	//	})
	//	// mapToString is an Operator[int, string]
	Operator[A, B any] = Kleisli[ReaderResult[A], B]

	// Endomorphism represents a function that transforms a value to the same type.
	// This is an alias for endomorphism.Endomorphism[A].
	//
	// Type Parameters:
	//   - A: The type of the value
	//
	// Signature:
	//
	//	type Endomorphism[A any] = func(A) A
	//
	// Example:
	//
	//	increment := func(x int) int { return x + 1 }
	//	// increment is an Endomorphism[int]
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Prism is an optic that focuses on a part of a data structure that may or may not be present.
	// This is an alias for prism.Prism[S, T].
	//
	// Type Parameters:
	//   - S: The source type
	//   - T: The target type
	//
	// Example:
	//
	//	// A prism that extracts an int from a string if it's a valid number
	//	intPrism := prism.Prism[string, int]{...}
	Prism[S, T any] = prism.Prism[S, T]

	// Lens is an optic that focuses on a part of a data structure that is always present.
	// This is an alias for lens.Lens[S, T].
	//
	// Type Parameters:
	//   - S: The source type
	//   - T: The target type
	//
	// Example:
	//
	//	// A lens that focuses on the Name field of a User
	//	nameLens := lens.Lens[User, string]{...}
	Lens[S, T any] = lens.Lens[S, T]

	// Trampoline represents a computation that can be executed in a stack-safe manner
	// using tail recursion elimination. This is an alias for tailrec.Trampoline[A, B].
	//
	// Type Parameters:
	//   - A: The input type
	//   - B: The output type
	//
	// Example:
	//
	//	// A tail-recursive factorial computation
	//	factorial := tailrec.Trampoline[int, int]{...}
	Trampoline[A, B any] = tailrec.Trampoline[A, B]

	// Predicate represents a function that tests a value and returns a boolean.
	// This is an alias for predicate.Predicate[A].
	//
	// Type Parameters:
	//   - A: The type of the value to test
	//
	// Signature:
	//
	//	type Predicate[A any] = func(A) bool
	//
	// Example:
	//
	//	isPositive := func(x int) bool { return x > 0 }
	//	// isPositive is a Predicate[int]
	Predicate[A any] = predicate.Predicate[A]

	// IO represents a side-effectful computation that produces a value of type A.
	// This is an alias for io.IO[A].
	//
	// IMPORTANT: IO operations have side effects (file I/O, network calls, etc.).
	// Combining IO with ReaderResult makes sense because ReaderResult is already effectful
	// due to its dependency on context.Context.
	//
	// Type Parameters:
	//   - A: The type of the value produced by the IO operation
	//
	// Signature:
	//
	//	type IO[A any] = func() A
	//
	// Example:
	//
	//	readConfig := func() Config {
	//	    // Side effect: read from file
	//	    data, _ := os.ReadFile("config.json")
	//	    return parseConfig(data)
	//	}
	//	// readConfig is an IO[Config]
	IO[A any] = io.IO[A]

	// IOResult represents a side-effectful computation that can fail with an error.
	// This combines IO (side effects) with Result (error handling).
	// This is an alias for ioresult.IOResult[A].
	//
	// IMPORTANT: IOResult operations have side effects and can fail.
	// Combining IOResult with ReaderResult makes sense because both are effectful.
	//
	// Type Parameters:
	//   - A: The type of the success value
	//
	// Signature:
	//
	//	type IOResult[A any] = func() Result[A]
	//
	// Example:
	//
	//	readConfig := func() result.Result[Config] {
	//	    // Side effect: read from file
	//	    data, err := os.ReadFile("config.json")
	//	    if err != nil {
	//	        return result.Error[Config](err)
	//	    }
	//	    return result.Of(parseConfig(data))
	//	}
	//	// readConfig is an IOResult[Config]
	IOResult[A any] = ioresult.IOResult[A]
)
