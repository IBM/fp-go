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

package generic

import (
	A "github.com/IBM/fp-go/v2/array/generic"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/utils"
	IO "github.com/IBM/fp-go/v2/iooption/generic"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

// Next returns the iterator for the next element in an iterator `Pair`
func Next[GU ~func() Option[Pair[GU, U]], U any](m Pair[GU, U]) GU {
	return P.Head(m)
}

// Current returns the current element in an iterator `Pair`
func Current[GU ~func() Option[Pair[GU, U]], U any](m Pair[GU, U]) U {
	return P.Tail(m)
}

// From constructs an array from a set of variadic arguments
func From[GU ~func() Option[Pair[GU, U]], U any](data ...U) GU {
	return FromArray[GU](data)
}

// Empty returns the empty iterator
func Empty[GU ~func() Option[Pair[GU, U]], U any]() GU {
	return IO.None[GU]()
}

// Of returns an iterator with one single element
func Of[GU ~func() Option[Pair[GU, U]], U any](a U) GU {
	return IO.Of[GU](P.MakePair(Empty[GU](), a))
}

// FromArray returns an iterator from multiple elements
func FromArray[GU ~func() Option[Pair[GU, U]], US ~[]U, U any](as US) GU {
	return A.MatchLeft(Empty[GU], func(head U, tail US) GU {
		return func() Option[Pair[GU, U]] {
			return O.Of(P.MakePair(FromArray[GU](tail), head))
		}
	})(as)
}

// reduce applies a function for each value of the iterator with a floating result
func reduce[GU ~func() Option[Pair[GU, U]], U, V any](as GU, f func(V, U) V, initial V) V {
	next, ok := O.Unwrap(as())
	current := initial
	for ok {
		// next (with bad side effect)
		current = f(current, Current(next))
		next, ok = O.Unwrap(Next(next)())
	}
	return current
}

// Reduce applies a function for each value of the iterator with a floating result
func Reduce[GU ~func() Option[Pair[GU, U]], U, V any](f func(V, U) V, initial V) func(GU) V {
	return F.Bind23of3(reduce[GU, U, V])(f, initial)
}

// ToArray converts the iterator to an array
func ToArray[GU ~func() Option[Pair[GU, U]], US ~[]U, U any](u GU) US {
	return Reduce[GU](A.Append[US], A.Empty[US]())(u)
}

func Map[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], FCT ~func(U) V, U, V any](f FCT) func(ma GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(Option[Pair[GU, U]]) Option[Pair[GV, V]]

	recurse := func(ma GU) GV {
		return F.Nullary2(
			ma,
			m,
		)
	}

	m = O.Map(P.BiMap(recurse, f))

	return recurse
}

func MonadMap[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](ma GU, f func(U) V) GV {
	return Map[GV, GU](f)(ma)
}

func concat[GU ~func() Option[Pair[GU, U]], U any](right, left GU) GU {
	var m func(ma Option[Pair[GU, U]]) Option[Pair[GU, U]]

	recurse := func(left GU) GU {
		return F.Nullary2(left, m)
	}

	m = O.Fold(
		right,
		F.Flow2(
			P.BiMap(recurse, F.Identity[U]),
			O.Some[Pair[GU, U]],
		))

	return recurse(left)
}

func Chain[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](f func(U) GV) func(GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(Option[Pair[GU, U]]) Option[Pair[GV, V]]

	recurse := func(ma GU) GV {
		return F.Nullary2(
			ma,
			m,
		)
	}
	m = O.Chain(
		F.Flow3(
			P.BiMap(recurse, f),
			P.Paired(concat[GV]),
			func(v GV) Option[Pair[GV, V]] {
				return v()
			},
		),
	)

	return recurse
}

func MonadChain[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](ma GU, f func(U) GV) GV {
	return Chain[GV, GU](f)(ma)
}

func MonadChainFirst[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](ma GU, f func(U) GV) GU {
	return C.MonadChainFirst(
		MonadChain[GU, GU, U, U],
		MonadMap[GU, GV, V, U],
		ma,
		f,
	)
}

func ChainFirst[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](f func(U) GV) func(GU) GU {
	return C.ChainFirst(
		Chain[GU, GU, U, U],
		Map[GU, GV, func(V) U, V, U],
		f,
	)
}

func Flatten[GV ~func() Option[Pair[GV, GU]], GU ~func() Option[Pair[GU, U]], U any](ma GV) GU {
	return MonadChain(ma, F.Identity[GU])
}

// MakeBy returns an [Iterator] with an infinite number of elements initialized with `f(i)`
func MakeBy[GU ~func() Option[Pair[GU, U]], FCT ~func(int) U, U any](f FCT) GU {

	var m func(int) Option[Pair[GU, U]]

	recurse := func(i int) GU {
		return F.Nullary2(
			F.Constant(i),
			m,
		)
	}

	m = F.Flow3(
		P.Of[int],
		P.BiMap(F.Flow2(
			utils.Inc,
			recurse),
			f),
		O.Of[Pair[GU, U]],
	)

	// bootstrap
	return recurse(0)
}

// Replicate creates an infinite [Iterator] containing a value.
func Replicate[GU ~func() Option[Pair[GU, U]], U any](a U) GU {
	return MakeBy[GU](F.Constant1[int](a))
}

// Repeat creates an [Iterator] containing a value repeated the specified number of times.
// Alias of [Replicate] combined with [Take]
func Repeat[GU ~func() Option[Pair[GU, U]], U any](n int, a U) GU {
	return F.Pipe2(
		a,
		Replicate[GU],
		Take[GU](n),
	)
}

// Count creates an [Iterator] containing a consecutive sequence of integers starting with the provided start value
func Count[GU ~func() Option[Pair[GU, int]]](start int) GU {
	return MakeBy[GU](N.Add(start))
}

func FilterMap[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], FCT ~func(U) Option[V], U, V any](f FCT) func(ma GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(Option[Pair[GU, U]]) Option[Pair[GV, V]]

	recurse := func(ma GU) GV {
		return F.Nullary2(
			ma,
			m,
		)
	}

	m = O.Fold(
		Empty[GV](),
		func(t Pair[GU, U]) Option[Pair[GV, V]] {
			r := recurse(Next(t))
			return O.MonadFold(f(Current(t)), r, F.Flow2(
				F.Bind1st(P.MakePair[GV, V], r),
				O.Some[Pair[GV, V]],
			))
		},
	)

	return recurse
}

func Filter[GU ~func() Option[Pair[GU, U]], FCT ~Predicate[U], U any](f FCT) func(ma GU) GU {
	return FilterMap[GU, GU](O.FromPredicate(f))
}

func Ap[GUV ~func() Option[Pair[GUV, func(U) V]], GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](ma GU) func(fab GUV) GV {
	return Chain[GV, GUV](F.Bind1st(MonadMap[GV, GU], ma))
}

func MonadAp[GUV ~func() Option[Pair[GUV, func(U) V]], GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], U, V any](fab GUV, ma GU) GV {
	return Ap[GUV, GV](ma)(fab)
}

func FilterChain[GVV ~func() Option[Pair[GVV, GV]], GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], FCT ~func(U) Option[GV], U, V any](f FCT) func(ma GU) GV {
	return F.Flow2(
		FilterMap[GVV, GU](f),
		Flatten[GVV],
	)
}

func FoldMap[GU ~func() Option[Pair[GU, U]], FCT ~func(U) V, U, V any](m M.Monoid[V]) func(FCT) func(ma GU) V {
	return func(f FCT) func(ma GU) V {
		return Reduce[GU](func(cur V, a U) V {
			return m.Concat(cur, f(a))
		}, m.Empty())
	}
}

func Fold[GU ~func() Option[Pair[GU, U]], U any](m M.Monoid[U]) func(ma GU) U {
	return Reduce[GU](m.Concat, m.Empty())
}
