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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	G "github.com/IBM/fp-go/v2/statereaderioeither/generic"
)

// Pointed returns the pointed operations for [StateReaderIOEither]
func Pointed[S, R, E, A any]() pointed.Pointed[A, StateReaderIOEither[S, R, E, A]] {
	return G.Pointed[StateReaderIOEither[S, R, E, A]]()
}

// Functor returns the functor operations for [StateReaderIOEither]
func Functor[S, R, E, A, B any]() functor.Functor[A, B, StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]] {
	return G.Functor[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]]()
}

// Applicative returns the applicative operations for [StateReaderIOEither]
func Applicative[S, R, E, A, B any]() applicative.Applicative[A, B, StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]] {
	return G.Applicative[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]]()
}

// Monad returns the monadic operations for [StateReaderIOEither]
func Monad[S, R, E, A, B any]() monad.Monad[A, B, StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]] {
	return G.Monad[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B], StateReaderIOEither[S, R, E, func(A) B]]()
}
