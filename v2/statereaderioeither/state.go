// Copyright (c) 2024 IBM Corp.
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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/statet"
	"github.com/IBM/fp-go/v2/readerioeither"
)

func Left[S, R, A, E any](e E) StateReaderIOEither[S, R, E, A] {
	return function.Constant1[S](readerioeither.Left[R, Pair[S, A]](e))
}

func Right[S, R, E, A any](a A) StateReaderIOEither[S, R, E, A] {
	return statet.Of[StateReaderIOEither[S, R, E, A]](readerioeither.Of[R, E, Pair[S, A]], a)
}

func Of[S, R, E, A any](a A) StateReaderIOEither[S, R, E, A] {
	return Right[S, R, E](a)
}

func MonadMap[S, R, E, A, B any](fa StateReaderIOEither[S, R, E, A], f func(A) B) StateReaderIOEither[S, R, E, B] {
	return statet.MonadMap[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]](
		readerioeither.MonadMap[R, E, Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

func Map[S, R, E, A, B any](f func(A) B) Operator[S, R, E, A, B] {
	return statet.Map[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]](
		readerioeither.Map[R, E, Pair[S, A], Pair[S, B]],
		f,
	)
}

func MonadChain[S, R, E, A, B any](fa StateReaderIOEither[S, R, E, A], f func(A) StateReaderIOEither[S, R, E, B]) StateReaderIOEither[S, R, E, B] {
	return statet.MonadChain(
		readerioeither.MonadChain[R, E, Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

func Chain[S, R, E, A, B any](f func(A) StateReaderIOEither[S, R, E, B]) Operator[S, R, E, A, B] {
	return statet.Chain[StateReaderIOEither[S, R, E, A]](
		readerioeither.Chain[R, E, Pair[S, A], Pair[S, B]],
		f,
	)
}

func MonadAp[B, S, R, E, A any](fab StateReaderIOEither[S, R, E, func(A) B], fa StateReaderIOEither[S, R, E, A]) StateReaderIOEither[S, R, E, B] {
	return statet.MonadAp[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]](
		readerioeither.MonadMap[R, E, Pair[S, A], Pair[S, B]],
		readerioeither.MonadChain[R, E, Pair[S, func(A) B], Pair[S, B]],
		fab,
		fa,
	)
}

func Ap[B, S, R, E, A any](fa StateReaderIOEither[S, R, E, A]) Operator[S, R, E, func(A) B, B] {
	return statet.Ap[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]](
		readerioeither.Map[R, E, Pair[S, A], Pair[S, B]],
		readerioeither.Chain[R, E, Pair[S, func(A) B], Pair[S, B]],
		fa,
	)
}

func FromReaderIOEither[S, R, E, A any](fa readerioeither.ReaderIOEither[R, E, A]) StateReaderIOEither[S, R, E, A] {
	return statet.FromF[StateReaderIOEither[S, R, E, A]](
		readerioeither.MonadMap[R, E, A],
		fa,
	)
}

func FromReaderEither[S, R, E, A any](fa ReaderEither[R, E, A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromReaderEither(fa))
}

func FromIOEither[S, R, E, A any](fa IOEither[E, A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromIOEither[R](fa))
}

func FromState[R, E, S, A any](sa State[S, A]) StateReaderIOEither[S, R, E, A] {
	return statet.FromState[StateReaderIOEither[S, R, E, A]](readerioeither.Of[R, E, Pair[S, A]], sa)
}

func FromIO[S, R, E, A any](fa IO[A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromIO[R, E](fa))
}

func FromReader[S, E, R, A any](fa Reader[R, A]) StateReaderIOEither[S, R, E, A] {
	return FromReaderIOEither[S](readerioeither.FromReader[E](fa))
}

func FromEither[S, R, E, A any](ma Either[E, A]) StateReaderIOEither[S, R, E, A] {
	return either.MonadFold(ma, Left[S, R, A, E], Right[S, R, E, A])
}

// Combinators

func Local[S, E, A, B, R1, R2 any](f func(R2) R1) func(StateReaderIOEither[S, R1, E, A]) StateReaderIOEither[S, R2, E, A] {
	return func(ma StateReaderIOEither[S, R1, E, A]) StateReaderIOEither[S, R2, E, A] {
		return function.Flow2(ma, readerioeither.Local[R1, R2, E, Pair[S, A]](f))
	}
}

func Asks[
	S, R, E, A any,
](f func(R) StateReaderIOEither[S, R, E, A]) StateReaderIOEither[S, R, E, A] {
	return func(s S) ReaderIOEither[R, E, Pair[S, A]] {
		return func(r R) IOEither[E, Pair[S, A]] {
			return f(r)(s)(r)
		}
	}
}

func FromEitherK[S, R, E, A, B any](f func(A) Either[E, B]) func(A) StateReaderIOEither[S, R, E, B] {
	return function.Flow2(
		f,
		FromEither[S, R, E, B],
	)
}

func FromIOK[S, R, E, A, B any](f func(A) IO[B]) func(A) StateReaderIOEither[S, R, E, B] {
	return function.Flow2(
		f,
		FromIO[S, R, E, B],
	)
}

func FromIOEitherK[
	S, R, E, A, B any,
](f func(A) IOEither[E, B]) func(A) StateReaderIOEither[S, R, E, B] {
	return function.Flow2(
		f,
		FromIOEither[S, R, E, B],
	)
}

func FromReaderIOEitherK[S, R, E, A, B any](f func(A) readerioeither.ReaderIOEither[R, E, B]) func(A) StateReaderIOEither[S, R, E, B] {
	return function.Flow2(
		f,
		FromReaderIOEither[S, R, E, B],
	)
}

func MonadChainReaderIOEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f func(A) readerioeither.ReaderIOEither[R, E, B]) StateReaderIOEither[S, R, E, B] {
	return MonadChain(ma, FromReaderIOEitherK[S](f))
}

func ChainReaderIOEitherK[S, R, E, A, B any](f func(A) readerioeither.ReaderIOEither[R, E, B]) Operator[S, R, E, A, B] {
	return Chain(FromReaderIOEitherK[S](f))
}

func MonadChainIOEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f func(A) IOEither[E, B]) StateReaderIOEither[S, R, E, B] {
	return MonadChain(ma, FromIOEitherK[S, R](f))
}

func ChainIOEitherK[S, R, E, A, B any](f func(A) IOEither[E, B]) Operator[S, R, E, A, B] {
	return Chain(FromIOEitherK[S, R](f))
}

func MonadChainEitherK[S, R, E, A, B any](ma StateReaderIOEither[S, R, E, A], f func(A) Either[E, B]) StateReaderIOEither[S, R, E, B] {
	return MonadChain(ma, FromEitherK[S, R](f))
}

func ChainEitherK[S, R, E, A, B any](f func(A) Either[E, B]) Operator[S, R, E, A, B] {
	return Chain(FromEitherK[S, R](f))
}
