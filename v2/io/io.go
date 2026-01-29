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

package io

import (
	"time"

	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	INTL "github.com/IBM/fp-go/v2/internal/lazy"
	"github.com/IBM/fp-go/v2/pair"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

// Of wraps a pure value in an IO context, creating a computation that returns that value.
// This is the monadic return operation for IO.
//
// Example:
//
//	greeting := io.Of("Hello, World!")
//	result := greeting() // returns "Hello, World!"
//
//go:inline
func Of[A any](a A) IO[A] {
	return F.Constant(a)
}

// FromIO is an identity function that returns the IO value unchanged.
// Useful for type conversions and maintaining consistency with other monad packages.
//
//go:inline
func FromIO[A any](a IO[A]) IO[A] {
	return a
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure[ANY ~func()](f ANY) IO[Void] {
	return func() Void {
		f()
		return function.VOID
	}
}

// MonadOf wraps a pure value in an IO context.
// This is an alias for Of, following the monadic naming convention.
//
//go:inline
func MonadOf[A any](a A) IO[A] {
	return F.Constant(a)
}

// MonadMap transforms the result of an IO computation by applying a function to it.
// The function is only applied when the IO is executed.
//
// Example:
//
//	doubled := io.MonadMap(io.Of(21), N.Mul(2))
//	result := doubled() // returns 42
//
//go:inline
func MonadMap[A, B any](fa IO[A], f func(A) B) IO[B] {
	//go:inline
	return func() B {
		return f(fa())
	}
}

// Map returns an operator that transforms the result of an IO computation.
// This is the curried version of MonadMap.
//
// Example:
//
//	double := io.Map(N.Mul(2))
//	doubled := double(io.Of(21))
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return F.Bind2nd(MonadMap[A, B], f)
}

// MonadMapTo replaces the result of an IO computation with a constant value.
// The original computation is still executed, but its result is discarded.
//
// Example:
//
//	always42 := io.MonadMapTo(sideEffect, 42)
//
//go:inline
func MonadMapTo[A, B any](fa IO[A], b B) IO[B] {
	return MonadMap(fa, F.Constant1[A](b))
}

// MapTo returns an operator that replaces the result with a constant value.
// This is the curried version of MonadMapTo.
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return Map(F.Constant1[A](b))
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
//
//go:inline
func MonadChain[A, B any](fa IO[A], f Kleisli[A, B]) IO[B] {
	//go:inline
	return func() B {
		return f(fa())()
	}
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return F.Bind2nd(MonadChain[A, B], f)
}

// MonadApSeq implements the applicative on a single thread by first executing mab and the ma
//
//go:inline
func MonadApSeq[A, B any](mab IO[func(A) B], ma IO[A]) IO[B] {
	return MonadChain(mab, F.Bind1st(MonadMap[A, B], ma))
}

// MonadApPar implements the applicative on two threads, the main thread executes mab and the actuall
// apply operation and the second thread computes ma. Communication between the threads happens via a channel
func MonadApPar[A, B any](mab IO[func(A) B], ma IO[A]) IO[B] {
	return func() B {
		c := make(chan A, 1)
		go func() {
			c <- ma()
			close(c)
		}()
		return mab()(<-c)
	}
}

// MonadAp implements the `ap` operation. Depending on a feature flag this will be sequential or parallel, the preferred implementation
// is parallel
//
//go:inline
func MonadAp[A, B any](mab IO[func(A) B], ma IO[A]) IO[B] {
	if useParallel {
		return MonadApPar(mab, ma)
	}
	return MonadApSeq(mab, ma)
}

// Ap returns an operator that applies a function wrapped in IO to a value wrapped in IO.
// This is the curried version of MonadAp and uses parallel execution by default.
//
// Example:
//
//	add := func(a int) func(int) int { return func(b int) int { return a + b } }
//	result := io.Ap(io.Of(2))(io.Of(add(3))) // parallel execution
//
//go:inline
func Ap[B, A any](ma IO[A]) Operator[func(A) B, B] {
	return F.Bind2nd(MonadAp[A, B], ma)
}

// ApSeq returns an operator that applies a function wrapped in IO to a value wrapped in IO sequentially.
// Unlike Ap, this executes the function and value computations in sequence.
//
//go:inline
func ApSeq[B, A any](ma IO[A]) Operator[func(A) B, B] {
	return Chain(F.Bind1st(MonadMap[A, B], ma))
}

// ApPar returns an operator that applies a function wrapped in IO to a value wrapped in IO in parallel.
// This explicitly uses parallel execution (same as Ap when useParallel is true).
//
//go:inline
func ApPar[B, A any](ma IO[A]) Operator[func(A) B, B] {
	return F.Bind2nd(MonadApPar[A, B], ma)
}

// Flatten removes one level of nesting from a nested IO computation.
// Converts IO[IO[A]] to IO[A].
//
// Example:
//
//	nested := io.Of(io.Of(42))
//	flattened := io.Flatten(nested)
//	result := flattened() // returns 42
//
//go:inline
func Flatten[A any](mma IO[IO[A]]) IO[A] {
	return MonadChain(mma, F.Identity)
}

// Memoize computes the value of the provided [IO] monad lazily but exactly once
func Memoize[A any](ma IO[A]) IO[A] {
	return INTL.Memoize(ma)
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[A, B any](fa IO[A], f Kleisli[A, B]) IO[A] {
	return chain.MonadChainFirst(MonadChain[A, A], MonadMap[B, A], fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(
		Chain[A, A],
		Map[B, A],
		f,
	)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, B any](first IO[A], second IO[B]) IO[A] {
	return apply.MonadApFirst(
		MonadAp[B, A],
		MonadMap[A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, B any](second IO[B]) Operator[A, A] {
	return apply.ApFirst(
		Ap[A, B],
		Map[A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, B any](first IO[A], second IO[B]) IO[B] {
	return apply.MonadApSecond(
		MonadAp[B, B],
		MonadMap[A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, B any](second IO[B]) Operator[A, B] {
	return apply.ApSecond(
		Ap[B, B],
		Map[A, func(B) B],

		second,
	)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[A, B any](fa IO[A], fb IO[B]) IO[B] {
	return MonadChain(fa, F.Constant1[A](fb))
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[A, B any](fb IO[B]) Operator[A, B] {
	return Chain(F.Constant1[A](fb))
}

// Now is an IO computation that returns the current timestamp when executed.
// Each execution returns the current time at that moment.
//
// Example:
//
//	timestamp := io.Now()
var Now IO[time.Time] = time.Now

// Defer creates an IO by creating a brand new IO via a generator function each time.
// This allows for dynamic creation of IO computations based on runtime conditions.
//
// Example:
//
//	deferred := io.Defer(func() io.IO[int] {
//	    if someCondition() {
//	        return io.Of(1)
//	    }
//	    return io.Of(2)
//	})
func Defer[A any](gen func() IO[A]) IO[A] {
	return func() A {
		return gen()()
	}
}

// MonadFlap applies a value to a function wrapped in IO.
// This is the reverse of Ap - instead of applying IO[func] to IO[value],
// it applies a pure value to IO[func].
//
// Example:
//
//	addFive := io.Of(N.Add(5))
//	result := io.MonadFlap(addFive, 10) // returns IO[15]
func MonadFlap[B, A any](fab IO[func(A) B], a A) IO[B] {
	return functor.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

// Flap returns an operator that applies a pure value to a function wrapped in IO.
// This is the curried version of MonadFlap.
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return functor.Flap(Map[func(A) B, B], a)
}

// Delay creates an operator that delays execution by the specified duration.
// The delay occurs before executing the wrapped computation.
//
// Example:
//
//	delayed := io.Delay(time.Second)(io.Of(42))
//	result := delayed() // waits 1 second, then returns 42
func Delay[A any](delay time.Duration) Operator[A, A] {
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

// After creates an operator that delays execution until after the given timestamp.
// If the timestamp is in the past, the computation executes immediately.
//
// Example:
//
//	future := time.Now().Add(5 * time.Second)
//	scheduled := io.After(future)(io.Of(42))
//	result := scheduled() // waits until future time, then returns 42
func After[A any](timestamp time.Time) Operator[A, A] {
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

// WithTime returns an IO that measures the start and end time.Time of the operation.
// Returns a Pair[Pair[time.Time, time.Time], A] where the head contains a nested pair of
// (start time, end time) and the tail contains the result. The result is placed in the tail
// position because that is the value that the pair monad operates on, allowing monadic
// operations to transform the result while preserving the timing information.
//
// Example:
//
//	timed := io.WithTime(expensiveComputation)
//	p := timed()
//	times := pair.Head(p)      // Pair[time.Time, time.Time]
//	result := pair.Tail(p)     // A
//	start := pair.Head(times)  // time.Time
//	end := pair.Tail(times)    // time.Time
func WithTime[A any](a IO[A]) IO[Pair[Pair[time.Time, time.Time], A]] {
	return func() Pair[Pair[time.Time, time.Time], A] {
		t0 := time.Now()
		res := a()
		t1 := time.Now()
		return pair.MakePair(pair.MakePair(t0, t1), res)
	}
}

// WithDuration returns an IO that measures the execution time.Duration of the operation.
// Returns a Pair with the duration as the head and the result as the tail.
// The result is placed in the tail position because that is the value that the pair monad
// operates on, allowing monadic operations to transform the result while preserving the duration.
//
// Example:
//
//	timed := io.WithDuration(expensiveComputation)
//	p := timed()
//	duration := pair.Head(p)
//	result := pair.Tail(p)
//	fmt.Printf("Took %v\n", duration)
func WithDuration[A any](a IO[A]) IO[Pair[time.Duration, A]] {
	return func() Pair[time.Duration, A] {
		t0 := time.Now()
		res := a()
		t1 := time.Now()
		return pair.MakePair(t1.Sub(t0), res)
	}
}
