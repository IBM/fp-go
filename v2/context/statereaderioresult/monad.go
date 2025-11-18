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

package statereaderioresult

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type stateReaderIOResultPointed[
	S, A any,
] struct{}

type stateReaderIOResultFunctor[
	S, A, B any,
] struct{}

type stateReaderIOResultApplicative[
	S, A, B any,
] struct{}

type stateReaderIOResultMonad[
	S, A, B any,
] struct{}

func (o *stateReaderIOResultPointed[S, A]) Of(a A) StateReaderIOResult[S, A] {
	return Of[S](a)
}

func (o *stateReaderIOResultMonad[S, A, B]) Of(a A) StateReaderIOResult[S, A] {
	return Of[S](a)
}

func (o *stateReaderIOResultApplicative[S, A, B]) Of(a A) StateReaderIOResult[S, A] {
	return Of[S](a)
}

func (o *stateReaderIOResultMonad[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *stateReaderIOResultApplicative[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *stateReaderIOResultFunctor[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *stateReaderIOResultMonad[S, A, B]) Chain(f Kleisli[S, A, B]) Operator[S, A, B] {
	return Chain(f)
}

func (o *stateReaderIOResultMonad[S, A, B]) Ap(fa StateReaderIOResult[S, A]) Operator[S, func(A) B, B] {
	return Ap[B](fa)
}

func (o *stateReaderIOResultApplicative[S, A, B]) Ap(fa StateReaderIOResult[S, A]) Operator[S, func(A) B, B] {
	return Ap[B](fa)
}

// Pointed implements the [pointed.Pointed] operations for [StateReaderIOResult]
func Pointed[
	S, A any,
]() pointed.Pointed[A, StateReaderIOResult[S, A]] {
	return &stateReaderIOResultPointed[S, A]{}
}

// Functor implements the [functor.Functor] operations for [StateReaderIOResult]
func Functor[
	S, A, B any,
]() functor.Functor[A, B, StateReaderIOResult[S, A], StateReaderIOResult[S, B]] {
	return &stateReaderIOResultFunctor[S, A, B]{}
}

// Applicative implements the [applicative.Applicative] operations for [StateReaderIOResult]
func Applicative[
	S, A, B any,
]() applicative.Applicative[A, B, StateReaderIOResult[S, A], StateReaderIOResult[S, B], StateReaderIOResult[S, func(A) B]] {
	return &stateReaderIOResultApplicative[S, A, B]{}
}

// Monad implements the [monad.Monad] operations for [StateReaderIOResult]
func Monad[
	S, A, B any,
]() monad.Monad[A, B, StateReaderIOResult[S, A], StateReaderIOResult[S, B], StateReaderIOResult[S, func(A) B]] {
	return &stateReaderIOResultMonad[S, A, B]{}
}
