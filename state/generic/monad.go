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

package generic

import (
	"github.com/IBM/fp-go/internal/applicative"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
	P "github.com/IBM/fp-go/pair"
)

type statePointed[GA ~func(S) P.Pair[A, S], S, A any] struct{}

type stateFunctor[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], S, A, B any] struct{}

type stateApplicative[GB ~func(S) P.Pair[B, S], GAB ~func(S) P.Pair[func(A) B, S], GA ~func(S) P.Pair[A, S], S, A, B any] struct{}

type stateMonad[GB ~func(S) P.Pair[B, S], GAB ~func(S) P.Pair[func(A) B, S], GA ~func(S) P.Pair[A, S], S, A, B any] struct{}

func (o *statePointed[GA, S, A]) Of(a A) GA {
	return Of[GA](a)
}

func (o *stateApplicative[GB, GAB, GA, S, A, B]) Of(a A) GA {
	return Of[GA](a)
}

func (o *stateMonad[GB, GAB, GA, S, A, B]) Of(a A) GA {
	return Of[GA](a)
}

func (o *stateFunctor[GB, GA, S, A, B]) Map(f func(A) B) func(GA) GB {
	return Map[GB, GA](f)
}

func (o *stateApplicative[GB, GAB, GA, S, A, B]) Map(f func(A) B) func(GA) GB {
	return Map[GB, GA](f)
}

func (o *stateMonad[GB, GAB, GA, S, A, B]) Map(f func(A) B) func(GA) GB {
	return Map[GB, GA](f)
}

func (o *stateMonad[GB, GAB, GA, S, A, B]) Chain(f func(A) GB) func(GA) GB {
	return Chain[GB, GA](f)
}

func (o *stateApplicative[GB, GAB, GA, S, A, B]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB, GA](fa)
}

func (o *stateMonad[GB, GAB, GA, S, A, B]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB, GA](fa)
}

// Pointed implements the pointed operations for [Writer]
func Pointed[GA ~func(S) P.Pair[A, S], S, A any]() pointed.Pointed[A, GA] {
	return &statePointed[GA, S, A]{}
}

// Functor implements the functor operations for [Writer]
func Functor[GB ~func(S) P.Pair[B, S], GA ~func(S) P.Pair[A, S], S, A, B any]() functor.Functor[A, B, GA, GB] {
	return &stateFunctor[GB, GA, S, A, B]{}
}

// Applicative implements the applicative operations for [Writer]
func Applicative[GB ~func(S) P.Pair[B, S], GAB ~func(S) P.Pair[func(A) B, S], GA ~func(S) P.Pair[A, S], S, A, B any]() applicative.Applicative[A, B, GA, GB, GAB] {
	return &stateApplicative[GB, GAB, GA, S, A, B]{}
}

// Monad implements the monadic operations for [Writer]
func Monad[GB ~func(S) P.Pair[B, S], GAB ~func(S) P.Pair[func(A) B, S], GA ~func(S) P.Pair[A, S], S, A, B any]() monad.Monad[A, B, GA, GB, GAB] {
	return &stateMonad[GB, GAB, GA, S, A, B]{}
}
