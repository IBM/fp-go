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

package either

import (
	"github.com/IBM/fp-go/v2/internal/monad"
)

type eitherMonad[E, A, B any] struct{}

func (o *eitherMonad[E, A, B]) Of(a A) Either[E, A] {
	return Of[E](a)
}

func (o *eitherMonad[E, A, B]) Map(f func(A) B) Operator[E, A, B] {
	return Map[E](f)
}

func (o *eitherMonad[E, A, B]) Chain(f func(A) Either[E, B]) Operator[E, A, B] {
	return Chain(f)
}

func (o *eitherMonad[E, A, B]) Ap(fa Either[E, A]) Operator[E, func(A) B, B] {
	return Ap[B](fa)
}

// Monad implements the monadic operations for Either.
// A monad combines the capabilities of Functor (Map), Applicative (Ap), and Chain (flatMap/bind).
// This allows for sequential composition of computations that may fail.
//
// Example:
//
//	m := either.Monad[error, int, string]()
//	result := m.Chain(func(x int) either.Either[error, string] {
//	    if x > 0 {
//	        return either.Right[error](strconv.Itoa(x))
//	    }
//	    return either.Left[string](errors.New("negative"))
//	})(either.Right[error](42))
//	// result is Right("42")
func Monad[E, A, B any]() monad.Monad[A, B, Either[E, A], Either[E, B], Either[E, func(A) B]] {
	return &eitherMonad[E, A, B]{}
}
