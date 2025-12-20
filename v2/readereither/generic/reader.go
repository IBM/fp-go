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

package generic

import (
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/eithert"
	FE "github.com/IBM/fp-go/v2/internal/fromeither"
	FR "github.com/IBM/fp-go/v2/internal/fromreader"
	FC "github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader/generic"
)

func MakeReaderEither[GEA ~func(E) ET.Either[L, A], L, E, A any](f func(E) ET.Either[L, A]) GEA {
	return f
}

func FromEither[GEA ~func(E) ET.Either[L, A], L, E, A any](e ET.Either[L, A]) GEA {
	return R.Of[GEA](e)
}

func RightReader[GA ~func(E) A, GEA ~func(E) ET.Either[L, A], L, E, A any](r GA) GEA {
	return eithert.RightF(R.MonadMap[GA, GEA, E, A, ET.Either[L, A]], r)
}

func LeftReader[GL ~func(E) L, GEA ~func(E) ET.Either[L, A], L, E, A any](l GL) GEA {
	return eithert.LeftF(R.MonadMap[GL, GEA, E, L, ET.Either[L, A]], l)
}

func Left[GEA ~func(E) ET.Either[L, A], L, E, A any](l L) GEA {
	return eithert.Left(R.Of[GEA, E, ET.Either[L, A]], l)
}

func Right[GEA ~func(E) ET.Either[L, A], L, E, A any](r A) GEA {
	return eithert.Right(R.Of[GEA, E, ET.Either[L, A]], r)
}

func FromReader[GA ~func(E) A, GEA ~func(E) ET.Either[L, A], L, E, A any](r GA) GEA {
	return RightReader[GA, GEA](r)
}

func MonadMap[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](fa GEA, f func(A) B) GEB {
	return readert.MonadMap[GEA, GEB](ET.MonadMap[L, A, B], fa, f)
}

func Map[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](f func(A) B) func(GEA) GEB {
	return readert.Map[GEA, GEB](ET.Map[L, A, B], f)
}

func MonadChain[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](ma GEA, f func(A) GEB) GEB {
	return readert.MonadChain(ET.MonadChain[L, A, B], ma, f)
}

func Chain[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](f func(A) GEB) func(GEA) GEB {
	return F.Bind2nd(MonadChain[GEA, GEB, L, E, A, B], f)
}

func MonadChainReaderK[
	GEA ~func(E) ET.Either[L, A],
	GEB ~func(E) ET.Either[L, B],
	GB ~func(E) B,
	L, E, A, B any](ma GEA, f func(A) GB) GEB {

	return MonadChain(ma, F.Flow2(f, FromReader[GB, GEB, L, E, B]))
}

func ChainReaderK[
	GEA ~func(E) ET.Either[L, A],
	GEB ~func(E) ET.Either[L, B],
	GB ~func(E) B,
	L, E, A, B any](f func(A) GB) func(GEA) GEB {
	return Chain[GEA](F.Flow2(f, FromReader[GB, GEB, L, E, B]))
}

func Of[GEA ~func(E) ET.Either[L, A], L, E, A any](a A) GEA {
	return readert.MonadOf[GEA](ET.Of[L, A], a)
}

func MonadAp[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], GEFAB ~func(E) ET.Either[L, func(A) B], L, E, A, B any](fab GEFAB, fa GEA) GEB {
	return readert.MonadAp[GEA, GEB, GEFAB, E, A](ET.MonadAp[B, L, A], fab, fa)
}

func Ap[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], GEFAB ~func(E) ET.Either[L, func(A) B], L, E, A, B any](fa GEA) func(GEFAB) GEB {
	return F.Bind2nd(MonadAp[GEA, GEB, GEFAB, L, E, A, B], fa)
}

func FromPredicate[GEA ~func(E) ET.Either[L, A], L, E, A any](pred func(A) bool, onFalse func(A) L) func(A) GEA {
	return FE.FromPredicate(FromEither[GEA, L, E, A], pred, onFalse)
}

func Fold[GEA ~func(E) ET.Either[L, A], GB ~func(E) B, E, L, A, B any](onLeft func(L) GB, onRight func(A) GB) func(GEA) GB {
	return eithert.MatchE(R.MonadChain[GEA, GB, E, ET.Either[L, A], B], onLeft, onRight)
}

func GetOrElse[GEA ~func(E) ET.Either[L, A], GA ~func(E) A, E, L, A any](onLeft func(L) GA) func(GEA) GA {
	return eithert.GetOrElse(R.MonadChain[GEA, GA, E, ET.Either[L, A], A], R.Of[GA, E, A], onLeft)
}

func OrElse[GEA1 ~func(E) ET.Either[L1, A], GEA2 ~func(E) ET.Either[L2, A], E, L1, A, L2 any](onLeft func(L1) GEA2) func(GEA1) GEA2 {
	return eithert.OrElse(R.MonadChain[GEA1, GEA2, E, ET.Either[L1, A], ET.Either[L2, A]], R.Of[GEA2, E, ET.Either[L2, A]], onLeft)
}

func OrLeft[GEA1 ~func(E) ET.Either[L1, A], GEA2 ~func(E) ET.Either[L2, A], GE2 ~func(E) L2, L1, E, L2, A any](onLeft func(L1) GE2) func(GEA1) GEA2 {
	return eithert.OrLeft(
		R.MonadChain[GEA1, GEA2, E, ET.Either[L1, A], ET.Either[L2, A]],
		R.MonadMap[GE2, GEA2, E, L2, ET.Either[L2, A]],
		R.Of[GEA2, E, ET.Either[L2, A]],
		onLeft,
	)
}

func Ask[GEE ~func(E) ET.Either[L, E], E, L any]() GEE {
	return FR.Ask(FromReader[func(E) E, GEE, L, E, E])()
}

func Asks[GA ~func(E) A, GEA ~func(E) ET.Either[L, A], E, L, A any](r GA) GEA {
	return FR.Asks(FromReader[GA, GEA, L, E, A])(r)
}

func MonadChainEitherK[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](ma GEA, f func(A) ET.Either[L, B]) GEB {
	return FE.MonadChainEitherK(
		MonadChain[GEA, GEB, L, E, A, B],
		FromEither[GEB, L, E, B],
		ma,
		f,
	)
}

func ChainEitherK[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](f func(A) ET.Either[L, B]) func(ma GEA) GEB {
	return F.Bind2nd(MonadChainEitherK[GEA, GEB, L, E, A, B], f)
}

func ChainOptionK[GEA ~func(E) ET.Either[L, A], GEB ~func(E) ET.Either[L, B], L, E, A, B any](onNone func() L) func(func(A) O.Option[B]) func(GEA) GEB {
	return FE.ChainOptionK(MonadChain[GEA, GEB, L, E, A, B], FromEither[GEB, L, E, B], onNone)
}

func Flatten[GEA ~func(E) ET.Either[L, A], GGA ~func(E) ET.Either[L, GEA], L, E, A any](mma GGA) GEA {
	return MonadChain(mma, F.Identity[GEA])
}

func MonadBiMap[GA ~func(E) ET.Either[E1, A], GB ~func(E) ET.Either[E2, B], E, E1, E2, A, B any](fa GA, f func(E1) E2, g func(A) B) GB {
	return eithert.MonadBiMap(R.MonadMap[GA, GB, E, ET.Either[E1, A], ET.Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[GA ~func(E) ET.Either[E1, A], GB ~func(E) ET.Either[E2, B], E, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(GA) GB {
	return eithert.BiMap(R.Map[GA, GB, E, ET.Either[E1, A], ET.Either[E2, B]], f, g)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[GA1 ~func(R1) ET.Either[E, A], GA2 ~func(R2) ET.Either[E, A], R2, R1, E, A any](f func(R2) R1) func(GA1) GA2 {
	return R.Local[GA1, GA2](f)
}

func MonadFlap[GEFAB ~func(E) ET.Either[L, func(A) B], GEB ~func(E) ET.Either[L, B], L, E, A, B any](fab GEFAB, a A) GEB {
	return FC.MonadFlap(MonadMap[GEFAB, GEB], fab, a)
}

func Flap[GEFAB ~func(E) ET.Either[L, func(A) B], GEB ~func(E) ET.Either[L, B], L, E, A, B any](a A) func(GEFAB) GEB {
	return FC.Flap(Map[GEFAB, GEB], a)
}

func MonadMapLeft[GA1 ~func(C) ET.Either[E1, A], GA2 ~func(C) ET.Either[E2, A], C, E1, E2, A any](fa GA1, f func(E1) E2) GA2 {
	return eithert.MonadMapLeft(R.MonadMap[GA1, GA2, C, ET.Either[E1, A], ET.Either[E2, A]], fa, f)
}

// MapLeft applies a mapping function to the error channel
func MapLeft[GA1 ~func(C) ET.Either[E1, A], GA2 ~func(C) ET.Either[E2, A], C, E1, E2, A any](f func(E1) E2) func(GA1) GA2 {
	return F.Bind2nd(MonadMapLeft[GA1, GA2, C, E1, E2, A], f)
}
