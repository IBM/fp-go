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

package readereither

import (
	ET "github.com/IBM/fp-go/v2/either"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
	G "github.com/IBM/fp-go/v2/readereither/generic"
)

type ReaderEither[E, L, A any] R.Reader[E, ET.Either[L, A]]

func MakeReaderEither[L, E, A any](f func(E) ET.Either[L, A]) ReaderEither[E, L, A] {
	return G.MakeReaderEither[ReaderEither[E, L, A]](f)
}

func FromEither[E, L, A any](e ET.Either[L, A]) ReaderEither[E, L, A] {
	return G.FromEither[ReaderEither[E, L, A]](e)
}

func RightReader[L, E, A any](r R.Reader[E, A]) ReaderEither[E, L, A] {
	return G.RightReader[R.Reader[E, A], ReaderEither[E, L, A]](r)
}

func LeftReader[A, E, L any](l R.Reader[E, L]) ReaderEither[E, L, A] {
	return G.LeftReader[R.Reader[E, L], ReaderEither[E, L, A]](l)
}

func Left[E, A, L any](l L) ReaderEither[E, L, A] {
	return G.Left[ReaderEither[E, L, A]](l)
}

func Right[E, L, A any](r A) ReaderEither[E, L, A] {
	return G.Right[ReaderEither[E, L, A]](r)
}

func FromReader[E, L, A any](r R.Reader[E, A]) ReaderEither[E, L, A] {
	return G.FromReader[R.Reader[E, A], ReaderEither[E, L, A]](r)
}

func MonadMap[E, L, A, B any](fa ReaderEither[E, L, A], f func(A) B) ReaderEither[E, L, B] {
	return G.MonadMap[ReaderEither[E, L, A], ReaderEither[E, L, B]](fa, f)
}

func Map[E, L, A, B any](f func(A) B) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return G.Map[ReaderEither[E, L, A], ReaderEither[E, L, B]](f)
}

func MonadChain[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) ReaderEither[E, L, B]) ReaderEither[E, L, B] {
	return G.MonadChain[ReaderEither[E, L, A], ReaderEither[E, L, B]](ma, f)
}

func Chain[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return G.Chain[ReaderEither[E, L, A], ReaderEither[E, L, B]](f)
}

func Of[E, L, A any](a A) ReaderEither[E, L, A] {
	return G.Of[ReaderEither[E, L, A]](a)
}

func MonadAp[E, L, A, B any](fab ReaderEither[E, L, func(A) B], fa ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return G.MonadAp[ReaderEither[E, L, A], ReaderEither[E, L, B], ReaderEither[E, L, func(A) B]](fab, fa)
}

func Ap[B, E, L, A any](fa ReaderEither[E, L, A]) func(ReaderEither[E, L, func(A) B]) ReaderEither[E, L, B] {
	return G.Ap[ReaderEither[E, L, A], ReaderEither[E, L, B], ReaderEither[E, L, func(A) B]](fa)
}

func FromPredicate[E, L, A any](pred func(A) bool, onFalse func(A) L) func(A) ReaderEither[E, L, A] {
	return G.FromPredicate[ReaderEither[E, L, A]](pred, onFalse)
}

func Fold[E, L, A, B any](onLeft func(L) R.Reader[E, B], onRight func(A) R.Reader[E, B]) func(ReaderEither[E, L, A]) R.Reader[E, B] {
	return G.Fold[ReaderEither[E, L, A]](onLeft, onRight)
}

func GetOrElse[E, L, A any](onLeft func(L) R.Reader[E, A]) func(ReaderEither[E, L, A]) R.Reader[E, A] {
	return G.GetOrElse[ReaderEither[E, L, A]](onLeft)
}

func OrElse[E, L1, A, L2 any](onLeft func(L1) ReaderEither[E, L2, A]) func(ReaderEither[E, L1, A]) ReaderEither[E, L2, A] {
	return G.OrElse[ReaderEither[E, L1, A]](onLeft)
}

func OrLeft[A, L1, E, L2 any](onLeft func(L1) R.Reader[E, L2]) func(ReaderEither[E, L1, A]) ReaderEither[E, L2, A] {
	return G.OrLeft[ReaderEither[E, L1, A], ReaderEither[E, L2, A]](onLeft)
}

func Ask[E, L any]() ReaderEither[E, L, E] {
	return G.Ask[ReaderEither[E, L, E]]()
}

func Asks[L, E, A any](r R.Reader[E, A]) ReaderEither[E, L, A] {
	return G.Asks[R.Reader[E, A], ReaderEither[E, L, A]](r)
}

func MonadChainEitherK[A, B, L, E any](ma ReaderEither[E, L, A], f func(A) ET.Either[L, B]) ReaderEither[E, L, B] {
	return G.MonadChainEitherK[ReaderEither[E, L, A], ReaderEither[E, L, B]](ma, f)
}

func ChainEitherK[A, B, L, E any](f func(A) ET.Either[L, B]) func(ma ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return G.ChainEitherK[ReaderEither[E, L, A], ReaderEither[E, L, B]](f)
}

func ChainOptionK[E, A, B, L any](onNone func() L) func(func(A) O.Option[B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return G.ChainOptionK[ReaderEither[E, L, A], ReaderEither[E, L, B]](onNone)
}

func Flatten[E, L, A any](mma ReaderEither[E, L, ReaderEither[E, L, A]]) ReaderEither[E, L, A] {
	return G.Flatten(mma)
}

func MonadBiMap[E, E1, E2, A, B any](fa ReaderEither[E, E1, A], f func(E1) E2, g func(A) B) ReaderEither[E, E2, B] {
	return G.MonadBiMap[ReaderEither[E, E1, A], ReaderEither[E, E2, B]](fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderEither[E, E1, A]) ReaderEither[E, E2, B] {
	return G.BiMap[ReaderEither[E, E1, A], ReaderEither[E, E2, B]](f, g)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[E, A, R2, R1 any](f func(R2) R1) func(ReaderEither[R1, E, A]) ReaderEither[R2, E, A] {
	return G.Local[ReaderEither[R1, E, A], ReaderEither[R2, E, A]](f)
}

// Read applies a context to a reader to obtain its value
func Read[E1, A, E any](e E) func(ReaderEither[E, E1, A]) ET.Either[E1, A] {
	return G.Read[ReaderEither[E, E1, A]](e)
}

func MonadFlap[L, E, A, B any](fab ReaderEither[L, E, func(A) B], a A) ReaderEither[L, E, B] {
	return G.MonadFlap[ReaderEither[L, E, func(A) B], ReaderEither[L, E, B]](fab, a)
}

func Flap[L, E, B, A any](a A) func(ReaderEither[L, E, func(A) B]) ReaderEither[L, E, B] {
	return G.Flap[ReaderEither[L, E, func(A) B], ReaderEither[L, E, B]](a)
}

func MonadMapLeft[C, E1, E2, A any](fa ReaderEither[C, E1, A], f func(E1) E2) ReaderEither[C, E2, A] {
	return G.MonadMapLeft[ReaderEither[C, E1, A], ReaderEither[C, E2, A]](fa, f)
}

// MapLeft applies a mapping function to the error channel
func MapLeft[C, E1, E2, A any](f func(E1) E2) func(ReaderEither[C, E1, A]) ReaderEither[C, E2, A] {
	return G.MapLeft[ReaderEither[C, E1, A], ReaderEither[C, E2, A]](f)
}
