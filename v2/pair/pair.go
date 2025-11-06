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
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/semigroup"
	"github.com/IBM/fp-go/v2/tuple"
)

// String prints some debug info for the object
//
//go:noinline
func pairString(s *pair) string {
	return fmt.Sprintf("Pair[%T, %T](%v, %v)", s.h, s.t, s.h, s.t)
}

// Format prints some debug info for the object
//
//go:noinline
func pairFormat(e *pair, f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, pairString(e))
	default:
		fmt.Fprint(f, pairString(e))
	}
}

// String prints some debug info for the object
func (s Pair[A, B]) String() string {
	return pairString((*pair)(&s))
}

// Format prints some debug info for the object
func (s Pair[A, B]) Format(f fmt.State, c rune) {
	pairFormat((*pair)(&s), f, c)
}

// Of creates a [Pair] with the same value in both the head and tail positions.
//
// Example:
//
//	p := pair.Of(42)  // Pair[int, int]{42, 42}
func Of[A any](value A) Pair[A, A] {
	return Pair[A, A]{h: value, t: value}
}

// FromTuple creates a [Pair] from a [tuple.Tuple2].
// The first element of the tuple becomes the head, and the second becomes the tail.
//
// Example:
//
//	t := tuple.MakeTuple2("hello", 42)
//	p := pair.FromTuple(t)  // Pair[string, int]{"hello", 42}
func FromTuple[A, B any](t tuple.Tuple2[A, B]) Pair[A, B] {
	return Pair[A, B]{h: t.F1, t: t.F2}
}

// ToTuple creates a [tuple.Tuple2] from a [Pair].
// The head becomes the first element, and the tail becomes the second element.
//
// Example:
//
//	p := pair.MakePair("hello", 42)
//	t := pair.ToTuple(p)  // tuple.Tuple2[string, int]{"hello", 42}
func ToTuple[A, B any](t Pair[A, B]) tuple.Tuple2[A, B] {
	return tuple.MakeTuple2(Head(t), Tail(t))
}

// MakePair creates a [Pair] from two values.
// The first value becomes the head, and the second becomes the tail.
//
// Example:
//
//	p := pair.MakePair("hello", 42)  // Pair[string, int]{"hello", 42}
func MakePair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{h: a, t: b}
}

// Head returns the head (first) value of the pair.
//
// Example:
//
//	p := pair.MakePair("hello", 42)
//	h := pair.Head(p)  // "hello"
func Head[A, B any](fa Pair[A, B]) A {
	return fa.h.(A)
}

// Tail returns the tail (second) value of the pair.
//
// Example:
//
//	p := pair.MakePair("hello", 42)
//	t := pair.Tail(p)  // 42
func Tail[A, B any](fa Pair[A, B]) B {
	return fa.t.(B)
}

// First returns the first value of the pair (alias for Head).
//
// Example:
//
//	p := pair.MakePair("hello", 42)
//	f := pair.First(p)  // "hello"
func First[A, B any](fa Pair[A, B]) A {
	return fa.h.(A)
}

// Second returns the second value of the pair (alias for Tail).
//
// Example:
//
//	p := pair.MakePair("hello", 42)
//	s := pair.Second(p)  // 42
func Second[A, B any](fa Pair[A, B]) B {
	return fa.t.(B)
}

// MonadMapHead maps a function over the head value of the pair, leaving the tail unchanged.
//
// Example:
//
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadMapHead(p, func(n int) string {
//	    return fmt.Sprintf("%d", n)
//	})  // Pair[string, string]{"5", "hello"}
func MonadMapHead[B, A, A1 any](fa Pair[A, B], f func(A) A1) Pair[A1, B] {
	return Pair[A1, B]{f(Head(fa)), fa.t}
}

// MonadMap maps a function over the head value of the pair (alias for MonadMapHead).
//
// Example:
//
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadMap(p, func(n int) string {
//	    return fmt.Sprintf("%d", n)
//	})  // Pair[string, string]{"5", "hello"}
func MonadMap[B, A, A1 any](fa Pair[A, B], f func(A) A1) Pair[A1, B] {
	return MonadMapHead(fa, f)
}

// MonadMapTail maps a function over the tail value of the pair, leaving the head unchanged.
//
// Example:
//
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadMapTail(p, func(s string) int {
//	    return len(s)
//	})  // Pair[int, int]{5, 5}
func MonadMapTail[A, B, B1 any](fa Pair[A, B], f func(B) B1) Pair[A, B1] {
	return Pair[A, B1]{fa.h, f(Tail(fa))}
}

// MonadBiMap maps functions over both the head and tail values of the pair.
//
// Example:
//
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadBiMap(p,
//	    func(n int) string { return fmt.Sprintf("%d", n) },
//	    func(s string) int { return len(s) },
//	)  // Pair[string, int]{"5", 5}
func MonadBiMap[A, B, A1, B1 any](fa Pair[A, B], f func(A) A1, g func(B) B1) Pair[A1, B1] {
	return Pair[A1, B1]{f(Head(fa)), g(Tail(fa))}
}

// Map returns a function that maps over the tail value of a pair (alias for MapTail).
// This is the curried version of MonadMapTail.
//
// Example:
//
//	mapper := pair.Map[int](func(s string) int { return len(s) })
//	p := pair.MakePair(5, "hello")
//	p2 := mapper(p)  // Pair[int, int]{5, 5}
func Map[A, B, B1 any](f func(B) B1) func(Pair[A, B]) Pair[A, B1] {
	return MapTail[A](f)
}

// MapHead returns a function that maps over the head value of a pair.
// This is the curried version of MonadMapHead.
//
// Example:
//
//	mapper := pair.MapHead[string](func(n int) string {
//	    return fmt.Sprintf("%d", n)
//	})
//	p := pair.MakePair(5, "hello")
//	p2 := mapper(p)  // Pair[string, string]{"5", "hello"}
func MapHead[B, A, A1 any](f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return F.Bind2nd(MonadMapHead[B, A, A1], f)
}

// MapTail returns a function that maps over the tail value of a pair.
// This is the curried version of MonadMapTail.
//
// Example:
//
//	mapper := pair.MapTail[int](func(s string) int { return len(s) })
//	p := pair.MakePair(5, "hello")
//	p2 := mapper(p)  // Pair[int, int]{5, 5}
func MapTail[A, B, B1 any](f func(B) B1) func(Pair[A, B]) Pair[A, B1] {
	return F.Bind2nd(MonadMapTail[A, B, B1], f)
}

// BiMap returns a function that maps over both values of a pair.
// This is the curried version of MonadBiMap.
//
// Example:
//
//	mapper := pair.BiMap(
//	    func(n int) string { return fmt.Sprintf("%d", n) },
//	    func(s string) int { return len(s) },
//	)
//	p := pair.MakePair(5, "hello")
//	p2 := mapper(p)  // Pair[string, int]{"5", 5}
func BiMap[A, B, A1, B1 any](f func(A) A1, g func(B) B1) func(Pair[A, B]) Pair[A1, B1] {
	return func(fa Pair[A, B]) Pair[A1, B1] {
		return MonadBiMap(fa, f, g)
	}
}

// MonadChainHead chains a function over the head value, combining tail values using a semigroup.
// The function receives the head value and returns a new pair. The tail values from both pairs
// are combined using the provided semigroup.
//
// Example:
//
//	import SG "github.com/IBM/fp-go/v2/semigroup"
//
//	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadChainHead(strConcat, p, func(n int) pair.Pair[string, string] {
//	    return pair.MakePair(fmt.Sprintf("%d", n), "!")
//	})  // Pair[string, string]{"5", "hello!"}
func MonadChainHead[B, A, A1 any](sg semigroup.Semigroup[B], fa Pair[A, B], f func(A) Pair[A1, B]) Pair[A1, B] {
	fb := f(Head(fa))
	return Pair[A1, B]{fb.h, sg.Concat(Tail(fa), Tail(fb))}
}

// MonadChainTail chains a function over the tail value, combining head values using a semigroup.
// The function receives the tail value and returns a new pair. The head values from both pairs
// are combined using the provided semigroup.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadChainTail(intSum, p, func(s string) pair.Pair[int, int] {
//	    return pair.MakePair(len(s), len(s) * 2)
//	})  // Pair[int, int]{10, 10}
func MonadChainTail[A, B, B1 any](sg semigroup.Semigroup[A], fb Pair[A, B], f func(B) Pair[A, B1]) Pair[A, B1] {
	fa := f(Tail(fb))
	return Pair[A, B1]{sg.Concat(Head(fb), Head(fa)), fa.t}
}

// MonadChain chains a function over the tail value (alias for MonadChainTail).
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	p := pair.MakePair(5, "hello")
//	p2 := pair.MonadChain(intSum, p, func(s string) pair.Pair[int, int] {
//	    return pair.MakePair(len(s), len(s) * 2)
//	})  // Pair[int, int]{10, 10}
func MonadChain[A, B, B1 any](sg semigroup.Semigroup[A], fa Pair[A, B], f func(B) Pair[A, B1]) Pair[A, B1] {
	return MonadChainTail(sg, fa, f)
}

// ChainHead returns a function that chains over the head value.
// This is the curried version of MonadChainHead.
//
// Example:
//
//	import SG "github.com/IBM/fp-go/v2/semigroup"
//
//	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
//	chain := pair.ChainHead(strConcat, func(n int) pair.Pair[string, string] {
//	    return pair.MakePair(fmt.Sprintf("%d", n), "!")
//	})
//	p := pair.MakePair(5, "hello")
//	p2 := chain(p)  // Pair[string, string]{"5", "hello!"}
func ChainHead[B, A, A1 any](sg semigroup.Semigroup[B], f func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B] {
	return func(fa Pair[A, B]) Pair[A1, B] {
		return MonadChainHead(sg, fa, f)
	}
}

// ChainTail returns a function that chains over the tail value.
// This is the curried version of MonadChainTail.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	chain := pair.ChainTail(intSum, func(s string) pair.Pair[int, int] {
//	    return pair.MakePair(len(s), len(s) * 2)
//	})
//	p := pair.MakePair(5, "hello")
//	p2 := chain(p)  // Pair[int, int]{10, 10}
func ChainTail[A, B, B1 any](sg semigroup.Semigroup[A], f func(B) Pair[A, B1]) func(Pair[A, B]) Pair[A, B1] {
	return func(fa Pair[A, B]) Pair[A, B1] {
		return MonadChainTail(sg, fa, f)
	}
}

// Chain returns a function that chains over the tail value (alias for ChainTail).
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	chain := pair.Chain(intSum, func(s string) pair.Pair[int, int] {
//	    return pair.MakePair(len(s), len(s) * 2)
//	})
//	p := pair.MakePair(5, "hello")
//	p2 := chain(p)  // Pair[int, int]{10, 10}
func Chain[A, B, B1 any](sg semigroup.Semigroup[A], f func(B) Pair[A, B1]) func(Pair[A, B]) Pair[A, B1] {
	return ChainTail(sg, f)
}

// MonadApHead applies a function wrapped in a pair to a value wrapped in a pair,
// operating on the head values and combining tail values using a semigroup.
//
// Example:
//
//	import SG "github.com/IBM/fp-go/v2/semigroup"
//
//	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
//	pf := pair.MakePair(func(n int) string { return fmt.Sprintf("%d", n) }, "!")
//	pv := pair.MakePair(42, "hello")
//	result := pair.MonadApHead(strConcat, pf, pv)  // Pair[string, string]{"42", "!hello"}
func MonadApHead[B, A, A1 any](sg semigroup.Semigroup[B], faa Pair[func(A) A1, B], fa Pair[A, B]) Pair[A1, B] {
	return Pair[A1, B]{Head(faa)(Head(fa)), sg.Concat(Tail(fa), Tail(faa))}
}

// MonadApTail applies a function wrapped in a pair to a value wrapped in a pair,
// operating on the tail values and combining head values using a semigroup.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	pf := pair.MakePair(10, func(s string) int { return len(s) })
//	pv := pair.MakePair(5, "hello")
//	result := pair.MonadApTail(intSum, pf, pv)  // Pair[int, int]{15, 5}
func MonadApTail[A, B, B1 any](sg semigroup.Semigroup[A], fbb Pair[A, func(B) B1], fb Pair[A, B]) Pair[A, B1] {
	return Pair[A, B1]{sg.Concat(Head(fb), Head(fbb)), Tail(fbb)(Tail(fb))}
}

// MonadAp applies a function wrapped in a pair to a value wrapped in a pair,
// operating on the tail values (alias for MonadApTail).
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	pf := pair.MakePair(10, func(s string) int { return len(s) })
//	pv := pair.MakePair(5, "hello")
//	result := pair.MonadAp(intSum, pf, pv)  // Pair[int, int]{15, 5}
func MonadAp[A, B, B1 any](sg semigroup.Semigroup[A], faa Pair[A, func(B) B1], fa Pair[A, B]) Pair[A, B1] {
	return MonadApTail(sg, faa, fa)
}

// ApHead returns a function that applies a function in a pair to a value in a pair,
// operating on head values. This is the curried version of MonadApHead.
//
// Example:
//
//	import SG "github.com/IBM/fp-go/v2/semigroup"
//
//	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
//	pv := pair.MakePair(42, "hello")
//	ap := pair.ApHead(strConcat, pv)
//	pf := pair.MakePair(func(n int) string { return fmt.Sprintf("%d", n) }, "!")
//	result := ap(pf)  // Pair[string, string]{"42", "!hello"}
func ApHead[B, A, A1 any](sg semigroup.Semigroup[B], fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return func(faa Pair[func(A) A1, B]) Pair[A1, B] {
		return MonadApHead(sg, faa, fa)
	}
}

// ApTail returns a function that applies a function in a pair to a value in a pair,
// operating on tail values. This is the curried version of MonadApTail.
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	pv := pair.MakePair(5, "hello")
//	ap := pair.ApTail(intSum, pv)
//	pf := pair.MakePair(10, func(s string) int { return len(s) })
//	result := ap(pf)  // Pair[int, int]{15, 5}
func ApTail[A, B, B1 any](sg semigroup.Semigroup[A], fb Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1] {
	return func(fbb Pair[A, func(B) B1]) Pair[A, B1] {
		return MonadApTail(sg, fbb, fb)
	}
}

// Ap returns a function that applies a function in a pair to a value in a pair,
// operating on tail values (alias for ApTail).
//
// Example:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	intSum := N.SemigroupSum[int]()
//	pv := pair.MakePair(5, "hello")
//	ap := pair.Ap(intSum, pv)
//	pf := pair.MakePair(10, func(s string) int { return len(s) })
//	result := ap(pf)  // Pair[int, int]{15, 5}
func Ap[A, B, B1 any](sg semigroup.Semigroup[A], fa Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1] {
	return ApTail[A, B, B1](sg, fa)
}

// Swap swaps the head and tail values of a pair.
//
// Example:
//
//	p := pair.MakePair("hello", 42)
//	swapped := pair.Swap(p)  // Pair[int, string]{42, "hello"}
func Swap[A, B any](fa Pair[A, B]) Pair[B, A] {
	return MakePair(Tail(fa), Head(fa))
}

// Paired converts a function with 2 parameters into a function taking a [Pair].
// The inverse function is [Unpaired].
//
// Example:
//
//	add := func(a, b int) int { return a + b }
//	pairedAdd := pair.Paired(add)
//	result := pairedAdd(pair.MakePair(3, 4))  // 7
func Paired[F ~func(T1, T2) R, T1, T2, R any](f F) func(Pair[T1, T2]) R {
	return func(t Pair[T1, T2]) R {
		return f(Head(t), Tail(t))
	}
}

// Unpaired converts a function with a [Pair] parameter into a function with 2 parameters.
// The inverse function is [Paired].
//
// Example:
//
//	pairedAdd := func(p pair.Pair[int, int]) int {
//	    return pair.Head(p) + pair.Tail(p)
//	}
//	add := pair.Unpaired(pairedAdd)
//	result := add(3, 4)  // 7
func Unpaired[F ~func(Pair[T1, T2]) R, T1, T2, R any](f F) func(T1, T2) R {
	return func(t1 T1, t2 T2) R {
		return f(MakePair(t1, t2))
	}
}

// Merge applies a curried function to a pair by applying the tail value first, then the head value.
//
// Example:
//
//	add := func(b int) func(a int) int {
//	    return func(a int) int { return a + b }
//	}
//	merge := pair.Merge(add)
//	result := merge(pair.MakePair(3, 4))  // 7 (applies 4 then 3)
func Merge[F ~func(B) func(A) R, A, B, R any](f F) func(Pair[A, B]) R {
	return func(p Pair[A, B]) R {
		return f(Tail(p))(Head(p))
	}
}
