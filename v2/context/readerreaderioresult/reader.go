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

package readerreaderioresult

import (
	"context"
	"time"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readeroption"
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
)

//go:inline
func FromReaderOption[R, A any](onNone Lazy[error]) Kleisli[R, ReaderOption[R, A], A] {
	return RRIOE.FromReaderOption[R, context.Context, A](onNone)
}

//go:inline
func FromReaderIOResult[R, A any](ma ReaderIOResult[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReaderIOEither[context.Context, error](ma)
}

//go:inline
func FromReaderIO[R, A any](ma ReaderIO[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReaderIO[context.Context, error](ma)
}

//go:inline
func RightReaderIO[R, A any](ma ReaderIO[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.RightReaderIO[context.Context, error](ma)
}

//go:inline
func LeftReaderIO[A, R any](me ReaderIO[R, error]) ReaderReaderIOResult[R, A] {
	return RRIOE.LeftReaderIO[context.Context, A](me)
}

//go:inline
func MonadMap[R, A, B any](fa ReaderReaderIOResult[R, A], f func(A) B) ReaderReaderIOResult[R, B] {
	return reader.MonadMap(fa, RIOE.Map(f))
}

//go:inline
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return reader.Map[R](RIOE.Map(f))
}

//go:inline
func MonadMapTo[R, A, B any](fa ReaderReaderIOResult[R, A], b B) ReaderReaderIOResult[R, B] {
	return reader.MonadMap(fa, RIOE.MapTo[A](b))
}

//go:inline
func MapTo[R, A, B any](b B) Operator[R, A, B] {
	return reader.Map[R](RIOE.MapTo[A](b))
}

//go:inline
func MonadChain[R, A, B any](fa ReaderReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return readert.MonadChain(
		RIOE.MonadChain[A, B],
		fa,
		f,
	)
}

//go:inline
func MonadChainFirst[R, A, B any](fa ReaderReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return chain.MonadChainFirst(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		fa,
		f)
}

//go:inline
func MonadTap[R, A, B any](fa ReaderReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirst(fa, f)
}

//go:inline
func MonadChainEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderReaderIOResult[R, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, A, B],
		FromEither[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, A, B],
		FromEither[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderReaderIOResult[R, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromEither[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, f)
}

//go:inline
func ChainFirstEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, A, A],
		Map[R, B, A],
		FromEither[R, B],
		f,
	)
}

//go:inline
func TapEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A] {
	return ChainFirstEitherK[R](f)
}

//go:inline
func MonadChainReaderK[R, A, B any](ma ReaderReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReader[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderK[R, A, B any](ma ReaderReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderK[R, A, B any](ma ReaderReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderK(ma, f)
}

//go:inline
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReader[R, B],
		f,
	)
}

//go:inline
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderK(f)
}

//go:inline
func MonadChainReaderIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReaderIO[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReaderIO[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReaderIO[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderIOK(ma, f)
}

//go:inline
func ChainFirstReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReaderIO[R, B],
		f,
	)
}

//go:inline
func TapReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderIOK(f)
}

//go:inline
func MonadChainReaderEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReaderEither[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReaderEither[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReaderEither[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderEitherK(ma, f)
}

//go:inline
func ChainFirstReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReaderEither[R, B],
		f,
	)
}

//go:inline
func TapReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return ChainFirstReaderEitherK(f)
}

//go:inline
func ChainReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, B] {
	return RRIOE.ChainReaderOptionK[R, context.Context, A, B](onNone)
}

func ChainFirstReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
	return RRIOE.ChainFirstReaderOptionK[R, context.Context, A, B](onNone)
}

//go:inline
func TapReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderOptionK[R, A, B](onNone)
}

//go:inline
func MonadChainIOEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f IOE.Kleisli[error, A, B]) ReaderReaderIOResult[R, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChain[R, A, B],
		FromIOEither[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainIOEitherK[R, A, B any](f IOE.Kleisli[error, A, B]) Operator[R, A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[R, A, B],
		FromIOEither[R, B],
		f,
	)
}

//go:inline
func MonadChainIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderReaderIOResult[R, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, A, B],
		FromIO[R, B],
		ma,
		f,
	)
}

//go:inline
func ChainIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, B] {
	return fromio.ChainIOK(
		Chain[R, A, B],
		FromIO[R, B],
		f,
	)
}

//go:inline
func MonadChainFirstIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromIO[R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstIOK(ma, f)
}

//go:inline
func ChainFirstIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, A, A],
		Map[R, B, A],
		FromIO[R, B],
		f,
	)
}

//go:inline
func TapIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOK[R](f)
}

//go:inline
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[R, A, B] {
	return fromeither.ChainOptionK(
		MonadChain[R, A, B],
		FromEither[R, B],
		onNone,
	)
}

//go:inline
func MonadAp[R, A, B any](fab ReaderReaderIOResult[R, func(A) B], fa ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, B] {
	return readert.MonadAp[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.MonadAp[B, A],
		fab,
		fa,
	)
}

//go:inline
func MonadApSeq[R, A, B any](fab ReaderReaderIOResult[R, func(A) B], fa ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, B] {
	return readert.MonadAp[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.MonadApSeq[B, A],
		fab,
		fa,
	)
}

//go:inline
func MonadApPar[R, A, B any](fab ReaderReaderIOResult[R, func(A) B], fa ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, B] {
	return readert.MonadAp[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.MonadApPar[B, A],
		fab,
		fa,
	)
}

//go:inline
func Ap[B, R, A any](fa ReaderReaderIOResult[R, A]) Operator[R, func(A) B, B] {
	return readert.Ap[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.Ap[B, A],
		fa,
	)
}

//go:inline
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return readert.Chain[ReaderReaderIOResult[R, A]](
		RIOE.Chain[A, B],
		f,
	)
}

//go:inline
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return chain.ChainFirst(
		Chain[R, A, A],
		Map[R, B, A],
		f)
}

//go:inline
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirst(f)
}

//go:inline
func Right[R, A any](a A) ReaderReaderIOResult[R, A] {
	return RRIOE.Right[R, context.Context, error](a)
}

//go:inline
func Left[R, A any](e error) ReaderReaderIOResult[R, A] {
	return RRIOE.Left[R, context.Context, A](e)
}

//go:inline
func Of[R, A any](a A) ReaderReaderIOResult[R, A] {
	return RRIOE.Of[R, context.Context, error](a)
}

//go:inline
func Flatten[R, A any](mma ReaderReaderIOResult[R, ReaderReaderIOResult[R, A]]) ReaderReaderIOResult[R, A] {
	return MonadChain(mma, function.Identity[ReaderReaderIOResult[R, A]])
}

//go:inline
func FromEither[R, A any](t Either[error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromEither[R, context.Context](t)
}

//go:inline
func FromResult[R, A any](t Result[A]) ReaderReaderIOResult[R, A] {
	return FromEither[R](t)
}

//go:inline
func RightReader[R, A any](ma Reader[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.RightReader[context.Context, error](ma)
}

//go:inline
func LeftReader[A, R any](ma Reader[R, error]) ReaderReaderIOResult[R, A] {
	return RRIOE.LeftReader[context.Context, A](ma)
}

//go:inline
func FromReader[R, A any](ma Reader[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReader[context.Context, error](ma)
}

//go:inline
func RightIO[R, A any](ma IO[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.RightIO[R, context.Context, error](ma)
}

//go:inline
func LeftIO[R, A any](ma IO[error]) ReaderReaderIOResult[R, A] {
	return RRIOE.LeftIO[R, context.Context, A](ma)
}

//go:inline
func FromIO[R, A any](ma IO[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromIO[R, context.Context, error](ma)
}

//go:inline
func FromIOEither[R, A any](ma IOEither[error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromIOEither[R, context.Context, error](ma)
}

//go:inline
func FromIOResult[R, A any](ma IOResult[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromIOEither[R, context.Context, error](ma)
}

//go:inline
func FromReaderEither[R, A any](ma RE.ReaderEither[R, error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReaderEither[R, context.Context, error](ma)
}

//go:inline
func Ask[R any]() ReaderReaderIOResult[R, R] {
	return RRIOE.Ask[R, context.Context, error]()
}

//go:inline
func Asks[R, A any](r Reader[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.Asks[context.Context, error](r)
}

//go:inline
func FromOption[R, A any](onNone Lazy[error]) func(Option[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromOption[R, context.Context, A](onNone)
}

//go:inline
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A] {
	return RRIOE.FromPredicate[R, context.Context, error](pred, onFalse)
}

//go:inline
func MonadAlt[R, A any](first ReaderReaderIOResult[R, A], second Lazy[ReaderReaderIOResult[R, A]]) ReaderReaderIOResult[R, A] {
	return RRIOE.MonadAlt(first, second)
}

//go:inline
func Alt[R, A any](second Lazy[ReaderReaderIOResult[R, A]]) Operator[R, A, A] {
	return RRIOE.Alt(second)
}

//go:inline
func MonadFlap[R, B, A any](fab ReaderReaderIOResult[R, func(A) B], a A) ReaderReaderIOResult[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}

//go:inline
func MonadMapLeft[R, A any](fa ReaderReaderIOResult[R, A], f Endmorphism[error]) ReaderReaderIOResult[R, A] {
	return RRIOE.MonadMapLeft[R, context.Context](fa, f)
}

//go:inline
func MapLeft[R, A any](f Endmorphism[error]) Operator[R, A, A] {
	return RRIOE.MapLeft[R, context.Context, A](f)
}

//go:inline
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.Local[context.Context, error, A](f)
}

//go:inline
func Read[A, R any](r R) func(ReaderReaderIOResult[R, A]) ReaderIOResult[context.Context, A] {
	return RRIOE.Read[context.Context, error, A](r)
}

//go:inline
func ReadIOEither[A, R any](rio IOEither[error, R]) func(ReaderReaderIOResult[R, A]) ReaderIOResult[context.Context, A] {
	return RRIOE.ReadIOEither[A, R, context.Context](rio)
}

//go:inline
func ReadIO[A, R any](rio IO[R]) func(ReaderReaderIOResult[R, A]) ReaderIOResult[context.Context, A] {
	return RRIOE.ReadIO[context.Context, error, A, R](rio)
}

//go:inline
func MonadChainLeft[R, A any](fa ReaderReaderIOResult[R, A], f Kleisli[R, error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.MonadChainLeft[R, context.Context, error, error, A](fa, f)
}

//go:inline
func ChainLeft[R, A any](f Kleisli[R, error, A]) func(ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.ChainLeft[R, context.Context, error, error, A](f)
}

//go:inline
func Delay[R, A any](delay time.Duration) Operator[R, A, A] {
	return reader.Map[R](RIOE.Delay[A](delay))
}
