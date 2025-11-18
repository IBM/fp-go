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

package state

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type statePointed[S, A any] struct{}

type stateFunctor[S, A, B any] struct{}

type stateApplicative[S, A, B any] struct{}

type stateMonad[S, A, B any] struct{}

func (o *statePointed[S, A]) Of(a A) State[S, A] {
	return Of[S](a)
}

func (o *stateApplicative[S, A, B]) Of(a A) State[S, A] {
	return Of[S](a)
}

func (o *stateMonad[S, A, B]) Of(a A) State[S, A] {
	return Of[S](a)
}

func (o *stateFunctor[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *stateApplicative[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *stateMonad[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *stateMonad[S, A, B]) Chain(f Kleisli[S, A, B]) Operator[S, A, B] {
	return Chain(f)
}

func (o *stateApplicative[S, A, B]) Ap(fa State[S, A]) func(State[S, func(A) B]) State[S, B] {
	return Ap[B](fa)
}

func (o *stateMonad[S, A, B]) Ap(fa State[S, A]) func(State[S, func(A) B]) State[S, B] {
	return Ap[B](fa)
}

// Pointed implements the pointed operations for [State]
func Pointed[S, A any]() pointed.Pointed[A, State[S, A]] {
	return &statePointed[S, A]{}
}

// Functor implements the functor operations for [State]
func Functor[S, A, B any]() functor.Functor[A, B, State[S, A], State[S, B]] {
	return &stateFunctor[S, A, B]{}
}

// Applicative implements the applicative operations for [State]
func Applicative[S, A, B any]() applicative.Applicative[A, B, State[S, A], State[S, B], State[S, func(A) B]] {
	return &stateApplicative[S, A, B]{}
}

// Monad implements the monadic operations for [State]
func Monad[S, A, B any]() monad.Monad[A, B, State[S, A], State[S, B], State[S, func(A) B]] {
	return &stateMonad[S, A, B]{}
}
