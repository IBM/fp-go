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

package stateless

import (
	G "github.com/IBM/fp-go/iterator/stateless/generic"
	L "github.com/IBM/fp-go/lazy"
	M "github.com/IBM/fp-go/monoid"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// Iterator represents a stateless, pure way to iterate over a sequence
type Iterator[U any] L.Lazy[O.Option[T.Tuple2[Iterator[U], U]]]

// Next returns the [Iterator] for the next element in an iterator `T.Tuple2`
func Next[U any](m T.Tuple2[Iterator[U], U]) Iterator[U] {
	return G.Next(m)
}

// Current returns the current element in an [Iterator] `T.Tuple2`
func Current[U any](m T.Tuple2[Iterator[U], U]) U {
	return G.Current(m)
}

// Empty returns the empty iterator
func Empty[U any]() Iterator[U] {
	return G.Empty[Iterator[U]]()
}

// Of returns an iterator with one single element
func Of[U any](a U) Iterator[U] {
	return G.Of[Iterator[U]](a)
}

// FromArray returns an iterator from multiple elements
func FromArray[U any](as []U) Iterator[U] {
	return G.FromArray[Iterator[U]](as)
}

// ToArray converts the iterator to an array
func ToArray[U any](u Iterator[U]) []U {
	return G.ToArray[Iterator[U], []U](u)
}

// Reduce applies a function for each value of the iterator with a floating result
func Reduce[U, V any](f func(V, U) V, initial V) func(Iterator[U]) V {
	return G.Reduce[Iterator[U]](f, initial)
}

// MonadMap transforms an [Iterator] of type [U] into an [Iterator] of type [V] via a mapping function
func MonadMap[U, V any](ma Iterator[U], f func(U) V) Iterator[V] {
	return G.MonadMap[Iterator[V], Iterator[U]](ma, f)
}

// Map transforms an [Iterator] of type [U] into an [Iterator] of type [V] via a mapping function
func Map[U, V any](f func(U) V) func(ma Iterator[U]) Iterator[V] {
	return G.Map[Iterator[V], Iterator[U]](f)
}

func MonadChain[U, V any](ma Iterator[U], f func(U) Iterator[V]) Iterator[V] {
	return G.MonadChain[Iterator[V], Iterator[U]](ma, f)
}

func Chain[U, V any](f func(U) Iterator[V]) func(Iterator[U]) Iterator[V] {
	return G.Chain[Iterator[V], Iterator[U]](f)
}

// Flatten converts an [Iterator] of [Iterator] into a simple [Iterator]
func Flatten[U any](ma Iterator[Iterator[U]]) Iterator[U] {
	return G.Flatten[Iterator[Iterator[U]], Iterator[U]](ma)
}

// From constructs an [Iterator] from a set of variadic arguments
func From[U any](data ...U) Iterator[U] {
	return G.From[Iterator[U]](data...)
}

// MakeBy returns an [Iterator] with an infinite number of elements initialized with `f(i)`
func MakeBy[FCT ~func(int) U, U any](f FCT) Iterator[U] {
	return G.MakeBy[Iterator[U]](f)
}

// Replicate creates an [Iterator] containing a value repeated an infinite number of times.
func Replicate[U any](a U) Iterator[U] {
	return G.Replicate[Iterator[U]](a)
}

// FilterMap filters and transforms the content of an iterator
func FilterMap[U, V any](f func(U) O.Option[V]) func(ma Iterator[U]) Iterator[V] {
	return G.FilterMap[Iterator[V], Iterator[U]](f)
}

// Filter filters the content of an iterator
func Filter[U any](f func(U) bool) func(ma Iterator[U]) Iterator[U] {
	return G.Filter[Iterator[U]](f)
}

// Ap is the applicative functor for iterators
func Ap[V, U any](ma Iterator[U]) func(Iterator[func(U) V]) Iterator[V] {
	return G.Ap[Iterator[func(U) V], Iterator[V]](ma)
}

// MonadAp is the applicative functor for iterators
func MonadAp[V, U any](fab Iterator[func(U) V], ma Iterator[U]) Iterator[V] {
	return G.MonadAp[Iterator[func(U) V], Iterator[V]](fab, ma)
}

// Repeat creates an [Iterator] containing a value repeated the specified number of times.
// Alias of [Replicate]
func Repeat[U any](n int, a U) Iterator[U] {
	return G.Repeat[Iterator[U]](n, a)
}

// Count creates an [Iterator] containing a consecutive sequence of integers starting with the provided start value
func Count(start int) Iterator[int] {
	return G.Count[Iterator[int]](start)
}

// FilterChain filters and transforms the content of an iterator
func FilterChain[U, V any](f func(U) O.Option[Iterator[V]]) func(ma Iterator[U]) Iterator[V] {
	return G.FilterChain[Iterator[Iterator[V]], Iterator[V], Iterator[U]](f)
}

// FoldMap maps and folds an iterator. Map the iterator passing each value to the iterating function. Then fold the results using the provided Monoid.
func FoldMap[U, V any](m M.Monoid[V]) func(func(U) V) func(ma Iterator[U]) V {
	return G.FoldMap[Iterator[U], func(U) V, U, V](m)
}

// Fold folds the iterator using the provided Monoid.
func Fold[U any](m M.Monoid[U]) func(Iterator[U]) U {
	return G.Fold[Iterator[U]](m)
}

func MonadChainFirst[U, V any](ma Iterator[U], f func(U) Iterator[V]) Iterator[U] {
	return G.MonadChainFirst[Iterator[V], Iterator[U], U, V](ma, f)
}

func ChainFirst[U, V any](f func(U) Iterator[V]) func(Iterator[U]) Iterator[U] {
	return G.ChainFirst[Iterator[V], Iterator[U], U, V](f)
}
