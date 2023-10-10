// Copyright (c) 2023 IBM Corp.
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
	F "github.com/IBM/fp-go/function"
	G "github.com/IBM/fp-go/identity/generic"
)

func MonadAp[B, A any](fab func(A) B, fa A) B {
	return G.MonadAp(fab, fa)
}

func Ap[B, A any](fa A) func(func(A) B) B {
	return G.Ap[func(A) B](fa)
}

func MonadMap[A, B any](fa A, f func(A) B) B {
	return G.MonadMap(fa, f)
}

func Map[A, B any](f func(A) B) func(A) B {
	return G.Map(f)
}

func MonadMapTo[A, B any](fa A, b B) B {
	return b
}

func MapTo[A, B any](b B) func(A) B {
	return F.Constant1[A](b)
}

func Of[A any](a A) A {
	return a
}

func MonadChain[A, B any](ma A, f func(A) B) B {
	return G.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) B) func(A) B {
	return G.Chain(f)
}

func MonadChainFirst[A, B any](fa A, f func(A) B) A {
	return G.MonadChainFirst(fa, f)
}

func ChainFirst[A, B any](f func(A) B) func(A) A {
	return G.ChainFirst(f)
}

func MonadFlap[B, A any](fab func(A) B, a A) B {
	return G.MonadFlap[func(A) B](fab, a)
}

func Flap[B, A any](a A) func(func(A) B) B {
	return G.Flap[func(A) B](a)
}
