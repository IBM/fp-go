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

package readerioresult

import (
	"context"
	"time"

	"github.com/IBM/fp-go/v2/context/readerio"
	"github.com/IBM/fp-go/v2/context/readerresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

// FromEither converts an [Either] into a [ReaderIOResult].
// The resulting computation ignores the context and immediately returns the Either value.
//
// Parameters:
//   - e: The Either value to lift into ReaderIOResult
//
// Returns a ReaderIOResult that produces the given Either value.
//
//go:inline
func FromEither[A any](e Either[A]) ReaderIOResult[A] {
	return RIOR.FromEither[context.Context](e)
}

// FromEither converts an [Either] into a [ReaderIOResult].
// The resulting computation ignores the context and immediately returns the Either value.
//
// Parameters:
//   - e: The Either value to lift into ReaderIOResult
//
// Returns a ReaderIOResult that produces the given Either value.
//
//go:inline
func FromResult[A any](e Result[A]) ReaderIOResult[A] {
	return RIOR.FromEither[context.Context](e)
}

// Left creates a [ReaderIOResult] that represents a failed computation with the given error.
//
// Parameters:
//   - l: The error value
//
// Returns a ReaderIOResult that always fails with the given error.
func Left[A any](l error) ReaderIOResult[A] {
	return RIOR.Left[context.Context, A](l)
}

// Right creates a [ReaderIOResult] that represents a successful computation with the given value.
//
// Parameters:
//   - r: The success value
//
// Returns a ReaderIOResult that always succeeds with the given value.
//
//go:inline
func Right[A any](r A) ReaderIOResult[A] {
	return RIOR.Right[context.Context](r)
}

// MonadMap transforms the success value of a [ReaderIOResult] using the provided function.
// If the computation fails, the error is propagated unchanged.
//
// Parameters:
//   - fa: The ReaderIOResult to transform
//   - f: The transformation function
//
// Returns a new ReaderIOResult with the transformed value.
//
//go:inline
func MonadMap[A, B any](fa ReaderIOResult[A], f func(A) B) ReaderIOResult[B] {
	return RIOR.MonadMap(fa, f)
}

// Map transforms the success value of a [ReaderIOResult] using the provided function.
// This is the curried version of [MonadMap], useful for composition.
//
// Parameters:
//   - f: The transformation function
//
// Returns a function that transforms a ReaderIOResult.
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return RIOR.Map[context.Context](f)
}

// MonadMapTo replaces the success value of a [ReaderIOResult] with a constant value.
// If the computation fails, the error is propagated unchanged.
//
// Parameters:
//   - fa: The ReaderIOResult to transform
//   - b: The constant value to use
//
// Returns a new ReaderIOResult with the constant value.
//
//go:inline
func MonadMapTo[A, B any](fa ReaderIOResult[A], b B) ReaderIOResult[B] {
	return RIOR.MonadMapTo(fa, b)
}

// MapTo replaces the success value of a [ReaderIOResult] with a constant value.
// This is the curried version of [MonadMapTo].
//
// Parameters:
//   - b: The constant value to use
//
// Returns a function that transforms a ReaderIOResult.
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return RIOR.MapTo[context.Context, A](b)
}

// MonadChain sequences two [ReaderIOResult] computations, where the second depends on the result of the first.
// If the first computation fails, the second is not executed.
//
// Parameters:
//   - ma: The first ReaderIOResult
//   - f: Function that produces the second ReaderIOResult based on the first's result
//
// Returns a new ReaderIOResult representing the sequenced computation.
//
//go:inline
func MonadChain[A, B any](ma ReaderIOResult[A], f Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChain(ma, WithContextK(f))
}

// Chain sequences two [ReaderIOResult] computations, where the second depends on the result of the first.
// This is the curried version of [MonadChain], useful for composition.
//
// Parameters:
//   - f: Function that produces the second ReaderIOResult based on the first's result
//
// Returns a function that sequences ReaderIOResult computations.
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return RIOR.Chain(WithContextK(f))
}

// MonadChainFirst sequences two [ReaderIOResult] computations but returns the result of the first.
// The second computation is executed for its side effects only.
//
// Parameters:
//   - ma: The first ReaderIOResult
//   - f: Function that produces the second ReaderIOResult
//
// Returns a ReaderIOResult with the result of the first computation.
//
//go:inline
func MonadChainFirst[A, B any](ma ReaderIOResult[A], f Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirst(ma, WithContextK(f))
}

//go:inline
func MonadTap[A, B any](ma ReaderIOResult[A], f Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTap(ma, WithContextK(f))
}

// ChainFirst sequences two [ReaderIOResult] computations but returns the result of the first.
// This is the curried version of [MonadChainFirst].
//
// Parameters:
//   - f: Function that produces the second ReaderIOResult
//
// Returns a function that sequences ReaderIOResult computations.
//
//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirst(WithContextK(f))
}

//go:inline
func Tap[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return RIOR.Tap(WithContextK(f))
}

// Of creates a [ReaderIOResult] that always succeeds with the given value.
// This is the same as [Right] and represents the monadic return operation.
//
// Parameters:
//   - a: The value to wrap
//
// Returns a ReaderIOResult that always succeeds with the given value.
//
//go:inline
func Of[A any](a A) ReaderIOResult[A] {
	return RIOR.Of[context.Context](a)
}

func withCancelCauseFunc[A any](cancel context.CancelCauseFunc, ma IOResult[A]) IOResult[A] {
	return function.Pipe3(
		ma,
		ioresult.Swap[A],
		ioeither.ChainFirstIOK[A](func(err error) func() Void {
			return io.FromImpure(func() { cancel(err) })
		}),
		ioeither.Swap[A],
	)
}

// MonadApPar implements parallel applicative application for [ReaderIOResult].
// It executes both computations in parallel and creates a sub-context that will be canceled
// if either operation fails. This provides automatic cancellation propagation.
//
// Parameters:
//   - fab: ReaderIOResult containing a function
//   - fa: ReaderIOResult containing a value
//
// Returns a ReaderIOResult with the function applied to the value.
func MonadApPar[B, A any](fab ReaderIOResult[func(A) B], fa ReaderIOResult[A]) ReaderIOResult[B] {
	// context sensitive input
	cfab := WithContext(fab)
	cfa := WithContext(fa)

	return func(ctx context.Context) IOResult[B] {
		// quick check for cancellation
		if ctx.Err() != nil {
			return ioeither.Left[B](context.Cause(ctx))
		}

		return func() Result[B] {
			// quick check for cancellation
			if ctx.Err() != nil {
				return either.Left[B](context.Cause(ctx))
			}

			// create sub-contexts for fa and fab, so they can cancel one other
			ctxSub, cancelSub := context.WithCancelCause(ctx)
			defer cancelSub(nil) // cancel has to be called in all paths

			fabIOE := withCancelCauseFunc(cancelSub, cfab(ctxSub))
			faIOE := withCancelCauseFunc(cancelSub, cfa(ctxSub))

			return ioresult.MonadApPar(fabIOE, faIOE)()
		}
	}
}

// MonadAp implements applicative application for [ReaderIOResult].
// By default, it uses parallel execution ([MonadApPar]) but can be configured to use
// sequential execution ([MonadApSeq]) via the useParallel constant.
//
// Parameters:
//   - fab: ReaderIOResult containing a function
//   - fa: ReaderIOResult containing a value
//
// Returns a ReaderIOResult with the function applied to the value.
func MonadAp[B, A any](fab ReaderIOResult[func(A) B], fa ReaderIOResult[A]) ReaderIOResult[B] {
	// dispatch to the configured version
	if useParallel {
		return MonadApPar(fab, fa)
	}
	return MonadApSeq(fab, fa)
}

// MonadApSeq implements sequential applicative application for [ReaderIOResult].
// It executes the function computation first, then the value computation.
//
// Parameters:
//   - fab: ReaderIOResult containing a function
//   - fa: ReaderIOResult containing a value
//
// Returns a ReaderIOResult with the function applied to the value.
//
//go:inline
func MonadApSeq[B, A any](fab ReaderIOResult[func(A) B], fa ReaderIOResult[A]) ReaderIOResult[B] {
	return RIOR.MonadApSeq(fab, fa)
}

// Ap applies a function wrapped in a [ReaderIOResult] to a value wrapped in a ReaderIOResult.
// This is the curried version of [MonadAp], using the default execution mode.
//
// Parameters:
//   - fa: ReaderIOResult containing a value
//
// Returns a function that applies a ReaderIOResult function to the value.
//
//go:inline
func Ap[B, A any](fa ReaderIOResult[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadAp[B, A], fa)
}

// ApSeq applies a function wrapped in a [ReaderIOResult] to a value sequentially.
// This is the curried version of [MonadApSeq].
//
// Parameters:
//   - fa: ReaderIOResult containing a value
//
// Returns a function that applies a ReaderIOResult function to the value sequentially.
//
//go:inline
func ApSeq[B, A any](fa ReaderIOResult[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApSeq[B, A], fa)
}

// ApPar applies a function wrapped in a [ReaderIOResult] to a value in parallel.
// This is the curried version of [MonadApPar].
//
// Parameters:
//   - fa: ReaderIOResult containing a value
//
// Returns a function that applies a ReaderIOResult function to the value in parallel.
//
//go:inline
func ApPar[B, A any](fa ReaderIOResult[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadApPar[B, A], fa)
}

// FromPredicate creates a [ReaderIOResult] from a predicate function.
// If the predicate returns true, the value is wrapped in Right; otherwise, Left with the error from onFalse.
//
// Parameters:
//   - pred: Predicate function to test the value
//   - onFalse: Function to generate an error when predicate fails
//
// Returns a function that converts a value to ReaderIOResult based on the predicate.
//
//go:inline
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return RIOR.FromPredicate[context.Context](pred, onFalse)
}

// OrElse provides an alternative [ReaderIOResult] computation if the first one fails.
// The alternative is only executed if the first computation results in a Left (error).
//
// Parameters:
//   - onLeft: Function that produces an alternative ReaderIOResult from the error
//
// Returns a function that provides fallback behavior for failed computations.
//
//go:inline
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A] {
	return RIOR.OrElse(onLeft)
}

// Ask returns a [ReaderIOResult] that provides access to the context.
// This is useful for accessing the [context.Context] within a computation.
//
// Returns a ReaderIOResult that produces the context.
//
//go:inline
func Ask() ReaderIOResult[context.Context] {
	return RIOR.Ask[context.Context]()
}

// MonadChainEitherK chains a function that returns an [Either] into a [ReaderIOResult] computation.
// This is useful for integrating pure Either-returning functions into ReaderIOResult workflows.
//
// Parameters:
//   - ma: The ReaderIOResult to chain from
//   - f: Function that produces an Either
//
// Returns a new ReaderIOResult with the chained computation.
//
//go:inline
func MonadChainEitherK[A, B any](ma ReaderIOResult[A], f either.Kleisli[error, A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainEitherK(ma, f)
}

// ChainEitherK chains a function that returns an [Either] into a [ReaderIOResult] computation.
// This is the curried version of [MonadChainEitherK].
//
// Parameters:
//   - f: Function that produces an Either
//
// Returns a function that chains the Either-returning function.
//
//go:inline
func ChainEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B] {
	return RIOR.ChainEitherK[context.Context](f)
}

//go:inline
func ChainResultK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B] {
	return RIOR.ChainEitherK[context.Context](f)
}

// MonadChainFirstEitherK chains a function that returns an [Either] but keeps the original value.
// The Either-returning function is executed for its validation/side effects only.
//
// Parameters:
//   - ma: The ReaderIOResult to chain from
//   - f: Function that produces an Either
//
// Returns a ReaderIOResult with the original value if both computations succeed.
//
//go:inline
func MonadChainFirstEitherK[A, B any](ma ReaderIOResult[A], f either.Kleisli[error, A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstEitherK(ma, f)
}

//go:inline
func MonadTapEitherK[A, B any](ma ReaderIOResult[A], f either.Kleisli[error, A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapEitherK(ma, f)
}

// ChainFirstEitherK chains a function that returns an [Either] but keeps the original value.
// This is the curried version of [MonadChainFirstEitherK].
//
// Parameters:
//   - f: Function that produces an Either
//
// Returns a function that chains the Either-returning function.
//
//go:inline
func ChainFirstEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return RIOR.ChainFirstEitherK[context.Context](f)
}

//go:inline
func TapEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, A] {
	return RIOR.TapEitherK[context.Context](f)
}

// ChainOptionK chains a function that returns an [Option] into a [ReaderIOResult] computation.
// If the Option is None, the provided error function is called.
//
// Parameters:
//   - onNone: Function to generate an error when Option is None
//
// Returns a function that chains Option-returning functions into ReaderIOResult.
//
//go:inline
func ChainOptionK[A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainOptionK[context.Context, A, B](onNone)
}

// FromIOEither converts an [IOResult] into a [ReaderIOResult].
// The resulting computation ignores the context.
//
// Parameters:
//   - t: The IOResult to convert
//
// Returns a ReaderIOResult that executes the IOResult.
//
//go:inline
func FromIOEither[A any](t IOResult[A]) ReaderIOResult[A] {
	return RIOR.FromIOEither[context.Context](t)
}

//go:inline
func FromIOResult[A any](t IOResult[A]) ReaderIOResult[A] {
	return RIOR.FromIOResult[context.Context](t)
}

// FromIO converts an [IO] into a [ReaderIOResult].
// The IO computation always succeeds, so it's wrapped in Right.
//
// Parameters:
//   - t: The IO to convert
//
// Returns a ReaderIOResult that executes the IO and wraps the result in Right.
//
//go:inline
func FromIO[A any](t IO[A]) ReaderIOResult[A] {
	return RIOR.FromIO[context.Context](t)
}

//go:inline
func FromReader[A any](t Reader[context.Context, A]) ReaderIOResult[A] {
	return RIOR.FromReader(t)
}

//go:inline
func FromReaderIO[A any](t ReaderIO[A]) ReaderIOResult[A] {
	return RIOR.FromReaderIO(t)
}

// FromLazy converts a [Lazy] computation into a [ReaderIOResult].
// The Lazy computation always succeeds, so it's wrapped in Right.
// This is an alias for [FromIO] since Lazy and IO have the same structure.
//
// Parameters:
//   - t: The Lazy computation to convert
//
// Returns a ReaderIOResult that executes the Lazy computation and wraps the result in Right.
//
//go:inline
func FromLazy[A any](t Lazy[A]) ReaderIOResult[A] {
	return RIOR.FromIO[context.Context](t)
}

// Never returns a [ReaderIOResult] that blocks indefinitely until the context is canceled.
// This is useful for creating computations that wait for external cancellation signals.
//
// Returns a ReaderIOResult that waits for context cancellation and returns the cancellation error.
func Never[A any]() ReaderIOResult[A] {
	return func(ctx context.Context) IOResult[A] {
		return func() Either[A] {
			<-ctx.Done()
			return either.Left[A](context.Cause(ctx))
		}
	}
}

// MonadChainIOK chains a function that returns an [IO] into a [ReaderIOResult] computation.
// The IO computation always succeeds, so it's wrapped in Right.
//
// Parameters:
//   - ma: The ReaderIOResult to chain from
//   - f: Function that produces an IO
//
// Returns a new ReaderIOResult with the chained IO computation.
//
//go:inline
func MonadChainIOK[A, B any](ma ReaderIOResult[A], f io.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainIOK(ma, f)
}

// ChainIOK chains a function that returns an [IO] into a [ReaderIOResult] computation.
// This is the curried version of [MonadChainIOK].
//
// Parameters:
//   - f: Function that produces an IO
//
// Returns a function that chains the IO-returning function.
//
//go:inline
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainIOK[context.Context](f)
}

// MonadChainFirstIOK chains a function that returns an [IO] but keeps the original value.
// The IO computation is executed for its side effects only.
//
// Parameters:
//   - ma: The ReaderIOResult to chain from
//   - f: Function that produces an IO
//
// Returns a ReaderIOResult with the original value after executing the IO.
//
//go:inline
func MonadChainFirstIOK[A, B any](ma ReaderIOResult[A], f io.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstIOK(ma, f)
}

//go:inline
func MonadTapIOK[A, B any](ma ReaderIOResult[A], f io.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapIOK(ma, f)
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
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstIOK[context.Context](f)
}

//go:inline
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapIOK[context.Context](f)
}

// ChainIOEitherK chains a function that returns an [IOResult] into a [ReaderIOResult] computation.
// This is useful for integrating IOResult-returning functions into ReaderIOResult workflows.
//
// Parameters:
//   - f: Function that produces an IOResult
//
// Returns a function that chains the IOResult-returning function.
//
//go:inline
func ChainIOEitherK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainIOEitherK[context.Context](f)
}

// Delay creates an operation that delays execution by the specified duration.
// The computation waits for either the delay to expire or the context to be canceled.
//
// Parameters:
//   - delay: The duration to wait before executing the computation
//
// Returns a function that delays a ReaderIOResult computation.
func Delay[A any](delay time.Duration) Operator[A, A] {
	return func(ma ReaderIOResult[A]) ReaderIOResult[A] {
		return func(ctx context.Context) IOResult[A] {
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
// Returns a ReaderIOResult that produces the current time after the delay.
func Timer(delay time.Duration) ReaderIOResult[time.Time] {
	return function.Pipe2(
		io.Now,
		FromIO[time.Time],
		Delay[time.Time](delay),
	)
}

// Defer creates a [ReaderIOResult] by lazily generating a new computation each time it's executed.
// This is useful for creating computations that should be re-evaluated on each execution.
//
// Parameters:
//   - gen: Lazy generator function that produces a ReaderIOResult
//
// Returns a ReaderIOResult that generates a fresh computation on each execution.
//
//go:inline
func Defer[A any](gen Lazy[ReaderIOResult[A]]) ReaderIOResult[A] {
	return RIOR.Defer(gen)
}

// TryCatch wraps a function that returns a tuple (value) into a [ReaderIOResult].
// This is the standard way to convert Go error-returning functions into ReaderIOResult.
//
// Parameters:
//   - f: Function that takes a context and returns a function producing (value)
//
// Returns a ReaderIOResult that wraps the error-returning function.
//
//go:inline
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOResult[A] {
	return RIOR.TryCatch(f, errors.Identity)
}

// MonadAlt provides an alternative [ReaderIOResult] if the first one fails.
// The alternative is lazily evaluated only if needed.
//
// Parameters:
//   - first: The primary ReaderIOResult to try
//   - second: Lazy alternative ReaderIOResult to use if first fails
//
// Returns a ReaderIOResult that tries the first, then the second if first fails.
//
//go:inline
func MonadAlt[A any](first ReaderIOResult[A], second Lazy[ReaderIOResult[A]]) ReaderIOResult[A] {
	return RIOR.MonadAlt(first, second)
}

// Alt provides an alternative [ReaderIOResult] if the first one fails.
// This is the curried version of [MonadAlt].
//
// Parameters:
//   - second: Lazy alternative ReaderIOResult to use if first fails
//
// Returns a function that provides fallback behavior.
//
//go:inline
func Alt[A any](second Lazy[ReaderIOResult[A]]) Operator[A, A] {
	return RIOR.Alt(second)
}

// Memoize computes the value of the provided [ReaderIOResult] monad lazily but exactly once.
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context.
//
// Parameters:
//   - rdr: The ReaderIOResult to memoize
//
// Returns a ReaderIOResult that caches its result after the first execution.
//
//go:inline
func Memoize[A any](rdr ReaderIOResult[A]) ReaderIOResult[A] {
	return RIOR.Memoize(rdr)
}

// Flatten converts a nested [ReaderIOResult] into a flat [ReaderIOResult].
// This is equivalent to [MonadChain] with the identity function.
//
// Parameters:
//   - rdr: The nested ReaderIOResult to flatten
//
// Returns a flattened ReaderIOResult.
//
//go:inline
func Flatten[A any](rdr ReaderIOResult[ReaderIOResult[A]]) ReaderIOResult[A] {
	return RIOR.Flatten(rdr)
}

// MonadFlap applies a value to a function wrapped in a [ReaderIOResult].
// This is the reverse of [MonadAp], useful in certain composition scenarios.
//
// Parameters:
//   - fab: ReaderIOResult containing a function
//   - a: The value to apply to the function
//
// Returns a ReaderIOResult with the function applied to the value.
//
//go:inline
func MonadFlap[B, A any](fab ReaderIOResult[func(A) B], a A) ReaderIOResult[B] {
	return RIOR.MonadFlap(fab, a)
}

// Flap applies a value to a function wrapped in a [ReaderIOResult].
// This is the curried version of [MonadFlap].
//
// Parameters:
//   - a: The value to apply to the function
//
// Returns a function that applies the value to a ReaderIOResult function.
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return RIOR.Flap[context.Context, B](a)
}

// Fold handles both success and error cases of a [ReaderIOResult] by providing handlers for each.
// Both handlers return ReaderIOResult, allowing for further composition.
//
// Parameters:
//   - onLeft: Handler for error case
//   - onRight: Handler for success case
//
// Returns a function that folds a ReaderIOResult into a new ReaderIOResult.
//
//go:inline
func Fold[A, B any](onLeft Kleisli[error, B], onRight Kleisli[A, B]) Operator[A, B] {
	return RIOR.Fold(function.Flow2(onLeft, WithContext), function.Flow2(onRight, WithContext))
}

// GetOrElse extracts the value from a [ReaderIOResult], providing a default via a function if it fails.
// The result is a [ReaderIO] that always succeeds.
//
// Parameters:
//   - onLeft: Function to provide a default value from the error
//
// Returns a function that converts a ReaderIOResult to a ReaderIO.
//
//go:inline
func GetOrElse[A any](onLeft readerio.Kleisli[error, A]) func(ReaderIOResult[A]) ReaderIO[A] {
	return RIOR.GetOrElse(onLeft)
}

// OrLeft transforms the error of a [ReaderIOResult] using the provided function.
// The success value is left unchanged.
//
// Parameters:
//   - onLeft: Function to transform the error
//
// Returns a function that transforms the error of a ReaderIOResult.
//
//go:inline
func OrLeft[A any](onLeft func(error) ReaderIO[error]) Operator[A, A] {
	return RIOR.OrLeft[A](onLeft)
}

//go:inline
func FromReaderEither[A any](ma ReaderEither[context.Context, error, A]) ReaderIOResult[A] {
	return RIOR.FromReaderEither(ma)
}

//go:inline
func FromReaderResult[A any](ma ReaderResult[A]) ReaderIOResult[A] {
	return RIOR.FromReaderEither(ma)
}

//go:inline
func FromReaderOption[A any](onNone Lazy[error]) Kleisli[ReaderOption[context.Context, A], A] {
	return RIOR.FromReaderOption[context.Context, A](onNone)
}

//go:inline
func MonadChainReaderK[A, B any](ma ReaderIOResult[A], f reader.Kleisli[context.Context, A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainReaderK(ma, f)
}

//go:inline
func ChainReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, B] {
	return RIOR.ChainReaderK(f)
}

//go:inline
func MonadChainFirstReaderK[A, B any](ma ReaderIOResult[A], f reader.Kleisli[context.Context, A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstReaderK(ma, f)
}

//go:inline
func MonadTapReaderK[A, B any](ma ReaderIOResult[A], f reader.Kleisli[context.Context, A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapReaderK(ma, f)
}

//go:inline
func ChainFirstReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, A] {
	return RIOR.ChainFirstReaderK(f)
}

//go:inline
func TapReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, A] {
	return RIOR.TapReaderK(f)
}

//go:inline
func MonadChainReaderResultK[A, B any](ma ReaderIOResult[A], f readerresult.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainReaderResultK(ma, f)
}

//go:inline
func ChainReaderResultK[A, B any](f readerresult.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainReaderResultK(f)
}

//go:inline
func MonadChainFirstReaderResultK[A, B any](ma ReaderIOResult[A], f readerresult.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstReaderResultK(ma, f)
}

//go:inline
func MonadTapReaderResultK[A, B any](ma ReaderIOResult[A], f readerresult.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapReaderResultK(ma, f)
}

//go:inline
func ChainFirstReaderResultK[A, B any](f readerresult.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstReaderResultK(f)
}

//go:inline
func TapReaderResultK[A, B any](f readerresult.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapReaderResultK(f)
}

//go:inline
func MonadChainReaderIOK[A, B any](ma ReaderIOResult[A], f readerio.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainReaderIOK(ma, f)
}

//go:inline
func ChainReaderIOK[A, B any](f readerio.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainReaderIOK(f)
}

//go:inline
func MonadChainFirstReaderIOK[A, B any](ma ReaderIOResult[A], f readerio.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstReaderIOK(ma, f)
}

//go:inline
func MonadTapReaderIOK[A, B any](ma ReaderIOResult[A], f readerio.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapReaderIOK(ma, f)
}

//go:inline
func ChainFirstReaderIOK[A, B any](f readerio.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstReaderIOK(f)
}

//go:inline
func TapReaderIOK[A, B any](f readerio.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapReaderIOK(f)
}

//go:inline
func ChainReaderOptionK[A, B any](onNone Lazy[error]) func(readeroption.Kleisli[context.Context, A, B]) Operator[A, B] {
	return RIOR.ChainReaderOptionK[context.Context, A, B](onNone)
}

//go:inline
func ChainFirstReaderOptionK[A, B any](onNone Lazy[error]) func(readeroption.Kleisli[context.Context, A, B]) Operator[A, A] {
	return RIOR.ChainFirstReaderOptionK[context.Context, A, B](onNone)
}

//go:inline
func TapReaderOptionK[A, B any](onNone Lazy[error]) func(readeroption.Kleisli[context.Context, A, B]) Operator[A, A] {
	return RIOR.TapReaderOptionK[context.Context, A, B](onNone)
}

//go:inline
func Read[A any](r context.Context) func(ReaderIOResult[A]) IOResult[A] {
	return RIOR.Read[A](r)
}

//go:inline
func ReadIO[A any](r IO[context.Context]) func(ReaderIOResult[A]) IOResult[A] {
	return RIOR.ReadIO[A](r)
}

//go:inline
func ReadIOEither[A any](r IOResult[context.Context]) func(ReaderIOResult[A]) IOResult[A] {
	return RIOR.ReadIOEither[A](r)
}

//go:inline
func ReadIOResult[A any](r IOResult[context.Context]) func(ReaderIOResult[A]) IOResult[A] {
	return RIOR.ReadIOResult[A](r)
}

// MonadChainLeft chains a computation on the left (error) side of a [ReaderIOResult].
// If the input is a Left value, it applies the function f to transform the error and potentially
// change the error type. If the input is a Right value, it passes through unchanged.
//
//go:inline
func MonadChainLeft[A any](fa ReaderIOResult[A], f Kleisli[error, A]) ReaderIOResult[A] {
	return RIOR.MonadChainLeft(fa, WithContextK(f))
}

// ChainLeft is the curried version of [MonadChainLeft].
// It returns a function that chains a computation on the left (error) side of a [ReaderIOResult].
//
//go:inline
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A] {
	return RIOR.ChainLeft(WithContextK(f))
}

// MonadChainFirstLeft chains a computation on the left (error) side but always returns the original error.
// If the input is a Left value, it applies the function f to the error and executes the resulting computation,
// but always returns the original Left error regardless of what f returns (Left or Right).
// If the input is a Right value, it passes through unchanged without calling f.
//
// This is useful for side effects on errors (like logging or metrics) where you want to perform an action
// when an error occurs but always propagate the original error, ensuring the error path is preserved.
//
//go:inline
func MonadChainFirstLeft[A, B any](ma ReaderIOResult[A], f Kleisli[error, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstLeft(ma, WithContextK(f))
}

//go:inline
func MonadTapLeft[A, B any](ma ReaderIOResult[A], f Kleisli[error, B]) ReaderIOResult[A] {
	return RIOR.MonadTapLeft(ma, WithContextK(f))
}

// ChainFirstLeft is the curried version of [MonadChainFirstLeft].
// It returns a function that chains a computation on the left (error) side while always preserving the original error.
//
// This is particularly useful for adding error handling side effects (like logging, metrics, or notifications)
// in a functional pipeline. The original error is always returned regardless of what f returns (Left or Right),
// ensuring the error path is preserved.
//
//go:inline
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return RIOR.ChainFirstLeft[A](WithContextK(f))
}

//go:inline
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return RIOR.TapLeft[A](WithContextK(f))
}

//go:inline
func ChainFirstLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A] {
	return RIOR.ChainFirstLeftIOK[A, context.Context](f)
}

//go:inline
func TapLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A] {
	return RIOR.TapLeftIOK[A, context.Context](f)
}

// Local transforms the context.Context environment before passing it to a ReaderIOResult computation.
//
// This is the Reader's local operation, which allows you to modify the environment
// for a specific computation without affecting the outer context. The transformation
// function receives the current context and returns a new context along with a
// cancel function. The cancel function is automatically called when the computation
// completes (via defer), ensuring proper cleanup of resources.
//
// The function checks for context cancellation before applying the transformation,
// returning an error immediately if the context is already cancelled.
//
// This is useful for:
//   - Adding timeouts or deadlines to specific operations
//   - Adding context values for nested computations
//   - Creating isolated context scopes
//   - Implementing context-based dependency injection
//
// Type Parameters:
//   - A: The value type of the ReaderIOResult
//
// Parameters:
//   - f: A function that transforms the context and returns a cancel function
//
// Returns:
//   - An Operator that runs the computation with the transformed context
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	// Add a custom value to the context
//	type key int
//	const userKey key = 0
//
//	addUser := readerioresult.Local[string](func(ctx context.Context) (context.Context, context.CancelFunc) {
//	    newCtx := context.WithValue(ctx, userKey, "Alice")
//	    return newCtx, func() {} // No-op cancel
//	})
//
//	getUser := readerioresult.FromReader(func(ctx context.Context) string {
//	    if user := ctx.Value(userKey); user != nil {
//	        return user.(string)
//	    }
//	    return "unknown"
//	})
//
//	result := F.Pipe1(
//	    getUser,
//	    addUser,
//	)
//	value, err := result(t.Context())()  // Returns ("Alice", nil)
//
// Timeout Example:
//
//	// Add a 5-second timeout to a specific operation
//	withTimeout := readerioresult.Local[Data](func(ctx context.Context) (context.Context, context.CancelFunc) {
//	    return context.WithTimeout(ctx, 5*time.Second)
//	})
//
//	result := F.Pipe1(
//	    fetchData,
//	    withTimeout,
//	)
func Local[A any](f func(context.Context) (context.Context, context.CancelFunc)) Operator[A, A] {
	return func(rr ReaderIOResult[A]) ReaderIOResult[A] {
		return func(ctx context.Context) IOResult[A] {
			return func() Result[A] {
				if ctx.Err() != nil {
					return result.Left[A](context.Cause(ctx))
				}
				otherCtx, otherCancel := f(ctx)
				defer otherCancel()
				return rr(otherCtx)()
			}
		}
	}
}

// WithTimeout adds a timeout to the context for a ReaderIOResult computation.
//
// This is a convenience wrapper around Local that uses context.WithTimeout.
// The computation must complete within the specified duration, or it will be
// cancelled. This is useful for ensuring operations don't run indefinitely
// and for implementing timeout-based error handling.
//
// The timeout is relative to when the ReaderIOResult is executed, not when
// WithTimeout is called. The cancel function is automatically called when
// the computation completes, ensuring proper cleanup. If the timeout expires,
// the computation will receive a context.DeadlineExceeded error.
//
// Type Parameters:
//   - A: The value type of the ReaderIOResult
//
// Parameters:
//   - timeout: The maximum duration for the computation
//
// Returns:
//   - An Operator that runs the computation with a timeout
//
// Example:
//
//	import (
//	    "time"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Fetch data with a 5-second timeout
//	fetchData := readerioresult.FromReader(func(ctx context.Context) Data {
//	    // Simulate slow operation
//	    select {
//	    case <-time.After(10 * time.Second):
//	        return Data{Value: "slow"}
//	    case <-ctx.Done():
//	        return Data{}
//	    }
//	})
//
//	result := F.Pipe1(
//	    fetchData,
//	    readerioresult.WithTimeout[Data](5*time.Second),
//	)
//	value, err := result(t.Context())()  // Returns (Data{}, context.DeadlineExceeded) after 5s
//
// Successful Example:
//
//	quickFetch := readerioresult.Right(Data{Value: "quick"})
//	result := F.Pipe1(
//	    quickFetch,
//	    readerioresult.WithTimeout[Data](5*time.Second),
//	)
//	value, err := result(t.Context())()  // Returns (Data{Value: "quick"}, nil)
func WithTimeout[A any](timeout time.Duration) Operator[A, A] {
	return Local[A](func(ctx context.Context) (context.Context, context.CancelFunc) {
		return context.WithTimeout(ctx, timeout)
	})
}

// WithDeadline adds an absolute deadline to the context for a ReaderIOResult computation.
//
// This is a convenience wrapper around Local that uses context.WithDeadline.
// The computation must complete before the specified time, or it will be
// cancelled. This is useful for coordinating operations that must finish
// by a specific time, such as request deadlines or scheduled tasks.
//
// The deadline is an absolute time, unlike WithTimeout which uses a relative
// duration. The cancel function is automatically called when the computation
// completes, ensuring proper cleanup. If the deadline passes, the computation
// will receive a context.DeadlineExceeded error.
//
// Type Parameters:
//   - A: The value type of the ReaderIOResult
//
// Parameters:
//   - deadline: The absolute time by which the computation must complete
//
// Returns:
//   - An Operator that runs the computation with a deadline
//
// Example:
//
//	import (
//	    "time"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Operation must complete by 3 PM
//	deadline := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)
//
//	fetchData := readerioresult.FromReader(func(ctx context.Context) Data {
//	    // Simulate operation
//	    select {
//	    case <-time.After(1 * time.Hour):
//	        return Data{Value: "done"}
//	    case <-ctx.Done():
//	        return Data{}
//	    }
//	})
//
//	result := F.Pipe1(
//	    fetchData,
//	    readerioresult.WithDeadline[Data](deadline),
//	)
//	value, err := result(t.Context())()  // Returns (Data{}, context.DeadlineExceeded) if past deadline
//
// Combining with Parent Context:
//
//	// If parent context already has a deadline, the earlier one takes precedence
//	parentCtx, cancel := context.WithDeadline(t.Context(), time.Now().Add(1*time.Hour))
//	defer cancel()
//
//	laterDeadline := time.Now().Add(2 * time.Hour)
//	result := F.Pipe1(
//	    fetchData,
//	    readerioresult.WithDeadline[Data](laterDeadline),
//	)
//	value, err := result(parentCtx)()  // Will use parent's 1-hour deadline
func WithDeadline[A any](deadline time.Time) Operator[A, A] {
	return Local[A](func(ctx context.Context) (context.Context, context.CancelFunc) {
		return context.WithDeadline(ctx, deadline)
	})
}
