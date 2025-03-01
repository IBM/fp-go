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

package state

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	G "github.com/IBM/fp-go/v2/state/generic"
)

// Pointed implements the pointed operations for [State]
func Pointed[S, A any]() pointed.Pointed[A, State[S, A]] {
	return G.Pointed[State[S, A], S, A]()
}

// Functor implements the pointed operations for [State]
func Functor[S, A, B any]() functor.Functor[A, B, State[S, A], State[S, B]] {
	return G.Functor[State[S, B], State[S, A], S, A, B]()
}

// Applicative implements the applicative operations for [State]
func Applicative[S, A, B any]() applicative.Applicative[A, B, State[S, A], State[S, B], State[S, func(A) B]] {
	return G.Applicative[State[S, B], State[S, func(A) B], State[S, A]]()
}

// Monad implements the monadic operations for [State]
func Monad[S, A, B any]() monad.Monad[A, B, State[S, A], State[S, B], State[S, func(A) B]] {
	return G.Monad[State[S, B], State[S, func(A) B], State[S, A]]()
}
