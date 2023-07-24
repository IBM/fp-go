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
	N "github.com/IBM/fp-go/number/integer"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

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

// Reduce applies a function for each value of the iterator with a floating result
func Reduce[GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](f func(V, U) V, initial V) func(GU) V {
	return func(as GU) V {
		next, ok := O.Unwrap(as())
		current := initial
		for ok {
			// next (with bad side effect)
			current = f(current, next.F2)
			next, ok = O.Unwrap(next.F1())
		}
		return current
	}
}

// ToArray converts the iterator to an array
func ToArray[GU ~func() O.Option[T.Tuple2[GU, U]], US ~[]U, U any](u GU) US {
	return Reduce[GU](A.Append[US], A.Empty[US]())(u)
}

func Map[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], U, V any](f func(U) V) func(ma GU) GV {
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

// MakeBy returns an [Iterator] with `n` elements initialized with `f(i)`
func MakeBy[GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(int) U, U any](n int, f FCT) GU {

	var m func(int) O.Option[T.Tuple2[GU, U]]

	recurse := func(i int) GU {
		return func() O.Option[T.Tuple2[GU, U]] {
			return F.Pipe1(
				i,
				m,
			)
		}
	}

	m = F.Flow2(
		O.FromPredicate(N.Between(0, n)),
		O.Map(F.Flow2(
			T.Replicate2[int],
			T.Map2(F.Flow2(
				utils.Inc,
				recurse),
				f),
		)),
	)

	return recurse(0)
}

// Replicate creates an [Iterator] containing a value repeated the specified number of times.
func Replicate[GU ~func() O.Option[T.Tuple2[GU, U]], U any](n int, a U) GU {
	return MakeBy[GU](n, F.Constant1[int](a))
}
