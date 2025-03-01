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

package lazy

import (
	"time"

	G "github.com/IBM/fp-go/v2/io/generic"
)

// Lazy represents a synchronous computation without side effects
type Lazy[A any] func() A

func MakeLazy[A any](f func() A) Lazy[A] {
	return G.MakeIO[Lazy[A]](f)
}

func Of[A any](a A) Lazy[A] {
	return G.Of[Lazy[A]](a)
}

func FromLazy[A any](a Lazy[A]) Lazy[A] {
	return G.FromIO(a)
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure(f func()) Lazy[any] {
	return G.FromImpure[Lazy[any]](f)
}

func MonadOf[A any](a A) Lazy[A] {
	return G.MonadOf[Lazy[A]](a)
}

func MonadMap[A, B any](fa Lazy[A], f func(A) B) Lazy[B] {
	return G.MonadMap[Lazy[A], Lazy[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(fa Lazy[A]) Lazy[B] {
	return G.Map[Lazy[A], Lazy[B]](f)
}

func MonadMapTo[A, B any](fa Lazy[A], b B) Lazy[B] {
	return G.MonadMapTo[Lazy[A], Lazy[B]](fa, b)
}

func MapTo[A, B any](b B) func(Lazy[A]) Lazy[B] {
	return G.MapTo[Lazy[A], Lazy[B]](b)
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
func MonadChain[A, B any](fa Lazy[A], f func(A) Lazy[B]) Lazy[B] {
	return G.MonadChain(fa, f)
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[A, B any](f func(A) Lazy[B]) func(Lazy[A]) Lazy[B] {
	return G.Chain[Lazy[A]](f)
}

func MonadAp[B, A any](mab Lazy[func(A) B], ma Lazy[A]) Lazy[B] {
	return G.MonadApSeq[Lazy[A], Lazy[B]](mab, ma)
}

func Ap[B, A any](ma Lazy[A]) func(Lazy[func(A) B]) Lazy[B] {
	return G.ApSeq[Lazy[B], Lazy[func(A) B], Lazy[A]](ma)
}

func Flatten[A any](mma Lazy[Lazy[A]]) Lazy[A] {
	return G.Flatten(mma)
}

// Memoize computes the value of the provided [Lazy] monad lazily but exactly once
func Memoize[A any](ma Lazy[A]) Lazy[A] {
	return G.Memoize(ma)
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[A, B any](fa Lazy[A], f func(A) Lazy[B]) Lazy[A] {
	return G.MonadChainFirst(fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[A, B any](f func(A) Lazy[B]) func(Lazy[A]) Lazy[A] {
	return G.ChainFirst[Lazy[A]](f)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, B any](first Lazy[A], second Lazy[B]) Lazy[A] {
	return G.MonadApFirst[Lazy[A], Lazy[B], Lazy[func(B) A]](first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, B any](second Lazy[B]) func(Lazy[A]) Lazy[A] {
	return G.ApFirst[Lazy[A], Lazy[B], Lazy[func(B) A]](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, B any](first Lazy[A], second Lazy[B]) Lazy[B] {
	return G.MonadApSecond[Lazy[A], Lazy[B], Lazy[func(B) B]](first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, B any](second Lazy[B]) func(Lazy[A]) Lazy[B] {
	return G.ApSecond[Lazy[A], Lazy[B], Lazy[func(B) B]](second)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[A, B any](fa Lazy[A], fb Lazy[B]) Lazy[B] {
	return G.MonadChainTo(fa, fb)
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[A, B any](fb Lazy[B]) func(Lazy[A]) Lazy[B] {
	return G.ChainTo[Lazy[A]](fb)
}

// Now returns the current timestamp
var Now = G.Now[Lazy[time.Time]]()

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() Lazy[A]) Lazy[A] {
	return G.Defer[Lazy[A]](gen)
}
