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

package readereither

import (
	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
)

func FromEither[E, L, A any](e Either[L, A]) ReaderEither[E, L, A] {
	return reader.Of[E](e)
}

func RightReader[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A] {
	return eithert.RightF(reader.MonadMap[E, A, Either[L, A]], r)
}

func LeftReader[A, E, L any](l Reader[E, L]) ReaderEither[E, L, A] {
	return eithert.LeftF(reader.MonadMap[E, L, Either[L, A]], l)
}

func Left[E, A, L any](l L) ReaderEither[E, L, A] {
	return eithert.Left(reader.Of[E, Either[L, A]], l)
}

func Right[E, L, A any](r A) ReaderEither[E, L, A] {
	return eithert.Right(reader.Of[E, Either[L, A]], r)
}

func FromReader[E, L, A any](r Reader[E, A]) ReaderEither[E, L, A] {
	return RightReader[L](r)
}

func MonadMap[E, L, A, B any](fa ReaderEither[E, L, A], f func(A) B) ReaderEither[E, L, B] {
	return readert.MonadMap[ReaderEither[E, L, A], ReaderEither[E, L, B]](ET.MonadMap[L, A, B], fa, f)
}

func Map[E, L, A, B any](f func(A) B) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return readert.Map[ReaderEither[E, L, A], ReaderEither[E, L, B]](ET.Map[L, A, B], f)
}

func MonadChain[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) ReaderEither[E, L, B]) ReaderEither[E, L, B] {
	return readert.MonadChain[ReaderEither[E, L, A], ReaderEither[E, L, B]](ET.MonadChain[L, A, B], ma, f)
}

func Chain[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return readert.Chain[ReaderEither[E, L, A], ReaderEither[E, L, B]](ET.Chain[L, A, B], f)
}

func Of[E, L, A any](a A) ReaderEither[E, L, A] {
	return readert.MonadOf[ReaderEither[E, L, A]](ET.Of[L, A], a)
}

func MonadAp[E, L, A, B any](fab ReaderEither[E, L, func(A) B], fa ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return readert.MonadAp[ReaderEither[E, L, A], ReaderEither[E, L, B], ReaderEither[E, L, func(A) B], E, A](ET.MonadAp[B, L, A], fab, fa)
}

func Ap[B, E, L, A any](fa ReaderEither[E, L, A]) func(ReaderEither[E, L, func(A) B]) ReaderEither[E, L, B] {
	return readert.Ap[ReaderEither[E, L, A], ReaderEither[E, L, B], ReaderEither[E, L, func(A) B], E, A](ET.Ap[B, L, A], fa)
}

func FromPredicate[E, L, A any](pred func(A) bool, onFalse func(A) L) func(A) ReaderEither[E, L, A] {
	return fromeither.FromPredicate(FromEither[E, L, A], pred, onFalse)
}

func Fold[E, L, A, B any](onLeft func(L) Reader[E, B], onRight func(A) Reader[E, B]) func(ReaderEither[E, L, A]) Reader[E, B] {
	return eithert.MatchE(reader.MonadChain[E, Either[L, A], B], onLeft, onRight)
}

func GetOrElse[E, L, A any](onLeft func(L) Reader[E, A]) func(ReaderEither[E, L, A]) Reader[E, A] {
	return eithert.GetOrElse(reader.MonadChain[E, Either[L, A], A], reader.Of[E, A], onLeft)
}

func OrElse[E, L1, A, L2 any](onLeft func(L1) ReaderEither[E, L2, A]) func(ReaderEither[E, L1, A]) ReaderEither[E, L2, A] {
	return eithert.OrElse(reader.MonadChain[E, Either[L1, A], Either[L2, A]], reader.Of[E, Either[L2, A]], onLeft)
}

func OrLeft[A, L1, E, L2 any](onLeft func(L1) Reader[E, L2]) func(ReaderEither[E, L1, A]) ReaderEither[E, L2, A] {
	return eithert.OrLeft(
		reader.MonadChain[E, Either[L1, A], Either[L2, A]],
		reader.MonadMap[E, L2, Either[L2, A]],
		reader.Of[E, Either[L2, A]],
		onLeft,
	)
}

func Ask[E, L any]() ReaderEither[E, L, E] {
	return fromreader.Ask(FromReader[E, L, E])()
}

func Asks[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A] {
	return fromreader.Asks(FromReader[E, L, A])(r)
}

func MonadChainEitherK[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) Either[L, B]) ReaderEither[E, L, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[E, L, A, B],
		FromEither[E, L, B],
		ma,
		f,
	)
}

func ChainEitherK[E, L, A, B any](f func(A) Either[L, B]) func(ma ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return fromeither.ChainEitherK(
		Chain[E, L, A, B],
		FromEither[E, L, B],
		f,
	)
}

func ChainOptionK[E, A, B, L any](onNone func() L) func(func(A) Option[B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return fromeither.ChainOptionK(MonadChain[E, L, A, B], FromEither[E, L, B], onNone)
}

func Flatten[E, L, A any](mma ReaderEither[E, L, ReaderEither[E, L, A]]) ReaderEither[E, L, A] {
	return MonadChain(mma, function.Identity[ReaderEither[E, L, A]])
}

func MonadBiMap[E, E1, E2, A, B any](fa ReaderEither[E, E1, A], f func(E1) E2, g func(A) B) ReaderEither[E, E2, B] {
	return eithert.MonadBiMap(reader.MonadMap[E, Either[E1, A], Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderEither[E, E1, A]) ReaderEither[E, E2, B] {
	return eithert.BiMap(reader.Map[E, Either[E1, A], Either[E2, B]], f, g)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[E, A, R2, R1 any](f func(R2) R1) func(ReaderEither[R1, E, A]) ReaderEither[R2, E, A] {
	return reader.Local[R2, R1, Either[E, A]](f)
}

// Read applies a context to a reader to obtain its value
func Read[E1, A, E any](e E) func(ReaderEither[E, E1, A]) Either[E1, A] {
	return reader.Read[E, Either[E1, A]](e)
}

func MonadFlap[L, E, A, B any](fab ReaderEither[L, E, func(A) B], a A) ReaderEither[L, E, B] {
	return functor.MonadFlap(MonadMap[L, E, func(A) B, B], fab, a)
}

func Flap[L, E, B, A any](a A) func(ReaderEither[L, E, func(A) B]) ReaderEither[L, E, B] {
	return functor.Flap(Map[L, E, func(A) B, B], a)
}

func MonadMapLeft[C, E1, E2, A any](fa ReaderEither[C, E1, A], f func(E1) E2) ReaderEither[C, E2, A] {
	return eithert.MonadMapLeft(reader.MonadMap[C, Either[E1, A], Either[E2, A]], fa, f)
}

// MapLeft applies a mapping function to the error channel
func MapLeft[C, E1, E2, A any](f func(E1) E2) func(ReaderEither[C, E1, A]) ReaderEither[C, E2, A] {
	return eithert.MapLeft(reader.Map[C, Either[E1, A], Either[E2, A]], f)
}
