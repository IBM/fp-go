// Copyright (c) 2024 - 2025 IBM Corp.
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
Package pair provides a strongly-typed data structure for holding two values with functional operations.

# Overview

A Pair is a data structure that holds exactly two values of potentially different types.
Unlike tuples which are position-based, Pair provides semantic operations for working with
"head" and "tail" (or "first" and "second") values, along with rich functional programming
capabilities including Functor, Applicative, and Monad operations.

The Pair type:

	type Pair[A, B any] struct {
	    h, t any  // head and tail (internal)
	}

# Basic Usage

Creating pairs:

	// Create a pair from two values
	p := pair.MakePair("hello", 42)  // Pair[string, int]

	// Create a pair with the same value in both positions
	p := pair.Of(5)  // Pair[int, int]{5, 5}

	// Convert from tuple
	t := tuple.MakeTuple2("world", 100)
	p := pair.FromTuple(t)  // Pair[string, int]

Accessing values:

	p := pair.MakePair("hello", 42)

	head := pair.Head(p)    // "hello"
	tail := pair.Tail(p)    // 42

	// Alternative names
	first := pair.First(p)  // "hello" (same as Head)
	second := pair.Second(p) // 42 (same as Tail)

# Transforming Pairs

Map operations transform one or both values:

	p := pair.MakePair(5, "hello")

	// Map the tail (second) value
	p2 := pair.MonadMapTail(p, func(s string) int {
	    return len(s)
	})  // Pair[int, int]{5, 5}

	// Map the head (first) value
	p3 := pair.MonadMapHead(p, func(n int) string {
	    return fmt.Sprintf("%d", n)
	})  // Pair[string, string]{"5", "hello"}

	// Map both values
	p4 := pair.MonadBiMap(p,
	    func(n int) string { return fmt.Sprintf("%d", n) },
	    S.Size,
	)  // Pair[string, int]{"5", 5}

Curried versions for composition:

	import F "github.com/IBM/fp-go/v2/function"

	// Create a mapper function
	doubleHead := pair.MapHead[string](func(n int) int {
	    return n * 2
	})

	p := pair.MakePair(5, "hello")
	result := doubleHead(p)  // Pair[int, string]{10, "hello"}

	// Compose multiple transformations
	transform := F.Flow2(
	    pair.MapHead[string](N.Mul(2)),
	    pair.MapTail[int](S.Size),
	)
	result := transform(p)  // Pair[int, int]{10, 5}

# Swapping

Swap exchanges the head and tail values:

	p := pair.MakePair("hello", 42)
	swapped := pair.Swap(p)  // Pair[int, string]{42, "hello"}

# Monadic Operations

Pair supports monadic operations on both head and tail, requiring a Semigroup
for combining values:

	import (
	    SG "github.com/IBM/fp-go/v2/semigroup"
	    N "github.com/IBM/fp-go/v2/number"
	)

	// Chain on the tail (requires semigroup for head)
	intSum := N.SemigroupSum[int]()

	p := pair.MakePair(5, "hello")
	result := pair.MonadChainTail(intSum, p, func(s string) pair.Pair[int, int] {
	    return pair.MakePair(len(s), len(s) * 2)
	})  // Pair[int, int]{10, 10} (5 + 5 from semigroup, 10 from function)

	// Chain on the head (requires semigroup for tail)
	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })

	p2 := pair.MakePair(5, "hello")
	result2 := pair.MonadChainHead(strConcat, p2, func(n int) pair.Pair[string, string] {
	    return pair.MakePair(fmt.Sprintf("%d", n), "!")
	})  // Pair[string, string]{"5", "hello!"}

Curried versions:

	intSum := N.SemigroupSum[int]()
	chain := pair.ChainTail(intSum, func(s string) pair.Pair[int, int] {
	    return pair.MakePair(len(s), len(s) * 2)
	})

	p := pair.MakePair(5, "hello")
	result := chain(p)  // Pair[int, int]{10, 10}

# Applicative Operations

Apply functions wrapped in pairs to values in pairs:

	import N "github.com/IBM/fp-go/v2/number"

	intSum := N.SemigroupSum[int]()

	// Function in a pair
	pf := pair.MakePair(10, S.Size)

	// Value in a pair
	pv := pair.MakePair(5, "hello")

	// Apply (on tail)
	result := pair.MonadApTail(intSum, pf, pv)
	// Pair[int, int]{15, 5} (10+5 from semigroup, len("hello") from function)

	// Apply (on head)
	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
	pf2 := pair.MakePair(func(n int) string { return fmt.Sprintf("%d", n) }, "!")
	pv2 := pair.MakePair(42, "hello")

	result2 := pair.MonadApHead(strConcat, pf2, pv2)
	// Pair[string, string]{"42", "!hello"}

# Function Conversion

Convert between regular functions and pair-taking functions:

	// Regular function
	add := func(a, b int) int { return a + b }

	// Convert to pair-taking function
	pairedAdd := pair.Paired(add)
	result := pairedAdd(pair.MakePair(3, 4))  // 7

	// Convert back
	unpairedAdd := pair.Unpaired(pairedAdd)
	result = unpairedAdd(3, 4)  // 7

Merge curried functions:

	// Curried function
	add := func(b int) func(a int) int {
	    return func(a int) int { return a + b }
	}

	// Apply to pair
	merge := pair.Merge(add)
	result := merge(pair.MakePair(3, 4))  // 7 (applies 4 then 3)

# Equality

Compare pairs for equality:

	import (
	    EQ "github.com/IBM/fp-go/v2/eq"
	)

	// Create equality for pairs
	pairEq := pair.Eq(
	    EQ.FromStrictEquals[string](),
	    EQ.FromStrictEquals[int](),
	)

	p1 := pair.MakePair("hello", 42)
	p2 := pair.MakePair("hello", 42)
	p3 := pair.MakePair("world", 42)

	pairEq.Equals(p1, p2)  // true
	pairEq.Equals(p1, p3)  // false

For comparable types:

	pairEq := pair.FromStrictEquals[string, int]()

	p1 := pair.MakePair("hello", 42)
	p2 := pair.MakePair("hello", 42)

	pairEq.Equals(p1, p2)  // true

# Tuple Conversion

Convert between Pair and Tuple2:

	import "github.com/IBM/fp-go/v2/tuple"

	// Pair to Tuple
	p := pair.MakePair("hello", 42)
	t := pair.ToTuple(p)  // Tuple2[string, int]

	// Tuple to Pair
	t := tuple.MakeTuple2("world", 100)
	p := pair.FromTuple(t)  // Pair[string, int]

# Type Class Instances

Pair provides type class instances for functional programming:

Functor - Map over values:

	import M "github.com/IBM/fp-go/v2/monoid"

	// Functor for tail
	functor := pair.FunctorTail[int, string, int]()
	mapper := functor.Map(S.Size)

	p := pair.MakePair(5, "hello")
	result := mapper(p)  // Pair[int, int]{5, 5}

Pointed - Wrap values:

	import M "github.com/IBM/fp-go/v2/monoid"

	// Pointed for tail (requires monoid for head)
	intSum := M.MonoidSum[int]()
	pointed := pair.PointedTail[string](intSum)

	p := pointed.Of("hello")  // Pair[int, string]{0, "hello"}

Applicative - Apply wrapped functions:

	import M "github.com/IBM/fp-go/v2/monoid"

	intSum := M.MonoidSum[int]()
	applicative := pair.ApplicativeTail[string, int, int](intSum)

	// Create a pair with a function
	pf := applicative.Of(S.Size)

	// Apply to a value
	pv := pair.MakePair(5, "hello")
	result := applicative.Ap(pv)(pf)  // Pair[int, int]{5, 5}

Monad - Chain operations:

	import M "github.com/IBM/fp-go/v2/monoid"

	intSum := M.MonoidSum[int]()
	monad := pair.MonadTail[string, int, int](intSum)

	p := monad.Of("hello")
	result := monad.Chain(func(s string) pair.Pair[int, int] {
	    return pair.MakePair(len(s), len(s) * 2)
	})(p)  // Pair[int, int]{5, 10}

# Practical Examples

Example 1: Accumulating Results with Context

	import (
	    N "github.com/IBM/fp-go/v2/number"
	    F "github.com/IBM/fp-go/v2/function"
	)

	// Process items while accumulating a count
	intSum := N.SemigroupSum[int]()

	processItem := func(item string) pair.Pair[int, string] {
	    return pair.MakePair(1, strings.ToUpper(item))
	}

	items := []string{"hello", "world", "foo"}
	initial := pair.MakePair(0, "")

	result := F.Pipe2(
	    items,
	    A.Map(processItem),
	    A.Reduce(func(acc, curr pair.Pair[int, string]) pair.Pair[int, string] {
	        return pair.MonadBiMap(
	            pair.MakePair(
	                pair.First(acc) + pair.First(curr),
	                pair.Second(acc) + " " + pair.Second(curr),
	            ),
	            F.Identity[int],
	            strings.TrimSpace,
	        )
	    }, initial),
	)
	// Result: Pair[int, string]{3, "HELLO WORLD FOO"}

Example 2: Tracking Computation Steps

	type Log []string

	logConcat := SG.MakeSemigroup(func(a, b Log) Log {
	    return append(a, b...)
	})

	compute := func(n int) pair.Pair[Log, int] {
	    return pair.MakePair(
	        Log{fmt.Sprintf("computed %d", n)},
	        n * 2,
	    )
	}

	p := pair.MakePair(Log{"start"}, 5)
	result := pair.MonadChainTail(logConcat, p, compute)
	// Pair[Log, int]{
	//     []string{"start", "computed 5"},
	//     10
	// }

Example 3: Writer Monad Pattern

	import M "github.com/IBM/fp-go/v2/monoid"

	// Use pair as a writer monad
	stringMonoid := M.MakeMonoid(
	    func(a, b string) string { return a + b },
	    "",
	)

	monad := pair.MonadTail[string, string, int](stringMonoid)

	// Log and compute
	logAndDouble := func(n int) pair.Pair[string, int] {
	    return pair.MakePair(
	        fmt.Sprintf("doubled %d; ", n),
	        n * 2,
	    )
	}

	logAndAdd := func(n int) pair.Pair[string, int] {
	    return pair.MakePair(
	        fmt.Sprintf("added 10; ", n),
	        n + 10,
	    )
	}

	result := F.Pipe2(
	    monad.Of(5),
	    monad.Chain(logAndDouble),
	    monad.Chain(logAndAdd),
	)
	// Pair[string, int]{"doubled 5; added 10; ", 20}

# Function Reference

Creation:
  - MakePair[A, B any](A, B) Pair[A, B] - Create a pair from two values
  - Of[A any](A) Pair[A, A] - Create a pair with same value in both positions
  - FromTuple[A, B any](Tuple2[A, B]) Pair[A, B] - Convert tuple to pair
  - ToTuple[A, B any](Pair[A, B]) Tuple2[A, B] - Convert pair to tuple

Access:
  - Head[A, B any](Pair[A, B]) A - Get the head (first) value
  - Tail[A, B any](Pair[A, B]) B - Get the tail (second) value
  - First[A, B any](Pair[A, B]) A - Get the first value (alias for Head)
  - Second[A, B any](Pair[A, B]) B - Get the second value (alias for Tail)

Transformations:
  - MonadMapHead[B, A, A1 any](Pair[A, B], func(A) A1) Pair[A1, B] - Map head value
  - MonadMapTail[A, B, B1 any](Pair[A, B], func(B) B1) Pair[A, B1] - Map tail value
  - MonadMap[B, A, A1 any](Pair[A, B], func(A) A1) Pair[A1, B] - Map head value (alias)
  - MonadBiMap[A, B, A1, B1 any](Pair[A, B], func(A) A1, func(B) B1) Pair[A1, B1] - Map both values
  - MapHead[B, A, A1 any](func(A) A1) func(Pair[A, B]) Pair[A1, B] - Curried map head
  - MapTail[A, B, B1 any](func(B) B1) func(Pair[A, B]) Pair[A, B1] - Curried map tail
  - Map[A, B, B1 any](func(B) B1) func(Pair[A, B]) Pair[A, B1] - Curried map tail (alias)
  - BiMap[A, B, A1, B1 any](func(A) A1, func(B) B1) func(Pair[A, B]) Pair[A1, B1] - Curried bimap
  - Swap[A, B any](Pair[A, B]) Pair[B, A] - Swap head and tail

Monadic Operations:
  - MonadChainHead[B, A, A1 any](Semigroup[B], Pair[A, B], func(A) Pair[A1, B]) Pair[A1, B]
  - MonadChainTail[A, B, B1 any](Semigroup[A], Pair[A, B], func(B) Pair[A, B1]) Pair[A, B1]
  - MonadChain[A, B, B1 any](Semigroup[A], Pair[A, B], func(B) Pair[A, B1]) Pair[A, B1]
  - ChainHead[B, A, A1 any](Semigroup[B], func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B]
  - ChainTail[A, B, B1 any](Semigroup[A], func(B) Pair[A, B1]) func(Pair[A, B]) Pair[A, B1]
  - Chain[A, B, B1 any](Semigroup[A], func(B) Pair[A, B1]) func(Pair[A, B]) Pair[A, B1]

Applicative Operations:
  - MonadApHead[B, A, A1 any](Semigroup[B], Pair[func(A) A1, B], Pair[A, B]) Pair[A1, B]
  - MonadApTail[A, B, B1 any](Semigroup[A], Pair[A, func(B) B1], Pair[A, B]) Pair[A, B1]
  - MonadAp[A, B, B1 any](Semigroup[A], Pair[A, func(B) B1], Pair[A, B]) Pair[A, B1]
  - ApHead[B, A, A1 any](Semigroup[B], Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B]
  - ApTail[A, B, B1 any](Semigroup[A], Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1]
  - Ap[A, B, B1 any](Semigroup[A], Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1]

Function Conversion:
  - Paired[F ~func(T1, T2) R, T1, T2, R any](F) func(Pair[T1, T2]) R
  - Unpaired[F ~func(Pair[T1, T2]) R, T1, T2, R any](F) func(T1, T2) R
  - Merge[F ~func(B) func(A) R, A, B, R any](F) func(Pair[A, B]) R

Equality:
  - Eq[A, B any](Eq[A], Eq[B]) Eq[Pair[A, B]] - Create equality for pairs
  - FromStrictEquals[A, B comparable]() Eq[Pair[A, B]] - Equality for comparable types

Type Classes:
  - MonadHead[A, B, A1 any](Monoid[B]) Monad[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]]
  - MonadTail[B, A, B1 any](Monoid[A]) Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
  - Monad[B, A, B1 any](Monoid[A]) Monad[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
  - PointedHead[A, B any](Monoid[B]) Pointed[A, Pair[A, B]]
  - PointedTail[B, A any](Monoid[A]) Pointed[B, Pair[A, B]]
  - Pointed[B, A any](Monoid[A]) Pointed[B, Pair[A, B]]
  - FunctorHead[A, B, A1 any]() Functor[A, A1, Pair[A, B], Pair[A1, B]]
  - FunctorTail[B, A, B1 any]() Functor[B, B1, Pair[A, B], Pair[A, B1]]
  - Functor[B, A, B1 any]() Functor[B, B1, Pair[A, B], Pair[A, B1]]
  - ApplicativeHead[A, B, A1 any](Monoid[B]) Applicative[A, A1, Pair[A, B], Pair[A1, B], Pair[func(A) A1, B]]
  - ApplicativeTail[B, A, B1 any](Monoid[A]) Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]
  - Applicative[B, A, B1 any](Monoid[A]) Applicative[B, B1, Pair[A, B], Pair[A, B1], Pair[A, func(B) B1]]

# Related Packages

  - github.com/IBM/fp-go/v2/tuple - Position-based heterogeneous tuples
  - github.com/IBM/fp-go/v2/semigroup - Associative binary operations
  - github.com/IBM/fp-go/v2/monoid - Semigroups with identity
  - github.com/IBM/fp-go/v2/eq - Equality type class
  - github.com/IBM/fp-go/v2/function - Function composition utilities
*/
package pair
