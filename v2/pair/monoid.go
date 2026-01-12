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
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
)

// ApplicativeMonoid creates a monoid for [Pair] using applicative functor operations on the tail.
//
// This is an alias for [ApplicativeMonoidTail], which lifts the right (tail) monoid into the
// Pair applicative functor. The left monoid provides the semigroup for combining head values
// during applicative operations.
//
// IMPORTANT: The three monoid constructors (ApplicativeMonoid/ApplicativeMonoidTail and
// ApplicativeMonoidHead) produce DIFFERENT results:
//   - ApplicativeMonoidTail: Combines head values in REVERSE order (right-to-left)
//   - ApplicativeMonoidHead: Combines tail values in REVERSE order (right-to-left)
//   - The "focused" component (tail for Tail, head for Head) combines in normal order (left-to-right)
//
// This difference is significant for non-commutative operations like string concatenation.
//
// Parameters:
//   - l: A monoid for the head (left) values of type L
//   - r: A monoid for the tail (right) values of type R
//
// Returns:
//   - A Monoid[Pair[L, R]] that combines pairs using applicative operations on the tail
//
// Example:
//
//	import (
//	    N "github.com/IBM/fp-go/v2/number"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	intAdd := N.MonoidSum[int]()
//	strConcat := S.Monoid
//
//	pairMonoid := pair.ApplicativeMonoid(intAdd, strConcat)
//
//	p1 := pair.MakePair(10, "foo")
//	p2 := pair.MakePair(20, "bar")
//
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[int, string]{30, "foobar"}
//	// Note: head combines normally (10+20), tail combines normally ("foo"+"bar")
//
//	empty := pairMonoid.Empty()
//	// empty is Pair[int, string]{0, ""}
//
//go:inline
func ApplicativeMonoid[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]] {
	return ApplicativeMonoidTail(l, r)
}

// ApplicativeMonoidTail creates a monoid for [Pair] by lifting the tail monoid into the applicative functor.
//
// This function constructs a monoid using the applicative structure of Pair, focusing on
// the tail (right) value. The head values are combined using the left monoid's semigroup
// operation during applicative application.
//
// CRITICAL BEHAVIOR: Due to the applicative functor implementation, the HEAD values are
// combined in REVERSE order (right-to-left), while TAIL values combine in normal order
// (left-to-right). This matters for non-commutative operations:
//
//	strConcat := S.Monoid
//	pairMonoid := pair.ApplicativeMonoidTail(strConcat, strConcat)
//	p1 := pair.MakePair("hello", "foo")
//	p2 := pair.MakePair(" world", "bar")
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[string, string]{" worldhello", "foobar"}
//	//                                 ^^^^^^^^^^^^^^  ^^^^^^
//	//                                 REVERSED!       normal
//
// The resulting monoid satisfies the standard monoid laws:
//   - Associativity: Concat(Concat(p1, p2), p3) = Concat(p1, Concat(p2, p3))
//   - Left identity: Concat(Empty(), p) = p
//   - Right identity: Concat(p, Empty()) = p
//
// Parameters:
//   - l: A monoid for the head (left) values of type L
//   - r: A monoid for the tail (right) values of type R
//
// Returns:
//   - A Monoid[Pair[L, R]] that combines pairs component-wise
//
// Example:
//
//	import (
//	    N "github.com/IBM/fp-go/v2/number"
//	    M "github.com/IBM/fp-go/v2/monoid"
//	)
//
//	intAdd := N.MonoidSum[int]()
//	intMul := N.MonoidProduct[int]()
//
//	pairMonoid := pair.ApplicativeMonoidTail(intAdd, intMul)
//
//	p1 := pair.MakePair(5, 3)
//	p2 := pair.MakePair(10, 4)
//
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[int, int]{15, 12}  (5+10, 3*4)
//	// Note: Addition is commutative, so order doesn't matter for head
//
//	empty := pairMonoid.Empty()
//	// empty is Pair[int, int]{0, 1}
//
// Example with different types:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	boolAnd := M.MakeMonoid(func(a, b bool) bool { return a && b }, true)
//	strConcat := S.Monoid
//
//	pairMonoid := pair.ApplicativeMonoidTail(boolAnd, strConcat)
//
//	p1 := pair.MakePair(true, "hello")
//	p2 := pair.MakePair(true, " world")
//
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[bool, string]{true, "hello world"}
//	// Note: Boolean AND is commutative, so order doesn't matter for head
//
//go:inline
func ApplicativeMonoidTail[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]] {
	return M.ApplicativeMonoid(
		FromHead[R](l.Empty()),
		MonadMapTail[L, R, func(R) R],
		F.Bind1of3(MonadApTail[L, R, R])(l),
		r)
}

// ApplicativeMonoidHead creates a monoid for [Pair] by lifting the head monoid into the applicative functor.
//
// This function constructs a monoid using the applicative structure of Pair, focusing on
// the head (left) value. The tail values are combined using the right monoid's semigroup
// operation during applicative application.
//
// This is the dual of [ApplicativeMonoidTail], operating on the head instead of the tail.
//
// CRITICAL BEHAVIOR: Due to the applicative functor implementation, the TAIL values are
// combined in REVERSE order (right-to-left), while HEAD values combine in normal order
// (left-to-right). This is the opposite of ApplicativeMonoidTail:
//
//	strConcat := S.Monoid
//	pairMonoid := pair.ApplicativeMonoidHead(strConcat, strConcat)
//	p1 := pair.MakePair("hello", "foo")
//	p2 := pair.MakePair(" world", "bar")
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[string, string]{"hello world", "barfoo"}
//	//                                 ^^^^^^^^^^^^  ^^^^^^^^
//	//                                 normal        REVERSED!
//
// The resulting monoid satisfies the standard monoid laws:
//   - Associativity: Concat(Concat(p1, p2), p3) = Concat(p1, Concat(p2, p3))
//   - Left identity: Concat(Empty(), p) = p
//   - Right identity: Concat(p, Empty()) = p
//
// Parameters:
//   - l: A monoid for the head (left) values of type L
//   - r: A monoid for the tail (right) values of type R
//
// Returns:
//   - A Monoid[Pair[L, R]] that combines pairs component-wise
//
// Example:
//
//	import (
//	    N "github.com/IBM/fp-go/v2/number"
//	    M "github.com/IBM/fp-go/v2/monoid"
//	)
//
//	intMul := N.MonoidProduct[int]()
//	intAdd := N.MonoidSum[int]()
//
//	pairMonoid := pair.ApplicativeMonoidHead(intMul, intAdd)
//
//	p1 := pair.MakePair(3, 5)
//	p2 := pair.MakePair(4, 10)
//
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[int, int]{12, 15}  (3*4, 5+10)
//	// Note: Both operations are commutative, so order doesn't matter
//
//	empty := pairMonoid.Empty()
//	// empty is Pair[int, int]{1, 0}
//
// Example comparing Head vs Tail with non-commutative operations:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	strConcat := S.Monoid
//
//	// Using ApplicativeMonoidHead - tail values REVERSED
//	headMonoid := pair.ApplicativeMonoidHead(strConcat, strConcat)
//	p1 := pair.MakePair("hello", "foo")
//	p2 := pair.MakePair(" world", "bar")
//	result := headMonoid.Concat(p1, p2)
//	// result is Pair[string, string]{"hello world", "barfoo"}
//
//	// Using ApplicativeMonoidTail - head values REVERSED
//	tailMonoid := pair.ApplicativeMonoidTail(strConcat, strConcat)
//	result2 := tailMonoid.Concat(p1, p2)
//	// result2 is Pair[string, string]{" worldhello", "foobar"}
//	// DIFFERENT result! Head and tail are swapped in their reversal behavior
//
//go:inline
func ApplicativeMonoidHead[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]] {
	return M.ApplicativeMonoid(
		FromTail[L](r.Empty()),
		MonadMapHead[R, L, func(L) L],
		F.Bind1of3(MonadApHead[R, L, L])(r),
		l)
}
