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
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
	"github.com/IBM/fp-go/internal/eithert"
	FE "github.com/IBM/fp-go/internal/fromeither"
	FIO "github.com/IBM/fp-go/internal/fromio"
	FIOE "github.com/IBM/fp-go/internal/fromioeither"
	FR "github.com/IBM/fp-go/internal/fromreader"
	FC "github.com/IBM/fp-go/internal/functor"
	IOE "github.com/IBM/fp-go/ioeither/generic"
	O "github.com/IBM/fp-go/option"
	RD "github.com/IBM/fp-go/reader/generic"
	G "github.com/IBM/fp-go/readerio/generic"
)

// MakeReader constructs an instance of a reader
func MakeReader[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](f func(R) GIOA) GEA {
	return f
}

func MonadAlt[LAZY ~func() GEA, GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](first GEA, second LAZY) GEA {
	return eithert.MonadAlt(
		G.Of[GEA],
		G.MonadChain[GEA, GEA],

		first,
		second,
	)
}

func Alt[LAZY ~func() GEA, GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](second LAZY) func(GEA) GEA {
	return F.Bind2nd(MonadAlt[LAZY], second)
}

func MonadMap[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](fa GEA, f func(A) B) GEB {
	return eithert.MonadMap(G.MonadMap[GEA, GEB, GIOA, GIOB, R, ET.Either[E, A], ET.Either[E, B]], fa, f)
}

func Map[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](f func(A) B) func(GEA) GEB {
	return F.Bind2nd(MonadMap[GEA, GEB, GIOA, GIOB, R, E, A, B], f)
}

func MonadMapTo[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](fa GEA, b B) GEB {
	return MonadMap[GEA, GEB](fa, F.Constant1[A](b))
}

func MapTo[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](b B) func(GEA) GEB {
	return Map[GEA, GEB](F.Constant1[A](b))
}

func MonadChain[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](fa GEA, f func(A) GEB) GEB {
	return eithert.MonadChain(
		G.MonadChain[GEA, GEB, GIOA, GIOB, R, ET.Either[E, A], ET.Either[E, B]],
		G.Of[GEB, GIOB, R, ET.Either[E, B]],
		fa,
		f)
}

func Chain[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](f func(A) GEB) func(fa GEA) GEB {
	return F.Bind2nd(MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B], f)
}

func MonadChainFirst[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](fa GEA, f func(A) GEB) GEA {
	return C.MonadChainFirst(
		MonadChain[GEA, GEA, GIOA, GIOA, R, E, A, A],
		MonadMap[GEB, GEA, GIOB, GIOA, R, E, B, A],
		fa,
		f)
}

func ChainFirst[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](f func(A) GEB) func(fa GEA) GEA {
	return F.Bind2nd(MonadChainFirst[GEA, GEB, GIOA, GIOB, R, E, A, B], f)
}

func MonadChainEitherK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](ma GEA, f func(A) ET.Either[E, B]) GEB {
	return FE.MonadChainEitherK(
		MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B],
		FromEither[GEB, GIOB, R, E, B],
		ma,
		f,
	)
}

func ChainEitherK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](f func(A) ET.Either[E, B]) func(ma GEA) GEB {
	return F.Bind2nd(MonadChainEitherK[GEA, GEB, GIOA, GIOB, R, E, A, B], f)
}

func MonadChainFirstEitherK[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A, B any](ma GEA, f func(A) ET.Either[E, B]) GEA {
	return FE.MonadChainFirstEitherK(
		MonadChain[GEA, GEA, GIOA, GIOA, R, E, A, A],
		MonadMap[func(R) func() ET.Either[E, B], GEA, func() ET.Either[E, B], GIOA, R, E, B, A],
		FromEither[func(R) func() ET.Either[E, B], func() ET.Either[E, B], R, E, B],
		ma,
		f,
	)
}

func ChainFirstEitherK[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A, B any](f func(A) ET.Either[E, B]) func(ma GEA) GEA {
	return F.Bind2nd(MonadChainFirstEitherK[GEA, GIOA, R, E, A, B], f)
}

func MonadChainFirstIOK[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GIO ~func() B, R, E, A, B any](ma GEA, f func(A) GIO) GEA {
	return FIO.MonadChainFirstIOK(
		MonadChain[GEA, GEA, GIOA, GIOA, R, E, A, A],
		MonadMap[func(R) func() ET.Either[E, B], GEA, func() ET.Either[E, B], GIOA, R, E, B, A],
		FromIO[func(R) func() ET.Either[E, B], func() ET.Either[E, B], GIO, R, E, B],
		ma,
		f,
	)
}

func ChainFirstIOK[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GIO ~func() B, R, E, A, B any](f func(A) GIO) func(GEA) GEA {
	return F.Bind2nd(MonadChainFirstIOK[GEA, GIOA, GIO, R, E, A, B], f)
}

func MonadChainReaderK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], GB ~func(R) B, R, E, A, B any](ma GEA, f func(A) GB) GEB {
	return FR.MonadChainReaderK(
		MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B],
		FromReader[GB, GEB, GIOB, R, E, B],
		ma,
		f,
	)
}

func ChainReaderK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], GB ~func(R) B, R, E, A, B any](f func(A) GB) func(GEA) GEB {
	return FR.ChainReaderK(
		MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B],
		FromReader[GB, GEB, GIOB, R, E, B],
		f,
	)
}

func MonadChainReaderIOK[GEA ~func(R) GIOEA, GEB ~func(R) GIOEB, GIOEA ~func() ET.Either[E, A], GIOEB ~func() ET.Either[E, B], GIOB ~func() B, GB ~func(R) GIOB, R, E, A, B any](ma GEA, f func(A) GB) GEB {
	return FR.MonadChainReaderK(
		MonadChain[GEA, GEB, GIOEA, GIOEB, R, E, A, B],
		RightReaderIO[GEB, GIOEB, GB, GIOB, R, E, B],
		ma,
		f,
	)
}

func ChainReaderIOK[GEA ~func(R) GIOEA, GEB ~func(R) GIOEB, GIOEA ~func() ET.Either[E, A], GIOEB ~func() ET.Either[E, B], GIOB ~func() B, GB ~func(R) GIOB, R, E, A, B any](f func(A) GB) func(GEA) GEB {
	return FR.ChainReaderK(
		MonadChain[GEA, GEB, GIOEA, GIOEB, R, E, A, B],
		RightReaderIO[GEB, GIOEB, GB, GIOB, R, E, B],
		f,
	)
}

func MonadChainIOEitherK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](ma GEA, f func(A) GIOB) GEB {
	return FIOE.MonadChainIOEitherK(
		MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B],
		FromIOEither[GEB, GIOB, R, E, B],
		ma,
		f,
	)
}

func ChainIOEitherK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](f func(A) GIOB) func(GEA) GEB {
	return F.Bind2nd(MonadChainIOEitherK[GEA, GEB, GIOA, GIOB, R, E, A, B], f)
}

func MonadChainIOK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], GIO ~func() B, R, E, A, B any](ma GEA, f func(A) GIO) GEB {
	return FIO.MonadChainIOK(
		MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B],
		FromIO[GEB, GIOB, GIO, R, E, B],
		ma,
		f,
	)
}

func ChainIOK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], GIO ~func() B, R, E, A, B any](f func(A) GIO) func(GEA) GEB {
	return F.Bind2nd(MonadChainIOK[GEA, GEB, GIOA, GIOB, GIO, R, E, A, B], f)
}

func ChainOptionK[GEA ~func(R) GIOA, GEB ~func(R) GIOB, GIOA ~func() ET.Either[E, A], GIOB ~func() ET.Either[E, B], R, E, A, B any](onNone func() E) func(func(A) O.Option[B]) func(GEA) GEB {
	return FE.ChainOptionK(MonadChain[GEA, GEB, GIOA, GIOB, R, E, A, B], FromEither[GEB, GIOB, R, E, B], onNone)
}

func MonadAp[
	GEA ~func(R) GIOA,
	GEB ~func(R) GIOB,
	GEFAB ~func(R) GIOFAB,
	GIOA ~func() ET.Either[E, A],
	GIOB ~func() ET.Either[E, B],
	GIOFAB ~func() ET.Either[E, func(A) B],
	R, E, A, B any](fab GEFAB, fa GEA) GEB {

	return eithert.MonadAp(
		G.MonadAp[GEA, GEB, func(R) func() func(ET.Either[E, A]) ET.Either[E, B], GIOA, GIOB, func() func(ET.Either[E, A]) ET.Either[E, B], R, ET.Either[E, A], ET.Either[E, B]],
		G.MonadMap[GEFAB, func(R) func() func(ET.Either[E, A]) ET.Either[E, B], GIOFAB, func() func(ET.Either[E, A]) ET.Either[E, B], R, ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		fab,
		fa,
	)
}

func Ap[
	GEA ~func(R) GIOA,
	GEB ~func(R) GIOB,
	GEFAB ~func(R) GIOFAB,
	GIOA ~func() ET.Either[E, A],
	GIOB ~func() ET.Either[E, B],
	GIOFAB ~func() ET.Either[E, func(A) B],
	R, E, A, B any](fa GEA) func(fab GEFAB) GEB {
	return F.Bind2nd(MonadAp[GEA, GEB, GEFAB, GIOA, GIOB, GIOFAB, R, E, A, B], fa)
}

func MonadApSeq[
	GEA ~func(R) GIOA,
	GEB ~func(R) GIOB,
	GEFAB ~func(R) GIOFAB,
	GIOA ~func() ET.Either[E, A],
	GIOB ~func() ET.Either[E, B],
	GIOFAB ~func() ET.Either[E, func(A) B],
	R, E, A, B any](fab GEFAB, fa GEA) GEB {

	return eithert.MonadAp(
		G.MonadApSeq[GEA, GEB, func(R) func() func(ET.Either[E, A]) ET.Either[E, B], GIOA, GIOB, func() func(ET.Either[E, A]) ET.Either[E, B], R, ET.Either[E, A], ET.Either[E, B]],
		G.MonadMap[GEFAB, func(R) func() func(ET.Either[E, A]) ET.Either[E, B], GIOFAB, func() func(ET.Either[E, A]) ET.Either[E, B], R, ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		fab,
		fa,
	)
}

func ApSeq[
	GEA ~func(R) GIOA,
	GEB ~func(R) GIOB,
	GEFAB ~func(R) GIOFAB,
	GIOA ~func() ET.Either[E, A],
	GIOB ~func() ET.Either[E, B],
	GIOFAB ~func() ET.Either[E, func(A) B],
	R, E, A, B any](fa GEA) func(fab GEFAB) GEB {
	return F.Bind2nd(MonadApSeq[GEA, GEB, GEFAB, GIOA, GIOB, GIOFAB, R, E, A, B], fa)
}

func MonadApPar[
	GEA ~func(R) GIOA,
	GEB ~func(R) GIOB,
	GEFAB ~func(R) GIOFAB,
	GIOA ~func() ET.Either[E, A],
	GIOB ~func() ET.Either[E, B],
	GIOFAB ~func() ET.Either[E, func(A) B],
	R, E, A, B any](fab GEFAB, fa GEA) GEB {

	return eithert.MonadAp(
		G.MonadApPar[GEA, GEB, func(R) func() func(ET.Either[E, A]) ET.Either[E, B], GIOA, GIOB, func() func(ET.Either[E, A]) ET.Either[E, B], R, ET.Either[E, A], ET.Either[E, B]],
		G.MonadMap[GEFAB, func(R) func() func(ET.Either[E, A]) ET.Either[E, B], GIOFAB, func() func(ET.Either[E, A]) ET.Either[E, B], R, ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		fab,
		fa,
	)
}

func ApPar[
	GEA ~func(R) GIOA,
	GEB ~func(R) GIOB,
	GEFAB ~func(R) GIOFAB,
	GIOA ~func() ET.Either[E, A],
	GIOB ~func() ET.Either[E, B],
	GIOFAB ~func() ET.Either[E, func(A) B],
	R, E, A, B any](fa GEA) func(fab GEFAB) GEB {
	return F.Bind2nd(MonadApPar[GEA, GEB, GEFAB, GIOA, GIOB, GIOFAB, R, E, A, B], fa)
}

func Right[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](a A) GEA {
	return eithert.Right(G.Of[GEA, GIOA, R, ET.Either[E, A]], a)
}

func Left[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](e E) GEA {
	return eithert.Left(G.Of[GEA, GIOA, R, ET.Either[E, A]], e)
}

func ThrowError[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](e E) GEA {
	return Left[GEA](e)
}

// Of returns a Reader with a fixed value
func Of[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](a A) GEA {
	return Right[GEA](a)
}

func Flatten[GEA ~func(R) GIOA, GGEA ~func(R) GIOEA, GIOA ~func() ET.Either[E, A], GIOEA ~func() ET.Either[E, GEA], R, E, A any](mma GGEA) GEA {
	return MonadChain(mma, F.Identity[GEA])
}

func FromIOEither[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](t GIOA) GEA {
	return RD.Of[GEA](t)
}

func FromEither[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](t ET.Either[E, A]) GEA {
	return G.Of[GEA](t)
}

func RightReader[GA ~func(R) A, GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](ma GA) GEA {
	return F.Flow2(ma, IOE.Right[GIOA, E, A])
}

func LeftReader[GE ~func(R) E, GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](ma GE) GEA {
	return F.Flow2(ma, IOE.Left[GIOA, E, A])
}

func FromReader[GA ~func(R) A, GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](ma GA) GEA {
	return RightReader[GA, GEA](ma)
}

func MonadFromReaderIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GRIO ~func(R) GIO, GIO ~func() A, R, E, A any](a A, f func(A) GRIO) GEA {
	return F.Pipe2(
		a,
		f,
		RightReaderIO[GEA, GIOA, GRIO, GIO, R, E, A],
	)
}

func FromReaderIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GRIO ~func(R) GIO, GIO ~func() A, R, E, A any](f func(A) GRIO) func(A) GEA {
	return F.Bind2nd(MonadFromReaderIO[GEA, GIOA, GRIO, GIO, R, E, A], f)
}

func RightReaderIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GRIO ~func(R) GIO, GIO ~func() A, R, E, A any](ma GRIO) GEA {
	return eithert.RightF(
		G.MonadMap[GRIO, GEA, GIO, GIOA, R, A, ET.Either[E, A]],
		ma,
	)
}

func LeftReaderIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GRIO ~func(R) GIO, GIO ~func() E, R, E, A any](me GRIO) GEA {
	return eithert.LeftF(
		G.MonadMap[GRIO, GEA, GIO, GIOA, R, E, ET.Either[E, A]],
		me,
	)
}

func RightIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GR ~func() A, R, E, A any](ma GR) GEA {
	return F.Pipe2(ma, IOE.RightIO[GIOA, GR, E, A], FromIOEither[GEA, GIOA, R, E, A])
}

func LeftIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GR ~func() E, R, E, A any](ma GR) GEA {
	return F.Pipe2(ma, IOE.LeftIO[GIOA, GR, E, A], FromIOEither[GEA, GIOA, R, E, A])
}

func FromIO[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], GR ~func() A, R, E, A any](ma GR) GEA {
	return RightIO[GEA](ma)
}

func FromReaderEither[GA ~func(R) ET.Either[E, A], GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](ma GA) GEA {
	return F.Flow2(ma, IOE.FromEither[GIOA, E, A])
}

func Ask[GER ~func(R) GIOR, GIOR ~func() ET.Either[E, R], R, E any]() GER {
	return FR.Ask(FromReader[func(R) R, GER, GIOR, R, E, R])()
}

func Asks[GA ~func(R) A, GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](r GA) GEA {
	return FR.Asks(FromReader[GA, GEA, GIOA, R, E, A])(r)
}

func FromOption[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](onNone func() E) func(O.Option[A]) GEA {
	return FE.FromOption(FromEither[GEA, GIOA, R, E, A], onNone)
}

func FromPredicate[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](pred func(A) bool, onFalse func(A) E) func(A) GEA {
	return FE.FromPredicate(FromEither[GEA, GIOA, R, E, A], pred, onFalse)
}

func Fold[GB ~func(R) GIOB, GEA ~func(R) GIOA, GIOB ~func() B, GIOA ~func() ET.Either[E, A], R, E, A, B any](onLeft func(E) GB, onRight func(A) GB) func(GEA) GB {
	return eithert.MatchE(G.MonadChain[GEA, GB, GIOA, GIOB, R, ET.Either[E, A], B], onLeft, onRight)
}

func GetOrElse[GA ~func(R) GIOB, GEA ~func(R) GIOA, GIOB ~func() A, GIOA ~func() ET.Either[E, A], R, E, A any](onLeft func(E) GA) func(GEA) GA {
	return eithert.GetOrElse(G.MonadChain[GEA, GA, GIOA, GIOB, R, ET.Either[E, A], A], G.Of[GA, GIOB, R, A], onLeft)
}

func OrElse[GEA1 ~func(R) GIOA1, GEA2 ~func(R) GIOA2, GIOA1 ~func() ET.Either[E1, A], GIOA2 ~func() ET.Either[E2, A], R, E1, A, E2 any](onLeft func(E1) GEA2) func(GEA1) GEA2 {
	return eithert.OrElse(G.MonadChain[GEA1, GEA2, GIOA1, GIOA2, R, ET.Either[E1, A], ET.Either[E2, A]], G.Of[GEA2, GIOA2, R, ET.Either[E2, A]], onLeft)
}

func OrLeft[GEA1 ~func(R) GIOA1, GE2 ~func(R) GIOE2, GEA2 ~func(R) GIOA2, GIOA1 ~func() ET.Either[E1, A], GIOE2 ~func() E2, GIOA2 ~func() ET.Either[E2, A], E1, R, E2, A any](onLeft func(E1) GE2) func(GEA1) GEA2 {
	return eithert.OrLeft(
		G.MonadChain[GEA1, GEA2, GIOA1, GIOA2, R, ET.Either[E1, A], ET.Either[E2, A]],
		G.MonadMap[GE2, GEA2, GIOE2, GIOA2, R, E2, ET.Either[E2, A]],
		G.Of[GEA2, GIOA2, R, ET.Either[E2, A]],
		onLeft,
	)
}

func MonadBiMap[GA ~func(R) GE1A, GB ~func(R) GE2B, GE1A ~func() ET.Either[E1, A], GE2B ~func() ET.Either[E2, B], R, E1, E2, A, B any](fa GA, f func(E1) E2, g func(A) B) GB {
	return eithert.MonadBiMap(G.MonadMap[GA, GB, GE1A, GE2B, R, ET.Either[E1, A], ET.Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[GA ~func(R) GE1A, GB ~func(R) GE2B, GE1A ~func() ET.Either[E1, A], GE2B ~func() ET.Either[E2, B], R, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(GA) GB {
	return eithert.BiMap(G.MonadMap[GA, GB, GE1A, GE2B, R, ET.Either[E1, A], ET.Either[E2, B]], f, g)
}

// Swap changes the order of type parameters
func Swap[GREA ~func(R) GEA, GRAE ~func(R) GAE, GEA ~func() ET.Either[E, A], GAE ~func() ET.Either[A, E], R, E, A any](val GREA) GRAE {
	return RD.MonadMap[GREA, GRAE, R, GEA, GAE](val, IOE.Swap[GEA, GAE])
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[GEA ~func(R) GA, GA ~func() ET.Either[E, A], R, E, A any](gen func() GEA) GEA {
	return G.Defer[GEA](gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[GEA ~func(R) GA, GA ~func() ET.Either[E, A], R, E, A any](f func(R) func() (A, error), onThrow func(error) E) GEA {
	return func(r R) GA {
		return IOE.TryCatch[GA](f(r), onThrow)
	}
}

// Memoize computes the value of the provided monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[
	GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](rdr GEA) GEA {
	return G.Memoize[GEA](rdr)
}

func MonadFlap[GREAB ~func(R) GEAB, GREB ~func(R) GEB, GEAB ~func() ET.Either[E, func(A) B], GEB ~func() ET.Either[E, B], R, E, B, A any](fab GREAB, a A) GREB {
	return FC.MonadFlap(MonadMap[GREAB, GREB], fab, a)
}

func Flap[GREAB ~func(R) GEAB, GREB ~func(R) GEB, GEAB ~func() ET.Either[E, func(A) B], GEB ~func() ET.Either[E, B], R, E, B, A any](a A) func(GREAB) GREB {
	return FC.Flap(MonadMap[GREAB, GREB], a)
}
