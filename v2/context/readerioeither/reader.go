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

package readerioeither

import (
	"context"
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/readerioeither"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

// FromEither converts an [Either] into a [ReaderIOEither].
// The resulting computation ignores the context and immediately returns the Either value.
//
// Parameters:
//   - e: The Either value to lift into ReaderIOEither
//
// Returns a ReaderIOEither that produces the given Either value.
func FromEither[A any](e Either[A]) ReaderIOEither[A] {
	return readerioeither.FromEither[context.Context](e)
}

// Left creates a [ReaderIOEither] that represents a failed computation with the given error.
//
// Parameters:
//   - l: The error value
//
// Returns a ReaderIOEither that always fails with the given error.
func Left[A any](l error) ReaderIOEither[A] {
	return readerioeither.Left[context.Context, A](l)
}

// Right creates a [ReaderIOEither] that represents a successful computation with the given value.
//
// Parameters:
//   - r: The success value
//
// Returns a ReaderIOEither that always succeeds with the given value.
func Right[A any](r A) ReaderIOEither[A] {
	return readerioeither.Right[context.Context, error](r)
}

// MonadMap transforms the success value of a [ReaderIOEither] using the provided function.
// If the computation fails, the error is propagated unchanged.
//
// Parameters:
//   - fa: The ReaderIOEither to transform
//   - f: The transformation function
//
// Returns a new ReaderIOEither with the transformed value.
func MonadMap[A, B any](fa ReaderIOEither[A], f func(A) B) ReaderIOEither[B] {
	return readerioeither.MonadMap(fa, f)
}

// Map transforms the success value of a [ReaderIOEither] using the provided function.
// This is the curried version of [MonadMap], useful for composition.
//
// Parameters:
//   - f: The transformation function
//
// Returns a function that transforms a ReaderIOEither.
func Map[A, B any](f func(A) B) Operator[A, B] {
	return readerioeither.Map[context.Context, error](f)
}

// MonadMapTo replaces the success value of a [ReaderIOEither] with a constant value.
// If the computation fails, the error is propagated unchanged.
//
// Parameters:
//   - fa: The ReaderIOEither to transform
//   - b: The constant value to use
//
// Returns a new ReaderIOEither with the constant value.
func MonadMapTo[A, B any](fa ReaderIOEither[A], b B) ReaderIOEither[B] {
	return readerioeither.MonadMapTo(fa, b)
}

// MapTo replaces the success value of a [ReaderIOEither] with a constant value.
// This is the curried version of [MonadMapTo].
//
// Parameters:
//   - b: The constant value to use
//
// Returns a function that transforms a ReaderIOEither.
func MapTo[A, B any](b B) Operator[A, B] {
	return readerioeither.MapTo[context.Context, error, A](b)
}

// MonadChain sequences two [ReaderIOEither] computations, where the second depends on the result of the first.
// If the first computation fails, the second is not executed.
//
// Parameters:
//   - ma: The first ReaderIOEither
//   - f: Function that produces the second ReaderIOEither based on the first's result
//
// Returns a new ReaderIOEither representing the sequenced computation.
func MonadChain[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[B] {
	return readerioeither.MonadChain(ma, f)
}

// Chain sequences two [ReaderIOEither] computations, where the second depends on the result of the first.
// This is the curried version of [MonadChain], useful for composition.
//
// Parameters:
//   - f: Function that produces the second ReaderIOEither based on the first's result
//
// Returns a function that sequences ReaderIOEither computations.
func Chain[A, B any](f func(A) ReaderIOEither[B]) Operator[A, B] {
	return readerioeither.Chain(f)
}

// MonadChainFirst sequences two [ReaderIOEither] computations but returns the result of the first.
// The second computation is executed for its side effects only.
//
// Parameters:
//   - ma: The first ReaderIOEither
//   - f: Function that produces the second ReaderIOEither
//
// Returns a ReaderIOEither with the result of the first computation.
func MonadChainFirst[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[A] {
	return readerioeither.MonadChainFirst(ma, f)
}

// ChainFirst sequences two [ReaderIOEither] computations but returns the result of the first.
// This is the curried version of [MonadChainFirst].
//
// Parameters:
//   - f: Function that produces the second ReaderIOEither
//
// Returns a function that sequences ReaderIOEither computations.
func ChainFirst[A, B any](f func(A) ReaderIOEither[B]) Operator[A, A] {
	return readerioeither.ChainFirst(f)
}

// Of creates a [ReaderIOEither] that always succeeds with the given value.
// This is the same as [Right] and represents the monadic return operation.
//
// Parameters:
//   - a: The value to wrap
//
// Returns a ReaderIOEither that always succeeds with the given value.
func Of[A any](a A) ReaderIOEither[A] {
	return readerioeither.Of[context.Context, error](a)
}

func withCancelCauseFunc[A any](cancel context.CancelCauseFunc, ma IOEither[A]) IOEither[A] {
	return function.Pipe3(
		ma,
		ioeither.Swap[error, A],
		ioeither.ChainFirstIOK[A](func(err error) func() any {
			return io.FromImpure(func() { cancel(err) })
		}),
		ioeither.Swap[A, error],
	)
}

// MonadApPar implements parallel applicative application for [ReaderIOEither].
// It executes both computations in parallel and creates a sub-context that will be canceled
// if either operation fails. This provides automatic cancellation propagation.
//
// Parameters:
//   - fab: ReaderIOEither containing a function
//   - fa: ReaderIOEither containing a value
//
// Returns a ReaderIOEither with the function applied to the value.
func MonadApPar[B, A any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	// context sensitive input
	cfab := WithContext(fab)
	cfa := WithContext(fa)

	return func(ctx context.Context) IOEither[B] {
		// quick check for cancellation
		if err := context.Cause(ctx); err != nil {
			return ioeither.Left[B](err)
		}

		return func() Either[B] {
			// quick check for cancellation
			if err := context.Cause(ctx); err != nil {
				return either.Left[B](err)
			}

			// create sub-contexts for fa and fab, so they can cancel one other
			ctxSub, cancelSub := context.WithCancelCause(ctx)
			defer cancelSub(nil) // cancel has to be called in all paths

			fabIOE := withCancelCauseFunc(cancelSub, cfab(ctxSub))
			faIOE := withCancelCauseFunc(cancelSub, cfa(ctxSub))

			return ioeither.MonadApPar(fabIOE, faIOE)()
		}
	}
}

// MonadAp implements applicative application for [ReaderIOEither].
// By default, it uses parallel execution ([MonadApPar]) but can be configured to use
// sequential execution ([MonadApSeq]) via the useParallel constant.
//
// Parameters:
//   - fab: ReaderIOEither containing a function
//   - fa: ReaderIOEither containing a value
//
// Returns a ReaderIOEither with the function applied to the value.
func MonadAp[B, A any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	// dispatch to the configured version
	if useParallel {
		return MonadApPar(fab, fa)
	}
	return MonadApSeq(fab, fa)
}

// MonadApSeq implements sequential applicative application for [ReaderIOEither].
// It executes the function computation first, then the value computation.
//
// Parameters:
//   - fab: ReaderIOEither containing a function
//   - fa: ReaderIOEither containing a value
//
// Returns a ReaderIOEither with the function applied to the value.
func MonadApSeq[B, A any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.MonadApSeq(fab, fa)
}

// Ap applies a function wrapped in a [ReaderIOEither] to a value wrapped in a ReaderIOEither.
// This is the curried version of [MonadAp], using the default execution mode.
//
// Parameters:
//   - fa: ReaderIOEither containing a value
//
// Returns a function that applies a ReaderIOEither function to the value.
func Ap[B, A any](fa ReaderIOEither[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadAp[B, A], fa)
}

// ApSeq applies a function wrapped in a [ReaderIOEither] to a value sequentially.
// This is the curried version of [MonadApSeq].
//
// Parameters:
//   - fa: ReaderIOEither containing a value
//
// Returns a function that applies a ReaderIOEither function to the value sequentially.
func ApSeq[B, A any](fa ReaderIOEither[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApSeq[B, A], fa)
}

// ApPar applies a function wrapped in a [ReaderIOEither] to a value in parallel.
// This is the curried version of [MonadApPar].
//
// Parameters:
//   - fa: ReaderIOEither containing a value
//
// Returns a function that applies a ReaderIOEither function to the value in parallel.
func ApPar[B, A any](fa ReaderIOEither[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApPar[B, A], fa)
}

// FromPredicate creates a [ReaderIOEither] from a predicate function.
// If the predicate returns true, the value is wrapped in Right; otherwise, Left with the error from onFalse.
//
// Parameters:
//   - pred: Predicate function to test the value
//   - onFalse: Function to generate an error when predicate fails
//
// Returns a function that converts a value to ReaderIOEither based on the predicate.
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) func(A) ReaderIOEither[A] {
	return readerioeither.FromPredicate[context.Context](pred, onFalse)
}

// OrElse provides an alternative [ReaderIOEither] computation if the first one fails.
// The alternative is only executed if the first computation results in a Left (error).
//
// Parameters:
//   - onLeft: Function that produces an alternative ReaderIOEither from the error
//
// Returns a function that provides fallback behavior for failed computations.
func OrElse[A any](onLeft func(error) ReaderIOEither[A]) Operator[A, A] {
	return readerioeither.OrElse[context.Context](onLeft)
}

// Ask returns a [ReaderIOEither] that provides access to the context.
// This is useful for accessing the [context.Context] within a computation.
//
// Returns a ReaderIOEither that produces the context.
func Ask() ReaderIOEither[context.Context] {
	return readerioeither.Ask[context.Context, error]()
}

// MonadChainEitherK chains a function that returns an [Either] into a [ReaderIOEither] computation.
// This is useful for integrating pure Either-returning functions into ReaderIOEither workflows.
//
// Parameters:
//   - ma: The ReaderIOEither to chain from
//   - f: Function that produces an Either
//
// Returns a new ReaderIOEither with the chained computation.
func MonadChainEitherK[A, B any](ma ReaderIOEither[A], f func(A) Either[B]) ReaderIOEither[B] {
	return readerioeither.MonadChainEitherK[context.Context](ma, f)
}

// ChainEitherK chains a function that returns an [Either] into a [ReaderIOEither] computation.
// This is the curried version of [MonadChainEitherK].
//
// Parameters:
//   - f: Function that produces an Either
//
// Returns a function that chains the Either-returning function.
func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainEitherK[context.Context](f)
}

// MonadChainFirstEitherK chains a function that returns an [Either] but keeps the original value.
// The Either-returning function is executed for its validation/side effects only.
//
// Parameters:
//   - ma: The ReaderIOEither to chain from
//   - f: Function that produces an Either
//
// Returns a ReaderIOEither with the original value if both computations succeed.
func MonadChainFirstEitherK[A, B any](ma ReaderIOEither[A], f func(A) Either[B]) ReaderIOEither[A] {
	return readerioeither.MonadChainFirstEitherK[context.Context](ma, f)
}

// ChainFirstEitherK chains a function that returns an [Either] but keeps the original value.
// This is the curried version of [MonadChainFirstEitherK].
//
// Parameters:
//   - f: Function that produces an Either
//
// Returns a function that chains the Either-returning function.
func ChainFirstEitherK[A, B any](f func(A) Either[B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return readerioeither.ChainFirstEitherK[context.Context](f)
}

// ChainOptionK chains a function that returns an [Option] into a [ReaderIOEither] computation.
// If the Option is None, the provided error function is called.
//
// Parameters:
//   - onNone: Function to generate an error when Option is None
//
// Returns a function that chains Option-returning functions into ReaderIOEither.
func ChainOptionK[A, B any](onNone func() error) func(func(A) Option[B]) Operator[A, B] {
	return readerioeither.ChainOptionK[context.Context, A, B](onNone)
}

// FromIOEither converts an [IOEither] into a [ReaderIOEither].
// The resulting computation ignores the context.
//
// Parameters:
//   - t: The IOEither to convert
//
// Returns a ReaderIOEither that executes the IOEither.
func FromIOEither[A any](t ioeither.IOEither[error, A]) ReaderIOEither[A] {
	return readerioeither.FromIOEither[context.Context](t)
}

// FromIO converts an [IO] into a [ReaderIOEither].
// The IO computation always succeeds, so it's wrapped in Right.
//
// Parameters:
//   - t: The IO to convert
//
// Returns a ReaderIOEither that executes the IO and wraps the result in Right.
func FromIO[A any](t IO[A]) ReaderIOEither[A] {
	return readerioeither.FromIO[context.Context, error](t)
}

// FromLazy converts a [Lazy] computation into a [ReaderIOEither].
// The Lazy computation always succeeds, so it's wrapped in Right.
// This is an alias for [FromIO] since Lazy and IO have the same structure.
//
// Parameters:
//   - t: The Lazy computation to convert
//
// Returns a ReaderIOEither that executes the Lazy computation and wraps the result in Right.
func FromLazy[A any](t Lazy[A]) ReaderIOEither[A] {
	return readerioeither.FromIO[context.Context, error](t)
}

// Never returns a [ReaderIOEither] that blocks indefinitely until the context is canceled.
// This is useful for creating computations that wait for external cancellation signals.
//
// Returns a ReaderIOEither that waits for context cancellation and returns the cancellation error.
func Never[A any]() ReaderIOEither[A] {
	return func(ctx context.Context) IOEither[A] {
		return func() Either[A] {
			<-ctx.Done()
			return either.Left[A](context.Cause(ctx))
		}
	}
}

// MonadChainIOK chains a function that returns an [IO] into a [ReaderIOEither] computation.
// The IO computation always succeeds, so it's wrapped in Right.
//
// Parameters:
//   - ma: The ReaderIOEither to chain from
//   - f: Function that produces an IO
//
// Returns a new ReaderIOEither with the chained IO computation.
func MonadChainIOK[A, B any](ma ReaderIOEither[A], f func(A) IO[B]) ReaderIOEither[B] {
	return readerioeither.MonadChainIOK(ma, f)
}

// ChainIOK chains a function that returns an [IO] into a [ReaderIOEither] computation.
// This is the curried version of [MonadChainIOK].
//
// Parameters:
//   - f: Function that produces an IO
//
// Returns a function that chains the IO-returning function.
func ChainIOK[A, B any](f func(A) IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainIOK[context.Context, error](f)
}

// MonadChainFirstIOK chains a function that returns an [IO] but keeps the original value.
// The IO computation is executed for its side effects only.
//
// Parameters:
//   - ma: The ReaderIOEither to chain from
//   - f: Function that produces an IO
//
// Returns a ReaderIOEither with the original value after executing the IO.
func MonadChainFirstIOK[A, B any](ma ReaderIOEither[A], f func(A) IO[B]) ReaderIOEither[A] {
	return readerioeither.MonadChainFirstIOK(ma, f)
}

// ChainFirstIOK chains a function that returns an [IO] but keeps the original value.
// This is the curried version of [MonadChainFirstIOK].
//
// Parameters:
//   - f: Function that produces an IO
//
// Returns a function that chains the IO-returning function.
func ChainFirstIOK[A, B any](f func(A) IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return readerioeither.ChainFirstIOK[context.Context, error](f)
}

// ChainIOEitherK chains a function that returns an [IOEither] into a [ReaderIOEither] computation.
// This is useful for integrating IOEither-returning functions into ReaderIOEither workflows.
//
// Parameters:
//   - f: Function that produces an IOEither
//
// Returns a function that chains the IOEither-returning function.
func ChainIOEitherK[A, B any](f func(A) ioeither.IOEither[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainIOEitherK[context.Context](f)
}

// Delay creates an operation that delays execution by the specified duration.
// The computation waits for either the delay to expire or the context to be canceled.
//
// Parameters:
//   - delay: The duration to wait before executing the computation
//
// Returns a function that delays a ReaderIOEither computation.
func Delay[A any](delay time.Duration) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return func(ma ReaderIOEither[A]) ReaderIOEither[A] {
		return func(ctx context.Context) IOEither[A] {
			return func() Either[A] {
				// manage the timeout
				timeoutCtx, cancelTimeout := context.WithTimeout(ctx, delay)
				defer cancelTimeout()
				// whatever comes first
				select {
				case <-timeoutCtx.Done():
					return ma(ctx)()
				case <-ctx.Done():
					return either.Left[A](context.Cause(ctx))
				}
			}
		}
	}
}

// Timer returns the current time after waiting for the specified delay.
// This is useful for creating time-based computations.
//
// Parameters:
//   - delay: The duration to wait before returning the time
//
// Returns a ReaderIOEither that produces the current time after the delay.
func Timer(delay time.Duration) ReaderIOEither[time.Time] {
	return function.Pipe2(
		io.Now,
		FromIO[time.Time],
		Delay[time.Time](delay),
	)
}

// Defer creates a [ReaderIOEither] by lazily generating a new computation each time it's executed.
// This is useful for creating computations that should be re-evaluated on each execution.
//
// Parameters:
//   - gen: Lazy generator function that produces a ReaderIOEither
//
// Returns a ReaderIOEither that generates a fresh computation on each execution.
func Defer[A any](gen Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return readerioeither.Defer(gen)
}

// TryCatch wraps a function that returns a tuple (value, error) into a [ReaderIOEither].
// This is the standard way to convert Go error-returning functions into ReaderIOEither.
//
// Parameters:
//   - f: Function that takes a context and returns a function producing (value, error)
//
// Returns a ReaderIOEither that wraps the error-returning function.
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOEither[A] {
	return readerioeither.TryCatch(f, errors.IdentityError)
}

// MonadAlt provides an alternative [ReaderIOEither] if the first one fails.
// The alternative is lazily evaluated only if needed.
//
// Parameters:
//   - first: The primary ReaderIOEither to try
//   - second: Lazy alternative ReaderIOEither to use if first fails
//
// Returns a ReaderIOEither that tries the first, then the second if first fails.
func MonadAlt[A any](first ReaderIOEither[A], second Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return readerioeither.MonadAlt(first, second)
}

// Alt provides an alternative [ReaderIOEither] if the first one fails.
// This is the curried version of [MonadAlt].
//
// Parameters:
//   - second: Lazy alternative ReaderIOEither to use if first fails
//
// Returns a function that provides fallback behavior.
func Alt[A any](second Lazy[ReaderIOEither[A]]) Operator[A, A] {
	return readerioeither.Alt(second)
}

// Memoize computes the value of the provided [ReaderIOEither] monad lazily but exactly once.
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context.
//
// Parameters:
//   - rdr: The ReaderIOEither to memoize
//
// Returns a ReaderIOEither that caches its result after the first execution.
func Memoize[A any](rdr ReaderIOEither[A]) ReaderIOEither[A] {
	return readerioeither.Memoize(rdr)
}

// Flatten converts a nested [ReaderIOEither] into a flat [ReaderIOEither].
// This is equivalent to [MonadChain] with the identity function.
//
// Parameters:
//   - rdr: The nested ReaderIOEither to flatten
//
// Returns a flattened ReaderIOEither.
func Flatten[A any](rdr ReaderIOEither[ReaderIOEither[A]]) ReaderIOEither[A] {
	return readerioeither.Flatten(rdr)
}

// MonadFlap applies a value to a function wrapped in a [ReaderIOEither].
// This is the reverse of [MonadAp], useful in certain composition scenarios.
//
// Parameters:
//   - fab: ReaderIOEither containing a function
//   - a: The value to apply to the function
//
// Returns a ReaderIOEither with the function applied to the value.
func MonadFlap[B, A any](fab ReaderIOEither[func(A) B], a A) ReaderIOEither[B] {
	return readerioeither.MonadFlap(fab, a)
}

// Flap applies a value to a function wrapped in a [ReaderIOEither].
// This is the curried version of [MonadFlap].
//
// Parameters:
//   - a: The value to apply to the function
//
// Returns a function that applies the value to a ReaderIOEither function.
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return readerioeither.Flap[context.Context, error, B](a)
}

// Fold handles both success and error cases of a [ReaderIOEither] by providing handlers for each.
// Both handlers return ReaderIOEither, allowing for further composition.
//
// Parameters:
//   - onLeft: Handler for error case
//   - onRight: Handler for success case
//
// Returns a function that folds a ReaderIOEither into a new ReaderIOEither.
func Fold[A, B any](onLeft func(error) ReaderIOEither[B], onRight func(A) ReaderIOEither[B]) Operator[A, B] {
	return readerioeither.Fold(onLeft, onRight)
}

// GetOrElse extracts the value from a [ReaderIOEither], providing a default via a function if it fails.
// The result is a [ReaderIO] that always succeeds.
//
// Parameters:
//   - onLeft: Function to provide a default value from the error
//
// Returns a function that converts a ReaderIOEither to a ReaderIO.
func GetOrElse[A any](onLeft func(error) ReaderIO[A]) func(ReaderIOEither[A]) ReaderIO[A] {
	return readerioeither.GetOrElse(onLeft)
}

// OrLeft transforms the error of a [ReaderIOEither] using the provided function.
// The success value is left unchanged.
//
// Parameters:
//   - onLeft: Function to transform the error
//
// Returns a function that transforms the error of a ReaderIOEither.
func OrLeft[A any](onLeft func(error) ReaderIO[error]) Operator[A, A] {
	return readerioeither.OrLeft[A](onLeft)
}
