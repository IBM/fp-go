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

package readerreaderioeither

import (
	"time"

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
	"github.com/IBM/fp-go/v2/ioeither"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/readeroption"
)

//go:inline
func FromReaderOption[R, C, A, E any](onNone Lazy[E]) Kleisli[R, C, E, ReaderOption[R, A], A] {
	return reader.Map[R](RIOE.FromOption[C, A](onNone))
}

//go:inline
func FromReaderIOEither[C, E, R, A any](ma ReaderIOEither[R, E, A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.MonadMap(ma, RIOE.FromIOEither[C])
}

//go:inline
func FromReaderIO[C, E, R, A any](ma ReaderIO[R, A]) ReaderReaderIOEither[R, C, E, A] {
	return RightReaderIO[C, E](ma)
}

//go:inline
func RightReaderIO[C, E, R, A any](ma ReaderIO[R, A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.MonadMap(ma, RIOE.RightIO[C, E, A])
}

//go:inline
func LeftReaderIO[C, A, R, E any](me ReaderIO[R, E]) ReaderReaderIOEither[R, C, E, A] {
	return reader.MonadMap(me, RIOE.LeftIO[C, A, E])
}

//go:inline
func MonadMap[R, C, E, A, B any](fa ReaderReaderIOEither[R, C, E, A], f func(A) B) ReaderReaderIOEither[R, C, E, B] {
	return reader.MonadMap(fa, RIOE.Map[C, E](f))
}

//go:inline
func Map[R, C, E, A, B any](f func(A) B) Operator[R, C, E, A, B] {
	return reader.Map[R](RIOE.Map[C, E](f))
}

//go:inline
func MonadMapTo[R, C, E, A, B any](fa ReaderReaderIOEither[R, C, E, A], b B) ReaderReaderIOEither[R, C, E, B] {
	return reader.MonadMap(fa, RIOE.MapTo[C, E, A](b))
}

//go:inline
func MapTo[R, C, E, A, B any](b B) Operator[R, C, E, A, B] {
	return reader.Map[R](RIOE.MapTo[C, E, A](b))
}

//go:inline
func MonadChain[R, C, E, A, B any](fa ReaderReaderIOEither[R, C, E, A], f Kleisli[R, C, E, A, B]) ReaderReaderIOEither[R, C, E, B] {
	return readert.MonadChain(
		RIOE.MonadChain[C, E, A, B],
		fa,
		f,
	)
}

//go:inline
func MonadChainFirst[R, C, E, A, B any](fa ReaderReaderIOEither[R, C, E, A], f Kleisli[R, C, E, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return chain.MonadChainFirst(
		MonadChain[R, C, E, A, A],
		MonadMap[R, C, E, B, A],
		fa,
		f)
}

//go:inline
func MonadTap[R, C, E, A, B any](fa ReaderReaderIOEither[R, C, E, A], f Kleisli[R, C, E, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChainFirst(fa, f)
}

//go:inline
func MonadChainEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f either.Kleisli[E, A, B]) ReaderReaderIOEither[R, C, E, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, C, E, A, B],
		FromEither[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func ChainEitherK[R, C, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, C, E, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, C, E, A, B],
		FromEither[R, C, E, B],
		f,
	)
}

//go:inline
func MonadChainFirstEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f either.Kleisli[E, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, C, E, A, A],
		MonadMap[R, C, E, B, A],
		FromEither[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f either.Kleisli[E, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChainFirstEitherK(ma, f)
}

//go:inline
func ChainFirstEitherK[R, C, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, C, E, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, C, E, A, A],
		Map[R, C, E, B, A],
		FromEither[R, C, E, B],
		f,
	)
}

//go:inline
func TapEitherK[R, C, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, C, E, A, A] {
	return ChainFirstEitherK[R, C](f)
}

//go:inline
func MonadChainReaderK[C, E, R, A, B any](ma ReaderReaderIOEither[R, C, E, A], f reader.Kleisli[R, A, B]) ReaderReaderIOEither[R, C, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, C, E, A, B],
		FromReader[C, E, R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderK[C, E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, C, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, C, E, A, B],
		FromReader[C, E, R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderK[C, E, R, A, B any](ma ReaderReaderIOEither[R, C, E, A], f reader.Kleisli[R, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, C, E, A, B],
		FromReader[C, E, R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderK[C, E, R, A, B any](ma ReaderReaderIOEither[R, C, E, A], f reader.Kleisli[R, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChainFirstReaderK(ma, f)
}

//go:inline
func ChainFirstReaderK[C, E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, C, E, A, B],
		FromReader[C, E, R, B],
		f,
	)
}

//go:inline
func TapReaderK[C, E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
	return ChainFirstReaderK[C, E](f)
}

//go:inline
func MonadChainReaderIOK[C, E, R, A, B any](ma ReaderReaderIOEither[R, C, E, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOEither[R, C, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, C, E, A, B],
		FromReaderIO[C, E, R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderIOK[C, E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, C, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, C, E, A, B],
		FromReaderIO[C, E, R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderIOK[C, E, R, A, B any](ma ReaderReaderIOEither[R, C, E, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, C, E, A, B],
		FromReaderIO[C, E, R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderIOK[C, E, R, A, B any](ma ReaderReaderIOEither[R, C, E, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChainFirstReaderIOK(ma, f)
}

//go:inline
func ChainFirstReaderIOK[C, E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, C, E, A, B],
		FromReaderIO[C, E, R, B],
		f,
	)
}

//go:inline
func TapReaderIOK[C, E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
	return ChainFirstReaderIOK[C, E](f)
}

//go:inline
func MonadChainReaderEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f RE.Kleisli[R, E, A, B]) ReaderReaderIOEither[R, C, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, C, E, A, B],
		FromReaderEither[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderEitherK[C, E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, C, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, C, E, A, B],
		FromReaderEither[R, C, E, B],
		f,
	)
}

//go:inline
func ChainReaderIOEitherK[C, R, E, A, B any](f RIOE.Kleisli[R, E, A, B]) Operator[R, C, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, C, E, A, B],
		FromReaderIOEither[C, E, R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f RE.Kleisli[R, E, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, C, E, A, B],
		FromReaderEither[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f RE.Kleisli[R, E, A, B]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChainFirstReaderEitherK(ma, f)
}

//go:inline
func ChainFirstReaderEitherK[C, E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, C, E, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, C, E, A, B],
		FromReaderEither[R, C, E, B],
		f,
	)
}

//go:inline
func TapReaderEitherK[C, E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, C, E, A, A] {
	return ChainFirstReaderEitherK[C](f)
}

func ChainReaderOptionK[R, C, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, B] {

	fro := FromReaderOption[R, C, B](onNone)

	return func(f readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, B] {
		return fromreader.ChainReaderK(
			Chain[R, C, E, A, B],
			fro,
			f,
		)

	}
}

func ChainFirstReaderOptionK[R, C, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
	fro := FromReaderOption[R, C, B](onNone)
	return func(f readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
		return fromreader.ChainFirstReaderK(
			ChainFirst[R, C, E, A, B],
			fro,
			f,
		)
	}
}

//go:inline
func TapReaderOptionK[R, C, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, C, E, A, A] {
	return ChainFirstReaderOptionK[R, C, A, B](onNone)
}

//go:inline
func MonadChainIOEitherK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f IOE.Kleisli[E, A, B]) ReaderReaderIOEither[R, C, E, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChain[R, C, E, A, B],
		FromIOEither[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func ChainIOEitherK[R, C, E, A, B any](f IOE.Kleisli[E, A, B]) Operator[R, C, E, A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[R, C, E, A, B],
		FromIOEither[R, C, E, B],
		f,
	)
}

//go:inline
func MonadChainIOK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f io.Kleisli[A, B]) ReaderReaderIOEither[R, C, E, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, C, E, A, B],
		FromIO[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func ChainIOK[R, C, E, A, B any](f io.Kleisli[A, B]) Operator[R, C, E, A, B] {
	return fromio.ChainIOK(
		Chain[R, C, E, A, B],
		FromIO[R, C, E, B],
		f,
	)
}

//go:inline
func MonadChainFirstIOK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f io.Kleisli[A, B]) ReaderReaderIOEither[R, C, E, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, C, E, A, A],
		MonadMap[R, C, E, B, A],
		FromIO[R, C, E, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapIOK[R, C, E, A, B any](ma ReaderReaderIOEither[R, C, E, A], f io.Kleisli[A, B]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChainFirstIOK(ma, f)
}

//go:inline
func ChainFirstIOK[R, C, E, A, B any](f io.Kleisli[A, B]) Operator[R, C, E, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, C, E, A, A],
		Map[R, C, E, B, A],
		FromIO[R, C, E, B],
		f,
	)
}

//go:inline
func TapIOK[R, C, E, A, B any](f io.Kleisli[A, B]) Operator[R, C, E, A, A] {
	return ChainFirstIOK[R, C, E](f)
}

//go:inline
func ChainOptionK[R, C, A, B, E any](onNone Lazy[E]) func(option.Kleisli[A, B]) Operator[R, C, E, A, B] {
	return fromeither.ChainOptionK(
		MonadChain[R, C, E, A, B],
		FromEither[R, C, E, B],
		onNone,
	)
}

//go:inline
func MonadAp[R, C, E, A, B any](fab ReaderReaderIOEither[R, C, E, func(A) B], fa ReaderReaderIOEither[R, C, E, A]) ReaderReaderIOEither[R, C, E, B] {
	return readert.MonadAp[
		ReaderReaderIOEither[R, C, E, A],
		ReaderReaderIOEither[R, C, E, B],
		ReaderReaderIOEither[R, C, E, func(A) B], R, A](
		RIOE.MonadAp[C, E, A, B],
		fab,
		fa,
	)
}

//go:inline
func MonadApSeq[R, C, E, A, B any](fab ReaderReaderIOEither[R, C, E, func(A) B], fa ReaderReaderIOEither[R, C, E, A]) ReaderReaderIOEither[R, C, E, B] {
	return readert.MonadAp[
		ReaderReaderIOEither[R, C, E, A],
		ReaderReaderIOEither[R, C, E, B],
		ReaderReaderIOEither[R, C, E, func(A) B], R, A](
		RIOE.MonadApSeq[C, E, A, B],
		fab,
		fa,
	)
}

//go:inline
func MonadApPar[R, C, E, A, B any](fab ReaderReaderIOEither[R, C, E, func(A) B], fa ReaderReaderIOEither[R, C, E, A]) ReaderReaderIOEither[R, C, E, B] {
	return readert.MonadAp[
		ReaderReaderIOEither[R, C, E, A],
		ReaderReaderIOEither[R, C, E, B],
		ReaderReaderIOEither[R, C, E, func(A) B], R, A](
		RIOE.MonadApPar[C, E, A, B],
		fab,
		fa,
	)
}

//go:inline
func Ap[B, R, C, E, A any](fa ReaderReaderIOEither[R, C, E, A]) Operator[R, C, E, func(A) B, B] {
	return readert.Ap[
		ReaderReaderIOEither[R, C, E, A],
		ReaderReaderIOEither[R, C, E, B],
		ReaderReaderIOEither[R, C, E, func(A) B], R, A](
		RIOE.Ap[B, C, E, A],
		fa,
	)
}

//go:inline
func Chain[R, C, E, A, B any](f Kleisli[R, C, E, A, B]) Operator[R, C, E, A, B] {
	return readert.Chain[ReaderReaderIOEither[R, C, E, A]](
		RIOE.Chain[C, E, A, B],
		f,
	)
}

//go:inline
func ChainFirst[R, C, E, A, B any](f Kleisli[R, C, E, A, B]) Operator[R, C, E, A, A] {
	return chain.ChainFirst(
		Chain[R, C, E, A, A],
		Map[R, C, E, B, A],
		f)
}

//go:inline
func Tap[R, C, E, A, B any](f Kleisli[R, C, E, A, B]) Operator[R, C, E, A, A] {
	return ChainFirst(f)
}

//go:inline
func Right[R, C, E, A any](a A) ReaderReaderIOEither[R, C, E, A] {
	return reader.Of[R](RIOE.Right[C, E](a))
}

//go:inline
func Left[R, C, A, E any](e E) ReaderReaderIOEither[R, C, E, A] {
	return reader.Of[R](RIOE.Left[C, A](e))
}

//go:inline
func Of[R, C, E, A any](a A) ReaderReaderIOEither[R, C, E, A] {
	return Right[R, C, E](a)
}

//go:inline
func Flatten[R, C, E, A any](mma ReaderReaderIOEither[R, C, E, ReaderReaderIOEither[R, C, E, A]]) ReaderReaderIOEither[R, C, E, A] {
	return MonadChain(mma, function.Identity[ReaderReaderIOEither[R, C, E, A]])
}

//go:inline
func FromEither[R, C, E, A any](t Either[E, A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.Of[R](RIOE.FromEither[C](t))
}

//go:inline
func RightReader[C, E, R, A any](ma Reader[R, A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.MonadMap(ma, RIOE.Right[C, E])
}

//go:inline
func LeftReader[C, A, R, E any](ma Reader[R, E]) ReaderReaderIOEither[R, C, E, A] {
	return reader.MonadMap(ma, RIOE.Left[C, A])
}

//go:inline
func FromReader[C, E, R, A any](ma Reader[R, A]) ReaderReaderIOEither[R, C, E, A] {
	return RightReader[C, E](ma)
}

//go:inline
func RightIO[R, C, E, A any](ma IO[A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.Of[R](RIOE.RightIO[C, E](ma))
}

//go:inline
func LeftIO[R, C, A, E any](ma IO[E]) ReaderReaderIOEither[R, C, E, A] {
	return reader.Of[R](RIOE.LeftIO[C, A](ma))
}

//go:inline
func FromIO[R, C, E, A any](ma IO[A]) ReaderReaderIOEither[R, C, E, A] {
	return RightIO[R, C, E](ma)
}

//go:inline
func FromIOEither[R, C, E, A any](ma IOEither[E, A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.Of[R](RIOE.FromIOEither[C](ma))
}

//go:inline
func FromReaderEither[R, C, E, A any](ma RE.ReaderEither[R, E, A]) ReaderReaderIOEither[R, C, E, A] {
	return reader.MonadMap(ma, RIOE.FromEither[C])
}

//go:inline
func Ask[R, C, E any]() ReaderReaderIOEither[R, C, E, R] {
	return fromreader.Ask(FromReader[C, E, R, R])()
}

//go:inline
func Asks[C, E, R, A any](r Reader[R, A]) ReaderReaderIOEither[R, C, E, A] {
	return fromreader.Asks(FromReader[C, E, R, A])(r)
}

//go:inline
func FromOption[R, C, A, E any](onNone Lazy[E]) func(Option[A]) ReaderReaderIOEither[R, C, E, A] {
	return fromeither.FromOption(FromEither[R, C, E, A], onNone)
}

//go:inline
func FromPredicate[R, C, E, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderReaderIOEither[R, C, E, A] {
	return fromeither.FromPredicate(FromEither[R, C, E, A], pred, onFalse)
}

//go:inline
func MonadAlt[R, C, E, A any](first ReaderReaderIOEither[R, C, E, A], second Lazy[ReaderReaderIOEither[R, C, E, A]]) ReaderReaderIOEither[R, C, E, A] {
	return func(r R) ReaderIOEither[C, E, A] {
		return RIOE.MonadAlt(first(r), func() ReaderIOEither[C, E, A] {
			return second()(r)
		})
	}
}

//go:inline
func Alt[R, C, E, A any](second Lazy[ReaderReaderIOEither[R, C, E, A]]) Operator[R, C, E, A, A] {
	return function.Bind2nd(MonadAlt, second)
}

//go:inline
func MonadFlap[R, C, E, B, A any](fab ReaderReaderIOEither[R, C, E, func(A) B], a A) ReaderReaderIOEither[R, C, E, B] {
	return functor.MonadFlap(MonadMap[R, C, E, func(A) B, B], fab, a)
}

//go:inline
func Flap[R, C, E, B, A any](a A) Operator[R, C, E, func(A) B, B] {
	return functor.Flap(Map[R, C, E, func(A) B, B], a)
}

//go:inline
func MonadMapLeft[R, C, E1, E2, A any](fa ReaderReaderIOEither[R, C, E1, A], f func(E1) E2) ReaderReaderIOEither[R, C, E2, A] {
	return reader.MonadMap(fa, RIOE.MapLeft[C, A](f))
}

//go:inline
func MapLeft[R, C, A, E1, E2 any](f func(E1) E2) func(ReaderReaderIOEither[R, C, E1, A]) ReaderReaderIOEither[R, C, E2, A] {
	return reader.Map[R](RIOE.MapLeft[C, A](f))
}

//go:inline
func Read[C, E, A, R any](r R) func(ReaderReaderIOEither[R, C, E, A]) ReaderIOEither[C, E, A] {
	return reader.Read[ReaderIOEither[C, E, A]](r)
}

//go:inline
func ReadIOEither[A, R, C, E any](rio IOEither[E, R]) func(ReaderReaderIOEither[R, C, E, A]) ReaderIOEither[C, E, A] {
	return func(rri ReaderReaderIOEither[R, C, E, A]) ReaderIOEither[C, E, A] {
		return func(c C) IOEither[E, A] {
			return function.Pipe1(
				rio,
				ioeither.Chain(func(r R) IOEither[E, A] {
					return rri(r)(c)
				}),
			)
		}
	}
}

//go:inline
func ReadIO[C, E, A, R any](rio IO[R]) func(ReaderReaderIOEither[R, C, E, A]) ReaderIOEither[C, E, A] {
	return func(rri ReaderReaderIOEither[R, C, E, A]) ReaderIOEither[C, E, A] {
		return func(c C) IOEither[E, A] {
			return function.Pipe1(
				rio,
				io.Chain(func(r R) IOEither[E, A] {
					return rri(r)(c)
				}),
			)
		}
	}
}

//go:inline
func MonadChainLeft[R, C, EA, EB, A any](fa ReaderReaderIOEither[R, C, EA, A], f Kleisli[R, C, EB, EA, A]) ReaderReaderIOEither[R, C, EB, A] {
	return readert.MonadChain(
		RIOE.MonadChainLeft[C, EA, EB, A],
		fa,
		f,
	)
}

//go:inline
func ChainLeft[R, C, EA, EB, A any](f Kleisli[R, C, EB, EA, A]) func(ReaderReaderIOEither[R, C, EA, A]) ReaderReaderIOEither[R, C, EB, A] {
	return readert.Chain[ReaderReaderIOEither[R, C, EA, A]](
		RIOE.ChainLeft[C, EA, EB, A],
		f,
	)
}

//go:inline
func Delay[R, C, E, A any](delay time.Duration) Operator[R, C, E, A, A] {
	return reader.Map[R](RIOE.Delay[C, E, A](delay))
}

//go:inline
func After[R, C, E, A any](timestamp time.Time) Operator[R, C, E, A, A] {
	return reader.Map[R](RIOE.After[C, E, A](timestamp))
}

func Defer[R, C, E, A any](fa Lazy[ReaderReaderIOEither[R, C, E, A]]) ReaderReaderIOEither[R, C, E, A] {
	return func(r R) ReaderIOEither[C, E, A] {
		return func(c C) RIOE.IOEither[E, A] {
			return func() IOE.Either[E, A] {
				return fa()(r)(c)()
			}
		}
	}
}
