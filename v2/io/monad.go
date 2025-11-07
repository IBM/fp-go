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

package io

import (
	"github.com/IBM/fp-go/v2/internal/monad"
)

type (
	ioMonad[A, B any] struct{}
	// IOMonad represents the monad type class for IO.
	// A monad combines the capabilities of Functor, Applicative, and Pointed
	// with the ability to chain computations (Chain/FlatMap).
	IOMonad[A, B any] = monad.Monad[A, B, IO[A], IO[B], IO[func(A) B]]
)

func (o *ioMonad[A, B]) Of(a A) IO[A] {
	return Of(a)
}

func (o *ioMonad[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

func (o *ioMonad[A, B]) Chain(f Kleisli[A, B]) Operator[A, B] {
	return Chain(f)
}

func (o *ioMonad[A, B]) Ap(fa IO[A]) Operator[func(A) B, B] {
	return Ap[B](fa)
}

// Monad returns an instance of the Monad type class for IO.
// This provides a structured way to access monadic operations (Of, Map, Chain, Ap)
// for IO computations.
//
// Example:
//
//	m := io.Monad[int, string]()
//	result := m.Chain(func(n int) io.IO[string] {
//	    return io.Of(strconv.Itoa(n))
//	})(m.Of(42))
func Monad[A, B any]() IOMonad[A, B] {
	return &ioMonad[A, B]{}
}
