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

package generic

import (
	"context"
	"time"

	E "github.com/IBM/fp-go/either"
	ER "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io/generic"
	IOE "github.com/IBM/fp-go/ioeither/generic"
	O "github.com/IBM/fp-go/option"
	RIE "github.com/IBM/fp-go/readerioeither/generic"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

func FromEither[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],
	A any](e E.Either[error, A]) GRA {
	return RIE.FromEither[GRA](e)
}

func RightReader[
	GRA ~func(context.Context) GIOA,
	GR ~func(context.Context) A,
	GIOA ~func() E.Either[error, A],
	A any](r GR) GRA {
	return RIE.RightReader[GR, GRA](r)
}

func LeftReader[
	GRA ~func(context.Context) GIOA,
	GR ~func(context.Context) error,
	GIOA ~func() E.Either[error, A],
	A any](l GR) GRA {
	return RIE.LeftReader[GR, GRA](l)
}

func Left[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],
	A any](l error) GRA {
	return RIE.Left[GRA](l)
}

func Right[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],
	A any](r A) GRA {
	return RIE.Right[GRA](r)
}

func FromReader[
	GRA ~func(context.Context) GIOA,
	GR ~func(context.Context) A,
	GIOA ~func() E.Either[error, A],
	A any](r GR) GRA {
	return RIE.FromReader[GR, GRA](r)
}

func MonadMap[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](fa GRA, f func(A) B) GRB {
	return RIE.MonadMap[GRA, GRB](fa, f)
}

func Map[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](f func(A) B) func(GRA) GRB {
	return RIE.Map[GRA, GRB](f)
}

func MonadMapTo[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](fa GRA, b B) GRB {
	return RIE.MonadMapTo[GRA, GRB](fa, b)
}

func MapTo[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](b B) func(GRA) GRB {
	return RIE.MapTo[GRA, GRB](b)
}

func MonadChain[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](ma GRA, f func(A) GRB) GRB {
	return RIE.MonadChain(ma, f)
}

func Chain[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](f func(A) GRB) func(GRA) GRB {
	return RIE.Chain[GRA](f)
}

func MonadChainFirst[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](ma GRA, f func(A) GRB) GRA {
	return RIE.MonadChainFirst(ma, f)
}

func ChainFirst[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	A, B any](f func(A) GRB) func(GRA) GRA {
	return RIE.ChainFirst[GRA](f)
}

func Of[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](a A) GRA {
	return RIE.Of[GRA](a)
}

// withCancelCauseFunc wraps an IOEither such that in case of an error the cancel function is invoked
func withCancelCauseFunc[
	GIOA ~func() E.Either[error, A],
	A any](cancel context.CancelCauseFunc, ma GIOA) GIOA {
	return F.Pipe3(
		ma,
		IOE.Swap[GIOA, func() E.Either[A, error]],
		IOE.ChainFirstIOK[func() E.Either[A, error], func() any](func(err error) func() any {
			return IO.MakeIO[func() any](func() any {
				cancel(err)
				return nil
			})
		}),
		IOE.Swap[func() E.Either[A, error], GIOA],
	)
}

// MonadApSeq implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
func MonadApSeq[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GRAB ~func(context.Context) GIOAB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],
	GIOAB ~func() E.Either[error, func(A) B],

	A, B any](fab GRAB, fa GRA) GRB {

	return RIE.MonadApSeq[GRA, GRB](fab, fa)
}

// MonadAp implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
func MonadApPar[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GRAB ~func(context.Context) GIOAB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],
	GIOAB ~func() E.Either[error, func(A) B],

	A, B any](fab GRAB, fa GRA) GRB {
	// context sensitive input
	cfab := WithContext(fab)
	cfa := WithContext(fa)

	return func(ctx context.Context) GIOB {
		// quick check for cancellation
		if err := context.Cause(ctx); err != nil {
			return IOE.Left[GIOB](err)
		}

		return func() E.Either[error, B] {
			// quick check for cancellation
			if err := context.Cause(ctx); err != nil {
				return E.Left[B](err)
			}

			// create sub-contexts for fa and fab, so they can cancel one other
			ctxSub, cancelSub := context.WithCancelCause(ctx)
			defer cancelSub(nil) // cancel has to be called in all paths

			fabIOE := withCancelCauseFunc(cancelSub, cfab(ctxSub))
			faIOE := withCancelCauseFunc(cancelSub, cfa(ctxSub))

			return IOE.MonadApPar[GIOB, GIOAB](fabIOE, faIOE)()
		}
	}
}

// MonadAp implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
func MonadAp[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GRAB ~func(context.Context) GIOAB,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],
	GIOAB ~func() E.Either[error, func(A) B],

	A, B any](fab GRAB, fa GRA) GRB {
	// dispatch to the configured version
	if useParallel {
		return MonadApPar[GRB](fab, fa)
	}
	return MonadApSeq[GRB](fab, fa)
}

func Ap[
	GRB ~func(context.Context) GIOB,
	GRAB ~func(context.Context) GIOAB,
	GRA ~func(context.Context) GIOA,

	GIOB ~func() E.Either[error, B],
	GIOAB ~func() E.Either[error, func(A) B],
	GIOA ~func() E.Either[error, A],

	A, B any](fa GRA) func(GRAB) GRB {
	return F.Bind2nd(MonadAp[GRB, GRA, GRAB], fa)
}

func ApSeq[
	GRB ~func(context.Context) GIOB,
	GRAB ~func(context.Context) GIOAB,
	GRA ~func(context.Context) GIOA,

	GIOB ~func() E.Either[error, B],
	GIOAB ~func() E.Either[error, func(A) B],
	GIOA ~func() E.Either[error, A],

	A, B any](fa GRA) func(GRAB) GRB {
	return F.Bind2nd(MonadApSeq[GRB, GRA, GRAB], fa)
}

func ApPar[
	GRB ~func(context.Context) GIOB,
	GRAB ~func(context.Context) GIOAB,
	GRA ~func(context.Context) GIOA,

	GIOB ~func() E.Either[error, B],
	GIOAB ~func() E.Either[error, func(A) B],
	GIOA ~func() E.Either[error, A],

	A, B any](fa GRA) func(GRAB) GRB {
	return F.Bind2nd(MonadApPar[GRB, GRA, GRAB], fa)
}

func FromPredicate[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](pred func(A) bool, onFalse func(A) error) func(A) GRA {
	return RIE.FromPredicate[GRA](pred, onFalse)
}

func Fold[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() B,

	A, B any](onLeft func(error) GRB, onRight func(A) GRB) func(GRA) GRB {
	return RIE.Fold[GRB, GRA](onLeft, onRight)
}

func GetOrElse[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() A,

	A any](onLeft func(error) GRB) func(GRA) GRB {
	return RIE.GetOrElse[GRB, GRA](onLeft)
}

func OrElse[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](onLeft func(error) GRA) func(GRA) GRA {
	return RIE.OrElse[GRA](onLeft)
}

func OrLeft[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() error,

	A any](onLeft func(error) GRB) func(GRA) GRA {
	return RIE.OrLeft[GRA, GRB, GRA](onLeft)
}

func Ask[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, context.Context],

]() GRA {
	return RIE.Ask[GRA]()
}

func Asks[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) A,
	GIOA ~func() E.Either[error, A],

	A any](r GRB) GRA {

	return RIE.Asks[GRB, GRA](r)
}

func MonadChainEitherK[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() E.Either[error, B],

	A, B any](ma GRA, f func(A) E.Either[error, B]) GRB {
	return RIE.MonadChainEitherK[GRA, GRB](ma, f)
}

func ChainEitherK[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() E.Either[error, B],

	A, B any](f func(A) E.Either[error, B]) func(ma GRA) GRB {
	return RIE.ChainEitherK[GRA, GRB](f)
}

func MonadChainFirstEitherK[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A, B any](ma GRA, f func(A) E.Either[error, B]) GRA {
	return RIE.MonadChainFirstEitherK[GRA](ma, f)
}

func ChainFirstEitherK[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A, B any](f func(A) E.Either[error, B]) func(ma GRA) GRA {
	return RIE.ChainFirstEitherK[GRA](f)
}

func ChainOptionK[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,

	GIOB ~func() E.Either[error, B],

	GIOA ~func() E.Either[error, A],

	A, B any](onNone func() error) func(func(A) O.Option[B]) func(GRA) GRB {
	return RIE.ChainOptionK[GRA, GRB](onNone)
}

func FromIOEither[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](t GIOA) GRA {
	return RIE.FromIOEither[GRA](t)
}

func FromIO[
	GRA ~func(context.Context) GIOA,
	GIOB ~func() A,

	GIOA ~func() E.Either[error, A],

	A any](t GIOB) GRA {
	return RIE.FromIO[GRA](t)
}

// Never returns a 'ReaderIOEither' that never returns, except if its context gets canceled
func Never[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any]() GRA {
	return func(ctx context.Context) GIOA {
		return IOE.MakeIO(func() E.Either[error, A] {
			<-ctx.Done()
			return E.Left[A](context.Cause(ctx))
		})
	}
}

func MonadChainIOK[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() E.Either[error, B],

	GIO ~func() B,

	A, B any](ma GRA, f func(A) GIO) GRB {
	return RIE.MonadChainIOK[GRA, GRB](ma, f)
}

func ChainIOK[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() E.Either[error, B],

	GIO ~func() B,

	A, B any](f func(A) GIO) func(ma GRA) GRB {
	return RIE.ChainIOK[GRA, GRB](f)
}

func MonadChainReaderIOK[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GRIO ~func(context.Context) GIO,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	GIO ~func() B,

	A, B any](ma GRA, f func(A) GRIO) GRB {
	return RIE.MonadChainReaderIOK[GRA, GRB](ma, f)
}

func ChainReaderIOK[
	GRB ~func(context.Context) GIOB,
	GRA ~func(context.Context) GIOA,
	GRIO ~func(context.Context) GIO,

	GIOA ~func() E.Either[error, A],
	GIOB ~func() E.Either[error, B],

	GIO ~func() B,

	A, B any](f func(A) GRIO) func(ma GRA) GRB {
	return RIE.ChainReaderIOK[GRA, GRB](f)
}

func MonadChainFirstIOK[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	GIO ~func() B,

	A, B any](ma GRA, f func(A) GIO) GRA {
	return RIE.MonadChainFirstIOK[GRA](ma, f)
}

func ChainFirstIOK[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	GIO ~func() B,

	A, B any](f func(A) GIO) func(ma GRA) GRA {
	return RIE.ChainFirstIOK[GRA](f)
}

func ChainIOEitherK[
	GRA ~func(context.Context) GIOA,
	GRB ~func(context.Context) GIOB,
	GIOA ~func() E.Either[error, A],

	GIOB ~func() E.Either[error, B],

	A, B any](f func(A) GIOB) func(ma GRA) GRB {
	return RIE.ChainIOEitherK[GRA, GRB](f)
}

// Delay creates an operation that passes in the value after some delay
func Delay[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](delay time.Duration) func(ma GRA) GRA {
	return func(ma GRA) GRA {
		return func(ctx context.Context) GIOA {
			return IOE.MakeIO(func() E.Either[error, A] {
				// manage the timeout
				timeoutCtx, cancelTimeout := context.WithTimeout(ctx, delay)
				defer cancelTimeout()
				// whatever comes first
				select {
				case <-timeoutCtx.Done():
					return ma(ctx)()
				case <-ctx.Done():
					return E.Left[A](context.Cause(ctx))
				}
			})
		}
	}
}

// Timer will return the current time after an initial delay
func Timer[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, time.Time],

](delay time.Duration) GRA {
	return F.Pipe2(
		IO.Now[func() time.Time](),
		FromIO[GRA, func() time.Time],
		Delay[GRA](delay),
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](gen func() GRA) GRA {
	return RIE.Defer[GRA](gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],

	A any](f func(context.Context) func() (A, error)) GRA {
	return RIE.TryCatch[GRA](f, ER.IdentityError)
}

func MonadAlt[LAZY ~func() GEA, GEA ~func(context.Context) GIOA, GIOA ~func() E.Either[error, A], A any](first GEA, second LAZY) GEA {
	return RIE.MonadAlt(first, second)
}

func Alt[LAZY ~func() GEA, GEA ~func(context.Context) GIOA, GIOA ~func() E.Either[error, A], A any](second LAZY) func(GEA) GEA {
	return RIE.Alt(second)
}

// Memoize computes the value of the provided monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],
	A any](rdr GRA) GRA {
	return RIE.Memoize[GRA](rdr)
}

func Flatten[
	GGRA ~func(context.Context) GGIOA,
	GGIOA ~func() E.Either[error, GRA],
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],
	A any](rdr GGRA) GRA {
	return RIE.Flatten[GRA](rdr)
}

func MonadFromReaderIO[
	GRIOEA ~func(context.Context) GIOEA,
	GIOEA ~func() E.Either[error, A],

	GRIOA ~func(context.Context) GIOA,
	GIOA ~func() A,

	A any](a A, f func(A) GRIOA) GRIOEA {
	return RIE.MonadFromReaderIO[GRIOEA](a, f)
}

func FromReaderIO[
	GRIOEA ~func(context.Context) GIOEA,
	GIOEA ~func() E.Either[error, A],

	GRIOA ~func(context.Context) GIOA,
	GIOA ~func() A,

	A any](f func(A) GRIOA) func(A) GRIOEA {
	return RIE.FromReaderIO[GRIOEA](f)
}

func RightReaderIO[
	GRIOEA ~func(context.Context) GIOEA,
	GIOEA ~func() E.Either[error, A],

	GRIOA ~func(context.Context) GIOA,
	GIOA ~func() A,

	A any](ma GRIOA) GRIOEA {
	return RIE.RightReaderIO[GRIOEA](ma)
}

func LeftReaderIO[
	GRIOEA ~func(context.Context) GIOEA,
	GIOEA ~func() E.Either[error, A],

	GRIOE ~func(context.Context) GIOE,
	GIOE ~func() error,

	A any](ma GRIOE) GRIOEA {
	return RIE.LeftReaderIO[GRIOEA](ma)
}

func MonadFlap[GREAB ~func(context.Context) GEAB, GREB ~func(context.Context) GEB, GEAB ~func() E.Either[error, func(A) B], GEB ~func() E.Either[error, B], B, A any](fab GREAB, a A) GREB {
	return RIE.MonadFlap[GREAB, GREB](fab, a)
}

func Flap[GREAB ~func(context.Context) GEAB, GREB ~func(context.Context) GEB, GEAB ~func() E.Either[error, func(A) B], GEB ~func() E.Either[error, B], B, A any](a A) func(GREAB) GREB {
	return RIE.Flap[GREAB, GREB](a)
}
