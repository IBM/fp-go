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

// Package readerio provides the ReaderIO monad, which combines the Reader and IO monads.
//
// ReaderIO[R, A] represents a computation that:
//   - Requires an environment of type R (Reader aspect)
//   - Performs side effects (IO aspect)
//   - Produces a value of type A
//
// This monad is particularly useful for dependency injection patterns and logging scenarios,
// where you need to:
//   - Access configuration or context throughout your application
//   - Perform side effects like logging, file I/O, or network calls
//   - Maintain functional composition and testability
//
// # Logging Use Case
// ReaderIO is especially well-suited for logging because it allows you to:
//   - Pass a logger through your computation chain without explicit parameter threading
//   - Compose logging operations with other side effects
//   - Test logging behavior by providing mock loggers in the environment
//
// Key functions for logging scenarios:
//   - [Ask]: Retrieve the entire environment (e.g., a logger instance)
//   - [Asks]: Extract a specific value from the environment (e.g., logger.Info method)
//   - [ChainIOK]: Chain logging operations that return IO effects
//   - [MonadChain]: Sequence multiple logging and computation steps
//
// Example logging usage:
//
//	type Env struct {
//	    Logger *log.Logger
//	}
//
//	// Log a message using the environment's logger
//	logInfo := func(msg string) readerio.ReaderIO[Env, func()] {
//	    return readerio.Asks(func(env Env) io.IO[func()] {
//	        return io.Of(func() { env.Logger.Println(msg) })
//	    })
//	}
//
//	// Compose logging with computation
//	computation := F.Pipe3(
//	    readerio.Of[Env](42),
//	    readerio.Chain(func(n int) readerio.ReaderIO[Env, int] {
//	        return F.Pipe1(
//	            logInfo(fmt.Sprintf("Processing: %d", n)),
//	            readerio.Map[Env](func(func()) int { return n * 2 }),
//	        )
//	    }),
//	    readerio.ChainIOK(func(result int) io.IO[int] {
//	        return io.Of(result)
//	    }),
//	)
//
//	// Execute with environment
//	env := Env{Logger: log.New(os.Stdout, "APP: ", log.LstdFlags)}
//	result := computation(env)() // Logs "Processing: 42" and returns 84
//
// # Core Operations
//
// The package provides standard monadic operations:
//   - [Of]: Lift a pure value into ReaderIO
//   - [Map]: Transform the result value
//   - [Chain]: Sequence dependent computations
//   - [Ap]: Apply a function in ReaderIO context
//
// # Integration
//
// Convert between different contexts:
//   - [FromIO]: Lift an IO action into ReaderIO
//   - [FromReader]: Lift a Reader into ReaderIO
//   - [ChainIOK]: Chain with IO-returning functions
//
// # Performance
//
//   - [Memoize]: Cache computation results (use with caution for context-dependent values)
//   - [Defer]: Ensure fresh computation on each execution
package readerio

import (
	"sync"
	"time"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
)

// FromIO converts an [IO] action to a [ReaderIO] that ignores the environment.
// This lifts a pure IO computation into the ReaderIO context.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Result type
//
// Parameters:
//   - t: The IO action to lift
//
// Returns:
//   - A ReaderIO that executes the IO action regardless of the environment
//
// Example:
//
//	ioAction := io.Of(42)
//	readerIO := readerio.FromIO[Config](ioAction)
//	result := readerIO(config)() // Returns 42
func FromIO[R, A any](t IO[A]) ReaderIO[R, A] {
	return reader.Of[R](t)
}

// FromReader converts a [Reader] to a [ReaderIO] by lifting the pure computation into IO.
// This allows you to use Reader computations in a ReaderIO context.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Result type
//
// Parameters:
//   - r: The Reader to convert
//
// Returns:
//   - A ReaderIO that wraps the Reader computation in IO
//
// Example:
//
//	reader := func(config Config) int { return config.Port }
//	readerIO := readerio.FromReader(reader)
//	result := readerIO(config)() // Returns config.Port
func FromReader[R, A any](r Reader[R, A]) ReaderIO[R, A] {
	return readert.MonadFromReader[Reader[R, A], ReaderIO[R, A]](io.Of[A], r)
}

// MonadMap applies a function to the value inside a ReaderIO context.
// This is the monadic version that takes the ReaderIO as the first parameter.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - fa: The ReaderIO containing the value to transform
//   - f: The transformation function
//
// Returns:
//   - A new ReaderIO with the transformed value
//
// Example:
//
//	rio := readerio.Of[Config](5)
//	doubled := readerio.MonadMap(rio, N.Mul(2))
//	result := doubled(config)() // Returns 10
func MonadMap[R, A, B any](fa ReaderIO[R, A], f func(A) B) ReaderIO[R, B] {
	return readert.MonadMap[ReaderIO[R, A], ReaderIO[R, B]](io.MonadMap[A, B], fa, f)
}

// MonadMapTo executes a ReaderIO computation, discards its result, and returns a constant value.
// This is the monadic version that takes both the ReaderIO and the constant value as parameters.
//
// IMPORTANT: ReaderIO represents a side-effectful computation (IO effects). For this reason,
// MonadMapTo WILL execute the original ReaderIO to allow any side effects to occur (such as
// logging, file I/O, network calls, etc.), then discard the result and return the constant value.
// The side effects are preserved even though the result value is discarded.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type (result will be discarded after execution)
//   - B: Output value type (constant to return)
//
// Parameters:
//   - fa: The ReaderIO to execute (side effects will occur, result discarded)
//   - b: The constant value to return after executing fa
//
// Returns:
//   - A new ReaderIO that executes fa for its side effects, then returns b
//
// Example:
//
//	logAndCompute := func(r Config) io.IO[int] {
//	    return io.Of(func() int {
//	        fmt.Println("Computing...") // Side effect
//	        return 42
//	    })
//	}
//	replaced := readerio.MonadMapTo(logAndCompute, "done")
//	result := replaced(config)() // Prints "Computing...", returns "done"
func MonadMapTo[R, A, B any](fa ReaderIO[R, A], b B) ReaderIO[R, B] {
	return MonadMap(fa, reader.Of[A](b))
}

// Map creates a function that applies a transformation to a ReaderIO value.
// This is the curried version suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - f: The transformation function
//
// Returns:
//   - An Operator that transforms ReaderIO[R, A] to ReaderIO[R, B]
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](5),
//	    readerio.Map[Config](N.Mul(2)),
//	)(config)() // Returns 10
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return readert.Map[ReaderIO[R, A], ReaderIO[R, B]](io.Map[A, B], f)
}

// MapTo creates an operator that executes a ReaderIO computation, discards its result,
// and returns a constant value. This is the curried version of [MonadMapTo], suitable for use in pipelines.
//
// IMPORTANT: ReaderIO represents a side-effectful computation (IO effects). For this reason,
// MapTo WILL execute the input ReaderIO to allow any side effects to occur (such as logging,
// file I/O, network calls, etc.), then discard the result and return the constant value.
// The side effects are preserved even though the result value is discarded.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type (result will be discarded after execution)
//   - B: Output value type (constant to return)
//
// Parameters:
//   - b: The constant value to return after executing the ReaderIO
//
// Returns:
//   - An Operator that executes a ReaderIO for its side effects, then returns b
//
// Example:
//
//	logStep := func(r Config) io.IO[int] {
//	    return io.Of(func() int {
//	        fmt.Println("Step executed") // Side effect
//	        return 42
//	    })
//	}
//	result := F.Pipe1(
//	    logStep,
//	    readerio.MapTo[Config, int]("complete"),
//	)(config)() // Prints "Step executed", returns "complete"
func MapTo[R, A, B any](b B) Operator[R, A, B] {
	return Map[R](reader.Of[A](b))
}

// MonadChain sequences two ReaderIO computations, where the second depends on the result of the first.
// This is the monadic bind operation for ReaderIO.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - ma: The first ReaderIO computation
//   - f: Function that takes the result of ma and returns a new ReaderIO
//
// Returns:
//   - A ReaderIO that sequences both computations
//
// Example:
//
//	rio1 := readerio.Of[Config](5)
//	result := readerio.MonadChain(rio1, func(n int) readerio.ReaderIO[Config, int] {
//	    return readerio.Of[Config](n * 2)
//	})
func MonadChain[R, A, B any](ma ReaderIO[R, A], f Kleisli[R, A, B]) ReaderIO[R, B] {
	return readert.MonadChain(io.MonadChain[A, B], ma, f)
}

// MonadChainFirst sequences two ReaderIO computations but returns the result of the first.
// The second computation is executed for its side effects only (e.g., logging, validation).
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Intermediate value type (discarded)
//
// Parameters:
//   - ma: The first ReaderIO computation
//   - f: Function that produces the second ReaderIO (for side effects)
//
// Returns:
//   - A ReaderIO with the result of the first computation
//
// Example:
//
//	rio := readerio.Of[Config](42)
//	result := readerio.MonadChainFirst(rio, func(n int) readerio.ReaderIO[Config, string] {
//	    // Log the value but don't change the result
//	    return readerio.Of[Config](fmt.Sprintf("Logged: %d", n))
//	})
//	value := result(config)() // Returns 42, but logging happened
func MonadChainFirst[R, A, B any](ma ReaderIO[R, A], f Kleisli[R, A, B]) ReaderIO[R, A] {
	return chain.MonadChainFirst(
		MonadChain,
		MonadMap,
		ma,
		f,
	)
}

// MonadTap executes a side-effect computation but returns the original value.
// This is an alias for [MonadChainFirst] and is useful for operations like logging
// or validation that should not affect the main computation flow.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Side effect value type (discarded)
//
// Parameters:
//   - ma: The ReaderIO to tap
//   - f: Function that produces a side-effect ReaderIO
//
// Returns:
//   - A ReaderIO with the original value after executing the side effect
//
// Example:
//
//	result := readerio.MonadTap(
//	    readerio.Of[Config](42),
//	    func(n int) readerio.ReaderIO[Config, func()] {
//	        return readerio.FromIO[Config](io.Of(func() { fmt.Println(n) }))
//	    },
//	)
func MonadTap[R, A, B any](ma ReaderIO[R, A], f Kleisli[R, A, B]) ReaderIO[R, A] {
	return MonadChainFirst(ma, f)
}

// Chain creates a function that sequences ReaderIO computations.
// This is the curried version suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - f: Function that takes a value and returns a ReaderIO
//
// Returns:
//   - An Operator that chains ReaderIO computations
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](5),
//	    readerio.Chain(func(n int) readerio.ReaderIO[Config, int] {
//	        return readerio.Of[Config](n * 2)
//	    }),
//	)(config)() // Returns 10
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return readert.Chain[ReaderIO[R, A]](io.Chain[A, B], f)
}

// ChainFirst creates a function that sequences ReaderIO computations but returns the first result.
// This is the curried version of [MonadChainFirst], suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Intermediate value type (discarded)
//
// Parameters:
//   - f: Function that produces a side-effect ReaderIO
//
// Returns:
//   - An Operator that sequences computations while preserving the original value
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](42),
//	    readerio.ChainFirst(func(n int) readerio.ReaderIO[Config, string] {
//	        return readerio.Of[Config](fmt.Sprintf("Logged: %d", n))
//	    }),
//	)(config)() // Returns 42
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return chain.ChainFirst(
		Chain[R, A, A],
		Map[R, B, A],
		f,
	)
}

// Tap creates a function that executes a side-effect computation but returns the original value.
// This is the curried version of [MonadTap], an alias for [ChainFirst].
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Side effect value type (discarded)
//
// Parameters:
//   - f: Function that produces a side-effect ReaderIO
//
// Returns:
//   - An Operator that taps ReaderIO computations
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](42),
//	    readerio.Tap(func(n int) readerio.ReaderIO[Config, func()] {
//	        return readerio.FromIO[Config](io.Of(func() { fmt.Println(n) }))
//	    }),
//	)(config)() // Returns 42, prints 42
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirst(f)
}

// Of creates a ReaderIO that returns a pure value, ignoring the environment.
// This is the monadic return/pure operation for ReaderIO.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Value type
//
// Parameters:
//   - a: The value to wrap
//
// Returns:
//   - A ReaderIO that always returns the given value
//
// Example:
//
//	rio := readerio.Of[Config](42)
//	result := rio(config)() // Returns 42
func Of[R, A any](a A) ReaderIO[R, A] {
	return readert.MonadOf[ReaderIO[R, A]](io.Of[A], a)
}

// MonadAp applies a function wrapped in a ReaderIO to a value wrapped in a ReaderIO.
// This is the applicative apply operation for ReaderIO.
//
// Type Parameters:
//   - B: Result type
//   - R: Reader environment type
//   - A: Input value type
//
// Parameters:
//   - fab: ReaderIO containing a function from A to B
//   - fa: ReaderIO containing a value of type A
//
// Returns:
//   - A ReaderIO containing the result of applying the function to the value
//
// Example:
//
//	fabIO := readerio.Of[Config](N.Mul(2))
//
//	faIO := readerio.Of[Config](5)
//	result := readerio.MonadAp(fabIO, faIO)(config)() // Returns 10
func MonadAp[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B] {
	return readert.MonadAp[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(A) B], R, A](io.MonadAp[A, B], fab, fa)
}

// MonadApSeq is like MonadAp but ensures sequential execution of effects.
//
// Type Parameters:
//   - B: Result type
//   - R: Reader environment type
//   - A: Input value type
//
// Parameters:
//   - fab: ReaderIO containing a function from A to B
//   - fa: ReaderIO containing a value of type A
//
// Returns:
//   - A ReaderIO containing the result, with sequential execution guaranteed
func MonadApSeq[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B] {
	return readert.MonadAp[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(A) B], R, A](io.MonadApSeq[A, B], fab, fa)
}

// MonadApPar is like MonadAp but allows parallel execution of effects where possible.
//
// Type Parameters:
//   - B: Result type
//   - R: Reader environment type
//   - A: Input value type
//
// Parameters:
//   - fab: ReaderIO containing a function from A to B
//   - fa: ReaderIO containing a value of type A
//
// Returns:
//   - A ReaderIO containing the result, with potential parallel execution
func MonadApPar[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B] {
	return readert.MonadAp[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(A) B], R, A](io.MonadApPar[A, B], fab, fa)
}

// Ap creates a function that applies a ReaderIO value to a ReaderIO function.
// This is the curried version suitable for use in pipelines.
//
// Type Parameters:
//   - B: Result type
//   - R: Reader environment type
//   - A: Input value type
//
// Parameters:
//   - fa: ReaderIO containing a value of type A
//
// Returns:
//   - An Operator that applies the value to a function
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](N.Mul(2)),
//	    readerio.Ap[int](readerio.Of[Config](5)),
//	)(config)() // Returns 10
func Ap[B, R, A any](fa ReaderIO[R, A]) Operator[R, func(A) B, B] {
	return function.Bind2nd(MonadAp[B, R, A], fa)
}

// Ask retrieves the current environment.
// This is the fundamental operation for accessing the Reader context.
//
// Type Parameters:
//   - R: Reader environment type
//
// Returns:
//   - A ReaderIO that returns the environment
//
// Example:
//
//	type Config struct { Port int }
//	rio := readerio.Ask[Config]()
//	config := Config{Port: 8080}
//	result := rio(config)() // Returns Config{Port: 8080}
func Ask[R any]() ReaderIO[R, R] {
	return fromreader.Ask(FromReader[R, R])()
}

// Asks retrieves a value derived from the environment using a Reader function.
// This allows you to extract specific information from the environment.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Result type
//
// Parameters:
//   - r: Function that extracts a value from the environment
//
// Returns:
//   - A ReaderIO that applies the function to the environment
//
// Example:
//
//	type Config struct { Port int }
//	rio := readerio.Asks(func(c Config) io.IO[int] {
//	    return io.Of(c.Port)
//	})
//	result := rio(Config{Port: 8080})() // Returns 8080
func Asks[R, A any](r Reader[R, A]) ReaderIO[R, A] {
	return fromreader.Asks(FromReader[R, A])(r)
}

// MonadChainIOK chains a ReaderIO with a function that returns an IO.
// This is useful for integrating IO operations into a ReaderIO pipeline.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - ma: The ReaderIO computation
//   - f: Function that takes a value and returns an IO
//
// Returns:
//   - A ReaderIO that sequences the computation with the IO operation
//
// Example:
//
//	rio := readerio.Of[Config](5)
//	result := readerio.MonadChainIOK(rio, func(n int) io.IO[int] {
//	    return io.Of(n * 2)
//	})
func MonadChainIOK[R, A, B any](ma ReaderIO[R, A], f io.Kleisli[A, B]) ReaderIO[R, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, A, B],
		FromIO[R, B],
		ma, f,
	)
}

// MonadChainFirstIOK chains a ReaderIO with an IO-returning function but keeps the original value.
// The IO computation is executed for its side effects only.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: IO result type (discarded)
//
// Parameters:
//   - ma: The ReaderIO computation
//   - f: Function that takes a value and returns an IO
//
// Returns:
//   - A ReaderIO with the original value after executing the IO
//
// Example:
//
//	rio := readerio.Of[Config](42)
//	result := readerio.MonadChainFirstIOK(rio, func(n int) io.IO[string] {
//	    return io.Of(fmt.Sprintf("Logged: %d", n))
//	})
//	value := result(config)() // Returns 42
func MonadChainFirstIOK[R, A, B any](ma ReaderIO[R, A], f io.Kleisli[A, B]) ReaderIO[R, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromIO[R, B],
		ma, f,
	)
}

// MonadTapIOK chains a ReaderIO with an IO-returning function but keeps the original value.
// This is an alias for [MonadChainFirstIOK] and is useful for side effects like logging.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: IO result type (discarded)
//
// Parameters:
//   - ma: The ReaderIO to tap
//   - f: Function that takes a value and returns an IO for side effects
//
// Returns:
//   - A ReaderIO with the original value after executing the IO
//
// Example:
//
//	result := readerio.MonadTapIOK(
//	    readerio.Of[Config](42),
//	    func(n int) io.IO[func()] {
//	        return io.Of(func() { fmt.Println(n) })
//	    },
//	)
func MonadTapIOK[R, A, B any](ma ReaderIO[R, A], f io.Kleisli[A, B]) ReaderIO[R, A] {
	return MonadChainFirstIOK(ma, f)
}

// ChainIOK creates a function that chains a ReaderIO with an IO operation.
// This is the curried version suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - f: Function that takes a value and returns an IO
//
// Returns:
//   - An Operator that chains ReaderIO with IO
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](5),
//	    readerio.ChainIOK(func(n int) io.IO[int] {
//	        return io.Of(n * 2)
//	    }),
//	)(config)() // Returns 10
func ChainIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, B] {
	return fromio.ChainIOK(
		Chain[R, A, B],
		FromIO[R, B],
		f,
	)
}

// ChainFirstIOK creates a function that chains a ReaderIO with an IO operation but keeps the original value.
// This is the curried version of [MonadChainFirstIOK], suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: IO result type (discarded)
//
// Parameters:
//   - f: Function that takes a value and returns an IO
//
// Returns:
//   - An Operator that chains with IO while preserving the original value
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](42),
//	    readerio.ChainFirstIOK(func(n int) io.IO[string] {
//	        return io.Of(fmt.Sprintf("Logged: %d", n))
//	    }),
//	)(config)() // Returns 42
func ChainFirstIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, A, A],
		Map[R, B, A],
		FromIO[R, B],
		f,
	)
}

// TapIOK creates a function that chains a ReaderIO with an IO operation but keeps the original value.
// This is the curried version of [MonadTapIOK], an alias for [ChainFirstIOK].
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: IO result type (discarded)
//
// Parameters:
//   - f: Function that takes a value and returns an IO for side effects
//
// Returns:
//   - An Operator that taps with IO-returning functions
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](42),
//	    readerio.TapIOK(func(n int) io.IO[func()] {
//	        return io.Of(func() { fmt.Println(n) })
//	    }),
//	)(config)() // Returns 42, prints 42
func TapIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOK[R](f)
}

// Defer creates a ReaderIO by calling a generator function each time it's executed.
// This allows for lazy evaluation and ensures a fresh computation on each invocation.
// Useful for operations that should not be cached or memoized.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Result type
//
// Parameters:
//   - gen: Generator function that creates a new ReaderIO on each call
//
// Returns:
//   - A ReaderIO that calls the generator function on each execution
//
// Example:
//
//	counter := 0
//	rio := readerio.Defer(func() readerio.ReaderIO[Config, int] {
//	    counter++
//	    return readerio.Of[Config](counter)
//	})
//	result1 := rio(config)() // Returns 1
//	result2 := rio(config)() // Returns 2 (fresh computation)
func Defer[R, A any](gen func() ReaderIO[R, A]) ReaderIO[R, A] {
	return func(r R) IO[A] {
		return func() A {
			return gen()(r)()
		}
	}
}

// Memoize computes the value of the provided [ReaderIO] monad lazily but exactly once.
// The first execution caches the result, and subsequent executions return the cached value.
//
// IMPORTANT: The context used to compute the value is the context of the first call.
// Do not use this method if the value has a functional dependency on the content of the context,
// as subsequent calls with different contexts will still return the memoized result from the first call.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Result type
//
// Parameters:
//   - rdr: The ReaderIO to memoize
//
// Returns:
//   - A ReaderIO that caches its result after the first execution
//
// Example:
//
//	expensive := readerio.Of[Config](computeExpensiveValue())
//	memoized := readerio.Memoize(expensive)
//	result1 := memoized(config)() // Computes the value
//	result2 := memoized(config)() // Returns cached value (no recomputation)
func Memoize[R, A any](rdr ReaderIO[R, A]) ReaderIO[R, A] {
	// synchronization primitives
	var once sync.Once
	var result A
	// callback
	gen := func(r R) func() {
		return func() {
			result = rdr(r)()
		}
	}
	// returns our memoized wrapper
	return func(r R) IO[A] {
		io := gen(r)
		return func() A {
			once.Do(io)
			return result
		}
	}
}

// Flatten removes one level of nesting from a ReaderIO structure.
// Converts ReaderIO[R, ReaderIO[R, A]] to ReaderIO[R, A].
// This is also known as "join" in monad terminology.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Result type
//
// Parameters:
//   - mma: A nested ReaderIO structure
//
// Returns:
//   - A flattened ReaderIO with one less level of nesting
//
// Example:
//
//	nested := readerio.Of[Config](readerio.Of[Config](42))
//	flattened := readerio.Flatten(nested)
//	result := flattened(config)() // Returns 42
func Flatten[R, A any](mma ReaderIO[R, ReaderIO[R, A]]) ReaderIO[R, A] {
	return MonadChain(mma, function.Identity[ReaderIO[R, A]])
}

// MonadFlap applies a value to a function wrapped in a ReaderIO.
// This is the "flipped" version of MonadAp, where the value comes second.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Result type
//
// Parameters:
//   - fab: ReaderIO containing a function from A to B
//   - a: The value to apply to the function
//
// Returns:
//   - A ReaderIO containing the result of applying the value to the function
//
// Example:
//
//	fabIO := readerio.Of[Config](N.Mul(2))
//	result := readerio.MonadFlap(fabIO, 5)(config)() // Returns 10
func MonadFlap[R, B, A any](fab ReaderIO[R, func(A) B], a A) ReaderIO[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap creates a function that applies a value to a ReaderIO function.
// This is the curried version of MonadFlap, suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Result type
//
// Parameters:
//   - a: The value to apply
//
// Returns:
//   - An Operator that applies the value to a ReaderIO function
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](N.Mul(2)),
//	    readerio.Flap[Config](5),
//	)(config)() // Returns 10
//
//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}

// MonadChainReaderK chains a ReaderIO with a function that returns a Reader.
// The Reader is lifted into the ReaderIO context, allowing composition of
// Reader and ReaderIO operations.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - ma: The ReaderIO to chain from
//   - f: Function that produces a Reader
//
// Returns:
//   - A new ReaderIO with the chained Reader computation
//
// Example:
//
//	rio := readerio.Of[Config](5)
//	result := readerio.MonadChainReaderK(rio, func(n int) reader.Reader[Config, int] {
//	    return func(c Config) int { return n + c.Value }
//	})
//
//go:inline
func MonadChainReaderK[R, A, B any](ma ReaderIO[R, A], f reader.Kleisli[R, A, B]) ReaderIO[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain,
		FromReader,
		ma,
		f,
	)
}

// ChainReaderK creates a function that chains a ReaderIO with a Reader-returning function.
// This is the curried version of [MonadChainReaderK], suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input value type
//   - B: Output value type
//
// Parameters:
//   - f: Function that produces a Reader
//
// Returns:
//   - An Operator that chains Reader-returning functions
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](5),
//	    readerio.ChainReaderK(func(n int) reader.Reader[Config, int] {
//	        return func(c Config) int { return n + c.Value }
//	    }),
//	)(config)()
//
//go:inline
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain,
		FromReader,
		f,
	)
}

// MonadChainFirstReaderK chains a function that returns a Reader but keeps the original value.
// The Reader computation is executed for its side effects only.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Reader result type (discarded)
//
// Parameters:
//   - ma: The ReaderIO to chain from
//   - f: Function that produces a Reader
//
// Returns:
//   - A ReaderIO with the original value after executing the Reader
//
// Example:
//
//	rio := readerio.Of[Config](42)
//	result := readerio.MonadChainFirstReaderK(rio, func(n int) reader.Reader[Config, string] {
//	    return func(c Config) string { return fmt.Sprintf("Logged: %d", n) }
//	})
//	value := result(config)() // Returns 42
//
//go:inline
func MonadChainFirstReaderK[R, A, B any](ma ReaderIO[R, A], f reader.Kleisli[R, A, B]) ReaderIO[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

// ChainFirstReaderK creates a function that chains a Reader but keeps the original value.
// This is the curried version of [MonadChainFirstReaderK], suitable for use in pipelines.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Reader result type (discarded)
//
// Parameters:
//   - f: Function that produces a Reader
//
// Returns:
//   - An Operator that chains Reader-returning functions while preserving the original value
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](42),
//	    readerio.ChainFirstReaderK(func(n int) reader.Reader[Config, string] {
//	        return func(c Config) string { return fmt.Sprintf("Logged: %d", n) }
//	    }),
//	)(config)() // Returns 42
//
//go:inline
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReader[R, B],
		f,
	)
}

// MonadTapReaderK chains a function that returns a Reader but keeps the original value.
// This is an alias for [MonadChainFirstReaderK] and is useful for side effects.
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Reader result type (discarded)
//
// Parameters:
//   - ma: The ReaderIO to tap
//   - f: Function that produces a Reader for side effects
//
// Returns:
//   - A ReaderIO with the original value after executing the Reader
//
// Example:
//
//	result := readerio.MonadTapReaderK(
//	    readerio.Of[Config](42),
//	    func(n int) reader.Reader[Config, func()] {
//	        return func(c Config) func() { return func() { fmt.Println(n) } }
//	    },
//	)
//
//go:inline
func MonadTapReaderK[R, A, B any](ma ReaderIO[R, A], f reader.Kleisli[R, A, B]) ReaderIO[R, A] {
	return MonadChainFirstReaderK(ma, f)
}

// TapReaderK creates a function that chains a Reader but keeps the original value.
// This is the curried version of [MonadTapReaderK], an alias for [ChainFirstReaderK].
//
// Type Parameters:
//   - R: Reader environment type
//   - A: Input and output value type
//   - B: Reader result type (discarded)
//
// Parameters:
//   - f: Function that produces a Reader for side effects
//
// Returns:
//   - An Operator that taps with Reader-returning functions
//
// Example:
//
//	result := F.Pipe1(
//	    readerio.Of[Config](42),
//	    readerio.TapReaderK(func(n int) reader.Reader[Config, func()] {
//	        return func(c Config) func() { return func() { fmt.Println(n) } }
//	    }),
//	)(config)() // Returns 42, prints 42
//
//go:inline
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderK(f)
}

// Read executes a ReaderIO with a given environment, returning the resulting IO.
// This is useful for providing the environment dependency and obtaining an IO action
// that can be executed later.
//
// Type Parameters:
//   - A: Result type
//   - R: Reader environment type
//
// Parameters:
//   - r: The environment to provide to the ReaderIO
//
// Returns:
//   - A function that converts a ReaderIO into an IO by applying the environment
//
// Example:
//
//	rio := readerio.Of[Config](42)
//	config := Config{Value: 10, Name: "test"}
//	ioAction := readerio.Read[int](config)(rio)
//	result := ioAction() // Returns 42
//
//go:inline
func Read[A, R any](r R) func(ReaderIO[R, A]) IO[A] {
	return reader.Read[IO[A]](r)
}

// ReadIO executes a ReaderIO computation by providing an environment wrapped in an IO effect.
// This is useful when the environment itself needs to be computed or retrieved through side effects.
//
// The function takes an IO[R] (an effectful computation that produces an environment) and returns
// a function that can execute a ReaderIO[R, A] to produce an IO[A].
//
// This is particularly useful in scenarios where:
//   - The environment needs to be loaded from a file, database, or network
//   - The environment requires initialization with side effects
//   - You want to compose environment retrieval with the computation that uses it
//
// The execution flow is:
//  1. Execute the IO[R] to get the environment R
//  2. Pass the environment to the ReaderIO[R, A] to get an IO[A]
//  3. Execute the resulting IO[A] to get the final result A
//
// Type Parameters:
//   - A: The result type of the ReaderIO computation
//   - R: The environment type required by the ReaderIO
//
// Parameters:
//   - r: An IO effect that produces the environment of type R
//
// Returns:
//   - A function that takes a ReaderIO[R, A] and returns an IO[A]
//
// Example:
//
//	type Config struct {
//	    DatabaseURL string
//	    Port        int
//	}
//
//	// Load config from file (side effect)
//	loadConfig := io.Of(Config{DatabaseURL: "localhost:5432", Port: 8080})
//
//	// A computation that uses the config
//	getConnectionString := readerio.Asks(func(c Config) io.IO[string] {
//	    return io.Of(c.DatabaseURL)
//	})
//
//	// Compose them together
//	result := readerio.ReadIO[string](loadConfig)(getConnectionString)
//	connectionString := result() // Executes both effects and returns "localhost:5432"
//
// Comparison with Read:
//   - [Read]: Takes a pure value R and executes the ReaderIO immediately
//   - [ReadIO]: Takes an IO[R] and chains the effects together
//
//go:inline
func ReadIO[A, R any](r IO[R]) func(ReaderIO[R, A]) IO[A] {
	return function.Flow2(
		io.Chain[R, A],
		Read[A](r),
	)
}

// Delay creates an operation that passes in the value after some delay
//
//go:inline
func Delay[R, A any](delay time.Duration) Operator[R, A, A] {
	return function.Bind2nd(function.Flow2[ReaderIO[R, A]], io.Delay[A](delay))
}

// After creates an operation that passes after the given [time.Time]
//
//go:inline
func After[R, A any](timestamp time.Time) Operator[R, A, A] {
	return function.Bind2nd(function.Flow2[ReaderIO[R, A]], io.After[A](timestamp))
}
