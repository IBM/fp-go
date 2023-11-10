// Copyright (c) 2023 IBM Corp.
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
	A "github.com/IBM/fp-go/array/generic"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	IO "github.com/IBM/fp-go/iooption/generic"
	M "github.com/IBM/fp-go/monoid"
	N "github.com/IBM/fp-go/number"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// Next returns the iterator for the next element in an iterator `T.Tuple2`
func Next[GU ~func() O.Option[T.Tuple2[GU, U]], U any](m T.Tuple2[GU, U]) GU {
	return T.First(m)
}

// Current returns the current element in an iterator `T.Tuple2`
func Current[GU ~func() O.Option[T.Tuple2[GU, U]], U any](m T.Tuple2[GU, U]) U {
	return T.Second(m)
}

// From constructs an array from a set of variadic arguments
func From[GU ~func() O.Option[T.Tuple2[GU, U]], U any](data ...U) GU {
	return FromArray[GU](data)
}

// Empty returns the empty iterator
func Empty[GU ~func() O.Option[T.Tuple2[GU, U]], U any]() GU {
	return IO.None[GU]()
}

// Of returns an iterator with one single element
func Of[GU ~func() O.Option[T.Tuple2[GU, U]], U any](a U) GU {
	return IO.Of[GU](T.MakeTuple2(Empty[GU](), a))
}

// FromArray returns an iterator from multiple elements
func FromArray[GU ~func() O.Option[T.Tuple2[GU, U]], US ~[]U, U any](as US) GU {
	return A.MatchLeft(Empty[GU], func(head U, tail US) GU {
		return func() O.Option[T.Tuple2[GU, U]] {
			return O.Of(T.MakeTuple2(FromArray[GU](tail), head))
		}
	})(as)
}

// reduce applies a function for each value of the iterator with a floating result
func reduce[GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](as GU, f func(V, U) V, initial V) V {
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
func Reduce[GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](f func(V, U) V, initial V) func(GU) V {
	return F.Bind23of3(reduce[GU, U, V])(f, initial)
}

// ToArray converts the iterator to an array
func ToArray[GU ~func() O.Option[T.Tuple2[GU, U]], US ~[]U, U any](u GU) US {
	return Reduce[GU](A.Append[US], A.Empty[US]())(u)
}

func Map[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(U) V, U, V any](f FCT) func(ma GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(O.Option[T.Tuple2[GU, U]]) O.Option[T.Tuple2[GV, V]]

	recurse := func(ma GU) GV {
		return F.Nullary2(
			ma,
			m,
		)
	}

	m = O.Map(T.Map2(recurse, f))

	return recurse
}

func MonadMap[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](ma GU, f func(U) V) GV {
	return Map[GV, GU](f)(ma)
}

func concat[GU ~func() O.Option[T.Tuple2[GU, U]], U any](right, left GU) GU {
	var m func(ma O.Option[T.Tuple2[GU, U]]) O.Option[T.Tuple2[GU, U]]

	recurse := func(left GU) GU {
		return F.Nullary2(left, m)
	}

	m = O.Fold(
		right,
		F.Flow2(
			T.Map2(recurse, F.Identity[U]),
			O.Some[T.Tuple2[GU, U]],
		))

	return recurse(left)
}

func Chain[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](f func(U) GV) func(GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(O.Option[T.Tuple2[GU, U]]) O.Option[T.Tuple2[GV, V]]

	recurse := func(ma GU) GV {
		return F.Nullary2(
			ma,
			m,
		)
	}
	m = O.Chain(
		F.Flow3(
			T.Map2(recurse, f),
			T.Tupled2(concat[GV]),
			func(v GV) O.Option[T.Tuple2[GV, V]] {
				return v()
			},
		),
	)

	return recurse
}

func MonadChain[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](ma GU, f func(U) GV) GV {
	return Chain[GV, GU](f)(ma)
}

func Flatten[GV ~func() O.Option[T.Tuple2[GV, GU]], GU ~func() O.Option[T.Tuple2[GU, U]], U any](ma GV) GU {
	return MonadChain(ma, F.Identity[GU])
}

// MakeBy returns an [Iterator] with an infinite number of elements initialized with `f(i)`
func MakeBy[GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(int) U, U any](f FCT) GU {

	var m func(int) O.Option[T.Tuple2[GU, U]]

	recurse := func(i int) GU {
		return F.Nullary2(
			F.Constant(i),
			m,
		)
	}

	m = F.Flow3(
		T.Replicate2[int],
		T.Map2(F.Flow2(
			utils.Inc,
			recurse),
			f),
		O.Of[T.Tuple2[GU, U]],
	)

	// bootstrap
	return recurse(0)
}

// Replicate creates an infinite [Iterator] containing a value.
func Replicate[GU ~func() O.Option[T.Tuple2[GU, U]], U any](a U) GU {
	return MakeBy[GU](F.Constant1[int](a))
}

// Repeat creates an [Iterator] containing a value repeated the specified number of times.
// Alias of [Replicate] combined with [Take]
func Repeat[GU ~func() O.Option[T.Tuple2[GU, U]], U any](n int, a U) GU {
	return F.Pipe2(
		a,
		Replicate[GU],
		Take[GU](n),
	)
}

// Count creates an [Iterator] containing a consecutive sequence of integers starting with the provided start value
func Count[GU ~func() O.Option[T.Tuple2[GU, int]]](start int) GU {
	return MakeBy[GU](N.Add(start))
}

func FilterMap[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(U) O.Option[V], U, V any](f FCT) func(ma GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(O.Option[T.Tuple2[GU, U]]) O.Option[T.Tuple2[GV, V]]

	recurse := func(ma GU) GV {
		return F.Nullary2(
			ma,
			m,
		)
	}

	m = O.Fold(
		Empty[GV](),
		func(t T.Tuple2[GU, U]) O.Option[T.Tuple2[GV, V]] {
			r := recurse(Next(t))
			return O.MonadFold(f(Current(t)), r, F.Flow2(
				F.Bind1st(T.MakeTuple2[GV, V], r),
				O.Some[T.Tuple2[GV, V]],
			))
		},
	)

	return recurse
}

func Filter[GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(U) bool, U any](f FCT) func(ma GU) GU {
	return FilterMap[GU, GU](O.FromPredicate(f))
}

func Ap[GUV ~func() O.Option[T.Tuple2[GUV, func(U) V]], GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](ma GU) func(fab GUV) GV {
	return Chain[GV, GUV](F.Bind1st(MonadMap[GV, GU], ma))
}

func MonadAp[GUV ~func() O.Option[T.Tuple2[GUV, func(U) V]], GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](fab GUV, ma GU) GV {
	return Ap[GUV, GV, GU](ma)(fab)
}

func FilterChain[GVV ~func() O.Option[T.Tuple2[GVV, GV]], GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(U) O.Option[GV], U, V any](f FCT) func(ma GU) GV {
	return F.Flow2(
		FilterMap[GVV, GU](f),
		Flatten[GVV],
	)
}

func FoldMap[GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(U) V, U, V any](m M.Monoid[V]) func(FCT) func(ma GU) V {
	return func(f FCT) func(ma GU) V {
		return Reduce[GU](func(cur V, a U) V {
			return m.Concat(cur, f(a))
		}, m.Empty())
	}
}

func Fold[GU ~func() O.Option[T.Tuple2[GU, U]], U any](m M.Monoid[U]) func(ma GU) U {
	return Reduce[GU](m.Concat, m.Empty())
}
