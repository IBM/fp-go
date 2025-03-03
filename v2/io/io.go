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

	F "github.com/IBM/fp-go/v2/function"
	INTA "github.com/IBM/fp-go/v2/internal/apply"
	INTC "github.com/IBM/fp-go/v2/internal/chain"
	INTF "github.com/IBM/fp-go/v2/internal/functor"
	INTL "github.com/IBM/fp-go/v2/internal/lazy"
	M "github.com/IBM/fp-go/v2/monoid"
	R "github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/semigroup"
	T "github.com/IBM/fp-go/v2/tuple"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

var (
	// undefined represents an undefined value
	undefined = struct{}{}
)

type (
	// IO represents a synchronous computation that cannot fail
	// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioltagt] for more details
	IO[A any] = func() A

	Mapper[A, B any] = R.Reader[IO[A], IO[B]]
	Monoid[A any]    = M.Monoid[IO[A]]
	Semigroup[A any] = S.Semigroup[IO[A]]
)

func Of[A any](a A) IO[A] {
	return F.Constant(a)
}

func FromIO[A any](a IO[A]) IO[A] {
	return a
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure(f func()) IO[any] {
	return func() any {
		f()
		return undefined
	}
}

func MonadOf[A any](a A) IO[A] {
	return F.Constant(a)
}

func MonadMap[A, B any](fa IO[A], f func(A) B) IO[B] {
	return func() B {
		return f(fa())
	}
}

func Map[A, B any](f func(A) B) Mapper[A, B] {
	return F.Bind2nd(MonadMap[A, B], f)
}

func MonadMapTo[A, B any](fa IO[A], b B) IO[B] {
	return MonadMap(fa, F.Constant1[A](b))
}

func MapTo[A, B any](b B) Mapper[A, B] {
	return Map(F.Constant1[A](b))
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
func MonadChain[A, B any](fa IO[A], f func(A) IO[B]) IO[B] {
	return func() B {
		return f(fa())()
	}
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[A, B any](f func(A) IO[B]) Mapper[A, B] {
	return F.Bind2nd(MonadChain[A, B], f)
}

// MonadApSeq implements the applicative on a single thread by first executing mab and the ma
func MonadApSeq[A, B any](mab IO[func(A) B], ma IO[A]) IO[B] {
	return MonadChain(mab, F.Bind1st(MonadMap[A, B], ma))
}

// MonadApPar implements the applicative on two threads, the main thread executes mab and the actuall
// apply operation and the second thread computes ma. Communication between the threads happens via a channel
func MonadApPar[A, B any](mab IO[func(A) B], ma IO[A]) IO[B] {
	return func() B {
		c := make(chan A)
		go func() {
			c <- ma()
			close(c)
		}()
		return mab()(<-c)
	}
}

// MonadAp implements the `ap` operation. Depending on a feature flag this will be sequential or parallel, the preferred implementation
// is parallel
func MonadAp[A, B any](mab IO[func(A) B], ma IO[A]) IO[B] {
	if useParallel {
		return MonadApPar(mab, ma)
	}
	return MonadApSeq(mab, ma)
}

func Ap[B, A any](ma IO[A]) Mapper[func(A) B, B] {
	return F.Bind2nd(MonadAp[A, B], ma)
}

func ApSeq[B, A any](ma IO[A]) Mapper[func(A) B, B] {
	return Chain(F.Bind1st(MonadMap[A, B], ma))
}

func ApPar[B, A any](ma IO[A]) Mapper[func(A) B, B] {
	return F.Bind2nd(MonadApPar[A, B], ma)
}

func Flatten[A any](mma IO[IO[A]]) IO[A] {
	return MonadChain(mma, F.Identity)
}

// Memoize computes the value of the provided [IO] monad lazily but exactly once
func Memoize[A any](ma IO[A]) IO[A] {
	return INTL.Memoize(ma)
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[A, B any](fa IO[A], f func(A) IO[B]) IO[A] {
	return INTC.MonadChainFirst(MonadChain[A, A], MonadMap[B, A], fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[A, B any](f func(A) IO[B]) Mapper[A, A] {
	return INTC.ChainFirst(
		Chain[A, A],
		Map[B, A],
		f,
	)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, B any](first IO[A], second IO[B]) IO[A] {
	return INTA.MonadApFirst(
		MonadAp[B, A],
		MonadMap[A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, B any](second IO[B]) Mapper[A, A] {
	return INTA.ApFirst(
		MonadAp[B, A],
		MonadMap[A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, B any](first IO[A], second IO[B]) IO[B] {
	return INTA.MonadApSecond(
		MonadAp[B, B],
		MonadMap[A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, B any](second IO[B]) Mapper[A, B] {
	return INTA.ApSecond(
		MonadAp[B, B],
		MonadMap[A, func(B) B],

		second,
	)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[A, B any](fa IO[A], fb IO[B]) IO[B] {
	return MonadChain(fa, F.Constant1[A](fb))
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[A, B any](fb IO[B]) Mapper[A, B] {
	return Chain(F.Constant1[A](fb))
}

// Now returns the current timestamp
var Now IO[time.Time] = time.Now

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() IO[A]) IO[A] {
	return func() A {
		return gen()()
	}
}

func MonadFlap[B, A any](fab IO[func(A) B], a A) IO[B] {
	return INTF.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

func Flap[B, A any](a A) Mapper[func(A) B, B] {
	return INTF.Flap(Map[func(A) B, B], a)
}

// Delay creates an operation that passes in the value after some delay
func Delay[A any](delay time.Duration) Mapper[A, A] {
	return func(ga IO[A]) IO[A] {
		return func() A {
			time.Sleep(delay)
			return ga()
		}
	}
}

func after(timestamp time.Time) func() {
	return func() {
		// check if we need to wait
		current := time.Now()
		if current.Before(timestamp) {
			time.Sleep(timestamp.Sub(current))
		}
	}
}

// After creates an operation that passes after the given timestamp
func After[A any](timestamp time.Time) Mapper[A, A] {
	aft := after(timestamp)
	return func(ga IO[A]) IO[A] {
		return func() A {
			// wait as long as necessary
			aft()
			// execute after wait
			return ga()
		}
	}
}

// WithTime returns an operation that measures the start and end [time.Time] of the operation
func WithTime[A any](a IO[A]) IO[T.Tuple3[A, time.Time, time.Time]] {
	return func() T.Tuple3[A, time.Time, time.Time] {
		t0 := time.Now()
		res := a()
		t1 := time.Now()
		return T.MakeTuple3(res, t0, t1)
	}
}

// WithDuration returns an operation that measures the [time.Duration]
func WithDuration[A any](a IO[A]) IO[T.Tuple2[A, time.Duration]] {
	return func() T.Tuple2[A, time.Duration] {
		t0 := time.Now()
		res := a()
		t1 := time.Now()
		return T.MakeTuple2(res, t1.Sub(t0))
	}
}
