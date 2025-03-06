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

	G "github.com/IBM/fp-go/v2/context/readerioeither/generic"
	"github.com/IBM/fp-go/v2/either"
	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RE "github.com/IBM/fp-go/v2/readerioeither"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

func FromEither[A any](e ET.Either[error, A]) ReaderIOEither[A] {
	return G.FromEither[ReaderIOEither[A]](e)
}

func Left[A any](l error) ReaderIOEither[A] {
	return G.Left[ReaderIOEither[A]](l)
}

func Right[A any](r A) ReaderIOEither[A] {
	return G.Right[ReaderIOEither[A]](r)
}

func MonadMap[A, B any](fa ReaderIOEither[A], f func(A) B) ReaderIOEither[B] {
	return G.MonadMap[ReaderIOEither[A], ReaderIOEither[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return G.Map[ReaderIOEither[A], ReaderIOEither[B]](f)
}

func MonadMapTo[A, B any](fa ReaderIOEither[A], b B) ReaderIOEither[B] {
	return G.MonadMapTo[ReaderIOEither[A], ReaderIOEither[B]](fa, b)
}

func MapTo[A, B any](b B) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return G.MapTo[ReaderIOEither[A], ReaderIOEither[B]](b)
}

func MonadChain[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[B] {
	return G.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderIOEither[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return G.Chain[ReaderIOEither[A]](f)
}

func MonadChainFirst[A, B any](ma ReaderIOEither[A], f func(A) ReaderIOEither[B]) ReaderIOEither[A] {
	return G.MonadChainFirst(ma, f)
}

func ChainFirst[A, B any](f func(A) ReaderIOEither[B]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return G.ChainFirst[ReaderIOEither[A]](f)
}

func Of[A any](a A) ReaderIOEither[A] {
	return G.Of[ReaderIOEither[A]](a)
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
	return RE.MonadApSeq(fab, fa)
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
	return RE.FromPredicate[context.Context](pred, onFalse)
}

func OrElse[A any](onLeft func(error) ReaderIOEither[A]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return RE.OrElse[context.Context](onLeft)
}

func Ask() ReaderIOEither[context.Context] {
	return RE.Ask[context.Context, error]()
}

func MonadChainEitherK[A, B any](ma ReaderIOEither[A], f func(A) ET.Either[error, B]) ReaderIOEither[B] {
	return RE.MonadChainEitherK[context.Context](ma, f)
}

func ChainEitherK[A, B any](f func(A) ET.Either[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return RE.ChainEitherK[context.Context](f)
}

func MonadChainFirstEitherK[A, B any](ma ReaderIOEither[A], f func(A) ET.Either[error, B]) ReaderIOEither[A] {
	return RE.MonadChainFirstEitherK[context.Context](ma, f)
}

func ChainFirstEitherK[A, B any](f func(A) ET.Either[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return RE.ChainFirstEitherK[context.Context](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) O.Option[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return RE.ChainOptionK[context.Context, A, B](onNone)
}

func FromIOEither[A any](t ioeither.IOEither[error, A]) ReaderIOEither[A] {
	return RE.FromIOEither[context.Context](t)
}

func FromIO[A any](t IO[A]) ReaderIOEither[A] {
	return RE.FromIO[context.Context, error](t)
}

func FromLazy[A any](t Lazy[A]) ReaderIOEither[A] {
	return RE.FromIO[context.Context, error](t)
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
	return RE.MonadChainIOK(ma, f)
}

func ChainIOK[A, B any](f func(A) IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return RE.ChainIOK[context.Context, error](f)
}

func MonadChainFirstIOK[A, B any](ma ReaderIOEither[A], f func(A) IO[B]) ReaderIOEither[A] {
	return RE.MonadChainFirstIOK(ma, f)
}

func ChainFirstIOK[A, B any](f func(A) IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return RE.ChainFirstIOK[context.Context, error](f)
}

func ChainIOEitherK[A, B any](f func(A) ioeither.IOEither[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return RE.ChainIOEitherK[context.Context](f)
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
	return RE.Defer(gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOEither[A] {
	return RE.TryCatch(f, errors.IdentityError)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[A any](first ReaderIOEither[A], second Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return RE.MonadAlt(first, second)
}

// Alt identifies an associative operation on a type constructor
func Alt[A any](second Lazy[ReaderIOEither[A]]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return RE.Alt(second)
}

// Memoize computes the value of the provided [ReaderIOEither] monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[A any](rdr ReaderIOEither[A]) ReaderIOEither[A] {
	return RE.Memoize(rdr)
}

// Flatten converts a nested [ReaderIOEither] into a [ReaderIOEither]
func Flatten[A any](rdr ReaderIOEither[ReaderIOEither[A]]) ReaderIOEither[A] {
	return RE.Flatten(rdr)
}

func MonadFlap[B, A any](fab ReaderIOEither[func(A) B], a A) ReaderIOEither[B] {
	return RE.MonadFlap(fab, a)
}

func Flap[B, A any](a A) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return RE.Flap[context.Context, error, B](a)
}

func Fold[A, B any](onLeft func(error) ReaderIOEither[B], onRight func(A) ReaderIOEither[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return RE.Fold(onLeft, onRight)
}

func GetOrElse[A any](onLeft func(error) RIO.ReaderIO[context.Context, A]) func(ReaderIOEither[A]) RIO.ReaderIO[context.Context, A] {
	return RE.GetOrElse(onLeft)
}

func OrLeft[A any](onLeft func(error) RIO.ReaderIO[context.Context, error]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return RE.OrLeft[A](onLeft)
}
