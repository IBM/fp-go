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

/*
Package identity implements the Identity monad, the simplest possible monad.

# Fantasy Land Specification

This implementation corresponds to the Fantasy Land Identity type:
https://github.com/fantasyland/fantasy-land

Implemented Fantasy Land algebras:
  - Functor: https://github.com/fantasyland/fantasy-land#functor
  - Apply: https://github.com/fantasyland/fantasy-land#apply
  - Applicative: https://github.com/fantasyland/fantasy-land#applicative
  - Chain: https://github.com/fantasyland/fantasy-land#chain
  - Monad: https://github.com/fantasyland/fantasy-land#monad

# Overview

The Identity monad is a trivial monad that simply wraps a value without adding
any computational context. It's the identity element in the category of monads,
meaning it doesn't add any effects or behavior - it just passes values through.

While seemingly useless, the Identity monad serves several important purposes:
  - As a baseline for understanding more complex monads
  - For testing monad transformers
  - As a default when no specific monad is needed
  - For generic code that works with any monad

In this implementation, Identity[A] is simply represented as type A itself,
making it a zero-cost abstraction.

# Core Concepts

The Identity monad implements the standard monadic operations:

  - Of: Wraps a value (identity function)
  - Map: Transforms the wrapped value
  - Chain (FlatMap): Chains computations
  - Ap: Applies a wrapped function to a wrapped value

Since Identity adds no context, all these operations reduce to simple function
application.

# Basic Usage

Creating and transforming Identity values:

	// Of wraps a value (but it's just the identity)
	x := identity.Of(42)
	// x is just 42

	// Map transforms the value
	doubled := identity.Map(func(n int) int {
	    return n * 2
	})(x)
	// doubled is 84

	// Chain for monadic composition
	result := identity.Chain(func(n int) int {
	    return n + 10
	})(doubled)
	// result is 94

# Functor Operations

Map transforms values:

	import F "github.com/IBM/fp-go/v2/function"

	// Simple mapping
	result := F.Pipe1(
	    5,
	    identity.Map(func(n int) int { return n * n }),
	)
	// result is 25

	// MapTo replaces with a constant
	result := F.Pipe1(
	    "ignored",
	    identity.MapTo[string, int](100),
	)
	// result is 100

# Applicative Operations

Ap applies wrapped functions:

	add := func(a int) func(int) int {
	    return func(b int) int {
	        return a + b
	    }
	}

	// Apply a curried function
	result := F.Pipe1(
	    add(10),
	    identity.Ap[int, int](5),
	)
	// result is 15

# Monad Operations

Chain for sequential composition:

	// Chain multiple operations
	result := F.Pipe2(
	    10,
	    identity.Chain(N.Mul(2)),
	    identity.Chain(N.Add(5)),
	)
	// result is 25

	// ChainFirst executes for side effects but keeps original value
	result := F.Pipe1(
	    42,
	    identity.ChainFirst(func(n int) string {
	        return fmt.Sprintf("Value: %d", n)
	    }),
	)
	// result is still 42

# Do Notation

The package provides "do notation" for imperative-style composition:

	type Result struct {
	    X int
	    Y int
	    Sum int
	}

	result := F.Pipe3(
	    identity.Do(Result{}),
	    identity.Bind(
	        func(r Result) func(int) Result {
	            return func(x int) Result {
	                r.X = x
	                return r
	            }
	        },
	        func(Result) int { return 10 },
	    ),
	    identity.Bind(
	        func(r Result) func(int) Result {
	            return func(y int) Result {
	                r.Y = y
	                return r
	            }
	        },
	        func(Result) int { return 20 },
	    ),
	    identity.Let(
	        func(r Result) func(int) Result {
	            return func(sum int) Result {
	                r.Sum = sum
	                return r
	            }
	        },
	        func(r Result) int { return r.X + r.Y },
	    ),
	)
	// result is Result{X: 10, Y: 20, Sum: 30}

# Sequence and Traverse

Convert tuples of Identity values:

	import T "github.com/IBM/fp-go/v2/tuple"

	// Sequence a tuple
	tuple := T.MakeTuple2(1, 2)
	result := identity.SequenceTuple2(tuple)
	// result is T.Tuple2[int, int]{1, 2}

	// Traverse with transformation
	tuple := T.MakeTuple2(1, 2)
	result := identity.TraverseTuple2(
	    N.Mul(2),
	    N.Mul(3),
	)(tuple)
	// result is T.Tuple2[int, int]{2, 6}

# Monad Interface

Get a monad instance for generic code:

	m := identity.Monad[int, string]()

	// Use monad operations
	value := m.Of(42)
	mapped := m.Map(func(n int) string {
	    return fmt.Sprintf("Number: %d", n)
	})(value)

# Why Identity?

The Identity monad might seem pointless, but it's useful for:

1. Testing: Test monad transformers with a simple base monad
2. Defaults: Provide a default when no specific monad is needed
3. Learning: Understand monad laws without additional complexity
4. Abstraction: Write generic code that works with any monad

Example of generic code:

	func ProcessWithMonad[M any](
	    monad monad.Monad[int, string, M, M, func(int) M],
	    value int,
	) M {
	    return F.Pipe2(
	        monad.Of(value),
	        monad.Map(N.Mul(2)),
	        monad.Map(func(n int) string { return fmt.Sprintf("%d", n) }),
	    )
	}

	// Works with Identity
	result := ProcessWithMonad(identity.Monad[int, string](), 21)
	// result is "42"

# Type Alias

The package defines:

	type Operator[A, B any] = func(A) B

This represents an Identity computation from A to B, which is just a function.

# Functions

Core operations:
  - Of[A any](A) A - Wrap a value (identity)
  - Map[A, B any](func(A) B) func(A) B - Transform value
  - Chain[A, B any](func(A) B) func(A) B - Monadic bind
  - Ap[B, A any](A) func(func(A) B) B - Apply function

Monad variants:
  - MonadMap, MonadChain, MonadAp - Uncurried versions

Additional operations:
  - MapTo[A, B any](B) func(A) B - Replace with constant
  - ChainFirst[A, B any](func(A) B) func(A) A - Execute for effect
  - Flap[B, A any](A) func(func(A) B) B - Flip application

Do notation:
  - Do[S any](S) S - Initialize context
  - Bind[S1, S2, T any] - Bind computation result
  - Let[S1, S2, T any] - Bind pure value
  - LetTo[S1, S2, B any] - Bind constant
  - BindTo[S1, T any] - Initialize from value
  - ApS[S1, S2, T any] - Apply in context

Sequence/Traverse:
  - SequenceT1-10 - Sequence tuples of size 1-10
  - SequenceTuple1-10 - Sequence tuple types
  - TraverseTuple1-10 - Traverse with transformations

Monad instance:
  - Monad[A, B any]() - Get monad interface

# Related Packages

  - function: Function composition utilities
  - monad: Monad interface definition
  - tuple: Tuple types for sequence operations
*/
package identity

//go:generate go run .. identity --count 10 --filename gen.go
