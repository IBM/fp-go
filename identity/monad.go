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

package identity

import (
	"github.com/IBM/fp-go/internal/monad"
)

type identityMonad[A, B any] struct{}

func (o *identityMonad[A, B]) Of(a A) A {
	return Of[A](a)
}

func (o *identityMonad[A, B]) Map(f func(A) B) func(A) B {
	return Map[A, B](f)
}

func (o *identityMonad[A, B]) Chain(f func(A) B) func(A) B {
	return Chain[A, B](f)
}

func (o *identityMonad[A, B]) Ap(fa A) func(func(A) B) B {
	return Ap[B, A](fa)
}

// Monad implements the monadic operations for [Option]
func Monad[A, B any]() monad.Monad[A, B, A, B, func(A) B] {
	return &identityMonad[A, B]{}
}
