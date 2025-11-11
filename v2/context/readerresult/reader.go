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

package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/readereither"
)

func FromEither[A any](e Either[A]) ReaderResult[A] {
	return readereither.FromEither[context.Context](e)
}

func Left[A any](l error) ReaderResult[A] {
	return readereither.Left[context.Context, A](l)
}

func Right[A any](r A) ReaderResult[A] {
	return readereither.Right[context.Context, error](r)
}

func MonadMap[A, B any](fa ReaderResult[A], f func(A) B) ReaderResult[B] {
	return readereither.MonadMap(fa, f)
}

func Map[A, B any](f func(A) B) Operator[A, B] {
	return readereither.Map[context.Context, error](f)
}

func MonadChain[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[B] {
	return readereither.MonadChain(ma, f)
}

func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return readereither.Chain(f)
}

func Of[A any](a A) ReaderResult[A] {
	return readereither.Of[context.Context, error](a)
}

func MonadAp[A, B any](fab ReaderResult[func(A) B], fa ReaderResult[A]) ReaderResult[B] {
	return readereither.MonadAp(fab, fa)
}

func Ap[A, B any](fa ReaderResult[A]) func(ReaderResult[func(A) B]) ReaderResult[B] {
	return readereither.Ap[B](fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return readereither.FromPredicate[context.Context](pred, onFalse)
}

func OrElse[A any](onLeft Kleisli[error, A]) Kleisli[ReaderResult[A], A] {
	return readereither.OrElse(onLeft)
}

func Ask() ReaderResult[context.Context] {
	return readereither.Ask[context.Context, error]()
}

func MonadChainEitherK[A, B any](ma ReaderResult[A], f func(A) Either[B]) ReaderResult[B] {
	return readereither.MonadChainEitherK(ma, f)
}

func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderResult[A]) ReaderResult[B] {
	return readereither.ChainEitherK[context.Context](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) Option[B]) Operator[A, B] {
	return readereither.ChainOptionK[context.Context, A, B](onNone)
}

func MonadFlap[B, A any](fab ReaderResult[func(A) B], a A) ReaderResult[B] {
	return readereither.MonadFlap(fab, a)
}

func Flap[B, A any](a A) Operator[func(A) B, B] {
	return readereither.Flap[context.Context, error, B](a)
}
