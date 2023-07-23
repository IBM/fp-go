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

	R "github.com/IBM/fp-go/context/reader"
	ET "github.com/IBM/fp-go/either"
	O "github.com/IBM/fp-go/option"
	RE "github.com/IBM/fp-go/readereither/generic"
)

func MakeReaderEither[A any](f func(context.Context) ET.Either[error, A]) ReaderEither[A] {
	return RE.MakeReaderEither[ReaderEither[A]](f)
}

func FromEither[A any](e ET.Either[error, A]) ReaderEither[A] {
	return RE.FromEither[ReaderEither[A]](e)
}

func RightReader[A any](r R.Reader[A]) ReaderEither[A] {
	return RE.RightReader[R.Reader[A], ReaderEither[A]](r)
}

func LeftReader[A any](l R.Reader[error]) ReaderEither[A] {
	return RE.LeftReader[R.Reader[error], ReaderEither[A]](l)
}

func Left[A any](l error) ReaderEither[A] {
	return RE.Left[ReaderEither[A]](l)
}

func Right[A any](r A) ReaderEither[A] {
	return RE.Right[ReaderEither[A]](r)
}

func FromReader[A any](r R.Reader[A]) ReaderEither[A] {
	return RE.FromReader[R.Reader[A], ReaderEither[A]](r)
}

func MonadMap[A, B any](fa ReaderEither[A], f func(A) B) ReaderEither[B] {
	return RE.MonadMap[ReaderEither[A], ReaderEither[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(ReaderEither[A]) ReaderEither[B] {
	return RE.Map[ReaderEither[A], ReaderEither[B]](f)
}

func MonadChain[A, B any](ma ReaderEither[A], f func(A) ReaderEither[B]) ReaderEither[B] {
	return RE.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderEither[B]) func(ReaderEither[A]) ReaderEither[B] {
	return RE.Chain[ReaderEither[A]](f)
}

func Of[A any](a A) ReaderEither[A] {
	return RE.Of[ReaderEither[A]](a)
}

func MonadAp[A, B any](fab ReaderEither[func(A) B], fa ReaderEither[A]) ReaderEither[B] {
	return RE.MonadAp[ReaderEither[A], ReaderEither[B]](fab, fa)
}

func Ap[A, B any](fa ReaderEither[A]) func(ReaderEither[func(A) B]) ReaderEither[B] {
	return RE.Ap[ReaderEither[A], ReaderEither[B], ReaderEither[func(A) B]](fa)
}

func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) func(A) ReaderEither[A] {
	return RE.FromPredicate[ReaderEither[A]](pred, onFalse)
}

func Fold[A, B any](onLeft func(error) R.Reader[B], onRight func(A) R.Reader[B]) func(ReaderEither[A]) R.Reader[B] {
	return RE.Fold[ReaderEither[A]](onLeft, onRight)
}

func GetOrElse[A any](onLeft func(error) R.Reader[A]) func(ReaderEither[A]) R.Reader[A] {
	return RE.GetOrElse[ReaderEither[A]](onLeft)
}

func OrElse[A any](onLeft func(error) ReaderEither[A]) func(ReaderEither[A]) ReaderEither[A] {
	return RE.OrElse[ReaderEither[A]](onLeft)
}

func OrLeft[A any](onLeft func(error) R.Reader[error]) func(ReaderEither[A]) ReaderEither[A] {
	return RE.OrLeft[ReaderEither[A], ReaderEither[A]](onLeft)
}

func Ask() ReaderEither[context.Context] {
	return RE.Ask[ReaderEither[context.Context]]()
}

func Asks[A any](r R.Reader[A]) ReaderEither[A] {
	return RE.Asks[R.Reader[A], ReaderEither[A]](r)
}

func MonadChainEitherK[A, B any](ma ReaderEither[A], f func(A) ET.Either[error, B]) ReaderEither[B] {
	return RE.MonadChainEitherK[ReaderEither[A], ReaderEither[B]](ma, f)
}

func ChainEitherK[A, B any](f func(A) ET.Either[error, B]) func(ma ReaderEither[A]) ReaderEither[B] {
	return RE.ChainEitherK[ReaderEither[A], ReaderEither[B]](f)
}

func ChainOptionK[A, B any](onNone func() error) func(func(A) O.Option[B]) func(ReaderEither[A]) ReaderEither[B] {
	return RE.ChainOptionK[ReaderEither[A], ReaderEither[B]](onNone)
}
