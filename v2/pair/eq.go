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

import "github.com/IBM/fp-go/v2/eq"

// Eq constructs an equality predicate for [Pair] from equality predicates for both components.
// Two pairs are considered equal if both their head values are equal and their tail values are equal.
// The head comparison is evaluated first, the tail comparison only if the heads are equal.
//
// Example:
//
//	import EQ "github.com/IBM/fp-go/v2/eq"
//
//	pairEq := pair.Eq(
//	    EQ.FromStrictEquals[string](),
//	    EQ.FromStrictEquals[int](),
//	)
//	p1 := pair.MakePair("hello", 42)
//	p2 := pair.MakePair("hello", 42)
//	p3 := pair.MakePair("world", 42)
//	pairEq.Equals(p1, p2)  // true
//	pairEq.Equals(p1, p3)  // false
func Eq[A, B any](a eq.Eq[A], b eq.Eq[B]) eq.Eq[Pair[A, B]] {
	return eq.Semigroup[Pair[A, B]]().Concat(
		eq.ContraMap(Head[A, B])(a),
		eq.ContraMap(Tail[A, B])(b),
	)
}

// FromStrictEquals constructs an [eq.Eq] for [Pair] using the built-in equality operator (==)
// for both components. This is only available when both type parameters are comparable.
//
// Example:
//
//	pairEq := pair.FromStrictEquals[string, int]()
//	p1 := pair.MakePair("hello", 42)
//	p2 := pair.MakePair("hello", 42)
//	pairEq.Equals(p1, p2)  // true
//
//go:inline
func FromStrictEquals[A, B comparable]() eq.Eq[Pair[A, B]] {
	return Eq(eq.FromStrictEquals[A](), eq.FromStrictEquals[B]())
}
