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

package result

import (
	"github.com/IBM/fp-go/v2/internal/monad"
)

type eitherMonad[A, B any] struct{}

func (o *eitherMonad[A, B]) Of(a A) Result[A] {
	return Of(a)
}

func (o *eitherMonad[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

func (o *eitherMonad[A, B]) Chain(f func(A) Result[B]) Operator[A, B] {
	return Chain(f)
}

func (o *eitherMonad[A, B]) Ap(fa Result[A]) Operator[func(A) B, B] {
	return Ap[B](fa)
}

// Monad implements the monadic operations for Either.
// A monad combines the capabilities of Functor (Map), Applicative (Ap), and Chain (flatMap/bind).
// This allows for sequential composition of computations that may fail.
//
// Example:
//
//	m := either.Monad[error, int, string]()
//	result := m.Chain(func(x int) either.Result[string] {
//	    if x > 0 {
//	        return either.Right[error](strconv.Itoa(x))
//	    }
//	    return either.Left[string](errors.New("negative"))
//	})(either.Right[error](42))
//	// result is Right("42")
func Monad[A, B any]() monad.Monad[A, B, Result[A], Result[B], Result[func(A) B]] {
	return &eitherMonad[A, B]{}
}
