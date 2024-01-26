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

package io

import (
	"time"

	G "github.com/IBM/fp-go/io/generic"
)

// IO represents a synchronous computation that cannot fail
// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioltagt] for more details
type IO[A any] func() A

func MakeIO[A any](f func() A) IO[A] {
	return G.MakeIO[IO[A]](f)
}

func Of[A any](a A) IO[A] {
	return G.Of[IO[A]](a)
}

func FromIO[A any](a IO[A]) IO[A] {
	return G.FromIO(a)
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure(f func()) IO[any] {
	return G.FromImpure[IO[any]](f)
}

func MonadOf[A any](a A) IO[A] {
	return G.MonadOf[IO[A]](a)
}

func MonadMap[A, B any](fa IO[A], f func(A) B) IO[B] {
	return G.MonadMap[IO[A], IO[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(fa IO[A]) IO[B] {
	return G.Map[IO[A], IO[B]](f)
}

func MonadMapTo[A, B any](fa IO[A], b B) IO[B] {
	return G.MonadMapTo[IO[A], IO[B]](fa, b)
}

func MapTo[A, B any](b B) func(IO[A]) IO[B] {
	return G.MapTo[IO[A], IO[B]](b)
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
func MonadChain[A, B any](fa IO[A], f func(A) IO[B]) IO[B] {
	return G.MonadChain(fa, f)
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[A, B any](f func(A) IO[B]) func(IO[A]) IO[B] {
	return G.Chain[IO[A]](f)
}

func MonadAp[B, A any](mab IO[func(A) B], ma IO[A]) IO[B] {
	return G.MonadAp[IO[A], IO[B]](mab, ma)
}

func Ap[B, A any](ma IO[A]) func(IO[func(A) B]) IO[B] {
	return G.Ap[IO[B], IO[func(A) B], IO[A]](ma)
}

func Flatten[A any](mma IO[IO[A]]) IO[A] {
	return G.Flatten(mma)
}

// Memoize computes the value of the provided [IO] monad lazily but exactly once
func Memoize[A any](ma IO[A]) IO[A] {
	return G.Memoize(ma)
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[A, B any](fa IO[A], f func(A) IO[B]) IO[A] {
	return G.MonadChainFirst(fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[A, B any](f func(A) IO[B]) func(IO[A]) IO[A] {
	return G.ChainFirst[IO[A]](f)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, B any](first IO[A], second IO[B]) IO[A] {
	return G.MonadApFirst[IO[A], IO[B], IO[func(B) A]](first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, B any](second IO[B]) func(IO[A]) IO[A] {
	return G.ApFirst[IO[A], IO[B], IO[func(B) A]](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, B any](first IO[A], second IO[B]) IO[B] {
	return G.MonadApSecond[IO[A], IO[B], IO[func(B) B]](first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, B any](second IO[B]) func(IO[A]) IO[B] {
	return G.ApSecond[IO[A], IO[B], IO[func(B) B]](second)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[A, B any](fa IO[A], fb IO[B]) IO[B] {
	return G.MonadChainTo(fa, fb)
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[A, B any](fb IO[B]) func(IO[A]) IO[B] {
	return G.ChainTo[IO[A]](fb)
}

// Now returns the current timestamp
var Now = G.Now[IO[time.Time]]()

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() IO[A]) IO[A] {
	return G.Defer[IO[A]](gen)
}

func MonadFlap[B, A any](fab IO[func(A) B], a A) IO[B] {
	return G.MonadFlap[func(A) B, IO[func(A) B], IO[B], A, B](fab, a)
}

func Flap[B, A any](a A) func(IO[func(A) B]) IO[B] {
	return G.Flap[func(A) B, IO[func(A) B], IO[B], A, B](a)
}

// Delay creates an operation that passes in the value after some delay
func Delay[A any](delay time.Duration) func(IO[A]) IO[A] {
	return G.Delay[IO[A]](delay)
}

// After creates an operation that passes after the given timestamp
func After[A any](timestamp time.Time) func(IO[A]) IO[A] {
	return G.After[IO[A]](timestamp)
}
