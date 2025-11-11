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

package generic

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type readerIOEitherPointed[R, E, A any, GRA ~func(R) GIOA, GIOA ~func() either.Either[E, A]] struct{}

type readerIOEitherMonad[R, E, A, B any, GRA ~func(R) GIOA, GRB ~func(R) GIOB, GRAB ~func(R) GIOAB, GIOA ~func() either.Either[E, A], GIOB ~func() either.Either[E, B], GIOAB ~func() either.Either[E, func(A) B]] struct{}

type readerIOEitherFunctor[R, E, A, B any, GRA ~func(R) GIOA, GRB ~func(R) GIOB, GIOA ~func() either.Either[E, A], GIOB ~func() either.Either[E, B]] struct{}

func (o *readerIOEitherPointed[R, E, A, GRA, GIOA]) Of(a A) GRA {
	return Of[GRA](a)
}

func (o *readerIOEitherMonad[R, E, A, B, GRA, GRB, GRAB, GIOA, GIOB, GIOAB]) Of(a A) GRA {
	return Of[GRA](a)
}

func (o *readerIOEitherMonad[R, E, A, B, GRA, GRB, GRAB, GIOA, GIOB, GIOAB]) Map(f func(A) B) func(GRA) GRB {
	return Map[GRA, GRB](f)
}

func (o *readerIOEitherMonad[R, E, A, B, GRA, GRB, GRAB, GIOA, GIOB, GIOAB]) Chain(f func(A) GRB) func(GRA) GRB {
	return Chain[GRA](f)
}

func (o *readerIOEitherMonad[R, E, A, B, GRA, GRB, GRAB, GIOA, GIOB, GIOAB]) Ap(fa GRA) func(GRAB) GRB {
	return Ap[GRA, GRB, GRAB](fa)
}

func (o *readerIOEitherFunctor[R, E, A, B, GRA, GRB, GIOA, GIOB]) Map(f func(A) B) func(GRA) GRB {
	return Map[GRA, GRB](f)
}

// Pointed implements the pointed operations for [ReaderIOEither]
func Pointed[R, E, A any, GRA ~func(R) GIOA, GIOA ~func() either.Either[E, A]]() pointed.Pointed[A, GRA] {
	return &readerIOEitherPointed[R, E, A, GRA, GIOA]{}
}

// Functor implements the monadic operations for [ReaderIOEither]
func Functor[R, E, A, B any, GRA ~func(R) GIOA, GRB ~func(R) GIOB, GIOA ~func() either.Either[E, A], GIOB ~func() either.Either[E, B]]() functor.Functor[A, B, GRA, GRB] {
	return &readerIOEitherFunctor[R, E, A, B, GRA, GRB, GIOA, GIOB]{}
}

// Monad implements the monadic operations for [ReaderIOEither]
func Monad[R, E, A, B any, GRA ~func(R) GIOA, GRB ~func(R) GIOB, GRAB ~func(R) GIOAB, GIOA ~func() either.Either[E, A], GIOB ~func() either.Either[E, B], GIOAB ~func() either.Either[E, func(A) B]]() monad.Monad[A, B, GRA, GRB, GRAB] {
	return &readerIOEitherMonad[R, E, A, B, GRA, GRB, GRAB, GIOA, GIOB, GIOAB]{}
}
