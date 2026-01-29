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

package lazy

import (
	"time"

	"github.com/IBM/fp-go/v2/io"
)

// Of creates a lazy computation that returns the given value.
// This is the most basic way to lift a value into the Lazy context.
//
// The computation is pure and will always return the same value when evaluated.
//
// Example:
//
//	computation := lazy.Of(42)
//	result := computation() // 42
func Of[A any](a A) Lazy[A] {
	return io.Of(a)
}

// FromLazy creates a lazy computation from another lazy computation.
// This is an identity function that can be useful for type conversions or
// making the intent explicit in code.
//
// Example:
//
//	original := func() int { return 42 }
//	wrapped := lazy.FromLazy(original)
//	result := wrapped() // 42
func FromLazy[A any](a Lazy[A]) Lazy[A] {
	return io.FromIO(a)
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure(f func()) Lazy[Void] {
	return io.FromImpure(f)
}

// MonadOf creates a lazy computation that returns the given value.
// This is an alias for Of, provided for consistency with monadic naming conventions.
//
// Example:
//
//	computation := lazy.MonadOf(42)
//	result := computation() // 42
func MonadOf[A any](a A) Lazy[A] {
	return io.MonadOf(a)
}

// MonadMap transforms the value inside a lazy computation using the provided function.
// The transformation is not applied until the lazy computation is evaluated.
//
// This is the monadic version of Map, taking the lazy computation as the first parameter.
//
// Example:
//
//	computation := lazy.Of(5)
//	doubled := lazy.MonadMap(computation, N.Mul(2))
//	result := doubled() // 10
func MonadMap[A, B any](fa Lazy[A], f func(A) B) Lazy[B] {
	return io.MonadMap(fa, f)
}

// Map transforms the value inside a lazy computation using the provided function.
// Returns a function that can be applied to a lazy computation.
//
// This is the curried version of MonadMap, useful for function composition.
//
// Example:
//
//	double := lazy.Map(N.Mul(2))
//	computation := lazy.Of(5)
//	result := double(computation)() // 10
//
//	// Or with pipe:
//	result := F.Pipe1(lazy.Of(5), double)() // 10
func Map[A, B any](f func(A) B) func(fa Lazy[A]) Lazy[B] {
	return io.Map(f)
}

// MonadMapTo replaces the value inside a lazy computation with a constant value.
// The original computation is still evaluated, but its result is discarded.
//
// This is useful when you want to sequence computations but only care about
// the side effects (though Lazy should represent pure computations).
//
// Example:
//
//	computation := lazy.Of("ignored")
//	replaced := lazy.MonadMapTo(computation, 42)
//	result := replaced() // 42
func MonadMapTo[A, B any](fa Lazy[A], b B) Lazy[B] {
	return io.MonadMapTo(fa, b)
}

// MapTo replaces the value inside a lazy computation with a constant value.
// Returns a function that can be applied to a lazy computation.
//
// This is the curried version of MonadMapTo.
//
// Example:
//
//	replaceWith42 := lazy.MapTo[string](42)
//	computation := lazy.Of("ignored")
//	result := replaceWith42(computation)() // 42
func MapTo[A, B any](b B) Kleisli[Lazy[A], B] {
	return io.MapTo[A](b)
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
func MonadChain[A, B any](fa Lazy[A], f Kleisli[A, B]) Lazy[B] {
	return io.MonadChain(fa, f)
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[A, B any](f Kleisli[A, B]) Kleisli[Lazy[A], B] {
	return io.Chain(f)
}

// MonadAp applies a lazy function to a lazy value.
// Both the function and the value are evaluated when the result is evaluated.
//
// This is the applicative functor operation, allowing you to apply functions
// that are themselves wrapped in a lazy context.
//
// Example:
//
//	lazyFunc := lazy.Of(N.Mul(2))
//	lazyValue := lazy.Of(5)
//	result := lazy.MonadAp(lazyFunc, lazyValue)() // 10
func MonadAp[B, A any](mab Lazy[func(A) B], ma Lazy[A]) Lazy[B] {
	return io.MonadApSeq(mab, ma)
}

// Ap applies a lazy function to a lazy value.
// Returns a function that takes a lazy function and returns a lazy result.
//
// This is the curried version of MonadAp, useful for function composition.
//
// Example:
//
//	lazyValue := lazy.Of(5)
//	applyTo5 := lazy.Ap[int](lazyValue)
//	lazyFunc := lazy.Of(N.Mul(2))
//	result := applyTo5(lazyFunc)() // 10
func Ap[B, A any](ma Lazy[A]) func(Lazy[func(A) B]) Lazy[B] {
	return io.ApSeq[B](ma)
}

func Flatten[A any](mma Lazy[Lazy[A]]) Lazy[A] {
	return io.Flatten(mma)
}

// Memoize computes the value of the provided [Lazy] monad lazily but exactly once
func Memoize[A any](ma Lazy[A]) Lazy[A] {
	return io.Memoize(ma)
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[A, B any](fa Lazy[A], f Kleisli[A, B]) Lazy[A] {
	return io.MonadChainFirst(fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[A, B any](f Kleisli[A, B]) Kleisli[Lazy[A], A] {
	return io.ChainFirst(f)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, B any](first Lazy[A], second Lazy[B]) Lazy[A] {
	return io.MonadApFirst(first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, B any](second Lazy[B]) Kleisli[Lazy[A], A] {
	return io.ApFirst[A](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, B any](first Lazy[A], second Lazy[B]) Lazy[B] {
	return io.MonadApSecond(first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, B any](second Lazy[B]) Kleisli[Lazy[A], B] {
	return io.ApSecond[A](second)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[A, B any](fa Lazy[A], fb Lazy[B]) Lazy[B] {
	return io.MonadChainTo(fa, fb)
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[A, B any](fb Lazy[B]) Kleisli[Lazy[A], B] {
	return io.ChainTo[A](fb)
}

// Now is a lazy computation that returns the current timestamp when evaluated.
// Each evaluation will return the current time at the moment of evaluation.
//
// Example:
//
//	time1 := lazy.Now()
//	// ... some time passes ...
//	time2 := lazy.Now()
//	// time1 and time2 will be different
var Now Lazy[time.Time] = io.Now

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() Lazy[A]) Lazy[A] {
	return io.Defer(gen)
}
