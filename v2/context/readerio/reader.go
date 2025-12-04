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
	"context"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

// MonadMap transforms the success value of a [ReaderIO] using the provided function.
// If the computation fails, the error is propagated unchanged.
//
// Parameters:
//   - fa: The ReaderIO to transform
//   - f: The transformation function
//
// Returns a new ReaderIO with the transformed value.
//
//go:inline
func MonadMap[A, B any](fa ReaderIO[A], f func(A) B) ReaderIO[B] {
	return RIO.MonadMap(fa, f)
}

// Map transforms the success value of a [ReaderIO] using the provided function.
// This is the curried version of [MonadMap], useful for composition.
//
// Parameters:
//   - f: The transformation function
//
// Returns a function that transforms a ReaderIO.
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return RIO.Map[context.Context](f)
}

// MonadMapTo replaces the success value of a [ReaderIO] with a constant value.
// If the computation fails, the error is propagated unchanged.
//
// Parameters:
//   - fa: The ReaderIO to transform
//   - b: The constant value to use
//
// Returns a new ReaderIO with the constant value.
//
//go:inline
func MonadMapTo[A, B any](fa ReaderIO[A], b B) ReaderIO[B] {
	return RIO.MonadMapTo(fa, b)
}

// MapTo replaces the success value of a [ReaderIO] with a constant value.
// This is the curried version of [MonadMapTo].
//
// Parameters:
//   - b: The constant value to use
//
// Returns a function that transforms a ReaderIO.
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return RIO.MapTo[context.Context, A](b)
}

// MonadChain sequences two [ReaderIO] computations, where the second depends on the result of the first.
// If the first computation fails, the second is not executed.
//
// Parameters:
//   - ma: The first ReaderIO
//   - f: Function that produces the second ReaderIO based on the first's result
//
// Returns a new ReaderIO representing the sequenced computation.
//
//go:inline
func MonadChain[A, B any](ma ReaderIO[A], f Kleisli[A, B]) ReaderIO[B] {
	return RIO.MonadChain(ma, f)
}

// Chain sequences two [ReaderIO] computations, where the second depends on the result of the first.
// This is the curried version of [MonadChain], useful for composition.
//
// Parameters:
//   - f: Function that produces the second ReaderIO based on the first's result
//
// Returns a function that sequences ReaderIO computations.
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return RIO.Chain(f)
}

// MonadChainFirst sequences two [ReaderIO] computations but returns the result of the first.
// The second computation is executed for its side effects only.
//
// Parameters:
//   - ma: The first ReaderIO
//   - f: Function that produces the second ReaderIO
//
// Returns a ReaderIO with the result of the first computation.
//
//go:inline
func MonadChainFirst[A, B any](ma ReaderIO[A], f Kleisli[A, B]) ReaderIO[A] {
	return RIO.MonadChainFirst(ma, f)
}

// MonadTap executes a side-effect computation but returns the original value.
// This is an alias for [MonadChainFirst] and is useful for operations like logging
// or validation that should not affect the main computation flow.
//
// Parameters:
//   - ma: The ReaderIO to tap
//   - f: Function that produces a side-effect ReaderIO
//
// Returns a ReaderIO with the original value after executing the side effect.
//
//go:inline
func MonadTap[A, B any](ma ReaderIO[A], f Kleisli[A, B]) ReaderIO[A] {
	return RIO.MonadTap(ma, f)
}

// ChainFirst sequences two [ReaderIO] computations but returns the result of the first.
// This is the curried version of [MonadChainFirst].
//
// Parameters:
//   - f: Function that produces the second ReaderIO
//
// Returns a function that sequences ReaderIO computations.
//
//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return RIO.ChainFirst(f)
}

// Tap executes a side-effect computation but returns the original value.
// This is the curried version of [MonadTap], an alias for [ChainFirst].
//
// Parameters:
//   - f: Function that produces a side-effect ReaderIO
//
// Returns a function that taps ReaderIO computations.
//
//go:inline
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return RIO.Tap(f)
}

// Of creates a [ReaderIO] that always succeeds with the given value.
// This is the same as [Right] and represents the monadic return operation.
//
// Parameters:
//   - a: The value to wrap
//
// Returns a ReaderIO that always succeeds with the given value.
//
//go:inline
func Of[A any](a A) ReaderIO[A] {
	return RIO.Of[context.Context](a)
}

// MonadApPar implements parallel applicative application for [ReaderIO].
// It executes the function and value computations in parallel where possible,
// potentially improving performance for independent operations.
//
// Parameters:
//   - fab: ReaderIO containing a function
//   - fa: ReaderIO containing a value
//
// Returns a ReaderIO with the function applied to the value.
//
//go:inline
func MonadApPar[B, A any](fab ReaderIO[func(A) B], fa ReaderIO[A]) ReaderIO[B] {
	return RIO.MonadApPar(fab, fa)
}

// MonadAp implements applicative application for [ReaderIO].
// By default, it uses parallel execution ([MonadApPar]) but can be configured to use
// sequential execution ([MonadApSeq]) via the useParallel constant.
//
// Parameters:
//   - fab: ReaderIO containing a function
//   - fa: ReaderIO containing a value
//
// Returns a ReaderIO with the function applied to the value.
//
//go:inline
func MonadAp[B, A any](fab ReaderIO[func(A) B], fa ReaderIO[A]) ReaderIO[B] {
	// dispatch to the configured version
	if useParallel {
		return MonadApPar(fab, fa)
	}
	return MonadApSeq(fab, fa)
}

// MonadApSeq implements sequential applicative application for [ReaderIO].
// It executes the function computation first, then the value computation.
//
// Parameters:
//   - fab: ReaderIO containing a function
//   - fa: ReaderIO containing a value
//
// Returns a ReaderIO with the function applied to the value.
//
//go:inline
func MonadApSeq[B, A any](fab ReaderIO[func(A) B], fa ReaderIO[A]) ReaderIO[B] {
	return RIO.MonadApSeq(fab, fa)
}

// Ap applies a function wrapped in a [ReaderIO] to a value wrapped in a ReaderIO.
// This is the curried version of [MonadAp], using the default execution mode.
//
// Parameters:
//   - fa: ReaderIO containing a value
//
// Returns a function that applies a ReaderIO function to the value.
//
//go:inline
func Ap[B, A any](fa ReaderIO[A]) Operator[func(A) B, B] {
	return RIO.Ap[B](fa)
}

// ApSeq applies a function wrapped in a [ReaderIO] to a value sequentially.
// This is the curried version of [MonadApSeq].
//
// Parameters:
//   - fa: ReaderIO containing a value
//
// Returns a function that applies a ReaderIO function to the value sequentially.
//
//go:inline
func ApSeq[B, A any](fa ReaderIO[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApSeq[B, A], fa)
}

// ApPar applies a function wrapped in a [ReaderIO] to a value in parallel.
// This is the curried version of [MonadApPar].
//
// Parameters:
//   - fa: ReaderIO containing a value
//
// Returns a function that applies a ReaderIO function to the value in parallel.
//
//go:inline
func ApPar[B, A any](fa ReaderIO[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApPar[B, A], fa)
}

// Ask returns a [ReaderIO] that provides access to the context.
// This is useful for accessing the [context.Context] within a computation.
//
// Returns a ReaderIO that produces the context.
//
//go:inline
func Ask() ReaderIO[context.Context] {
	return RIO.Ask[context.Context]()
}

// FromIO converts an [IO] into a [ReaderIO].
// The IO computation always succeeds, so it's wrapped in Right.
//
// Parameters:
//   - t: The IO to convert
//
// Returns a ReaderIO that executes the IO and wraps the result in Right.
//
//go:inline
func FromIO[A any](t IO[A]) ReaderIO[A] {
	return RIO.FromIO[context.Context](t)
}

// FromReader converts a [Reader] into a [ReaderIO].
// The Reader computation is lifted into the IO context, allowing it to be
// composed with other ReaderIO operations.
//
// Parameters:
//   - t: The Reader to convert
//
// Returns a ReaderIO that executes the Reader and wraps the result in IO.
//
//go:inline
func FromReader[A any](t Reader[context.Context, A]) ReaderIO[A] {
	return RIO.FromReader(t)
}

// FromLazy converts a [Lazy] computation into a [ReaderIO].
// The Lazy computation always succeeds, so it's wrapped in Right.
// This is an alias for [FromIO] since Lazy and IO have the same structure.
//
// Parameters:
//   - t: The Lazy computation to convert
//
// Returns a ReaderIO that executes the Lazy computation and wraps the result in Right.
//
//go:inline
func FromLazy[A any](t Lazy[A]) ReaderIO[A] {
	return RIO.FromIO[context.Context](t)
}

// MonadChainIOK chains a function that returns an [IO] into a [ReaderIO] computation.
// The IO computation always succeeds, so it's wrapped in Right.
//
// Parameters:
//   - ma: The ReaderIO to chain from
//   - f: Function that produces an IO
//
// Returns a new ReaderIO with the chained IO computation.
//
//go:inline
func MonadChainIOK[A, B any](ma ReaderIO[A], f func(A) IO[B]) ReaderIO[B] {
	return RIO.MonadChainIOK(ma, f)
}

// ChainIOK chains a function that returns an [IO] into a [ReaderIO] computation.
// This is the curried version of [MonadChainIOK].
//
// Parameters:
//   - f: Function that produces an IO
//
// Returns a function that chains the IO-returning function.
//
//go:inline
func ChainIOK[A, B any](f func(A) IO[B]) Operator[A, B] {
	return RIO.ChainIOK[context.Context](f)
}

// MonadChainFirstIOK chains a function that returns an [IO] but keeps the original value.
// The IO computation is executed for its side effects only.
//
// Parameters:
//   - ma: The ReaderIO to chain from
//   - f: Function that produces an IO
//
// Returns a ReaderIO with the original value after executing the IO.
//
//go:inline
func MonadChainFirstIOK[A, B any](ma ReaderIO[A], f func(A) IO[B]) ReaderIO[A] {
	return RIO.MonadChainFirstIOK(ma, f)
}

// MonadTapIOK chains a function that returns an [IO] but keeps the original value.
// This is an alias for [MonadChainFirstIOK] and is useful for side effects like logging.
//
// Parameters:
//   - ma: The ReaderIO to tap
//   - f: Function that produces an IO for side effects
//
// Returns a ReaderIO with the original value after executing the IO.
//
//go:inline
func MonadTapIOK[A, B any](ma ReaderIO[A], f func(A) IO[B]) ReaderIO[A] {
	return RIO.MonadTapIOK(ma, f)
}

// ChainFirstIOK chains a function that returns an [IO] but keeps the original value.
// This is the curried version of [MonadChainFirstIOK].
//
// Parameters:
//   - f: Function that produces an IO
//
// Returns a function that chains the IO-returning function.
//
//go:inline
func ChainFirstIOK[A, B any](f func(A) IO[B]) Operator[A, A] {
	return RIO.ChainFirstIOK[context.Context](f)
}

// TapIOK chains a function that returns an [IO] but keeps the original value.
// This is the curried version of [MonadTapIOK], an alias for [ChainFirstIOK].
//
// Parameters:
//   - f: Function that produces an IO for side effects
//
// Returns a function that taps with IO-returning functions.
//
//go:inline
func TapIOK[A, B any](f func(A) IO[B]) Operator[A, A] {
	return RIO.TapIOK[context.Context](f)
}

// Defer creates a [ReaderIO] by lazily generating a new computation each time it's executed.
// This is useful for creating computations that should be re-evaluated on each execution.
//
// Parameters:
//   - gen: Lazy generator function that produces a ReaderIO
//
// Returns a ReaderIO that generates a fresh computation on each execution.
//
//go:inline
func Defer[A any](gen Lazy[ReaderIO[A]]) ReaderIO[A] {
	return RIO.Defer(gen)
}

// Memoize computes the value of the provided [ReaderIO] monad lazily but exactly once.
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context.
//
// Parameters:
//   - rdr: The ReaderIO to memoize
//
// Returns a ReaderIO that caches its result after the first execution.
//
//go:inline
func Memoize[A any](rdr ReaderIO[A]) ReaderIO[A] {
	return RIO.Memoize(rdr)
}

// Flatten converts a nested [ReaderIO] into a flat [ReaderIO].
// This is equivalent to [MonadChain] with the identity function.
//
// Parameters:
//   - rdr: The nested ReaderIO to flatten
//
// Returns a flattened ReaderIO.
//
//go:inline
func Flatten[A any](rdr ReaderIO[ReaderIO[A]]) ReaderIO[A] {
	return RIO.Flatten(rdr)
}

// MonadFlap applies a value to a function wrapped in a [ReaderIO].
// This is the reverse of [MonadAp], useful in certain composition scenarios.
//
// Parameters:
//   - fab: ReaderIO containing a function
//   - a: The value to apply to the function
//
// Returns a ReaderIO with the function applied to the value.
//
//go:inline
func MonadFlap[B, A any](fab ReaderIO[func(A) B], a A) ReaderIO[B] {
	return RIO.MonadFlap(fab, a)
}

// Flap applies a value to a function wrapped in a [ReaderIO].
// This is the curried version of [MonadFlap].
//
// Parameters:
//   - a: The value to apply to the function
//
// Returns a function that applies the value to a ReaderIO function.
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return RIO.Flap[context.Context, B](a)
}

// MonadChainReaderK chains a [ReaderIO] with a function that returns a [Reader].
// The Reader is lifted into the ReaderIO context, allowing composition of
// Reader and ReaderIO operations.
//
// Parameters:
//   - ma: The ReaderIO to chain from
//   - f: Function that produces a Reader
//
// Returns a new ReaderIO with the chained Reader computation.
//
//go:inline
func MonadChainReaderK[A, B any](ma ReaderIO[A], f reader.Kleisli[context.Context, A, B]) ReaderIO[B] {
	return RIO.MonadChainReaderK(ma, f)
}

// ChainReaderK chains a [ReaderIO] with a function that returns a [Reader].
// This is the curried version of [MonadChainReaderK].
//
// Parameters:
//   - f: Function that produces a Reader
//
// Returns a function that chains Reader-returning functions.
//
//go:inline
func ChainReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, B] {
	return RIO.ChainReaderK(f)
}

// MonadChainFirstReaderK chains a function that returns a [Reader] but keeps the original value.
// The Reader computation is executed for its side effects only.
//
// Parameters:
//   - ma: The ReaderIO to chain from
//   - f: Function that produces a Reader
//
// Returns a ReaderIO with the original value after executing the Reader.
//
//go:inline
func MonadChainFirstReaderK[A, B any](ma ReaderIO[A], f reader.Kleisli[context.Context, A, B]) ReaderIO[A] {
	return RIO.MonadChainFirstReaderK(ma, f)
}

// MonadTapReaderK chains a function that returns a [Reader] but keeps the original value.
// This is an alias for [MonadChainFirstReaderK] and is useful for side effects.
//
// Parameters:
//   - ma: The ReaderIO to tap
//   - f: Function that produces a Reader for side effects
//
// Returns a ReaderIO with the original value after executing the Reader.
//
//go:inline
func MonadTapReaderK[A, B any](ma ReaderIO[A], f reader.Kleisli[context.Context, A, B]) ReaderIO[A] {
	return RIO.MonadTapReaderK(ma, f)
}

// ChainFirstReaderK chains a function that returns a [Reader] but keeps the original value.
// This is the curried version of [MonadChainFirstReaderK].
//
// Parameters:
//   - f: Function that produces a Reader
//
// Returns a function that chains Reader-returning functions while preserving the original value.
//
//go:inline
func ChainFirstReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, A] {
	return RIO.ChainFirstReaderK(f)
}

// TapReaderK chains a function that returns a [Reader] but keeps the original value.
// This is the curried version of [MonadTapReaderK], an alias for [ChainFirstReaderK].
//
// Parameters:
//   - f: Function that produces a Reader for side effects
//
// Returns a function that taps with Reader-returning functions.
//
//go:inline
func TapReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, A] {
	return RIO.TapReaderK(f)
}

// Read executes a [ReaderIO] with a given context, returning the resulting [IO].
// This is useful for providing the context dependency and obtaining an IO action
// that can be executed later.
//
// Parameters:
//   - r: The context to provide to the ReaderIO
//
// Returns a function that converts a ReaderIO into an IO by applying the context.
//
//go:inline
func Read[A any](r context.Context) func(ReaderIO[A]) IO[A] {
	return RIO.Read[A](r)
}
