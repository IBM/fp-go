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
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
)

func MonadFromReaderIO[R, E, A any](a A, f func(A) ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return function.Pipe2(
		a,
		f,
		RightReaderIO[R, E, A],
	)
}

func FromReaderIO[R, E, A any](f func(A) ReaderIO[R, A]) func(A) ReaderIOEither[R, E, A] {
	return function.Bind2nd(MonadFromReaderIO[R, E, A], f)
}

func RightReaderIO[R, E, A any](ma ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return eithert.RightF(
		readerio.MonadMap[R, A, either.Either[E, A]],
		ma,
	)
}

func LeftReaderIO[A, R, E any](me ReaderIO[R, E]) ReaderIOEither[R, E, A] {
	return eithert.LeftF(
		readerio.MonadMap[R, E, either.Either[E, A]],
		me,
	)
}

func MonadMap[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) B) ReaderIOEither[R, E, B] {
	return eithert.MonadMap(readerio.MonadMap[R, either.Either[E, A], either.Either[E, B]], fa, f)
}

func Map[R, E, A, B any](f func(A) B) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.Map(readerio.Map[R, either.Either[E, A], either.Either[E, B]], f)
}

func MonadMapTo[R, E, A, B any](fa ReaderIOEither[R, E, A], b B) ReaderIOEither[R, E, B] {
	return MonadMap(fa, function.Constant1[A](b))
}

func MapTo[R, E, A, B any](b B) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return Map[R, E](function.Constant1[A](b))
}

func MonadChain[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) ReaderIOEither[R, E, B]) ReaderIOEither[R, E, B] {
	return eithert.MonadChain(
		readerio.MonadChain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		fa,
		f)
}

func MonadChainFirst[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) ReaderIOEither[R, E, B]) ReaderIOEither[R, E, A] {
	return chain.MonadChainFirst(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		fa,
		f)
}

func MonadChainEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) either.Either[E, B]) ReaderIOEither[R, E, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		ma,
		f,
	)
}

func ChainEitherK[R, E, A, B any](f func(A) either.Either[E, B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return fromeither.ChainEitherK(
		Chain[R, E, A, B],
		FromEither[R, E, B],
		f,
	)
}

func MonadChainFirstEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) either.Either[E, B]) ReaderIOEither[R, E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		FromEither[R, E, B],
		ma,
		f,
	)
}

func ChainFirstEitherK[R, E, A, B any](f func(A) either.Either[E, B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		FromEither[R, E, B],
		f,
	)
}

func MonadChainReaderK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) Reader[R, B]) ReaderIOEither[R, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, E, A, B],
		FromReader[E, R, B],
		ma,
		f,
	)
}

func ChainReaderK[E, R, A, B any](f func(A) Reader[R, B]) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return fromreader.ChainReaderK(
		MonadChain[R, E, A, B],
		FromReader[E, R, B],
		f,
	)
}

func MonadChainIOEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) IOE.IOEither[E, B]) ReaderIOEither[R, E, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChain[R, E, A, B],
		FromIOEither[R, E, B],
		ma,
		f,
	)
}

func ChainIOEitherK[R, E, A, B any](f func(A) IOE.IOEither[E, B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return fromioeither.ChainIOEitherK(
		Chain[R, E, A, B],
		FromIOEither[R, E, B],
		f,
	)
}

func MonadChainIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) io.IO[B]) ReaderIOEither[R, E, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, E, A, B],
		FromIO[R, E, B],
		ma,
		f,
	)
}

func ChainIOK[R, E, A, B any](f func(A) io.IO[B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return fromio.ChainIOK(
		Chain[R, E, A, B],
		FromIO[R, E, B],
		f,
	)
}

func MonadChainFirstIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) io.IO[B]) ReaderIOEither[R, E, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		FromIO[R, E, B],
		ma,
		f,
	)
}

func ChainFirstIOK[R, E, A, B any](f func(A) io.IO[B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return fromio.ChainFirstIOK(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		FromIO[R, E, B],
		f,
	)
}

func ChainOptionK[R, A, B, E any](onNone func() E) func(func(A) O.Option[B]) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return fromeither.ChainOptionK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		onNone,
	)
}

func MonadAp[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadAp[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

func MonadApSeq[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadApSeq[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

func MonadApPar[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadApPar[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(fab ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return function.Bind2nd(MonadAp[R, E, A, B], fa)
}

func Chain[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.Chain(
		readerio.Chain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		f)
}

func ChainFirst[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return chain.ChainFirst(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		f)
}

func Right[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return eithert.Right(readerio.Of[R, Either[E, A]], a)
}

func Left[R, A, E any](e E) ReaderIOEither[R, E, A] {
	return eithert.Left(readerio.Of[R, Either[E, A]], e)
}

func ThrowError[R, A, E any](e E) ReaderIOEither[R, E, A] {
	return Left[R, A](e)
}

// Of returns a Reader with a fixed value
func Of[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return Right[R, E](a)
}

func Flatten[R, E, A any](mma ReaderIOEither[R, E, ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return MonadChain(mma, function.Identity[ReaderIOEither[R, E, A]])
}

func FromEither[R, E, A any](t either.Either[E, A]) ReaderIOEither[R, E, A] {
	return readerio.Of[R](t)
}

func RightReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, ioeither.Right[E, A])
}

func LeftReader[A, R, E any](ma Reader[R, E]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, ioeither.Left[A, E])
}

func FromReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A] {
	return RightReader[E](ma)
}

func RightIO[R, E, A any](ma io.IO[A]) ReaderIOEither[R, E, A] {
	return function.Pipe2(ma, ioeither.RightIO[E, A], FromIOEither[R, E, A])
}

func LeftIO[R, A, E any](ma io.IO[E]) ReaderIOEither[R, E, A] {
	return function.Pipe2(ma, ioeither.LeftIO[A, E], FromIOEither[R, E, A])
}

func FromIO[R, E, A any](ma io.IO[A]) ReaderIOEither[R, E, A] {
	return RightIO[R, E](ma)
}

func FromIOEither[R, E, A any](ma IOE.IOEither[E, A]) ReaderIOEither[R, E, A] {
	return reader.Of[R](ma)
}

func FromReaderEither[R, E, A any](ma RE.ReaderEither[R, E, A]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, ioeither.FromEither[E, A])
}

func Ask[R, E any]() ReaderIOEither[R, E, R] {
	return fromreader.Ask(FromReader[E, R, R])()
}

func Asks[E, R, A any](r Reader[R, A]) ReaderIOEither[R, E, A] {
	return fromreader.Asks(FromReader[E, R, A])(r)
}

func FromOption[R, A, E any](onNone func() E) func(O.Option[A]) ReaderIOEither[R, E, A] {
	return fromeither.FromOption(FromEither[R, E, A], onNone)
}

func FromPredicate[R, E, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderIOEither[R, E, A] {
	return fromeither.FromPredicate(FromEither[R, E, A], pred, onFalse)
}

func Fold[R, E, A, B any](onLeft func(E) ReaderIO[R, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOEither[R, E, A]) ReaderIO[R, B] {
	return eithert.MatchE(readerio.MonadChain[R, either.Either[E, A], B], onLeft, onRight)
}

func GetOrElse[R, E, A any](onLeft func(E) ReaderIO[R, A]) func(ReaderIOEither[R, E, A]) ReaderIO[R, A] {
	return eithert.GetOrElse(readerio.MonadChain[R, either.Either[E, A], A], readerio.Of[R, A], onLeft)
}

func OrElse[R, E1, A, E2 any](onLeft func(E1) ReaderIOEither[R, E2, A]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.OrElse(readerio.MonadChain[R, either.Either[E1, A], either.Either[E2, A]], readerio.Of[R, either.Either[E2, A]], onLeft)
}

func OrLeft[A, E1, R, E2 any](onLeft func(E1) ReaderIO[R, E2]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.OrLeft(
		readerio.MonadChain[R, either.Either[E1, A], either.Either[E2, A]],
		readerio.MonadMap[R, E2, either.Either[E2, A]],
		readerio.Of[R, either.Either[E2, A]],
		onLeft,
	)
}

func MonadBiMap[R, E1, E2, A, B any](fa ReaderIOEither[R, E1, A], f func(E1) E2, g func(A) B) ReaderIOEither[R, E2, B] {
	return eithert.MonadBiMap(
		readerio.MonadMap[R, either.Either[E1, A], either.Either[E2, B]],
		fa, f, g,
	)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[R, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, B] {
	return eithert.BiMap(readerio.Map[R, either.Either[E1, A], either.Either[E2, B]], f, g)
}

// Swap changes the order of type parameters
func Swap[R, E, A any](val ReaderIOEither[R, E, A]) ReaderIOEither[R, A, E] {
	return reader.MonadMap(val, ioeither.Swap[E, A])
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[R, E, A any](gen L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return readerio.Defer(gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[R, E, A any](f func(R) func() (A, error), onThrow func(error) E) ReaderIOEither[R, E, A] {
	return func(r R) IOEither[E, A] {
		return ioeither.TryCatch(f(r), onThrow)
	}
}

// MonadAlt identifies an associative operation on a type constructor.
func MonadAlt[R, E, A any](first ReaderIOEither[R, E, A], second L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return eithert.MonadAlt(
		readerio.Of[R, Either[E, A]],
		readerio.MonadChain[R, Either[E, A], Either[E, A]],

		first,
		second,
	)
}

// Alt identifies an associative operation on a type constructor.
func Alt[R, E, A any](second L.Lazy[ReaderIOEither[R, E, A]]) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return eithert.Alt(
		readerio.Of[R, Either[E, A]],
		readerio.MonadChain[R, Either[E, A], Either[E, A]],

		second,
	)
}

// Memoize computes the value of the provided [ReaderIOEither] monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[
	R, E, A any](rdr ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return readerio.Memoize(rdr)
}

func MonadFlap[R, E, B, A any](fab ReaderIOEither[R, E, func(A) B], a A) ReaderIOEither[R, E, B] {
	return functor.MonadFlap(MonadMap[R, E, func(A) B, B], fab, a)
}

func Flap[R, E, B, A any](a A) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return functor.Flap(Map[R, E, func(A) B, B], a)
}

func MonadMapLeft[R, E1, E2, A any](fa ReaderIOEither[R, E1, A], f func(E1) E2) ReaderIOEither[R, E2, A] {
	return eithert.MonadMapLeft(readerio.MonadMap[R, Either[E1, A], Either[E2, A]], fa, f)
}

// MapLeft applies a mapping function to the error channel
func MapLeft[R, A, E1, E2 any](f func(E1) E2) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.MapLeft(readerio.Map[R, Either[E1, A], Either[E2, A]], f)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[R1, R2, E, A any](f func(R2) R1) func(ReaderIOEither[R1, E, A]) ReaderIOEither[R2, E, A] {
	return reader.Local[R2, R1, IOEither[E, A]](f)
}
