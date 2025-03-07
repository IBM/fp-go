// Copyright (c) 2023 IBM Corp.
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

func FromEither[A any](e Either[A]) ReaderIOEither[A] {
	return readerioeither.FromEither[context.Context](e)
}

func Left[A any](l error) ReaderIOEither[A] {
	return readerioeither.Left[context.Context, A](l)
}

func Right[A any](r A) ReaderIOEither[A] {
	return readerioeither.Right[context.Context, error](r)
}

func MonadMap[A, B any](fa ReaderIOEither[A], f func(A) B) ReaderIOEither[B] {
	return readerioeither.MonadMap(fa, f)
}

func Map[A, B any](f func(A) B) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.Map[context.Context, error](f)
}

func MonadMapTo[A, B any](fa ReaderIOEither[A], b B) ReaderIOEither[B] {
	return readerioeither.MonadMapTo(fa, b)
}

func MapTo[A, B any](b B) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.MapTo[context.Context, error, A](b)
}

func MonadChain[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[B] {
	return readerioeither.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderIOEither[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.Chain(f)
}

func MonadChainFirst[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[A] {
	return readerioeither.MonadChainFirst(ma, f)
}

func ChainFirst[A, B any](f func(A) ReaderIOEither[B]) Operator[A, A] {
	return readerioeither.ChainFirst(f)
}

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

// MonadApPar implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
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

// MonadAp implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
func MonadAp[B, A any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	// dispatch to the configured version
	if useParallel {
		return MonadApPar(fab, fa)
	}
	return MonadApSeq(fab, fa)
}

func MonadApSeq[B, A any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.MonadApSeq(fab, fa)
}

func Ap[B, A any](fa ReaderIOEither[A]) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return function.Bind2nd(MonadAp[B, A], fa)
}

func ApSeq[B, A any](fa ReaderIOEither[A]) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return function.Bind2nd(MonadApSeq[B, A], fa)
}

func ApPar[B, A any](fa ReaderIOEither[A]) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return function.Bind2nd(MonadApPar[B, A], fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) func(A) ReaderIOEither[A] {
	return readerioeither.FromPredicate[context.Context](pred, onFalse)
}

func OrElse[A any](onLeft func(error) ReaderIOEither[A]) Operator[A, A] {
	return readerioeither.OrElse[context.Context](onLeft)
}

func Ask() ReaderIOEither[context.Context] {
	return readerioeither.Ask[context.Context, error]()
}

func MonadChainEitherK[A, B any](ma ReaderIOEither[A], f func(A) Either[B]) ReaderIOEither[B] {
	return readerioeither.MonadChainEitherK[context.Context](ma, f)
}

func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainEitherK[context.Context](f)
}

func MonadChainFirstEitherK[A, B any](ma ReaderIOEither[A], f func(A) Either[B]) ReaderIOEither[A] {
	return readerioeither.MonadChainFirstEitherK[context.Context](ma, f)
}

func ChainFirstEitherK[A, B any](f func(A) Either[B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return readerioeither.ChainFirstEitherK[context.Context](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) Option[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainOptionK[context.Context, A, B](onNone)
}

func FromIOEither[A any](t ioeither.IOEither[error, A]) ReaderIOEither[A] {
	return readerioeither.FromIOEither[context.Context](t)
}

func FromIO[A any](t IO[A]) ReaderIOEither[A] {
	return readerioeither.FromIO[context.Context, error](t)
}

func FromLazy[A any](t Lazy[A]) ReaderIOEither[A] {
	return readerioeither.FromIO[context.Context, error](t)
}

// Never returns a 'ReaderIOEither' that never returns, except if its context gets canceled
func Never[A any]() ReaderIOEither[A] {
	return func(ctx context.Context) IOEither[A] {
		return func() Either[A] {
			<-ctx.Done()
			return either.Left[A](context.Cause(ctx))
		}
	}
}

func MonadChainIOK[A, B any](ma ReaderIOEither[A], f func(A) IO[B]) ReaderIOEither[B] {
	return readerioeither.MonadChainIOK(ma, f)
}

func ChainIOK[A, B any](f func(A) IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainIOK[context.Context, error](f)
}

func MonadChainFirstIOK[A, B any](ma ReaderIOEither[A], f func(A) IO[B]) ReaderIOEither[A] {
	return readerioeither.MonadChainFirstIOK(ma, f)
}

func ChainFirstIOK[A, B any](f func(A) IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return readerioeither.ChainFirstIOK[context.Context, error](f)
}

func ChainIOEitherK[A, B any](f func(A) ioeither.IOEither[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.ChainIOEitherK[context.Context](f)
}

// Delay creates an operation that passes in the value after some delay
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

// Timer will return the current time after an initial delay
func Timer(delay time.Duration) ReaderIOEither[time.Time] {
	return function.Pipe2(
		io.Now,
		FromIO[time.Time],
		Delay[time.Time](delay),
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return readerioeither.Defer(gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOEither[A] {
	return readerioeither.TryCatch(f, errors.IdentityError)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[A any](first ReaderIOEither[A], second Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return readerioeither.MonadAlt(first, second)
}

// Alt identifies an associative operation on a type constructor
func Alt[A any](second Lazy[ReaderIOEither[A]]) Operator[A, A] {
	return readerioeither.Alt(second)
}

// Memoize computes the value of the provided [ReaderIOEither] monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[A any](rdr ReaderIOEither[A]) ReaderIOEither[A] {
	return readerioeither.Memoize(rdr)
}

// Flatten converts a nested [ReaderIOEither] into a [ReaderIOEither]
func Flatten[A any](rdr ReaderIOEither[ReaderIOEither[A]]) ReaderIOEither[A] {
	return readerioeither.Flatten(rdr)
}

func MonadFlap[B, A any](fab ReaderIOEither[func(A) B], a A) ReaderIOEither[B] {
	return readerioeither.MonadFlap(fab, a)
}

func Flap[B, A any](a A) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return readerioeither.Flap[context.Context, error, B](a)
}

func Fold[A, B any](onLeft func(error) ReaderIOEither[B], onRight func(A) ReaderIOEither[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return readerioeither.Fold(onLeft, onRight)
}

func GetOrElse[A any](onLeft func(error) ReaderIO[A]) func(ReaderIOEither[A]) ReaderIO[A] {
	return readerioeither.GetOrElse(onLeft)
}

func OrLeft[A any](onLeft func(error) ReaderIO[error]) Operator[A, A] {
	return readerioeither.OrLeft[A](onLeft)
}
