// Copyright (c) 2023 - 2025 IBM Corp.
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
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
)

func MonadAp[B, A any](fab func(A) B, fa A) B {
	return fab(fa)
}

func Ap[B, A any](fa A) Operator[func(A) B, B] {
	return function.Bind2nd(MonadAp[B, A], fa)
}

func MonadMap[A, B any](fa A, f func(A) B) B {
	return f(fa)
}

func Map[A, B any](f func(A) B) Operator[A, B] {
	return f
}

func MonadMapTo[A, B any](_ A, b B) B {
	return b
}

func MapTo[A, B any](b B) func(A) B {
	return function.Constant1[A](b)
}

func Of[A any](a A) A {
	return a
}

func MonadChain[A, B any](ma A, f Kleisli[A, B]) B {
	return f(ma)
}

func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return f
}

func MonadChainFirst[A, B any](fa A, f Kleisli[A, B]) A {
	return chain.MonadChainFirst(MonadChain[A, A], MonadMap[B, A], fa, f)
}

func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(Chain[A, A], Map[B, A], f)
}

func MonadFlap[B, A any](fab func(A) B, a A) B {
	return functor.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

func Flap[B, A any](a A) Operator[func(A) B, B] {
	return functor.Flap(Map[func(A) B, B], a)
}
