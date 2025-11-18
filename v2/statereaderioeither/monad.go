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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type stateReaderIOEitherPointed[
	S, R, E, A any,
] struct{}

type stateReaderIOEitherFunctor[
	S, R, E, A, B any,
] struct{}

type stateReaderIOEitherApplicative[
	S, R, E, A, B any,
] struct{}

type stateReaderIOEitherMonad[
	S, R, E, A, B any,
] struct{}

func (o *stateReaderIOEitherPointed[S, R, E, A]) Of(a A) StateReaderIOEither[S, R, E, A] {
	return Of[S, R, E](a)
}

func (o *stateReaderIOEitherMonad[S, R, E, A, B]) Of(a A) StateReaderIOEither[S, R, E, A] {
	return Of[S, R, E](a)
}

func (o *stateReaderIOEitherApplicative[S, R, E, A, B]) Of(a A) StateReaderIOEither[S, R, E, A] {
	return Of[S, R, E](a)
}

func (o *stateReaderIOEitherMonad[S, R, E, A, B]) Map(f func(A) B) Operator[S, R, E, A, B] {
	return Map[S, R, E](f)
}

func (o *stateReaderIOEitherApplicative[S, R, E, A, B]) Map(f func(A) B) Operator[S, R, E, A, B] {
	return Map[S, R, E](f)
}

func (o *stateReaderIOEitherFunctor[S, R, E, A, B]) Map(f func(A) B) Operator[S, R, E, A, B] {
	return Map[S, R, E](f)
}

func (o *stateReaderIOEitherMonad[S, R, E, A, B]) Chain(f Kleisli[S, R, E, A, B]) Operator[S, R, E, A, B] {
	return Chain(f)
}

func (o *stateReaderIOEitherMonad[S, R, E, A, B]) Ap(fa StateReaderIOEither[S, R, E, A]) Operator[S, R, E, func(A) B, B] {
	return Ap[B](fa)
}

func (o *stateReaderIOEitherApplicative[S, R, E, A, B]) Ap(fa StateReaderIOEither[S, R, E, A]) Operator[S, R, E, func(A) B, B] {
	return Ap[B](fa)
}

// Pointed implements the [pointed.Pointed] operations for [StateReaderIOEither]
func Pointed[
	S, R, E, A any,
]() pointed.Pointed[A, StateReaderIOEither[S, R, E, A]] {
	return &stateReaderIOEitherPointed[S, R, E, A]{}
}

// Functor implements the [functor.Functor] operations for [StateReaderIOEither]
func Functor[
	S, R, E, A, B any,
]() functor.Functor[A, B, StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]] {
	return &stateReaderIOEitherFunctor[S, R, E, A, B]{}
}

// Applicative implements the [applicative.Applicative] operations for [StateReaderIOEither]
func Applicative[
	S, R, E, A, B any,
]() applicative.Applicative[A, B, StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]] {
	return &stateReaderIOEitherApplicative[S, R, E, A, B]{}
}

// Monad implements the [monad.Monad] operations for [StateReaderIOEither]
func Monad[
	S, R, E, A, B any,
]() monad.Monad[A, B, StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]] {
	return &stateReaderIOEitherMonad[S, R, E, A, B]{}
}
