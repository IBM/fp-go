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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
)

// Pointed returns the pointed operations for [ReaderIOEither]
//
//go:inline
func Pointed[R, E, A any]() pointed.Pointed[A, ReaderIOEither[R, E, A]] {
	return G.Pointed[R, E, A, ReaderIOEither[R, E, A]]()
}

// Functor returns the functor operations for [ReaderIOEither]
//
//go:inline
func Functor[R, E, A, B any]() functor.Functor[A, B, ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]] {
	return G.Functor[R, E, A, B, ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]]()
}

// Monad returns the monadic operations for [ReaderIOEither]
//
//go:inline
func Monad[R, E, A, B any]() monad.Monad[A, B, ReaderIOEither[R, E, A], ReaderIOEither[R, E, B], ReaderIOEither[R, E, func(A) B]] {
	return G.Monad[R, E, A, B, ReaderIOEither[R, E, A], ReaderIOEither[R, E, B], ReaderIOEither[R, E, func(A) B]]()
}
