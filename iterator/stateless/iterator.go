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
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// Iterator represents a stateless, pure way to iterate over a sequence
type Iterator[U any] L.Lazy[O.Option[T.Tuple2[Iterator[U], U]]]

// Empty returns the empty iterator
func Empty[U any]() Iterator[U] {
	return G.Empty[Iterator[U]]()
}

// Of returns an iterator with one single element
func Of[GU ~func() O.Option[T.Tuple2[GU, U]], U any](a U) Iterator[U] {
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

// MakeBy returns an [Iterator] with `n` elements initialized with `f(i)`
func MakeBy[FCT ~func(int) U, U any](n int, f FCT) Iterator[U] {
	return G.MakeBy[Iterator[U]](n, f)
}

// Replicate creates an [Iterator] containing a value repeated the specified number of times.
func Replicate[U any](n int, a U) Iterator[U] {
	return G.Replicate[Iterator[U]](n, a)
}
