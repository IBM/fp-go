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

package generic

import (
	"time"

	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	FC "github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/internal/lazy"
	P "github.com/IBM/fp-go/v2/pair"
	T "github.com/IBM/fp-go/v2/tuple"
)

var (
	// undefined represents an undefined value
	undefined = struct{}{}
)

// type IO[A any] = func() A

func MakeIO[GA ~func() A, A any](f func() A) GA {
	return f
}

func Of[GA ~func() A, A any](a A) GA {
	return MakeIO[GA](F.Constant(a))
}

func FromIO[GA ~func() A, A any](a GA) GA {
	return a
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure[GA ~func() any, IMP ~func()](f IMP) GA {
	return MakeIO[GA](func() any {
		f()
		return undefined
	})
}

func MonadOf[GA ~func() A, A any](a A) GA {
	return MakeIO[GA](F.Constant(a))
}

func MonadMap[GA ~func() A, GB ~func() B, A, B any](fa GA, f func(A) B) GB {
	return MakeIO[GB](func() B {
		return f(fa())
	})
}

func Map[GA ~func() A, GB ~func() B, A, B any](f func(A) B) func(GA) GB {
	return F.Bind2nd(MonadMap[GA, GB, A, B], f)
}

func MonadMapTo[GA ~func() A, GB ~func() B, A, B any](fa GA, b B) GB {
	return MonadMap[GA, GB](fa, F.Constant1[A](b))
}

func MapTo[GA ~func() A, GB ~func() B, A, B any](b B) func(GA) GB {
	return Map[GA, GB](F.Constant1[A](b))
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
func MonadChain[GA ~func() A, GB ~func() B, A, B any](fa GA, f func(A) GB) GB {
	return MakeIO[GB](func() B {
		return f(fa())()
	})
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[GA ~func() A, GB ~func() B, A, B any](f func(A) GB) func(GA) GB {
	return F.Bind2nd(MonadChain[GA, GB, A, B], f)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[GA ~func() A, GB ~func() B, A, B any](fa GA, fb GB) GB {
	return MonadChain(fa, F.Constant1[A](fb))
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[GA ~func() A, GB ~func() B, A, B any](fb GB) func(GA) GB {
	return Chain[GA](F.Constant1[A](fb))
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[GA ~func() A, GB ~func() B, A, B any](fa GA, f func(A) GB) GA {
	return C.MonadChainFirst(MonadChain[GA, GA, A, A], MonadMap[GB, GA, B, A], fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[GA ~func() A, GB ~func() B, A, B any](f func(A) GB) func(GA) GA {
	return C.ChainFirst(
		Chain[GA, GA, A, A],
		Map[GB, GA, B, A],
		f,
	)
}

func ApSeq[GB ~func() B, GAB ~func() func(A) B, GA ~func() A, B, A any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadApSeq[GA, GB, GAB, A, B], ma)
}

func ApPar[GB ~func() B, GAB ~func() func(A) B, GA ~func() A, B, A any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadApPar[GA, GB, GAB, A, B], ma)
}

func Ap[GB ~func() B, GAB ~func() func(A) B, GA ~func() A, B, A any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadAp[GA, GB, GAB, A, B], ma)
}

func Flatten[GA ~func() A, GAA ~func() GA, A any](mma GAA) GA {
	return mma()
}

// Memoize computes the value of the provided IO monad lazily but exactly once
func Memoize[GA ~func() A, A any](ma GA) GA {
	return L.Memoize(ma)
}

// Delay creates an operation that passes in the value after some delay
func Delay[GA ~func() A, A any](delay time.Duration) func(GA) GA {
	return func(ga GA) GA {
		return MakeIO[GA](func() A {
			time.Sleep(delay)
			return ga()
		})
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
func After[GA ~func() A, A any](timestamp time.Time) func(GA) GA {
	aft := after(timestamp)
	return func(ga GA) GA {
		return MakeIO[GA](func() A {
			// wait as long as necessary
			aft()
			// execute after wait
			return ga()
		})
	}
}

// Now returns the current timestamp
func Now[GA ~func() time.Time]() GA {
	return MakeIO[GA](time.Now)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[GA ~func() A, A any](gen func() GA) GA {
	return MakeIO[GA](func() A {
		return gen()()
	})
}

func MonadFlap[FAB ~func(A) B, GFAB ~func() FAB, GB ~func() B, A, B any](fab GFAB, a A) GB {
	return FC.MonadFlap(MonadMap[GFAB, GB, FAB, B], fab, a)
}

func Flap[FAB ~func(A) B, GFAB ~func() FAB, GB ~func() B, A, B any](a A) func(GFAB) GB {
	return FC.Flap(Map[GFAB, GB, FAB, B], a)
}

// WithTime returns an operation that measures the start and end timestamp of the operation
func WithTime[GTA ~func() T.Tuple3[A, time.Time, time.Time], GA ~func() A, A any](a GA) GTA {
	return MakeIO[GTA](func() T.Tuple3[A, time.Time, time.Time] {
		t0 := time.Now()
		res := a()
		t1 := time.Now()
		return T.MakeTuple3(res, t0, t1)
	})
}

// WithDuration returns an operation that measures the duration of the operation
func WithDuration[GTA ~func() P.Pair[time.Duration, A], GA ~func() A, A any](a GA) GTA {
	return MakeIO[GTA](func() P.Pair[time.Duration, A] {
		t0 := time.Now()
		res := a()
		t1 := time.Now()
		return P.MakePair(t1.Sub(t0), res)
	})
}
