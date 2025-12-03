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

package readerio

import (
	"sync"

	"github.com/IBM/fp-go/v2/function"
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
//	doubled := readerio.MonadMap(rio, func(n int) int { return n * 2 })
//	result := doubled(config)() // Returns 10
func MonadMap[R, A, B any](fa ReaderIO[R, A], f func(A) B) ReaderIO[R, B] {
	return readert.MonadMap[ReaderIO[R, A], ReaderIO[R, B]](io.MonadMap[A, B], fa, f)
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
//	    readerio.Map[Config](func(n int) int { return n * 2 }),
//	)(config)() // Returns 10
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return readert.Map[ReaderIO[R, A], ReaderIO[R, B]](io.Map[A, B], f)
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
func MonadChain[R, A, B any](ma ReaderIO[R, A], f func(A) ReaderIO[R, B]) ReaderIO[R, B] {
	return readert.MonadChain(io.MonadChain[A, B], ma, f)
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
func Chain[R, A, B any](f func(A) ReaderIO[R, B]) Operator[R, A, B] {
	return readert.Chain[ReaderIO[R, A]](io.Chain[A, B], f)
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
//	fabIO := readerio.Of[Config](func(n int) int { return n * 2 })
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
//	    readerio.Of[Config](func(n int) int { return n * 2 }),
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
func MonadChainIOK[R, A, B any](ma ReaderIO[R, A], f func(A) IO[B]) ReaderIO[R, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, A, B],
		FromIO[R, B],
		ma, f,
	)
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
func ChainIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, B] {
	return fromio.ChainIOK(
		Chain[R, A, B],
		FromIO[R, B],
		f,
	)
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
//	fabIO := readerio.Of[Config](func(n int) int { return n * 2 })
//	result := readerio.MonadFlap(fabIO, 5)(config)() // Returns 10
func MonadFlap[R, A, B any](fab ReaderIO[R, func(A) B], a A) ReaderIO[R, B] {
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
//	    readerio.Of[Config](func(n int) int { return n * 2 }),
//	    readerio.Flap[Config](5),
//	)(config)() // Returns 10
func Flap[R, A, B any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}
