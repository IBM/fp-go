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
	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	P "github.com/IBM/fp-go/v2/pair"
)

type stateReaderIOEitherPointed[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any,
] struct{}

type stateReaderIOEitherFunctor[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
] struct{}

type stateReaderIOEitherApplicative[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	SRIOEAB ~func(S) RIOEAB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	RIOEAB ~func(R) IOEAB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEAB ~func() ET.Either[E, P.Pair[func(A) B, S]],
	S, R, E, A, B any,
] struct{}

type stateReaderIOEitherMonad[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	SRIOEAB ~func(S) RIOEAB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	RIOEAB ~func(R) IOEAB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEAB ~func() ET.Either[E, P.Pair[func(A) B, S]],
	S, R, E, A, B any,
] struct{}

func (o *stateReaderIOEitherPointed[SRIOEA, RIOEA, IOEA, S, R, E, A]) Of(a A) SRIOEA {
	return Of[SRIOEA](a)
}

func (o *stateReaderIOEitherMonad[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Of(a A) SRIOEA {
	return Of[SRIOEA](a)
}

func (o *stateReaderIOEitherApplicative[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Of(a A) SRIOEA {
	return Of[SRIOEA](a)
}

func (o *stateReaderIOEitherMonad[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Map(f func(A) B) func(SRIOEA) SRIOEB {
	return Map[SRIOEA, SRIOEB](f)
}

func (o *stateReaderIOEitherApplicative[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Map(f func(A) B) func(SRIOEA) SRIOEB {
	return Map[SRIOEA, SRIOEB](f)
}

func (o *stateReaderIOEitherFunctor[SRIOEA, SRIOEB, RIOEA, RIOEB, IOEA, IOEB, S, R, E, A, B]) Map(f func(A) B) func(SRIOEA) SRIOEB {
	return Map[SRIOEA, SRIOEB](f)
}

func (o *stateReaderIOEitherMonad[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Chain(f func(A) SRIOEB) func(SRIOEA) SRIOEB {
	return Chain[SRIOEA, SRIOEB](f)
}

func (o *stateReaderIOEitherMonad[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Ap(fa SRIOEA) func(SRIOEAB) SRIOEB {
	return Ap[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B](fa)
}

func (o *stateReaderIOEitherApplicative[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]) Ap(fa SRIOEA) func(SRIOEAB) SRIOEB {
	return Ap[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B](fa)
}

// Pointed implements the pointed operations for [StateReaderIOEither]
func Pointed[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any,
]() pointed.Pointed[A, SRIOEA] {
	return &stateReaderIOEitherPointed[SRIOEA, RIOEA, IOEA, S, R, E, A]{}
}

// Functor implements the functor operations for [StateReaderIOEither]
func Functor[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
]() functor.Functor[A, B, SRIOEA, SRIOEB] {
	return &stateReaderIOEitherFunctor[SRIOEA, SRIOEB, RIOEA, RIOEB, IOEA, IOEB, S, R, E, A, B]{}
}

// Applicative implements the applicative operations for [StateReaderIOEither]
func Applicative[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	SRIOEAB ~func(S) RIOEAB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	RIOEAB ~func(R) IOEAB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEAB ~func() ET.Either[E, P.Pair[func(A) B, S]],
	S, R, E, A, B any,
]() applicative.Applicative[A, B, SRIOEA, SRIOEB, SRIOEAB] {
	return &stateReaderIOEitherApplicative[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]{}
}

// Monad implements the monadic operations for [StateReaderIOEither]
func Monad[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	SRIOEAB ~func(S) RIOEAB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	RIOEAB ~func(R) IOEAB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEAB ~func() ET.Either[E, P.Pair[func(A) B, S]],
	S, R, E, A, B any,
]() monad.Monad[A, B, SRIOEA, SRIOEB, SRIOEAB] {
	return &stateReaderIOEitherMonad[SRIOEA, SRIOEB, SRIOEAB, RIOEA, RIOEB, RIOEAB, IOEA, IOEB, IOEAB, S, R, E, A, B]{}
}
