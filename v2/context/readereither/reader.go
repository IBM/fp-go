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
	"context"

	"github.com/IBM/fp-go/v2/readereither"
)

func FromEither[A any](e Either[A]) ReaderEither[A] {
	return readereither.FromEither[context.Context](e)
}

func Left[A any](l error) ReaderEither[A] {
	return readereither.Left[context.Context, A](l)
}

func Right[A any](r A) ReaderEither[A] {
	return readereither.Right[context.Context, error](r)
}

func MonadMap[A, B any](fa ReaderEither[A], f func(A) B) ReaderEither[B] {
	return readereither.MonadMap(fa, f)
}

func Map[A, B any](f func(A) B) Operator[A, B] {
	return readereither.Map[context.Context, error](f)
}

func MonadChain[A, B any](ma ReaderEither[A], f func(A) ReaderEither[B]) ReaderEither[B] {
	return readereither.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderEither[B]) Operator[A, B] {
	return readereither.Chain(f)
}

func Of[A any](a A) ReaderEither[A] {
	return readereither.Of[context.Context, error](a)
}

func MonadAp[A, B any](fab ReaderEither[func(A) B], fa ReaderEither[A]) ReaderEither[B] {
	return readereither.MonadAp(fab, fa)
}

func Ap[A, B any](fa ReaderEither[A]) func(ReaderEither[func(A) B]) ReaderEither[B] {
	return readereither.Ap[B](fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) func(A) ReaderEither[A] {
	return readereither.FromPredicate[context.Context](pred, onFalse)
}

func OrElse[A any](onLeft func(error) ReaderEither[A]) func(ReaderEither[A]) ReaderEither[A] {
	return readereither.OrElse(onLeft)
}

func Ask() ReaderEither[context.Context] {
	return readereither.Ask[context.Context, error]()
}

func MonadChainEitherK[A, B any](ma ReaderEither[A], f func(A) Either[B]) ReaderEither[B] {
	return readereither.MonadChainEitherK(ma, f)
}

func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderEither[A]) ReaderEither[B] {
	return readereither.ChainEitherK[context.Context](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) Option[B]) Operator[A, B] {
	return readereither.ChainOptionK[context.Context, A, B](onNone)
}

func MonadFlap[B, A any](fab ReaderEither[func(A) B], a A) ReaderEither[B] {
	return readereither.MonadFlap(fab, a)
}

func Flap[B, A any](a A) Operator[func(A) B, B] {
	return readereither.Flap[context.Context, error, B](a)
}
