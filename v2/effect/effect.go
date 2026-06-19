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

package effect

import (
	"context"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
)

// FromThunk lifts a Thunk (context-independent IO computation with error handling) into an Effect.
// This allows you to integrate computations that don't need the effect's context type C
// into effect chains. The Thunk will be executed with the runtime context when the effect runs.
//
// # Type Parameters
//
//   - C: The context type required by the effect (not used by the thunk)
//   - A: The type of the success value
//
// # Parameters
//
//   - f: A Thunk[A] that performs IO with error handling
//
// # Returns
//
//   - Effect[C, A]: An effect that ignores its context and executes the thunk
//
// # Example
//
//	thunk := func(ctx context.Context) io.IO[result.Result[int]] {
//	    return func() result.Result[int] {
//	        // Perform IO operation
//	        return result.Of(42)
//	    }
//	}
//
//	eff := effect.FromThunk[MyContext](thunk)
//	// eff can be used in any context but executes the thunk
//
//go:inline
func FromThunk[C, A any](f Thunk[A]) Effect[C, A] {
	return reader.Of[C](f)
}

//go:inline
func FromIO[C, A any](f IO[A]) Effect[C, A] {
	return readerreaderioresult.FromIO[C](f)
}

//go:inline
func FromResult[C, A any](r Result[A]) Effect[C, A] {
	return readerreaderioresult.FromEither[C](r)
}

// Succeed creates a successful Effect that produces the given value.
// This is the primary way to lift a pure value into the Effect context.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The type of the success value
//
// # Parameters
//
//   - a: The value to wrap in a successful effect
//
// # Returns
//
//   - Effect[C, A]: An effect that always succeeds with the given value
//
// # Example
//
//	eff := effect.Succeed[MyContext](42)
//	result, err := runEffect(eff, myContext)
//	// result == 42, err == nil
func Succeed[C, A any](a A) Effect[C, A] {
	return readerreaderioresult.Of[C](a)
}

// Fail creates a failed Effect with the given error.
// This is used to represent computations that have failed.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The type of the success value (never produced)
//
// # Parameters
//
//   - err: The error that caused the failure
//
// # Returns
//
//   - Effect[C, A]: An effect that always fails with the given error
//
// # Example
//
//	eff := effect.Fail[MyContext, int](errors.New("failed"))
//	_, err := runEffect(eff, myContext)
//	// err == errors.New("failed")
func Fail[C, A any](err error) Effect[C, A] {
	return readerreaderioresult.Left[C, A](err)
}

// Of creates a successful Effect that produces the given value.
// This is an alias for Succeed and follows the pointed functor convention.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The type of the success value
//
// # Parameters
//
//   - a: The value to wrap in a successful effect
//
// # Returns
//
//   - Effect[C, A]: An effect that always succeeds with the given value
//
// # Example
//
//	eff := effect.Of[MyContext]("hello")
//	result, err := runEffect(eff, myContext)
//	// result == "hello", err == nil
func Of[C, A any](a A) Effect[C, A] {
	return readerreaderioresult.Of[C](a)
}

// Map transforms the success value of an Effect using the provided function.
// If the effect fails, the error is propagated unchanged.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: The transformation function to apply to the success value
//
// # Returns
//
//   - Operator[C, A, B]: A function that transforms Effect[C, A] to Effect[C, B]
//
// # Example
//
//	eff := effect.Of[MyContext](42)
//	mapped := effect.Map[MyContext](func(x int) string {
//		return strconv.Itoa(x)
//	})(eff)
//	// mapped produces "42"
func Map[C, A, B any](f func(A) B) Operator[C, A, B] {
	return readerreaderioresult.Map[C](f)
}

// Chain sequences two effects, where the second effect depends on the result of the first.
// This is the monadic bind operation (flatMap) for effects.
// If the first effect fails, the second is not executed.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: A function that takes the result of the first effect and returns a new effect
//
// # Returns
//
//   - Operator[C, A, B]: A function that transforms Effect[C, A] to Effect[C, B]
//
// # Example
//
//	eff := effect.Of[MyContext](42)
//	chained := effect.Chain[MyContext](func(x int) Effect[MyContext, string] {
//		return effect.Of[MyContext](strconv.Itoa(x * 2))
//	})(eff)
//	// chained produces "84"
//
//go:inline
func Chain[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, B] {
	return readerreaderioresult.Chain(f)
}

//go:inline
func ChainFirst[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, A] {
	return readerreaderioresult.ChainFirst(f)
}

// ChainFirstThunkK chains an effect with a function that returns a Thunk,
// but discards the result and returns the original value.
// This is useful for performing side effects (like logging or IO operations) that don't
// need the effect's context, without changing the value flowing through the computation.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The value type (preserved)
//   - B: The type produced by the Thunk (discarded)
//
// # Parameters
//
//   - f: A function that takes A and returns Thunk[B] for side effects
//
// # Returns
//
//   - Operator[C, A, A]: A function that executes the Thunk but preserves the original value
//
// # Example
//
//	logToFile := func(n int) readerioresult.ReaderIOResult[any] {
//	    return func(ctx context.Context) io.IO[result.Result[any]] {
//	        return func() result.Result[any] {
//	            // Perform IO operation that doesn't need effect context
//	            fmt.Printf("Logging: %d\n", n)
//	            return result.Of[any](nil)
//	        }
//	    }
//	}
//
//	eff := effect.Of[MyContext](42)
//	logged := effect.ChainFirstThunkK[MyContext](logToFile)(eff)
//	// Prints "Logging: 42" but still produces 42
//
// # See Also
//
//   - ChainThunkK: Chains with a Thunk and uses its result
//   - TapThunkK: Alias for ChainFirstThunkK
//   - ChainFirstIOK: Similar but for IO operations
//
//go:inline
func ChainFirstThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[C, A, B],
		FromThunk[C, B],
		f,
	)
}

// TapThunkK is an alias for ChainFirstThunkK.
// It chains an effect with a function that returns a Thunk for side effects,
// but preserves the original value. This is useful for logging, debugging, or
// performing IO operations that don't need the effect's context.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The value type (preserved)
//   - B: The type produced by the Thunk (discarded)
//
// # Parameters
//
//   - f: A function that takes A and returns Thunk[B] for side effects
//
// # Returns
//
//   - Operator[C, A, A]: A function that executes the Thunk but preserves the original value
//
// # Example
//
//	performSideEffect := func(n int) readerioresult.ReaderIOResult[any] {
//	    return func(ctx context.Context) io.IO[result.Result[any]] {
//	        return func() result.Result[any] {
//	            // Perform context-independent IO operation
//	            log.Printf("Processing value: %d", n)
//	            return result.Of[any](nil)
//	        }
//	    }
//	}
//
//	eff := effect.Of[MyContext](42)
//	tapped := effect.TapThunkK[MyContext](performSideEffect)(eff)
//	// Logs "Processing value: 42" but still produces 42
//
// # See Also
//
//   - ChainFirstThunkK: The underlying implementation
//   - TapIOK: Similar but for IO operations
//   - Tap: Similar but for full effects
//
//go:inline
func TapThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, A] {
	return ChainFirstThunkK[C](f)
}

// ChainIOK chains an effect with a function that returns an IO action.
// This is useful for integrating IO-based computations (synchronous side effects)
// into effect chains. The IO action is automatically lifted into the Effect context.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: A function that takes A and returns IO[B]
//
// # Returns
//
//   - Operator[C, A, B]: A function that chains the IO-returning function with the effect
//
// # Example
//
//	performIO := func(n int) io.IO[string] {
//	    return func() string {
//	        // Perform synchronous side effect
//	        return fmt.Sprintf("Value: %d", n)
//	    }
//	}
//
//	eff := effect.Of[MyContext](42)
//	chained := effect.ChainIOK[MyContext](performIO)(eff)
//	// chained produces "Value: 42"
//
//go:inline
func ChainIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainIOK[C](f)
}

// ChainFirstIOK chains an effect with a function that returns an IO action,
// but discards the result and returns the original value.
// This is useful for performing side effects (like logging) without changing the value.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The value type (preserved)
//   - B: The type produced by the IO action (discarded)
//
// # Parameters
//
//   - f: A function that takes A and returns IO[B] for side effects
//
// # Returns
//
//   - Operator[C, A, A]: A function that executes the IO action but preserves the original value
//
// # Example
//
//	logValue := func(n int) io.IO[any] {
//	    return func() any {
//	        fmt.Printf("Processing: %d\n", n)
//	        return nil
//	    }
//	}
//
//	eff := effect.Of[MyContext](42)
//	logged := effect.ChainFirstIOK[MyContext](logValue)(eff)
//	// Prints "Processing: 42" but still produces 42
//
//go:inline
func ChainFirstIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.ChainFirstIOK[C](f)
}

// TapIOK is an alias for ChainFirstIOK.
// It chains an effect with a function that returns an IO action for side effects,
// but preserves the original value. This is useful for logging, debugging, or
// performing actions without changing the result.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The value type (preserved)
//   - B: The type produced by the IO action (discarded)
//
// # Parameters
//
//   - f: A function that takes A and returns IO[B] for side effects
//
// # Returns
//
//   - Operator[C, A, A]: A function that executes the IO action but preserves the original value
//
// # Example
//
//	logValue := func(n int) io.IO[any] {
//	    return func() any {
//	        fmt.Printf("Value: %d\n", n)
//	        return nil
//	    }
//	}
//
//	eff := effect.Of[MyContext](42)
//	tapped := effect.TapIOK[MyContext](logValue)(eff)
//	// Prints "Value: 42" but still produces 42
//
//go:inline
func TapIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.ChainFirstIOK[C](f)
}

// Ap applies a function wrapped in an Effect to a value wrapped in an Effect.
// This is the applicative apply operation, useful for applying effects in parallel.
//
// # Type Parameters
//
//   - B: The output value type
//   - C: The context type required by the effects
//   - A: The input value type
//
// # Parameters
//
//   - fa: The effect containing the value to apply the function to
//
// # Returns
//
//   - Operator[C, func(A) B, B]: A function that applies the function effect to the value effect
//
// # Example
//
//	fnEff := effect.Of[MyContext](N.Mul(2))
//	valEff := effect.Of[MyContext](21)
//	result := effect.Ap[int](valEff)(fnEff)
//	// result produces 42
func Ap[B, C, A any](fa Effect[C, A]) Operator[C, func(A) B, B] {
	return readerreaderioresult.Ap[B](fa)
}

// Suspend delays the evaluation of an effect until it is run.
// This is useful for recursive effects or when you need lazy evaluation.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The type of the success value
//
// # Parameters
//
//   - fa: A lazy computation that produces an effect
//
// # Returns
//
//   - Effect[C, A]: An effect that evaluates the lazy computation when run
//
// # Example
//
//	var recursiveEff func(int) Effect[MyContext, int]
//	recursiveEff = func(n int) Effect[MyContext, int] {
//		if n <= 0 {
//			return effect.Of[MyContext](0)
//		}
//		return effect.Suspend(func() Effect[MyContext, int] {
//			return effect.Map[MyContext](func(x int) int {
//				return x + n
//			})(recursiveEff(n - 1))
//		})
//	}
func Suspend[C, A any](fa Lazy[Effect[C, A]]) Effect[C, A] {
	return readerreaderioresult.Defer(fa)
}

// Tap executes a side effect for its effect, but returns the original value.
// This is useful for logging, debugging, or performing actions without changing the result.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The value type
//   - ANY: The type produced by the side effect (ignored)
//
// # Parameters
//
//   - f: A function that performs a side effect based on the value
//
// # Returns
//
//   - Operator[C, A, A]: A function that executes the side effect but preserves the original value
//
// # Example
//
//	eff := effect.Of[MyContext](42)
//	tapped := effect.Tap[MyContext](func(x int) Effect[MyContext, any] {
//		fmt.Println("Value:", x)
//		return effect.Of[MyContext, any](nil)
//	})(eff)
//	// Prints "Value: 42" but still produces 42
func Tap[C, A, ANY any](f Kleisli[C, A, ANY]) Operator[C, A, A] {
	return readerreaderioresult.Tap(f)
}

// Ternary creates a conditional effect based on a predicate.
// If the predicate returns true, onTrue is executed; otherwise, onFalse is executed.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - pred: A predicate function to test the input value
//   - onTrue: The effect to execute if the predicate is true
//   - onFalse: The effect to execute if the predicate is false
//
// # Returns
//
//   - Kleisli[C, A, B]: A function that conditionally executes one of two effects
//
// # Example
//
//	kleisli := effect.Ternary(
//		func(x int) bool { return x > 10 },
//		func(x int) Effect[MyContext, string] {
//			return effect.Of[MyContext]("large")
//		},
//		func(x int) Effect[MyContext, string] {
//			return effect.Of[MyContext]("small")
//		},
//	)
//	result := kleisli(15) // produces "large"
func Ternary[C, A, B any](pred Predicate[A], onTrue, onFalse Kleisli[C, A, B]) Kleisli[C, A, B] {
	return function.Ternary(pred, onTrue, onFalse)
}

// ChainResultK chains an effect with a function that returns a Result.
// This is useful for integrating Result-based computations into effect chains.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: A function that takes A and returns Result[B]
//
// # Returns
//
//   - Operator[C, A, B]: A function that chains the Result-returning function with the effect
//
// # Example
//
//	parseIntResult := result.Eitherize1(strconv.Atoi)
//	eff := effect.Of[MyContext]("42")
//	chained := effect.ChainResultK[MyContext](parseIntResult)(eff)
//	// chained produces 42 as an int
//
//go:inline
func ChainResultK[C, A, B any](f result.Kleisli[A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainResultK[C](f)
}

// ChainReaderK chains an effect with a function that returns a Reader.
// This is useful for integrating Reader-based computations (pure context-dependent functions)
// into effect chains. The Reader is automatically lifted into the Effect context.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: A function that takes A and returns Reader[C, B]
//
// # Returns
//
//   - Operator[C, A, B]: A function that chains the Reader-returning function with the effect
//
// # Example
//
//	type Config struct { Multiplier int }
//
//	getMultiplied := func(n int) reader.Reader[Config, int] {
//	    return func(cfg Config) int {
//	        return n * cfg.Multiplier
//	    }
//	}
//
//	eff := effect.Of[Config](5)
//	chained := effect.ChainReaderK[Config](getMultiplied)(eff)
//	// With Config{Multiplier: 3}, produces 15
//
//go:inline
func ChainReaderK[C, A, B any](f reader.Kleisli[C, A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainReaderK(f)
}

// ChainThunkK chains an effect with a function that returns a Thunk.
// This is useful for integrating Thunk-based computations (context-independent IO with error handling)
// into effect chains. The Thunk is automatically lifted into the Effect context.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: A function that takes A and returns Thunk[B] (readerioresult.Kleisli[A, B])
//
// # Returns
//
//   - Operator[C, A, B]: A function that chains the Thunk-returning function with the effect
//
// # Example
//
//	performIO := func(n int) readerioresult.ReaderIOResult[string] {
//	    return func(ctx context.Context) io.IO[result.Result[string]] {
//	        return func() result.Result[string] {
//	            // Perform IO operation that doesn't need effect context
//	            return result.Of(fmt.Sprintf("Processed: %d", n))
//	        }
//	    }
//	}
//
//	eff := effect.Of[MyContext](42)
//	chained := effect.ChainThunkK[MyContext](performIO)(eff)
//	// chained produces "Processed: 42"
//
//go:inline
func ChainThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, B] {
	return fromreader.ChainReaderK(
		Chain[C, A, B],
		FromThunk[C, B],
		f,
	)
}

// ChainReaderIOK chains an effect with a function that returns a ReaderIO.
// This is useful for integrating ReaderIO-based computations (context-dependent IO operations)
// into effect chains. The ReaderIO is automatically lifted into the Effect context.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input value type
//   - B: The output value type
//
// # Parameters
//
//   - f: A function that takes A and returns ReaderIO[C, B]
//
// # Returns
//
//   - Operator[C, A, B]: A function that chains the ReaderIO-returning function with the effect
//
// # Example
//
//	type Config struct { LogPrefix string }
//
//	logAndDouble := func(n int) readerio.ReaderIO[Config, int] {
//	    return func(cfg Config) io.IO[int] {
//	        return func() int {
//	            fmt.Printf("%s: %d\n", cfg.LogPrefix, n)
//	            return n * 2
//	        }
//	    }
//	}
//
//	eff := effect.Of[Config](21)
//	chained := effect.ChainReaderIOK[Config](logAndDouble)(eff)
//	// Logs "prefix: 21" and produces 42
//
//go:inline
func ChainReaderIOK[C, A, B any](f readerio.Kleisli[C, A, B]) Operator[C, A, B] {
	return readerreaderioresult.ChainReaderIOK(f)
}

// Read provides a context to an effect, partially applying it.
// This converts an Effect[C, A] to a Thunk[A] by supplying the required context.
//
// # Type Parameters
//
//   - A: The type of the success value
//   - C: The context type
//
// # Parameters
//
//   - c: The context to provide to the effect
//
// # Returns
//
//   - func(Effect[C, A]) Thunk[A]: A function that converts an effect to a thunk
//
// # Example
//
//	ctx := MyContext{Value: "test"}
//	eff := effect.Of[MyContext](42)
//	thunk := effect.Read[int](ctx)(eff)
//	// thunk is now a Thunk[int] that can be run without context
//
//go:inline
func Read[A, C any](c C) func(Effect[C, A]) Thunk[A] {
	return readerreaderioresult.Read[A](c)
}

// ReadIO provides a context from an IO computation to an effect, partially applying it.
// This converts an Effect[C, A] to a Thunk[A] by supplying the required context through
// an IO action. This is useful when the context itself needs to be computed or retrieved
// through side effects.
//
// Type Parameters:
//   - A: The type of the success value
//   - C: The context type
//
// Parameters:
//   - c: An IO computation that produces the context
//
// Returns:
//   - func(Effect[C, A]) Thunk[A]: A function that converts an effect to a thunk
//
// See Also:
//   - Read: Provides a pure context value instead of an IO computation
//   - Asks: Projects a value from the context
//
//go:inline
func ReadIO[A, C any](c IO[C]) func(Effect[C, A]) Thunk[A] {
	return readerreaderioresult.ReadIO[A](c)
}

// Asks creates an Effect that projects a value from the context using a Reader function.
// This is useful for extracting specific fields or computing derived values from the context.
// It's essentially a lifted version of the Reader pattern into the Effect context.
//
// # Type Parameters
//
//   - C: The context type
//   - A: The type of the projected value
//
// # Parameters
//
//   - r: A Reader function that extracts or computes a value from the context
//
// # Returns
//
//   - Effect[C, A]: An effect that succeeds with the projected value
//
// # Example
//
//	type Config struct {
//		Host string
//		Port int
//	}
//
//	// Extract a specific field
//	getHost := effect.Asks[Config](func(cfg Config) string {
//		return cfg.Host
//	})
//
//	// Compute a derived value
//	getURL := effect.Asks[Config](func(cfg Config) string {
//		return fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
//	})
//
//	result, err := runEffect(getHost, Config{Host: "localhost", Port: 8080})
//	// result == "localhost", err == nil
//
// # See Also
//
// See Also:
//
//   - Ask: Returns the entire context as the value
//   - Map: Transforms the value after extraction
//
//go:inline
func Asks[C, A any](r Reader[C, A]) Effect[C, A] {
	return readerreaderioresult.Asks(r)
}

// Paired converts an [Effect] into a single-argument function that accepts a [Pair]
// bundling the context.Context (head) and the outer environment R (tail).
//
// This is a thin wrapper over [readerreaderioresult.Paired]; see that function for
// the full rationale. In short: R sits in the tail because the tail is the primary
// value that [pair.Map] and other functor operations act on, while context.Context
// sits in the head as auxiliary threading data that passes through unchanged.
//
// # Type Parameters
//
//   - R: The context type required by the effect (outer environment)
//   - A: The type of the success value
//
// # Parameters
//
//   - f: The effect to convert
//
// # Returns
//
//   - func(Pair[context.Context, R]) IOResult[A]: A function that accepts a bundled pair
//     and runs the effect, equivalent to f(pair.Tail(p))(pair.Head(p))
//
// # Example
//
//	type Config struct{ BaseURL string }
//
//	fetch := effect.Of[Config]("hello")
//	paired := effect.Paired(fetch)
//
//	p := pair.MakePair[context.Context, Config](ctx, Config{BaseURL: "http://example.com"})
//	res := paired(p)()  // Result[string]
func Paired[R, A any](f Effect[R, A]) ioresult.Kleisli[Pair[context.Context, R], A] {
	return readerreaderioresult.Paired(f)
}

// MonadChainLeft handles errors by chaining a recovery computation.
// If the effect fails, the error is passed to f which can produce a recovery effect.
// If the effect succeeds, its value is returned unchanged.
// This is the monadic version that takes the computation as the first parameter.
//
// # Type Parameters
//
//   - R: The context type required by the effects
//   - A: The value type
//
// # Parameters
//
//   - fa: The effect that may fail
//   - f: A function that takes an error and returns a recovery effect
//
// # Returns
//
//   - Effect[R, A]: An effect that either succeeds with the original value or recovers from the error
//
// # Example
//
//	type Config struct{ RetryCount int }
//
//	fetchData := effect.Fail[Config, string](errors.New("network error"))
//
//	recover := func(err error) effect.Effect[Config, string] {
//	    return effect.Of[Config]("fallback data")
//	}
//
//	result := effect.MonadChainLeft(fetchData, recover)
//	// result produces "fallback data" instead of failing
//
// # See Also
//
//   - ChainLeft: The curried version that returns an operator
//   - MonadAlt: Alternative composition without error inspection
func MonadChainLeft[R, A any](fa Effect[R, A], f Kleisli[R, error, A]) Effect[R, A] {
	return readerreaderioresult.MonadChainLeft(fa, f)
}

// ChainLeft handles errors by chaining a recovery computation.
// If the effect fails, the error is passed to f which can produce a recovery effect.
// If the effect succeeds, its value is returned unchanged.
// This is the curried version that returns an operator.
//
// # Type Parameters
//
//   - R: The context type required by the effects
//   - A: The value type
//
// # Parameters
//
//   - f: A function that takes an error and returns a recovery effect
//
// # Returns
//
//   - Operator[R, A, A]: A function that transforms a failing effect into a recovered one
//
// # Example
//
//	type Config struct{ MaxRetries int }
//
//	recoverFromError := func(err error) effect.Effect[Config, int] {
//	    return effect.Asks[Config](func(cfg Config) int {
//	        return cfg.MaxRetries
//	    })
//	}
//
//	pipeline := F.Pipe1(
//	    effect.Fail[Config, int](errors.New("operation failed")),
//	    effect.ChainLeft[Config](recoverFromError),
//	)
//	// With Config{MaxRetries: 3}, produces 3
//
// # See Also
//
//   - MonadChainLeft: The monadic version that takes the computation first
//   - Alt: Alternative composition without error inspection
func ChainLeft[R, A any](f Kleisli[R, error, A]) Operator[R, A, A] {
	return readerreaderioresult.ChainLeft(f)
}

// MonadAlt provides alternative/fallback behavior for effects.
// If the first effect fails, it tries the second effect (which is lazy-evaluated).
// If the first effect succeeds, its value is returned and the second effect is never evaluated.
// This is the monadic version that takes both effects as parameters.
//
// # Type Parameters
//
//   - R: The context type required by the effects
//   - A: The value type
//
// # Parameters
//
//   - first: The primary effect to try
//   - second: A lazy computation that produces the fallback effect (only evaluated if first fails)
//
// # Returns
//
//   - Effect[R, A]: An effect that succeeds with the first successful result, or fails if both fail
//
// # Example
//
//	type Config struct{ PrimaryURL, FallbackURL string }
//
//	fetchFromPrimary := effect.Fail[Config, string](errors.New("primary unavailable"))
//	fetchFromFallback := func() effect.Effect[Config, string] {
//	    return effect.Of[Config]("data from fallback")
//	}
//
//	result := effect.MonadAlt(fetchFromPrimary, fetchFromFallback)
//	// result produces "data from fallback"
//
// # See Also
//
//   - Alt: The curried version that returns an operator
//   - MonadChainLeft: Similar but allows error inspection
func MonadAlt[R, A any](first Effect[R, A], second Lazy[Effect[R, A]]) Effect[R, A] {
	return readerreaderioresult.MonadAlt(first, second)
}

// Alt provides alternative/fallback behavior for effects.
// If the first effect fails, it tries the second effect (which is lazy-evaluated).
// If the first effect succeeds, its value is returned and the second effect is never evaluated.
// This is the curried version that returns an operator.
//
// # Type Parameters
//
//   - R: The context type required by the effects
//   - A: The value type
//
// # Parameters
//
//   - second: A lazy computation that produces the fallback effect (only evaluated if first fails)
//
// # Returns
//
//   - Operator[R, A, A]: A function that provides fallback behavior for an effect
//
// # Example
//
//	type Config struct{ Endpoints []string }
//
//	tryEndpoint := func(url string) effect.Effect[Config, string] {
//	    if url == "backup.api" {
//	        return effect.Of[Config]("success from backup")
//	    }
//	    return effect.Fail[Config, string](fmt.Errorf("%s failed", url))
//	}
//
//	pipeline := F.Pipe2(
//	    tryEndpoint("primary.api"),
//	    effect.Alt[Config](func() effect.Effect[Config, string] {
//	        return tryEndpoint("backup.api")
//	    }),
//	)
//	// pipeline produces "success from backup"
//
// # See Also
//
//   - MonadAlt: The monadic version that takes both effects
//   - ChainLeft: Similar but allows error inspection
func Alt[R, A any](second Lazy[Effect[R, A]]) Operator[R, A, A] {
	return readerreaderioresult.Alt(second)
}

// ChainFirstLeft chains a computation on the error path but preserves the original value.
// If the effect succeeds, the original value is returned unchanged.
// If it fails, the error handler f is executed, and its result determines the final outcome.
// This is the curried version that returns an operator.
//
// Use cases:
//   - Error logging without changing the error
//   - Error recovery with fallback logic
//   - Side effects on error path
//
// See Also:
//   - TapLeft: Alias for this function
//   - MonadChainFirstLeft: Monadic version
//   - ChainLeft: Similar but replaces the error
func ChainFirstLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstLeft[A](f)
}

// MonadChainFirstLeft chains a computation on the error path but preserves the original value.
// If the effect succeeds, the original value is returned unchanged.
// If it fails, the error handler f is executed, and its result determines the final outcome.
// This is the monadic version that takes the effect as the first parameter.
//
// See Also:
//   - MonadTapLeft: Alias for this function
//   - ChainFirstLeft: Curried version
func MonadChainFirstLeft[R, A, B any](ma Effect[R, A], f Kleisli[R, error, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstLeft(ma, f)
}

// TapLeft is an alias for ChainFirstLeft.
// Executes a side effect on the error path while preserving the original value or error.
//
// Common use cases:
//   - Logging errors without modifying them
//   - Sending error notifications
//   - Recording error metrics
//
// See Also:
//   - ChainFirstLeft: The underlying implementation
//   - MonadTapLeft: Monadic version
func TapLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return readerreaderioresult.TapLeft[A](f)
}

// MonadTapLeft is an alias for MonadChainFirstLeft.
// Executes a side effect on the error path while preserving the original value or error.
// This is the monadic version that takes the effect as the first parameter.
//
// See Also:
//   - TapLeft: Curried version
//   - MonadChainFirstLeft: The underlying implementation
func MonadTapLeft[R, A, B any](ma Effect[R, A], f Kleisli[R, error, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapLeft(ma, f)
}

// ChainFirstLeftIOK chains an IO computation on the error path but preserves the original value.
// The IO computation is automatically lifted into Effect.
// This is the curried version that returns an operator.
//
// Use cases:
//   - Logging errors to console or file
//   - Sending error notifications via IO
//   - Recording metrics on error
//
// See Also:
//   - TapLeftIOK: Alias for this function
//   - MonadChainFirstLeftIOK: Monadic version
func ChainFirstLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstLeftIOK[A, R](f)
}

// MonadChainFirstLeftIOK chains an IO computation on the error path but preserves the original value.
// The IO computation is automatically lifted into Effect.
// This is the monadic version that takes the effect as the first parameter.
//
// See Also:
//   - MonadTapLeftIOK: Alias for this function
//   - ChainFirstLeftIOK: Curried version
func MonadChainFirstLeftIOK[R, A, B any](ma Effect[R, A], f io.Kleisli[error, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstLeftIOK(ma, f)
}

// TapLeftIOK is an alias for ChainFirstLeftIOK.
// Executes an IO side effect on the error path while preserving the original value or error.
//
// Common use cases:
//   - Logging errors to console: func(e error) io.IO[Void] { return func() Void { fmt.Println(e); return VOID } }
//   - Writing errors to file
//   - Sending error notifications
//
// See Also:
//   - ChainFirstLeftIOK: The underlying implementation
//   - MonadTapLeftIOK: Monadic version
func TapLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A] {
	return readerreaderioresult.TapLeftIOK[A, R](f)
}

// MonadTapLeftIOK is an alias for MonadChainFirstLeftIOK.
// Executes an IO side effect on the error path while preserving the original value or error.
// This is the monadic version that takes the effect as the first parameter.
//
// See Also:
//   - TapLeftIOK: Curried version
//   - MonadChainFirstLeftIOK: The underlying implementation
func MonadTapLeftIOK[R, A, B any](ma Effect[R, A], f io.Kleisli[error, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapLeftIOK(ma, f)
}

// ChainFirstLeftThunkK chains a Thunk computation on the error path but preserves the original value.
// If the effect succeeds, the original value is returned unchanged.
// If it fails, the error handler f is executed with the runtime context, and its result determines the final outcome.
//
// This function is similar to ChainFirstLeft but accepts a Thunk-based Kleisli arrow instead of a full Effect Kleisli.
// A Thunk is a context-independent computation that only needs the runtime context.Context, making it useful for
// error handlers that don't need access to the effect's context type C.
//
// The key difference from ChainFirstLeftIOK is that Thunk computations have access to context.Context,
// enabling cancellation, timeouts, and context values, while IO computations do not.
//
// Type Parameters:
//   - C: The context type required by the effect
//   - A: The success type of the effect
//   - B: The result type of the error handler (typically discarded)
//
// Parameters:
//   - f: A Thunk Kleisli arrow that takes an error and returns a Thunk[B]
//
// Returns:
//   - An Operator that preserves the original value or error after executing the handler
//
// Example with error logging:
//
//	logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
//	    return func(ctx context.Context) io.IO[result.Result[F.Void]] {
//	        return func() result.Result[F.Void] {
//	            slog.ErrorContext(ctx, "Operation failed", "error", err)
//	            return result.Of(F.VOID)
//	        }
//	    }
//	}
//
//	pipeline := F.Pipe2(
//	    fetchData[Config](id),
//	    ChainFirstLeftThunkK[Config, Data](logError),
//	    Map(processData),
//	)
//
// Example with error recovery:
//
//	recordError := func(err error) readerioresult.ReaderIOResult[F.Void] {
//	    return func(ctx context.Context) io.IO[result.Result[F.Void]] {
//	        return func() result.Result[F.Void] {
//	            if dbErr := recordToDatabase(ctx, err); dbErr != nil {
//	                return result.Left[F.Void](dbErr)
//	            }
//	            return result.Of(F.VOID)
//	        }
//	    }
//	}
//
//	pipeline := F.Pipe2(
//	    performOperation[Config](data),
//	    ChainFirstLeftThunkK[Config, Result](recordError),
//	    OrElse(fallbackOperation),
//	)
//
// Use Cases:
//   - Error logging with context (cancellation, request IDs)
//   - Recording errors to external systems (databases, monitoring)
//   - Sending error notifications with timeout handling
//   - Error recovery with context-aware operations
//
// See Also:
//   - TapLeftThunkK: Alias for this function
//   - ChainFirstLeft: Similar but requires full Effect context
//   - ChainFirstLeftIOK: Similar but without context.Context access
//   - TapLeft: For error handlers that need the effect's context type
//
//go:inline
func ChainFirstLeftThunkK[C, A, B any](f thunk.Kleisli[error, B]) Operator[C, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirstLeft[A, C, B],
		FromThunk[C, B],
		f,
	)
}

// TapLeftThunkK is an alias for ChainFirstLeftThunkK.
// Executes a Thunk side effect on the error path while preserving the original value or error.
//
// This function is ideal for error handling scenarios where you need access to context.Context
// but don't need the effect's context type C. Common use cases include logging with context,
// recording errors to external systems, and sending notifications with timeout handling.
//
// The key advantage over TapLeftIOK is access to context.Context, enabling:
//   - Cancellation and timeout handling
//   - Request-scoped values (trace IDs, user info)
//   - Deadline propagation
//
// Type Parameters:
//   - C: The context type required by the effect
//   - A: The success type of the effect
//   - B: The result type of the error handler (typically F.Void)
//
// Parameters:
//   - f: A Thunk Kleisli arrow that takes an error and returns a Thunk[B]
//
// Returns:
//   - An Operator that preserves the original value or error after executing the handler
//
// Use Cases:
//   - Logging errors with context-aware loggers
//   - Recording errors with cancellation support
//   - Sending notifications with timeout handling
//   - Error metrics with request tracing
//
// See Also:
//   - ChainFirstLeftThunkK: The underlying implementation
//   - TapLeft: For error handlers that need the effect's context type
//   - TapLeftIOK: For simpler error handlers without context.Context
//
//go:inline
func TapLeftThunkK[C, A, B any](f thunk.Kleisli[error, B]) Operator[C, A, A] {
	return ChainFirstLeftThunkK[C, A](f)
}
