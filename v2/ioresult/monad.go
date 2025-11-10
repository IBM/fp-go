// Copyright (c) 2024 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package ioresult

import (
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type (
	ioEitherPointed[A any] struct{}

	ioEitherMonad[A, B any] struct{}

	ioEitherFunctor[A, B any] struct{}
)

func (o *ioEitherPointed[A]) Of(a A) IOResult[A] {
	return Of(a)
}

func (o *ioEitherMonad[A, B]) Of(a A) IOResult[A] {
	return Of(a)
}

func (o *ioEitherMonad[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

func (o *ioEitherMonad[A, B]) Chain(f Kleisli[A, B]) Operator[A, B] {
	return Chain(f)
}

func (o *ioEitherMonad[A, B]) Ap(fa IOResult[A]) Operator[func(A) B, B] {
	return Ap[B](fa)
}

func (o *ioEitherFunctor[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

// Pointed implements the pointed operations for [IOResult]
func Pointed[A any]() pointed.Pointed[A, IOResult[A]] {
	return &ioEitherPointed[A]{}
}

// Functor implements the monadic operations for [IOResult]
func Functor[A, B any]() functor.Functor[A, B, IOResult[A], IOResult[B]] {
	return &ioEitherFunctor[A, B]{}
}

// Monad implements the monadic operations for [IOResult]
func Monad[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]] {
	return &ioEitherMonad[A, B]{}
}
