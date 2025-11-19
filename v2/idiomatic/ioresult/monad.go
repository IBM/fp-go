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

package ioresult

import (
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type (
	ioEitherPointed[A any] struct {
		fof Kleisli[A, A]
	}

	ioEitherFunctor[A, B any] struct {
		fmap func(func(A) B) Operator[A, B]
	}

	ioEitherApply[A, B any] struct {
		ioEitherFunctor[A, B]
		fap func(IOResult[A]) Operator[func(A) B, B]
	}

	ioEitherChainable[A, B any] struct {
		ioEitherApply[A, B]
		fchain func(Kleisli[A, B]) Operator[A, B]
	}

	ioEitherMonad[A, B any] struct {
		ioEitherPointed[A]
		ioEitherChainable[A, B]
	}
)

// Of implements the Pointed interface for IOResult.
func (o *ioEitherPointed[A]) Of(a A) IOResult[A] {
	return o.fof(a)
}

// Map implements the Monad interface's Map operation.
func (o *ioEitherFunctor[A, B]) Map(f func(A) B) Operator[A, B] {
	return o.fmap(f)
}

// Chain implements the Monad interface's Chain operation.
func (o *ioEitherChainable[A, B]) Chain(f Kleisli[A, B]) Operator[A, B] {
	return o.fchain(f)
}

// Ap implements the Monad interface's Ap operation.
func (o *ioEitherApply[A, B]) Ap(fa IOResult[A]) Operator[func(A) B, B] {
	return o.fap(fa)
}

// Pointed implements the pointed operations for [IOEither]
// Pointed returns a Pointed instance for IOResult.
// Pointed provides the ability to lift pure values into the IOResult context.
func Pointed[A any]() pointed.Pointed[A, IOResult[A]] {
	return &ioEitherPointed[A]{
		Of[A],
	}
}

// Functor implements the monadic operations for [IOEither]
// Functor returns a Functor instance for IOResult.
// Functor provides the Map operation for transforming values.
func Functor[A, B any]() functor.Functor[A, B, IOResult[A], IOResult[B]] {
	return &ioEitherFunctor[A, B]{
		Map[A, B],
	}
}

// Monad implements the monadic operations for [IOEither]
// Monad returns a Monad instance for IOResult.
// Monad provides the full monadic interface including Map, Chain, and Ap.
func Monad[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]] {
	return MonadPar[A, B]()
}

// Monad implements the monadic operations for [IOEither]
// Monad returns a Monad instance for IOResult.
// Monad provides the full monadic interface including Map, Chain, and Ap.
func MonadPar[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]] {
	return &ioEitherMonad[A, B]{
		ioEitherPointed[A]{
			Of[A],
		},
		ioEitherChainable[A, B]{
			ioEitherApply[A, B]{
				ioEitherFunctor[A, B]{
					Map[A, B],
				},
				ApPar[B, A],
			},
			Chain[A, B],
		},
	}
}

// Monad implements the monadic operations for [IOEither]
// Monad returns a Monad instance for IOResult.
// Monad provides the full monadic interface including Map, Chain, and Ap.
func MonadSeq[A, B any]() monad.Monad[A, B, IOResult[A], IOResult[B], IOResult[func(A) B]] {
	return &ioEitherMonad[A, B]{
		ioEitherPointed[A]{
			Of[A],
		},
		ioEitherChainable[A, B]{
			ioEitherApply[A, B]{
				ioEitherFunctor[A, B]{
					Map[A, B],
				},
				ApSeq[B, A],
			},
			Chain[A, B],
		},
	}
}
