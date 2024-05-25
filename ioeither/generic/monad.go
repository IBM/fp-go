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
	ET "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/monad"
	"github.com/IBM/fp-go/internal/pointed"
)

type ioEitherPointed[E, A any, GA ~func() ET.Either[E, A]] struct{}

type ioEitherMonad[E, A, B any, GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B]] struct{}

type ioEitherFunctor[E, A, B any, GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B]] struct{}

func (o *ioEitherPointed[E, A, GA]) Of(a A) GA {
	return Of[GA, E, A](a)
}

func (o *ioEitherMonad[E, A, B, GA, GB, GAB]) Of(a A) GA {
	return Of[GA, E, A](a)
}

func (o *ioEitherMonad[E, A, B, GA, GB, GAB]) Map(f func(A) B) func(GA) GB {
	return Map[GA, GB, E, A, B](f)
}

func (o *ioEitherMonad[E, A, B, GA, GB, GAB]) Chain(f func(A) GB) func(GA) GB {
	return Chain[GA, GB, E, A, B](f)
}

func (o *ioEitherMonad[E, A, B, GA, GB, GAB]) Ap(fa GA) func(GAB) GB {
	return Ap[GB, GAB, GA, E, A, B](fa)
}

func (o *ioEitherFunctor[E, A, B, GA, GB]) Map(f func(A) B) func(GA) GB {
	return Map[GA, GB, E, A, B](f)
}

// Pointed implements the pointed operations for [IOEither]
func Pointed[E, A any, GA ~func() ET.Either[E, A]]() pointed.Pointed[A, GA] {
	return &ioEitherPointed[E, A, GA]{}
}

// Functor implements the monadic operations for [IOEither]
func Functor[E, A, B any, GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B]]() functor.Functor[A, B, GA, GB] {
	return &ioEitherFunctor[E, A, B, GA, GB]{}
}

// Monad implements the monadic operations for [IOEither]
func Monad[E, A, B any, GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B]]() monad.Monad[A, B, GA, GB, GAB] {
	return &ioEitherMonad[E, A, B, GA, GB, GAB]{}
}
