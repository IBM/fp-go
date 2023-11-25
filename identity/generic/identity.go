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

package generic

import (
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
	FC "github.com/IBM/fp-go/internal/functor"
)

func MonadAp[GAB ~func(A) B, B, A any](fab GAB, fa A) B {
	return fab(fa)
}

func Ap[GAB ~func(A) B, B, A any](fa A) func(GAB) B {
	return F.Bind2nd(MonadAp[GAB, B, A], fa)
}

func MonadMap[GAB ~func(A) B, A, B any](fa A, f GAB) B {
	return f(fa)
}

func Map[GAB ~func(A) B, A, B any](f GAB) func(A) B {
	return f
}

func MonadChain[GAB ~func(A) B, A, B any](ma A, f GAB) B {
	return f(ma)
}

func Chain[GAB ~func(A) B, A, B any](f GAB) func(A) B {
	return f
}

func MonadChainFirst[GAB ~func(A) B, A, B any](fa A, f GAB) A {
	return C.MonadChainFirst(MonadChain[func(A) A, A, A], MonadMap[func(B) A, B, A], fa, f)
}

func ChainFirst[GAB ~func(A) B, A, B any](f GAB) func(A) A {
	return C.ChainFirst(MonadChain[func(A) A, A, A], MonadMap[func(B) A, B, A], f)
}

func MonadFlap[GAB ~func(A) B, A, B any](fab GAB, a A) B {
	return FC.MonadFlap(MonadMap[func(GAB) B, GAB, B], fab, a)
}

func Flap[GAB ~func(A) B, B, A any](a A) func(GAB) B {
	return F.Bind2nd(MonadFlap[GAB, A, B], a)
}
