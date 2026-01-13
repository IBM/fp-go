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

// MonadAp applies a function to a value in the Identity monad context.
// Since Identity has no computational context, this is just function application.
//
// This is the uncurried version of Ap.
//
// Implements the Fantasy Land Apply specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#apply
//
// Example:
//
//	result := identity.MonadAp(func(n int) int { return n * 2 }, 21)
//	// result is 42
func MonadAp[B, A any](fab func(A) B, fa A) B {
	return fab(fa)
}

// Ap applies a wrapped function to a wrapped value.
// Returns a function that takes a function and applies the value to it.
//
// This is the curried version of MonadAp, useful for composition with Pipe.
//
// Implements the Fantasy Land Apply specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#apply
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	double := func(n int) int { return n * 2 }
//	result := F.Pipe1(double, identity.Ap[int](21))
//	// result is 42
func Ap[B, A any](fa A) Operator[func(A) B, B] {
	return function.Bind2nd(MonadAp[B, A], fa)
}

// MonadMap transforms a value using a function in the Identity monad context.
// Since Identity has no computational context, this is just function application.
//
// This is the uncurried version of Map.
//
// Implements the Fantasy Land Functor specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#functor
//
// Example:
//
//	result := identity.MonadMap(21, func(n int) int { return n * 2 })
//	// result is 42
func MonadMap[A, B any](fa A, f func(A) B) B {
	return f(fa)
}

// Map transforms a value using a function.
// Returns the function itself since Identity adds no context.
//
// This is the curried version of MonadMap, useful for composition with Pipe.
//
// Implements the Fantasy Land Functor specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#functor
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe1(21, identity.Map(func(n int) int { return n * 2 }))
//	// result is 42
func Map[A, B any](f func(A) B) Operator[A, B] {
	return f
}

// MonadMapTo replaces a value with a constant, ignoring the input.
//
// This is the uncurried version of MapTo.
//
// Example:
//
//	result := identity.MonadMapTo("ignored", 42)
//	// result is 42
func MonadMapTo[A, B any](_ A, b B) B {
	return b
}

// MapTo replaces any value with a constant value.
// Returns a function that ignores its input and returns the constant.
//
// This is the curried version of MonadMapTo, useful for composition with Pipe.
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe1("ignored", identity.MapTo[string](42))
//	// result is 42
func MapTo[A, B any](b B) func(A) B {
	return function.Constant1[A](b)
}

// Of wraps a value in the Identity monad.
// Since Identity has no computational context, this is just the identity function.
//
// This is the Pointed/Applicative "pure" operation.
//
// Implements the Fantasy Land Applicative specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#applicative
//
// Example:
//
//	value := identity.Of(42)
//	// value is 42
//
//go:inline
func Of[A any](a A) A {
	return a
}

// MonadChain applies a Kleisli arrow to a value in the Identity monad context.
// Since Identity has no computational context, this is just function application.
//
// This is the uncurried version of Chain, also known as "bind" or "flatMap".
//
// Implements the Fantasy Land Chain specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#chain
//
// Example:
//
//	result := identity.MonadChain(21, func(n int) int { return n * 2 })
//	// result is 42
func MonadChain[A, B any](ma A, f Kleisli[A, B]) B {
	return f(ma)
}

// Chain applies a Kleisli arrow to a value.
// Returns the function itself since Identity adds no context.
//
// This is the curried version of MonadChain, also known as "bind" or "flatMap".
// Useful for composition with Pipe.
//
// Implements the Fantasy Land Chain specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#chain
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe1(21, identity.Chain(func(n int) int { return n * 2 }))
//	// result is 42
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return f
}

// MonadChainFirst executes a computation for its effect but returns the original value.
// Useful for side effects like logging while preserving the original value.
//
// This is the uncurried version of ChainFirst.
//
// Example:
//
//	result := identity.MonadChainFirst(42, func(n int) string {
//	    fmt.Printf("Value: %d\n", n)
//	    return "logged"
//	})
//	// result is 42 (original value preserved)
func MonadChainFirst[A, B any](fa A, f Kleisli[A, B]) A {
	return chain.MonadChainFirst(MonadChain[A, A], MonadMap[B, A], fa, f)
}

// ChainFirst executes a computation for its effect but returns the original value.
// Useful for side effects like logging while preserving the original value.
//
// This is the curried version of MonadChainFirst, useful for composition with Pipe.
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe1(
//	    42,
//	    identity.ChainFirst(func(n int) string {
//	        fmt.Printf("Value: %d\n", n)
//	        return "logged"
//	    }),
//	)
//	// result is 42 (original value preserved)
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(Chain[A, A], Map[B, A], f)
}

// MonadFlap applies a value to a function, flipping the normal application order.
// Instead of applying a function to a value, it applies a value to a function.
//
// This is the uncurried version of Flap.
//
// Example:
//
//	double := func(n int) int { return n * 2 }
//	result := identity.MonadFlap(double, 21)
//	// result is 42
func MonadFlap[B, A any](fab func(A) B, a A) B {
	return functor.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

// Flap applies a value to a function, flipping the normal application order.
// Returns a function that takes a function and applies the value to it.
//
// This is the curried version of MonadFlap, useful for composition with Pipe.
// Useful when you have a value and want to apply it to multiple functions.
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	double := func(n int) int { return n * 2 }
//	result := F.Pipe1(double, identity.Flap[int](21))
//	// result is 42
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return functor.Flap(Map[func(A) B, B], a)
}

// Extract extracts the value from the Identity monad.
// Since Identity has no computational context, this is just the identity function.
//
// This is the Comonad "extract" operation.
//
// Implements the Fantasy Land Comonad specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#comonad
//
// Example:
//
//	value := identity.Extract(42)
//	// value is 42
//
//go:inline
func Extract[A any](a A) A {
	return a
}

// Extend extends a computation over the Identity monad.
// Since Identity has no computational context, this is just function application.
//
// This is the Comonad "extend" operation, also known as "cobind".
//
// Implements the Fantasy Land Extend specification:
// https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#extend
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	result := F.Pipe1(21, identity.Extend(func(n int) int { return n * 2 }))
//	// result is 42
//
//go:inline
func Extend[A, B any](f func(A) B) Operator[A, B] {
	return f
}
