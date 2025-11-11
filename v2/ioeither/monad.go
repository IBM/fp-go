// Copyright (c) 2024 - 2025 IBM Corp.
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

package ioeither

import (
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type (
	ioEitherPointed[E, A any] struct{}

	ioEitherMonad[E, A, B any] struct{}

	ioEitherFunctor[E, A, B any] struct{}
)

func (o *ioEitherPointed[E, A]) Of(a A) IOEither[E, A] {
	return Of[E](a)
}

func (o *ioEitherMonad[E, A, B]) Of(a A) IOEither[E, A] {
	return Of[E](a)
}

func (o *ioEitherMonad[E, A, B]) Map(f func(A) B) Operator[E, A, B] {
	return Map[E](f)
}

func (o *ioEitherMonad[E, A, B]) Chain(f Kleisli[E, A, B]) Operator[E, A, B] {
	return Chain(f)
}

func (o *ioEitherMonad[E, A, B]) Ap(fa IOEither[E, A]) Operator[E, func(A) B, B] {
	return Ap[B](fa)
}

func (o *ioEitherFunctor[E, A, B]) Map(f func(A) B) Operator[E, A, B] {
	return Map[E](f)
}

// Pointed implements the pointed operations for [IOEither]
func Pointed[E, A any]() pointed.Pointed[A, IOEither[E, A]] {
	return &ioEitherPointed[E, A]{}
}

// Functor implements the monadic operations for [IOEither]
func Functor[E, A, B any]() functor.Functor[A, B, IOEither[E, A], IOEither[E, B]] {
	return &ioEitherFunctor[E, A, B]{}
}

// Monad implements the monadic operations for [IOEither]
func Monad[E, A, B any]() monad.Monad[A, B, IOEither[E, A], IOEither[E, B], IOEither[E, func(A) B]] {
	return &ioEitherMonad[E, A, B]{}
}
