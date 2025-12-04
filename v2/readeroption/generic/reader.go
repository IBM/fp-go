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
	F "github.com/IBM/fp-go/v2/function"
	FO "github.com/IBM/fp-go/v2/internal/fromoption"
	FR "github.com/IBM/fp-go/v2/internal/fromreader"
	FC "github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/optiont"
	"github.com/IBM/fp-go/v2/internal/readert"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader/generic"
)

//go:inline
func MakeReaderOption[GEA ~func(E) O.Option[A], E, A any](f func(E) O.Option[A]) GEA {
	return f
}

//go:inline
func FromOption[GEA ~func(E) O.Option[A], E, A any](e O.Option[A]) GEA {
	return R.Of[GEA](e)
}

func SomeReader[GA ~func(E) A, GEA ~func(E) O.Option[A], E, A any](r GA) GEA {
	return optiont.SomeF(R.MonadMap[GA, GEA, E, A, O.Option[A]], r)
}

func Some[GEA ~func(E) O.Option[A], E, A any](r A) GEA {
	return optiont.Of(R.Of[GEA, E, O.Option[A]], r)
}

//go:inline
func FromReader[GA ~func(E) A, GEA ~func(E) O.Option[A], E, A any](r GA) GEA {
	return SomeReader[GA, GEA](r)
}

func MonadMap[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], E, A, B any](fa GEA, f func(A) B) GEB {
	return readert.MonadMap[GEA, GEB](O.MonadMap[A, B], fa, f)
}

func Map[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], E, A, B any](f func(A) B) func(GEA) GEB {
	return readert.Map[GEA, GEB](O.Map[A, B], f)
}

func MonadChain[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], E, A, B any](ma GEA, f func(A) GEB) GEB {
	return readert.MonadChain(O.MonadChain[A, B], ma, f)
}

func Chain[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], E, A, B any](f func(A) GEB) func(GEA) GEB {
	return F.Bind2nd(MonadChain[GEA, GEB, E, A, B], f)
}

func Of[GEA ~func(E) O.Option[A], E, A any](a A) GEA {
	return readert.MonadOf[GEA](O.Of[A], a)
}

func MonadAp[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], GEFAB ~func(E) O.Option[func(A) B], E, A, B any](fab GEFAB, fa GEA) GEB {
	return readert.MonadAp[GEA, GEB, GEFAB, E, A](O.MonadAp[B, A], fab, fa)
}

func Ap[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], GEFAB ~func(E) O.Option[func(A) B], E, A, B any](fa GEA) func(GEFAB) GEB {
	return F.Bind2nd(MonadAp[GEA, GEB, GEFAB, E, A, B], fa)
}

func FromPredicate[GEA ~func(E) O.Option[A], E, A any](pred func(A) bool) func(A) GEA {
	return FO.FromPredicate(FromOption[GEA, E, A], pred)
}

func Fold[GEA ~func(E) O.Option[A], GB ~func(E) B, E, A, B any](onNone func() GB, onRight func(A) GB) func(GEA) GB {
	return optiont.MatchE(R.Chain[GEA, GB, E, O.Option[A], B], onNone, onRight)
}

func GetOrElse[GEA ~func(E) O.Option[A], GA ~func(E) A, E, A any](onNone func() GA) func(GEA) GA {
	return optiont.GetOrElse(R.Chain[GEA, GA, E, O.Option[A], A], onNone, R.Of[GA, E, A])
}

func Ask[GEE ~func(E) O.Option[E], E any]() GEE {
	return FR.Ask(FromReader[func(E) E, GEE, E, E])()
}

func Asks[GA ~func(E) A, GEA ~func(E) O.Option[A], E, A any](r GA) GEA {
	return FR.Asks(FromReader[GA, GEA, E, A])(r)
}

func MonadChainOptionK[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], E, A, B any](ma GEA, f func(A) O.Option[B]) GEB {
	return FO.MonadChainOptionK(
		MonadChain[GEA, GEB, E, A, B],
		FromOption[GEB, E, B],
		ma,
		f,
	)
}

func ChainOptionK[GEA ~func(E) O.Option[A], GEB ~func(E) O.Option[B], E, A, B any](f func(A) O.Option[B]) func(ma GEA) GEB {
	return F.Bind2nd(MonadChainOptionK[GEA, GEB, E, A, B], f)
}

func Flatten[GEA ~func(E) O.Option[A], GGA ~func(E) O.Option[GEA], E, A any](mma GGA) GEA {
	return MonadChain(mma, F.Identity[GEA])
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[GA1 ~func(R1) O.Option[A], GA2 ~func(R2) O.Option[A], R2, R1, E, A any](f func(R2) R1) func(GA1) GA2 {
	return R.Local[GA1, GA2](f)
}

func MonadFlap[GEFAB ~func(E) O.Option[func(A) B], GEB ~func(E) O.Option[B], E, A, B any](fab GEFAB, a A) GEB {
	return FC.MonadFlap(MonadMap[GEFAB, GEB], fab, a)
}

func Flap[GEFAB ~func(E) O.Option[func(A) B], GEB ~func(E) O.Option[B], E, A, B any](a A) func(GEFAB) GEB {
	return FC.Flap(Map[GEFAB, GEB], a)
}
