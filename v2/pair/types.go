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

package pair

import (
	"github.com/IBM/fp-go/v2/semigroup"
	"github.com/IBM/fp-go/v2/tuple"
)

type (
	// Semigroup is a type alias for semigroup.Semigroup.
	// A Semigroup represents an algebraic structure with an associative binary operation.
	// It is used in various Pair operations to combine values when chaining or applying
	// operations that need to merge the head or tail values.
	//
	// Example:
	//
	//	import N "github.com/IBM/fp-go/v2/number"
	//	intSum := N.SemigroupSum[int]()  // Combines integers by addition
	Semigroup[A any] = semigroup.Semigroup[A]

	// Tuple2 is a type alias for tuple.Tuple2.
	// It represents a 2-tuple with two values. Pairs can be converted to and from Tuple2
	// using the FromTuple and ToTuple functions.
	//
	// Example:
	//
	//	t := tuple.MakeTuple2("hello", 42)
	//	p := pair.FromTuple(t)  // Convert Tuple2 to Pair
	Tuple2[A, B any] = tuple.Tuple2[A, B]

	// Pair defines a data structure that holds two strongly typed values.
	// The first value is called the "head" or "left" (L), and the second is called
	// the "tail" or "right" (R).
	//
	// Pair provides a foundation for functional programming patterns, supporting operations
	// like mapping, chaining, and applicative application on either or both values.
	//
	// Example:
	//
	//	p := pair.MakePair("hello", 42)
	//	head := pair.Head(p)  // "hello"
	//	tail := pair.Tail(p)  // 42
	Pair[L, R any] struct {
		r R
		l L
	}

	// Kleisli represents a Kleisli arrow for Pair.
	// It is a function that takes a value of type R1 and returns a Pair[L, R2].
	//
	// Kleisli arrows are used for monadic composition, allowing you to chain operations
	// that produce Pairs. They are particularly useful with Chain and ChainTail operations.
	//
	// Example:
	//
	//	// A Kleisli arrow that takes a string and returns a Pair
	//	lengthPair := func(s string) pair.Pair[int, string] {
	//	    return pair.MakePair(len(s), s + "!")
	//	}
	//	// Type: Kleisli[int, string, string]
	Kleisli[L, R1, R2 any] = func(R1) Pair[L, R2]

	// Operator represents an endomorphism on Pair that transforms the tail value.
	// It is a function that takes a Pair[L, R1] and returns a Pair[L, R2], preserving
	// the head type L while potentially transforming the tail from R1 to R2.
	//
	// Operators are commonly used as the return type of curried functions like Map, Chain,
	// and Ap, enabling function composition and point-free style programming.
	//
	// Example:
	//
	//	// An operator that maps the tail value
	//	toLengths := pair.Map[string](func(s string) int {
	//	    return len(s)
	//	})
	//	// Type: Operator[string, string, int]
	//	// Usage: p2 := toLengths(pair.MakePair("key", "value"))
	Operator[L, R1, R2 any] = func(Pair[L, R1]) Pair[L, R2]
)
