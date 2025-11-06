// Copyright (c) 2024 IBM Corp.
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

	F "github.com/IBM/fp-go/function"
	Sg "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
)

type (
	pair struct {
		h, t any
	}

	// Pair defines a data structure that holds two strongly typed values
	Pair[A, B any] pair
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

// Of creates a [Pair] with the same value to to both fields
func Of[A any](value A) Pair[A, A] {
	return Pair[A, A]{h: value, t: value}
}

// FromTuple creates a [Pair] from a [T.Tuple2]
func FromTuple[A, B any](t T.Tuple2[A, B]) Pair[A, B] {
	return Pair[A, B]{h: t.F1, t: t.F2}
}

// ToTuple creates a [T.Tuple2] from a [Pair]
func ToTuple[A, B any](t Pair[A, B]) T.Tuple2[A, B] {
	return T.MakeTuple2(Head(t), Tail(t))
}

// MakePair creates a [Pair] from two values
func MakePair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{h: a, t: b}
}

// Head returns the head value of the pair
func Head[A, B any](fa Pair[A, B]) A {
	return fa.h.(A)
}

// Tail returns the head value of the pair
func Tail[A, B any](fa Pair[A, B]) B {
	return fa.t.(B)
}

// MonadMapHead maps the head value
func MonadMapHead[B, A, A1 any](fa Pair[A, B], f func(A) A1) Pair[A1, B] {
	return Pair[A1, B]{f(Head(fa)), fa.t}
}

// MonadMap maps the head value
func MonadMap[B, A, A1 any](fa Pair[A, B], f func(A) A1) Pair[A1, B] {
	return MonadMapHead(fa, f)
}

// MonadMapTail maps the Tail value
func MonadMapTail[A, B, B1 any](fa Pair[A, B], f func(B) B1) Pair[A, B1] {
	return Pair[A, B1]{fa.h, f(Tail(fa))}
}

// MonadBiMap maps both values
func MonadBiMap[A, B, A1, B1 any](fa Pair[A, B], f func(A) A1, g func(B) B1) Pair[A1, B1] {
	return Pair[A1, B1]{f(Head(fa)), g(Tail(fa))}
}

// Map maps the head value
func Map[B, A, A1 any](f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return MapHead[B, A, A1](f)
}

// MapHead maps the head value
func MapHead[B, A, A1 any](f func(A) A1) func(Pair[A, B]) Pair[A1, B] {
	return F.Bind2nd(MonadMapHead[B, A, A1], f)
}

// MapTail maps the Tail value
func MapTail[A, B, B1 any](f func(B) B1) func(Pair[A, B]) Pair[A, B1] {
	return F.Bind2nd(MonadMapTail[A, B, B1], f)
}

// BiMap maps both values
func BiMap[A, B, A1, B1 any](f func(A) A1, g func(B) B1) func(Pair[A, B]) Pair[A1, B1] {
	return func(fa Pair[A, B]) Pair[A1, B1] {
		return MonadBiMap(fa, f, g)
	}
}

// MonadChainHead chains on the head value
func MonadChainHead[B, A, A1 any](sg Sg.Semigroup[B], fa Pair[A, B], f func(A) Pair[A1, B]) Pair[A1, B] {
	fb := f(Head(fa))
	return Pair[A1, B]{fb.h, sg.Concat(Tail(fa), Tail(fb))}
}

// MonadChainTail chains on the Tail value
func MonadChainTail[A, B, B1 any](sg Sg.Semigroup[A], fb Pair[A, B], f func(B) Pair[A, B1]) Pair[A, B1] {
	fa := f(Tail(fb))
	return Pair[A, B1]{sg.Concat(Head(fb), Head(fa)), fa.t}
}

// MonadChain chains on the head value
func MonadChain[B, A, A1 any](sg Sg.Semigroup[B], fa Pair[A, B], f func(A) Pair[A1, B]) Pair[A1, B] {
	return MonadChainHead(sg, fa, f)
}

// ChainHead chains on the head value
func ChainHead[B, A, A1 any](sg Sg.Semigroup[B], f func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B] {
	return func(fa Pair[A, B]) Pair[A1, B] {
		return MonadChainHead(sg, fa, f)
	}
}

// ChainTail chains on the Tail value
func ChainTail[A, B, B1 any](sg Sg.Semigroup[A], f func(B) Pair[A, B1]) func(Pair[A, B]) Pair[A, B1] {
	return func(fa Pair[A, B]) Pair[A, B1] {
		return MonadChainTail(sg, fa, f)
	}
}

// Chain chains on the head value
func Chain[B, A, A1 any](sg Sg.Semigroup[B], f func(A) Pair[A1, B]) func(Pair[A, B]) Pair[A1, B] {
	return ChainHead[B, A, A1](sg, f)
}

// MonadApHead applies on the head value
func MonadApHead[B, A, A1 any](sg Sg.Semigroup[B], faa Pair[func(A) A1, B], fa Pair[A, B]) Pair[A1, B] {
	return Pair[A1, B]{Head(faa)(Head(fa)), sg.Concat(Tail(fa), Tail(faa))}
}

// MonadApTail applies on the Tail value
func MonadApTail[A, B, B1 any](sg Sg.Semigroup[A], fbb Pair[A, func(B) B1], fb Pair[A, B]) Pair[A, B1] {
	return Pair[A, B1]{sg.Concat(Head(fb), Head(fbb)), Tail(fbb)(Tail(fb))}
}

// MonadAp applies on the head value
func MonadAp[B, A, A1 any](sg Sg.Semigroup[B], faa Pair[func(A) A1, B], fa Pair[A, B]) Pair[A1, B] {
	return MonadApHead(sg, faa, fa)
}

// ApHead applies on the head value
func ApHead[B, A, A1 any](sg Sg.Semigroup[B], fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return func(faa Pair[func(A) A1, B]) Pair[A1, B] {
		return MonadApHead(sg, faa, fa)
	}
}

// ApTail applies on the Tail value
func ApTail[A, B, B1 any](sg Sg.Semigroup[A], fb Pair[A, B]) func(Pair[A, func(B) B1]) Pair[A, B1] {
	return func(fbb Pair[A, func(B) B1]) Pair[A, B1] {
		return MonadApTail(sg, fbb, fb)
	}
}

// Ap applies on the head value
func Ap[B, A, A1 any](sg Sg.Semigroup[B], fa Pair[A, B]) func(Pair[func(A) A1, B]) Pair[A1, B] {
	return ApHead[B, A, A1](sg, fa)
}

// Swap swaps the two channels
func Swap[A, B any](fa Pair[A, B]) Pair[B, A] {
	return MakePair(Tail(fa), Head(fa))
}

// Paired converts a function with 2 parameters into a function taking a [Pair]
// The inverse function is [Unpaired]
func Paired[F ~func(T1, T2) R, T1, T2, R any](f F) func(Pair[T1, T2]) R {
	return func(t Pair[T1, T2]) R {
		return f(Head(t), Tail(t))
	}
}

// Unpaired converts a function with a [Pair] parameter into a function with 2 parameters
// The inverse function is [Paired]
func Unpaired[F ~func(Pair[T1, T2]) R, T1, T2, R any](f F) func(T1, T2) R {
	return func(t1 T1, t2 T2) R {
		return f(MakePair(t1, t2))
	}
}

func Merge[F ~func(B) func(A) R, A, B, R any](f F) func(Pair[A, B]) R {
	return func(p Pair[A, B]) R {
		return f(Tail(p))(Head(p))
	}
}
