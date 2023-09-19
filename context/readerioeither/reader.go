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

	R "github.com/IBM/fp-go/context/reader"
	RIO "github.com/IBM/fp-go/context/readerio"
	G "github.com/IBM/fp-go/context/readerioeither/generic"
	ET "github.com/IBM/fp-go/either"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
	L "github.com/IBM/fp-go/lazy"
	O "github.com/IBM/fp-go/option"
)

func FromEither[A any](e ET.Either[error, A]) ReaderIOEither[A] {
	return G.FromEither[ReaderIOEither[A]](e)
}

func RightReader[A any](r R.Reader[A]) ReaderIOEither[A] {
	return G.RightReader[ReaderIOEither[A]](r)
}

func LeftReader[A any](l R.Reader[error]) ReaderIOEither[A] {
	return G.LeftReader[ReaderIOEither[A]](l)
}

func Left[A any](l error) ReaderIOEither[A] {
	return G.Left[ReaderIOEither[A]](l)
}

func Right[A any](r A) ReaderIOEither[A] {
	return G.Right[ReaderIOEither[A]](r)
}

func FromReader[A any](r R.Reader[A]) ReaderIOEither[A] {
	return G.FromReader[ReaderIOEither[A]](r)
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

// MonadAp implements the `Ap` function for a reader with context. It creates a sub-context that will
// be canceled if any of the input operations errors out or
func MonadAp[B, A any](fab ReaderIOEither[func(A) B], fa ReaderIOEither[A]) ReaderIOEither[B] {
	return G.MonadAp[ReaderIOEither[B]](fab, fa)
}

func Ap[B, A any](fa ReaderIOEither[A]) func(ReaderIOEither[func(A) B]) ReaderIOEither[B] {
	return G.Ap[ReaderIOEither[B], ReaderIOEither[func(A) B]](fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) func(A) ReaderIOEither[A] {
	return G.FromPredicate[ReaderIOEither[A]](pred, onFalse)
}

func Fold[A, B any](onLeft func(error) RIO.ReaderIO[B], onRight func(A) RIO.ReaderIO[B]) func(ReaderIOEither[A]) RIO.ReaderIO[B] {
	return G.Fold[RIO.ReaderIO[B], ReaderIOEither[A]](onLeft, onRight)
}

func GetOrElse[A any](onLeft func(error) RIO.ReaderIO[A]) func(ReaderIOEither[A]) RIO.ReaderIO[A] {
	return G.GetOrElse[RIO.ReaderIO[A], ReaderIOEither[A]](onLeft)
}

func OrElse[A any](onLeft func(error) ReaderIOEither[A]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return G.OrElse[ReaderIOEither[A]](onLeft)
}

func OrLeft[A any](onLeft func(error) RIO.ReaderIO[error]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return G.OrLeft[ReaderIOEither[A], RIO.ReaderIO[error]](onLeft)
}

func Ask() ReaderIOEither[context.Context] {
	return G.Ask[ReaderIOEither[context.Context]]()
}

func Asks[A any](r R.Reader[A]) ReaderIOEither[A] {
	return G.Asks[ReaderIOEither[A]](r)
}

func MonadChainEitherK[A, B any](ma ReaderIOEither[A], f func(A) ET.Either[error, B]) ReaderIOEither[B] {
	return G.MonadChainEitherK[ReaderIOEither[A], ReaderIOEither[B]](ma, f)
}

func ChainEitherK[A, B any](f func(A) ET.Either[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return G.ChainEitherK[ReaderIOEither[A], ReaderIOEither[B]](f)
}

func MonadChainFirstEitherK[A, B any](ma ReaderIOEither[A], f func(A) ET.Either[error, B]) ReaderIOEither[A] {
	return G.MonadChainFirstEitherK[ReaderIOEither[A]](ma, f)
}

func ChainFirstEitherK[A, B any](f func(A) ET.Either[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return G.ChainFirstEitherK[ReaderIOEither[A]](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) O.Option[B]) func(ReaderIOEither[A]) ReaderIOEither[B] {
	return G.ChainOptionK[ReaderIOEither[A], ReaderIOEither[B]](onNone)
}

func FromIOEither[A any](t IOE.IOEither[error, A]) ReaderIOEither[A] {
	return G.FromIOEither[ReaderIOEither[A]](t)
}

func FromIO[A any](t IO.IO[A]) ReaderIOEither[A] {
	return G.FromIO[ReaderIOEither[A]](t)
}

func FromLazy[A any](t L.Lazy[A]) ReaderIOEither[A] {
	return G.FromIO[ReaderIOEither[A]](t)
}

// Never returns a 'ReaderIOEither' that never returns, except if its context gets canceled
func Never[A any]() ReaderIOEither[A] {
	return G.Never[ReaderIOEither[A]]()
}

func MonadChainIOK[A, B any](ma ReaderIOEither[A], f func(A) IO.IO[B]) ReaderIOEither[B] {
	return G.MonadChainIOK[ReaderIOEither[B], ReaderIOEither[A]](ma, f)
}

func ChainIOK[A, B any](f func(A) IO.IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return G.ChainIOK[ReaderIOEither[B], ReaderIOEither[A]](f)
}

func MonadChainFirstIOK[A, B any](ma ReaderIOEither[A], f func(A) IO.IO[B]) ReaderIOEither[A] {
	return G.MonadChainFirstIOK[ReaderIOEither[A]](ma, f)
}

func ChainFirstIOK[A, B any](f func(A) IO.IO[B]) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return G.ChainFirstIOK[ReaderIOEither[A]](f)
}

func ChainIOEitherK[A, B any](f func(A) IOE.IOEither[error, B]) func(ma ReaderIOEither[A]) ReaderIOEither[B] {
	return G.ChainIOEitherK[ReaderIOEither[A], ReaderIOEither[B]](f)
}

// Delay creates an operation that passes in the value after some delay
func Delay[A any](delay time.Duration) func(ma ReaderIOEither[A]) ReaderIOEither[A] {
	return G.Delay[ReaderIOEither[A]](delay)
}

// Timer will return the current time after an initial delay
func Timer(delay time.Duration) ReaderIOEither[time.Time] {
	return G.Timer[ReaderIOEither[time.Time]](delay)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen L.Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return G.Defer[ReaderIOEither[A]](gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[A any](f func(context.Context) func() (A, error)) ReaderIOEither[A] {
	return G.TryCatch[ReaderIOEither[A]](f)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[A any](first ReaderIOEither[A], second L.Lazy[ReaderIOEither[A]]) ReaderIOEither[A] {
	return G.MonadAlt(first, second)
}

// Alt identifies an associative operation on a type constructor
func Alt[A any](second L.Lazy[ReaderIOEither[A]]) func(ReaderIOEither[A]) ReaderIOEither[A] {
	return G.Alt(second)
}
