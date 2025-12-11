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

// Package lazy provides a functional programming abstraction for synchronous computations
// without side effects. It represents deferred computations that are evaluated only when
// their result is needed.
//
// # Fantasy Land Specification
//
// This implementation corresponds to the Fantasy Land IO type (for pure computations):
// https://github.com/fantasyland/fantasy-land
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//
// # Overview
//
// A Lazy[A] is simply a function that takes no arguments and returns a value of type A:
//
//	type Lazy[A any] = func() A
//
// This allows you to defer the evaluation of a computation until it's actually needed,
// which is useful for:
//   - Avoiding unnecessary computations
//   - Creating infinite data structures
//   - Implementing memoization
//   - Composing computations in a pure functional style
//
// # Core Concepts
//
// The lazy package implements several functional programming patterns:
//
// **Functor**: Transform values inside a Lazy context using Map
//
// **Applicative**: Combine multiple Lazy computations using Ap and ApS
//
// **Monad**: Chain dependent computations using Chain and Bind
//
// **Memoization**: Cache computation results using Memoize
//
// # Basic Usage
//
// Creating and evaluating lazy computations:
//
//	import (
//	    "fmt"
//	    "github.com/IBM/fp-go/v2/lazy"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Create a lazy computation
//	computation := lazy.Of(42)
//
//	// Transform it
//	doubled := F.Pipe1(
//	    computation,
//	    lazy.Map(N.Mul(2)),
//	)
//
//	// Evaluate when needed
//	result := doubled() // 84
//
// # Memoization
//
// Lazy computations can be memoized to ensure they're evaluated only once:
//
//	import "math/rand"
//
//	// Without memoization - generates different values each time
//	random := lazy.FromLazy(rand.Int)
//	value1 := random() // e.g., 12345
//	value2 := random() // e.g., 67890 (different)
//
//	// With memoization - caches the first result
//	memoized := lazy.Memoize(rand.Int)
//	value1 := memoized() // e.g., 12345
//	value2 := memoized() // 12345 (same as value1)
//
// # Chaining Computations
//
// Use Chain to compose dependent computations:
//
//	getUserId := lazy.Of(123)
//
//	getUser := F.Pipe1(
//	    getUserId,
//	    lazy.Chain(func(id int) lazy.Lazy[User] {
//	        return lazy.Of(fetchUser(id))
//	    }),
//	)
//
//	user := getUser()
//
// # Do-Notation Style
//
// The package supports do-notation style composition using Bind and ApS:
//
//	type Config struct {
//	    Host string
//	    Port int
//	}
//
//	result := F.Pipe2(
//	    lazy.Do(Config{}),
//	    lazy.Bind(
//	        func(host string) func(Config) Config {
//	            return func(c Config) Config { c.Host = host; return c }
//	        },
//	        func(c Config) lazy.Lazy[string] {
//	            return lazy.Of("localhost")
//	        },
//	    ),
//	    lazy.Bind(
//	        func(port int) func(Config) Config {
//	            return func(c Config) Config { c.Port = port; return c }
//	        },
//	        func(c Config) lazy.Lazy[int] {
//	            return lazy.Of(8080)
//	        },
//	    ),
//	)
//
//	config := result() // Config{Host: "localhost", Port: 8080}
//
// # Traverse and Sequence
//
// Transform collections of values into lazy computations:
//
//	// Transform array elements
//	numbers := []int{1, 2, 3}
//	doubled := F.Pipe1(
//	    numbers,
//	    lazy.TraverseArray(func(x int) lazy.Lazy[int] {
//	        return lazy.Of(x * 2)
//	    }),
//	)
//	result := doubled() // []int{2, 4, 6}
//
//	// Sequence array of lazy computations
//	computations := []lazy.Lazy[int]{
//	    lazy.Of(1),
//	    lazy.Of(2),
//	    lazy.Of(3),
//	}
//	result := lazy.SequenceArray(computations)() // []int{1, 2, 3}
//
// # Retry Logic
//
// The package includes retry functionality for computations that may fail:
//
//	import (
//	    R "github.com/IBM/fp-go/v2/retry"
//	    "time"
//	)
//
//	policy := R.CapDelay(
//	    2*time.Second,
//	    R.Monoid.Concat(
//	        R.ExponentialBackoff(10),
//	        R.LimitRetries(5),
//	    ),
//	)
//
//	action := func(status R.RetryStatus) lazy.Lazy[string] {
//	    return lazy.Of(fetchData())
//	}
//
//	check := func(value string) bool {
//	    return value == "" // retry if empty
//	}
//
//	result := lazy.Retrying(policy, action, check)()
//
// # Algebraic Structures
//
// The package provides algebraic structures for combining lazy computations:
//
// **Semigroup**: Combine two lazy values using a semigroup operation
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	intAddSemigroup := lazy.ApplySemigroup(M.MonoidSum[int]())
//	result := intAddSemigroup.Concat(lazy.Of(5), lazy.Of(10))() // 15
//
// **Monoid**: Combine lazy values with an identity element
//
//	intAddMonoid := lazy.ApplicativeMonoid(M.MonoidSum[int]())
//	empty := intAddMonoid.Empty()() // 0
//	result := intAddMonoid.Concat(lazy.Of(5), lazy.Of(10))() // 15
//
// # Comparison
//
// Compare lazy computations by evaluating and comparing their results:
//
//	import EQ "github.com/IBM/fp-go/v2/eq"
//
//	eq := lazy.Eq(EQ.FromEquals[int]())
//	result := eq.Equals(lazy.Of(42), lazy.Of(42)) // true
//
// # Key Functions
//
// **Creation**:
//   - Of: Create a lazy computation from a value
//   - FromLazy: Create a lazy computation from another lazy computation
//   - FromImpure: Convert a side effect into a lazy computation
//   - Defer: Create a lazy computation from a generator function
//
// **Transformation**:
//   - Map: Transform the value inside a lazy computation
//   - MapTo: Replace the value with a constant
//   - Chain: Chain dependent computations
//   - ChainFirst: Chain computations but keep the first result
//   - Flatten: Flatten nested lazy computations
//
// **Combination**:
//   - Ap: Apply a lazy function to a lazy value
//   - ApFirst: Combine two computations, keeping the first result
//   - ApSecond: Combine two computations, keeping the second result
//
// **Memoization**:
//   - Memoize: Cache the result of a computation
//
// **Do-Notation**:
//   - Do: Start a do-notation context
//   - Bind: Bind a computation result to a context
//   - Let: Attach a pure value to a context
//   - LetTo: Attach a constant to a context
//   - BindTo: Initialize a context from a value
//   - ApS: Attach a value using applicative style
//
// **Lens-Based Operations**:
//   - BindL: Bind using a lens
//   - LetL: Let using a lens
//   - LetToL: LetTo using a lens
//   - ApSL: ApS using a lens
//
// **Collections**:
//   - TraverseArray: Transform array elements into lazy computations
//   - SequenceArray: Convert array of lazy computations to lazy array
//   - TraverseRecord: Transform record values into lazy computations
//   - SequenceRecord: Convert record of lazy computations to lazy record
//
// **Tuples**:
//   - SequenceT1, SequenceT2, SequenceT3, SequenceT4: Combine lazy computations into tuples
//
// **Retry**:
//   - Retrying: Retry a computation according to a policy
//
// **Algebraic**:
//   - ApplySemigroup: Create a semigroup for lazy values
//   - ApplicativeMonoid: Create a monoid for lazy values
//   - Eq: Create an equality predicate for lazy values
//
// # Relationship to IO
//
// The lazy package is built on top of the io package and shares the same underlying
// implementation. The key difference is conceptual:
//   - lazy.Lazy[A] represents a pure, synchronous computation without side effects
//   - io.IO[A] represents a computation that may have side effects
//
// In practice, they are the same type, but the lazy package provides a more focused
// API for pure computations.
package lazy
