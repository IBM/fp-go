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
// # See Also
//
//   - FromIO: Similar but for infallible IO operations
//   - ChainThunkK: Chains a Thunk-returning function over an existing Effect
//
//go:inline
func FromThunk[C, A any](f Thunk[A]) Effect[C, A] {
	return reader.Of[C](f)
}

// FromIO lifts an infallible IO action into an Effect.
// The IO is not executed until the effect is run. Unlike FromThunk, the IO action
// cannot fail and does not see the runtime context.Context.
//
// # Type Parameters
//
//   - C: The context type required by the effect (not used by the IO)
//   - A: The type of the value produced by the IO
//
// # Parameters
//
//   - f: An IO[A] action to lift
//
// # Returns
//
//   - Effect[C, A]: An effect that always succeeds with the value produced by f
//
// # See Also
//
//   - FromThunk: Similar but the action can fail and sees context.Context
//   - FromResult: Lifts an already-computed Result
//
//go:inline
func FromIO[C, A any](f IO[A]) Effect[C, A] {
	return readerreaderioresult.FromIO[C](f)
}

// FromResult lifts an already-computed Result into an Effect.
// A Right value becomes a successful Effect; a Left value becomes a failed Effect.
//
// # Type Parameters
//
//   - C: The context type required by the effect (not used)
//   - A: The type of the success value
//
// # Parameters
//
//   - r: The Result to lift
//
// # Returns
//
//   - Effect[C, A]: An effect that succeeds or fails according to r
//
// # See Also
//
//   - FromIO: Lifts an IO computation that cannot fail
//   - FromReaderResult: Lifts a context-dependent fallible computation
//
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
// # See Also
//
//   - Of: Alias for Succeed following the pointed-functor convention
//   - Fail: Constructs a failing Effect
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
// # See Also
//
//   - ChainLeft: Recover from a failure
//   - Alt: Fall back to another Effect on failure
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
// # See Also
//
//   - Succeed: Synonym for Of
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
// # See Also
//
//   - Chain: Map where the transformation itself returns an Effect
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
// # See Also
//
//   - Map: Chain where the transformation is a pure function
//   - ChainFirst: Chain that discards the result of the inner effect
//
//go:inline
func Chain[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, B] {
	return readerreaderioresult.Chain(f)
}

// ChainFirst sequences an effect for its side effect, discarding the inner result
// and returning the original value unchanged.
// If the inner effect fails, the failure propagates; if the outer effect fails,
// f is never called.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The value type (preserved)
//   - B: The type produced by the inner effect (discarded)
//
// # Parameters
//
//   - f: A Kleisli arrow executed for its side effect
//
// # Returns
//
//   - Operator[C, A, A]: A function that executes f but preserves the original value
//
// # See Also
//
//   - Tap: Alias for ChainFirst with a name that signals observation intent
//   - Chain: Like ChainFirst but uses the inner result
//
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
// # See Also
//
//   - ChainThunkK: Similar but the IO may fail and sees context.Context
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
// # See Also
//
//   - TapIOK: Alias for ChainFirstIOK
//   - ChainFirstThunkK: Similar but the IO may fail and sees context.Context
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
// # See Also
//
//   - ChainFirstIOK: The underlying implementation
//   - TapThunkK: Similar but the IO may fail and sees context.Context
//
//go:inline
func TapIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A] {
	return readerreaderioresult.ChainFirstIOK[C](f)
}

// Ap applies a function wrapped in an Effect to a value wrapped in an Effect.
// This is the applicative apply operation. Because neither effect depends on the
// other's result, both are run in parallel; the result fails if either one fails.
//
// B comes first in the type parameter list so it can be given explicitly
// (e.g. Ap[int](fa)) while C and A are inferred from fa.
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
// # See Also
//
//   - Chain: When the second effect depends on the result of the first
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
// # See Also
//
//   - Of: Constructs an immediately-resolved Effect
func Suspend[C, A any](fa Lazy[Effect[C, A]]) Effect[C, A] {
	return readerreaderioresult.Defer(fa)
}

// Tap runs an effect derived from the current value for its side effect only and
// returns the original value unchanged. If the tapped effect fails, the failure propagates.
// This is useful for logging, debugging, or auditing without changing the result.
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
// # See Also
//
//   - ChainFirst: Alias for Tap
//   - TapIOK: Similar but for IO-returning functions
//   - TapThunkK: Similar but for Thunk-returning functions
func Tap[C, A, ANY any](f Kleisli[C, A, ANY]) Operator[C, A, A] {
	return readerreaderioresult.Tap(f)
}

// Ternary returns a Kleisli arrow that dispatches to onTrue when pred holds for
// the input and to onFalse otherwise.
//
// Deprecated: Use predicate.Fold instead. Note that Fold takes the branches in
// the opposite order and the predicate last: predicate.Fold(onFalse, onTrue)(pred).
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
// # See Also
//
//   - FromResult: Lifts a single Result without an upstream Effect
//   - ChainThunkK: Similar but the computation performs IO and sees context.Context
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
// # See Also
//
//   - Asks: Read a value from the context without an upstream value
//   - ChainReaderIOK: Similar but the computation also performs IO
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
// # See Also
//
//   - ChainFirstThunkK: Like ChainThunkK but discards the Thunk result
//   - ChainIOK: Similar but the IO cannot fail and does not see context.Context
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
// # See Also
//
//   - ChainReaderK: Similar but the computation does not perform IO
//   - ChainThunkK: Similar but the computation may fail
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
// # See Also
//
//   - ReadIO: Similar but the context is produced by an IO computation
//   - Provide: Synonym used in the run pipeline
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
// # Type Parameters
//
//   - A: The type of the success value
//   - C: The context type
//
// # Parameters
//
//   - c: An IO computation that produces the context
//
// # Returns
//
//   - func(Effect[C, A]) Thunk[A]: A function that converts an effect to a thunk
//
// # See Also
//
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
// # See Also
//
//   - Ask: Returns the entire context as the value
//   - FromReader: Canonical lifting function when you already have a Reader
//   - Map: Transforms the projected value
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
// value that pair.Map and other functor operations act on, while context.Context
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
// # See Also
//
//   - Read: Provides the context as a plain value instead of a Pair
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
//   - Effect[R, A]: An effect that succeeds with the first successful result, or fails with
//     the fallback's error if both fail
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
// # See Also
//
//   - MonadAlt: The monadic version that takes both effects
//   - ChainLeft: Similar but allows error inspection
func Alt[R, A any](second Lazy[Effect[R, A]]) Operator[R, A, A] {
	return readerreaderioresult.Alt(second)
}

// ChainFirstLeft runs an effect on the error path for its side effect only.
// If the upstream effect succeeds, f is not called and the value passes through.
// If it fails, f is called with the error and its effect is executed, but the
// original error is always what propagates: f cannot recover from the failure
// (use ChainLeft for that), and a failure of f itself is discarded.
// This is the curried version that returns an operator.
//
// Because A only appears in the result type, it must be given explicitly,
// e.g. ChainFirstLeft[int](logError).
//
// # Type Parameters
//
//   - A: The success type of the effect (must be specified)
//   - R: The context type required by the effects
//   - B: The result type of the handler (discarded)
//
// # Parameters
//
//   - f: A Kleisli arrow from the error to an effect executed for its side effect
//
// # Returns
//
//   - Operator[R, A, A]: An operator that preserves the original value or error
//
// # See Also
//
//   - TapLeft: Alias for this function
//   - MonadChainFirstLeft: Monadic version
//   - ChainLeft: Recovers from or replaces the error
//   - ChainFirstLeftThunkK: Similar but the handler is a Thunk
func ChainFirstLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstLeft[A](f)
}

// MonadChainFirstLeft runs an effect on the error path for its side effect only.
// It is the uncurried form of ChainFirstLeft with the same semantics: f runs only
// when ma fails, and the original error always propagates, even if f itself fails.
//
// # See Also
//
//   - MonadTapLeft: Alias for this function
//   - ChainFirstLeft: Curried version
func MonadChainFirstLeft[R, A, B any](ma Effect[R, A], f Kleisli[R, error, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstLeft(ma, f)
}

// TapLeft is an alias for ChainFirstLeft.
// Executes a side effect on the error path while preserving the original value or error.
//
// # See Also
//
//   - ChainFirstLeft: The underlying implementation
//   - MonadTapLeft: Monadic version
func TapLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return readerreaderioresult.TapLeft[A](f)
}

// MonadTapLeft is an alias for MonadChainFirstLeft.
// Executes a side effect on the error path while preserving the original value or error.
// This is the monadic version that takes the effect as the first parameter.
//
// # See Also
//
//   - TapLeft: Curried version
//   - MonadChainFirstLeft: The underlying implementation
func MonadTapLeft[R, A, B any](ma Effect[R, A], f Kleisli[R, error, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapLeft(ma, f)
}

// ChainFirstLeftIOK runs an IO action on the error path for its side effect only.
// The action receives the error; the original value or error passes through unchanged.
// This is the curried version that returns an operator.
//
// Neither A nor R can be inferred from f, so both must be given explicitly,
// e.g. ChainFirstLeftIOK[int, Config](logError).
//
// # See Also
//
//   - TapLeftIOK: Alias for this function
//   - MonadChainFirstLeftIOK: Monadic version
//   - ChainFirstLeftThunkK: Similar but the IO may fail and sees context.Context
func ChainFirstLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A] {
	return readerreaderioresult.ChainFirstLeftIOK[A, R](f)
}

// MonadChainFirstLeftIOK chains an IO computation on the error path but preserves the original value.
// The IO computation is automatically lifted into Effect.
// This is the monadic version that takes the effect as the first parameter.
//
// # See Also
//
//   - MonadTapLeftIOK: Alias for this function
//   - ChainFirstLeftIOK: Curried version
func MonadChainFirstLeftIOK[R, A, B any](ma Effect[R, A], f io.Kleisli[error, B]) Effect[R, A] {
	return readerreaderioresult.MonadChainFirstLeftIOK(ma, f)
}

// TapLeftIOK is an alias for ChainFirstLeftIOK.
// Executes an IO side effect on the error path while preserving the original value or error.
//
// # See Also
//
//   - ChainFirstLeftIOK: The underlying implementation
//   - MonadTapLeftIOK: Monadic version
func TapLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A] {
	return readerreaderioresult.TapLeftIOK[A, R](f)
}

// MonadTapLeftIOK is an alias for MonadChainFirstLeftIOK.
// Executes an IO side effect on the error path while preserving the original value or error.
// This is the monadic version that takes the effect as the first parameter.
//
// # See Also
//
//   - TapLeftIOK: Curried version
//   - MonadChainFirstLeftIOK: The underlying implementation
func MonadTapLeftIOK[R, A, B any](ma Effect[R, A], f io.Kleisli[error, B]) Effect[R, A] {
	return readerreaderioresult.MonadTapLeftIOK(ma, f)
}

// ChainFirstLeftThunkK chains a Thunk computation on the error path but preserves the original value.
// If the effect succeeds, the original value is returned unchanged.
// If it fails, f is executed with the runtime context, but the original error always
// propagates: f cannot recover from the failure, and a failure of f itself is discarded.
//
// This function is similar to ChainFirstLeft but accepts a Thunk-based Kleisli arrow instead of a full Effect Kleisli.
// A Thunk is a context-independent computation that only needs the runtime context.Context, making it useful for
// error handlers that don't need access to the effect's context type C.
//
// The key difference from ChainFirstLeftIOK is that Thunk computations have access to context.Context,
// enabling cancellation, timeouts, and context values, while IO computations do not.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The success type of the effect
//   - B: The result type of the error handler (typically discarded)
//
// # Parameters
//
//   - f: A Thunk Kleisli arrow that takes an error and returns a Thunk[B]
//
// # Returns
//
//   - Operator[C, A, A]: An operator that preserves the original value or error after executing the handler
//
// # See Also
//
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
// The key advantage over TapLeftIOK is access to context.Context, enabling
// cancellation, timeouts, and request-scoped values (trace IDs, deadlines).
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The success type of the effect
//   - B: The result type of the error handler (typically F.Void)
//
// # Parameters
//
//   - f: A Thunk Kleisli arrow that takes an error and returns a Thunk[B]
//
// # Returns
//
//   - Operator[C, A, A]: An operator that preserves the original value or error after executing the handler
//
// # See Also
//
//   - ChainFirstLeftThunkK: The underlying implementation
//   - TapLeft: For error handlers that need the effect's context type
//   - TapLeftIOK: For simpler error handlers without context.Context
//
//go:inline
func TapLeftThunkK[C, A, B any](f thunk.Kleisli[error, B]) Operator[C, A, A] {
	return ChainFirstLeftThunkK[C, A](f)
}

// FromReader lifts a Reader into an Effect, wrapping the Reader's result as a
// success value. Unlike Asks, which emphasises the "ask the environment" pattern,
// FromReader is the canonical lifting function used when you already have a Reader
// and want to treat it as an Effect without additional transformation.
//
// # Type Parameters
//
//   - R: The outer context (environment) type consumed by the Reader
//   - A: The type of the value produced by the Reader
//
// # Parameters
//
//   - ma: A Reader[R, A] that reads from environment R and returns A
//
// # Returns
//
//   - Effect[R, A]: An effect that always succeeds with the value produced by ma
//
// # See Also
//
//   - Asks: Alternative constructor that emphasises environment projection
//   - FromResult: Lifts a pre-computed Result into an Effect
//   - FromIO: Lifts a context-independent IO computation into an Effect
//
//go:inline
func FromReader[R, A any](ma Reader[R, A]) Effect[R, A] {
	return readerreaderioresult.FromReader(ma)
}

// FromReaderResult lifts a ReaderResult into an Effect. A ReaderResult is a
// function from an environment R to a Result[A], so it can already encode both
// success and failure. FromReaderResult preserves that structure: a Right value
// inside the ReaderResult becomes a successful Effect, and a Left value becomes
// a failed Effect carrying the original error.
//
// This is the canonical way to integrate existing fallible, environment-reading
// computations (e.g. result.Eitherize1 applied to a func(R) (A, error)) into
// an effect pipeline.
//
// # Type Parameters
//
//   - R: The environment type consumed by the ReaderResult
//   - A: The type of the success value produced on a Right result
//
// # Parameters
//
//   - ma: A ReaderResult[R, A] (i.e. func(R) Result[A]) to lift
//
// # Returns
//
//   - Effect[R, A]: An effect that evaluates ma against its environment and
//     propagates the result as a success or failure
//
// # See Also
//
//   - FromReader: Lifts an always-successful Reader into an Effect
//   - ChainResultK: Chains a result-returning function over an existing Effect
//   - FromResult: Lifts a pre-computed Result into an Effect
//
//go:inline
func FromReaderResult[R, A any](ma ReaderResult[R, A]) Effect[R, A] {
	return readerreaderioresult.FromReaderResult(ma)
}
