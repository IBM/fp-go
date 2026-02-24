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
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/fromreader"
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
func Chain[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, B] {
	return readerreaderioresult.Chain(f)
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
//	fnEff := effect.Of[MyContext](func(x int) int { return x * 2 })
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
