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

package stateio

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type StateIOPointed[
	S, A any,
] struct{}

type StateIOFunctor[
	S, A, B any,
] struct{}

type StateIOApplicative[
	S, A, B any,
] struct{}

type StateIOMonad[
	S, A, B any,
] struct{}

func (o *StateIOPointed[S, A]) Of(a A) StateIO[S, A] {
	return Of[S, A](a)
}

func (o *StateIOMonad[S, A, B]) Of(a A) StateIO[S, A] {
	return Of[S](a)
}

func (o *StateIOApplicative[S, A, B]) Of(a A) StateIO[S, A] {
	return Of[S](a)
}

func (o *StateIOMonad[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *StateIOApplicative[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *StateIOFunctor[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

func (o *StateIOMonad[S, A, B]) Chain(f Kleisli[S, A, B]) Operator[S, A, B] {
	return Chain(f)
}

func (o *StateIOMonad[S, A, B]) Ap(fa StateIO[S, A]) Operator[S, func(A) B, B] {
	return Ap[B](fa)
}

func (o *StateIOApplicative[S, A, B]) Ap(fa StateIO[S, A]) Operator[S, func(A) B, B] {
	return Ap[B](fa)
}

// Pointed implements the [pointed.Pointed] operations for [StateIO]
func Pointed[
	S, A any,
]() pointed.Pointed[A, StateIO[S, A]] {
	return &StateIOPointed[S, A]{}
}

// Functor implements the [functor.Functor] operations for [StateIO]
func Functor[
	S, A, B any,
]() functor.Functor[A, B, StateIO[S, A], StateIO[S, B]] {
	return &StateIOFunctor[S, A, B]{}
}

// Applicative implements the [applicative.Applicative] operations for [StateIO]
func Applicative[
	S, A, B any,
]() applicative.Applicative[A, B, StateIO[S, A], StateIO[S, B], StateIO[S, func(A) B]] {
	return &StateIOApplicative[S, A, B]{}
}

// Monad implements the [monad.Monad] operations for [StateIO]
func Monad[
	S, A, B any,
]() monad.Monad[A, B, StateIO[S, A], StateIO[S, B], StateIO[S, func(A) B]] {
	return &StateIOMonad[S, A, B]{}
}
