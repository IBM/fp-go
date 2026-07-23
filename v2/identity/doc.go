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
meaning it doesn't add any effects or behavior — it just passes values through.

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

Of wraps a value (but it is just the identity), Map transforms it, and Chain
sequences computations. See the ExampleOf, ExampleMap, and ExampleChain
functions for runnable demonstrations.

# Functor Operations

Map transforms values. MapTo replaces a value with a constant.
See ExampleMap and ExampleMapTo for runnable demonstrations.

# Applicative Operations

Ap applies wrapped functions. See ExampleAp for a runnable demonstration.

# Monad Operations

Chain composes computations sequentially. ChainFirst executes a computation
for its side effect and returns the original value unchanged.
See ExampleChain and ExampleChainFirst for runnable demonstrations.

# Do Notation

The package provides "do notation" for imperative-style composition using
Do, Bind, Let, LetTo, BindTo, and ApS.
See ExampleBind, ExampleLet, ExampleBindTo, and ExampleApS for runnable
demonstrations.

# Sequence and Traverse

SequenceTuple and TraverseTuple convert tuples of Identity values.
See ExampleSequenceTuple2 and ExampleTraverseTuple2 for runnable demonstrations.

# Monad Interface

Monad returns a monad.Monad instance for use in generic code that works
with any monad.

# Why Identity?

The Identity monad might seem pointless, but it is useful for:

 1. Testing: Test monad transformers with a simple base monad
 2. Defaults: Provide a default when no specific monad is needed
 3. Learning: Understand monad laws without additional complexity
 4. Abstraction: Write generic code that works with any monad

# Type Alias

The package defines:

	type Operator[A, B any] = func(A) B

This represents an Identity computation from A to B, which is just a function.

# Functions

Core operations:
  - Of[A any](A) A — Wrap a value (identity)
  - Map[A, B any](func(A) B) func(A) B — Transform value
  - Chain[A, B any](func(A) B) func(A) B — Monadic bind
  - Ap[B, A any](A) func(func(A) B) B — Apply function

Monad variants:
  - MonadMap, MonadChain, MonadAp — Uncurried versions

Additional operations:
  - MapTo[A, B any](B) func(A) B — Replace with constant
  - ChainFirst[A, B any](func(A) B) func(A) A — Execute for effect
  - Flap[B, A any](A) func(func(A) B) B — Flip application

Do notation:
  - Do[S any](S) S — Initialize context
  - Bind[S1, S2, T any] — Bind computation result
  - Let[S1, S2, T any] — Bind pure value
  - LetTo[S1, S2, B any] — Bind constant
  - BindTo[S1, T any] — Initialize from value
  - ApS[S1, S2, T any] — Apply in context

Sequence/Traverse:
  - SequenceT1-10 — Sequence tuples of size 1–10
  - SequenceTuple1-10 — Sequence tuple types
  - TraverseTuple1-10 — Traverse with transformations

Monad instance:
  - Monad[A, B any]() — Get monad interface

# Related Packages

  - function: Function composition utilities
  - monad: Monad interface definition
  - tuple: Tuple types for sequence operations
*/
package identity

//go:generate go run .. identity --count 10 --filename gen.go
