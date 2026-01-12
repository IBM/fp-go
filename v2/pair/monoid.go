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

// Monoid creates a simple component-wise monoid for [Pair].
//
// This function creates a monoid that combines pairs by independently combining their
// head and tail components using the provided monoids. Both components are combined
// in NORMAL left-to-right order.
//
// IMPORTANT: This is DIFFERENT from [ApplicativeMonoidTail] and [ApplicativeMonoidHead],
// which use applicative functor operations and reverse the order of the non-focused component.
//
// Use this function when you want:
//   - Simple, predictable left-to-right combination for both components
//   - Behavior that matches intuition for non-commutative operations
//   - Direct component-wise combination without applicative functor semantics
//
// Use [ApplicativeMonoidTail] or [ApplicativeMonoidHead] when you need applicative
// functor semantics for lifting monoid operations into the Pair context.
//
// Parameters:
//   - l: A monoid for the head (left) values of type L
//   - r: A monoid for the tail (right) values of type R
//
// Returns:
//   - A Monoid[Pair[L, R]] that combines both components left-to-right
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
//	pairMonoid := pair.Monoid(intAdd, strConcat)
//
//	p1 := pair.MakePair(5, "hello")
//	p2 := pair.MakePair(10, " world")
//
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[int, string]{15, "hello world"}
//	// Both components combine left-to-right: (5+10, "hello"+" world")
//
//	empty := pairMonoid.Empty()
//	// empty is Pair[int, string]{0, ""}
//
// Comparison with ApplicativeMonoidTail:
//
//	strConcat := S.Monoid
//
//	// Simple component-wise monoid
//	simpleMonoid := pair.Monoid(strConcat, strConcat)
//	p1 := pair.MakePair("A", "1")
//	p2 := pair.MakePair("B", "2")
//	result1 := simpleMonoid.Concat(p1, p2)
//	// result1 is Pair[string, string]{"AB", "12"}
//	// Both components: left-to-right
//
//	// Applicative monoid
//	appMonoid := pair.ApplicativeMonoidTail(strConcat, strConcat)
//	result2 := appMonoid.Concat(p1, p2)
//	// result2 is Pair[string, string]{"BA", "12"}
//	// Head: reversed, Tail: normal
//
//go:inline
func Monoid[L, R any](l M.Monoid[L], r M.Monoid[R]) M.Monoid[Pair[L, R]] {
	return M.MakeMonoid(
		func(pl, pr Pair[L, R]) Pair[L, R] {
			return MakePair(l.Concat(Head(pl), Head(pr)), r.Concat(Tail(pl), Tail(pr)))
		},
		MakePair(l.Empty(), r.Empty()),
	)
}

// ApplicativeMonoid creates a monoid for [Pair] using applicative functor operations on the tail.
//
// This is an alias for [ApplicativeMonoidTail], which lifts the right (tail) monoid into the
// Pair applicative functor. The left monoid provides the semigroup for combining head values
// during applicative operations.
//
// IMPORTANT BEHAVIORAL NOTE: The applicative implementation causes the HEAD component to be
// combined in REVERSE order (right-to-left) while the TAIL combines normally (left-to-right).
// This differs from Haskell's standard Applicative instance for pairs, which combines the
// first component left-to-right. This matters for non-commutative operations like string
// concatenation.
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
// CRITICAL BEHAVIORAL NOTE: The HEAD values are combined in REVERSE order (right-to-left),
// while TAIL values combine in normal order (left-to-right). This is due to how the
// applicative `ap` operation is implemented for Pair.
//
// NOTE: This differs from Haskell's standard Applicative instance for (,) which combines
// the first component left-to-right. The reversal occurs because MonadApTail implements:
//
//	MakePair(sg.Concat(second.head, first.head), ...)
//
// Example showing the reversal with non-commutative operations:
//
//	strConcat := S.Monoid
//	pairMonoid := pair.ApplicativeMonoidTail(strConcat, strConcat)
//	p1 := pair.MakePair("hello", "foo")
//	p2 := pair.MakePair(" world", "bar")
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[string, string]{" worldhello", "foobar"}
//	//                                 ^^^^^^^^^^^^^^  ^^^^^^
//	//                                 REVERSED!       normal order
//
// In Haskell's Applicative for (,), this would give ("hellohello world", "foobar")
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
// CRITICAL BEHAVIORAL NOTE: The TAIL values are combined in REVERSE order (right-to-left),
// while HEAD values combine in normal order (left-to-right). This is the opposite behavior
// of ApplicativeMonoidTail. The reversal occurs because MonadApHead implements:
//
//	MakePair(..., sg.Concat(second.tail, first.tail))
//
// Example showing the reversal with non-commutative operations:
//
//	strConcat := S.Monoid
//	pairMonoid := pair.ApplicativeMonoidHead(strConcat, strConcat)
//	p1 := pair.MakePair("hello", "foo")
//	p2 := pair.MakePair(" world", "bar")
//	result := pairMonoid.Concat(p1, p2)
//	// result is Pair[string, string]{"hello world", "barfoo"}
//	//                                 ^^^^^^^^^^^^  ^^^^^^^^
//	//                                 normal order  REVERSED!
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
