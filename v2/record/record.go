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

package record

import (
	EM "github.com/IBM/fp-go/v2/endomorphism"
	Mg "github.com/IBM/fp-go/v2/magma"
	Mo "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	G "github.com/IBM/fp-go/v2/record/generic"
)

// IsEmpty tests if a map is empty
func IsEmpty[K comparable, V any](r map[K]V) bool {
	return G.IsEmpty(r)
}

// IsNonEmpty tests if a map is not empty
func IsNonEmpty[K comparable, V any](r map[K]V) bool {
	return G.IsNonEmpty(r)
}

// Keys returns the key in a map
func Keys[K comparable, V any](r map[K]V) []K {
	return G.Keys[map[K]V, []K](r)
}

// Values returns the values in a map
func Values[K comparable, V any](r map[K]V) []V {
	return G.Values[map[K]V, []V](r)
}

// Collect applies a collector function to the key value pairs in a map and returns the result as an array
func Collect[K comparable, V, R any](f func(K, V) R) func(map[K]V) []R {
	return G.Collect[map[K]V, []R](f)
}

// CollectOrd applies a collector function to the key value pairs in a map and returns the result as an array
func CollectOrd[V, R any, K comparable](o ord.Ord[K]) func(func(K, V) R) func(map[K]V) []R {
	return G.CollectOrd[map[K]V, []R](o)
}

// Reduce reduces a map to a single value by applying a reducer function to each value
func Reduce[K comparable, V, R any](f func(R, V) R, initial R) func(map[K]V) R {
	return G.Reduce[map[K]V](f, initial)
}

// ReduceWithIndex reduces a map to a single value by applying a reducer function to each key-value pair
func ReduceWithIndex[K comparable, V, R any](f func(K, R, V) R, initial R) func(map[K]V) R {
	return G.ReduceWithIndex[map[K]V](f, initial)
}

// ReduceRef reduces a map to a single value by applying a reducer function to each value reference
func ReduceRef[K comparable, V, R any](f func(R, *V) R, initial R) func(map[K]V) R {
	return G.ReduceRef[map[K]V](f, initial)
}

// ReduceRefWithIndex reduces a map to a single value by applying a reducer function to each key-value pair with value references
func ReduceRefWithIndex[K comparable, V, R any](f func(K, R, *V) R, initial R) func(map[K]V) R {
	return G.ReduceRefWithIndex[map[K]V](f, initial)
}

// MonadMap transforms each value in a map using the provided function
func MonadMap[K comparable, V, R any](r map[K]V, f func(V) R) map[K]R {
	return G.MonadMap[map[K]V, map[K]R](r, f)
}

// MonadMapWithIndex transforms each key-value pair in a map using the provided function
func MonadMapWithIndex[K comparable, V, R any](r map[K]V, f func(K, V) R) map[K]R {
	return G.MonadMapWithIndex[map[K]V, map[K]R](r, f)
}

// MonadMapRefWithIndex transforms each key-value pair in a map using the provided function with value references
func MonadMapRefWithIndex[K comparable, V, R any](r map[K]V, f func(K, *V) R) map[K]R {
	return G.MonadMapRefWithIndex[map[K]V, map[K]R](r, f)
}

// MonadMapRef transforms each value in a map using the provided function with value references
func MonadMapRef[K comparable, V, R any](r map[K]V, f func(*V) R) map[K]R {
	return G.MonadMapRef[map[K]V, map[K]R](r, f)
}

// Map returns a function that transforms each value in a map using the provided function
func Map[K comparable, V, R any](f func(V) R) Operator[K, V, R] {
	return G.Map[map[K]V, map[K]R](f)
}

// MapRef returns a function that transforms each value in a map using the provided function with value references
func MapRef[K comparable, V, R any](f func(*V) R) Operator[K, V, R] {
	return G.MapRef[map[K]V, map[K]R](f)
}

// MapWithIndex returns a function that transforms each key-value pair in a map using the provided function
func MapWithIndex[K comparable, V, R any](f func(K, V) R) Operator[K, V, R] {
	return G.MapWithIndex[map[K]V, map[K]R](f)
}

// MapRefWithIndex returns a function that transforms each key-value pair in a map using the provided function with value references
func MapRefWithIndex[K comparable, V, R any](f func(K, *V) R) Operator[K, V, R] {
	return G.MapRefWithIndex[map[K]V, map[K]R](f)
}

// Lookup returns the entry for a key in a map if it exists
func Lookup[V any, K comparable](k K) func(map[K]V) O.Option[V] {
	return G.Lookup[map[K]V](k)
}

// MonadLookup returns the entry for a key in a map if it exists
func MonadLookup[V any, K comparable](m map[K]V, k K) O.Option[V] {
	return G.MonadLookup(m, k)
}

// Has tests if a key is contained in a map
func Has[K comparable, V any](k K, r map[K]V) bool {
	return G.Has(k, r)
}

// Union combines two maps using the provided Magma to resolve conflicts for duplicate keys
func Union[K comparable, V any](m Mg.Magma[V]) func(map[K]V) Operator[K, V, V] {
	return G.Union[map[K]V](m)
}

// Merge combines two maps giving the values in the right one precedence. Also refer to [MergeMonoid]
func Merge[K comparable, V any](right map[K]V) Operator[K, V, V] {
	return G.Merge(right)
}

// Empty creates an empty map
func Empty[K comparable, V any]() map[K]V {
	return G.Empty[map[K]V]()
}

// Size returns the number of elements in a map
func Size[K comparable, V any](r map[K]V) int {
	return G.Size(r)
}

// ToArray converts a map to an array of key-value pairs
func ToArray[K comparable, V any](r map[K]V) Entries[K, V] {
	return G.ToArray[map[K]V, Entries[K, V]](r)
}

// ToEntries converts a map to an array of key-value pairs (alias for ToArray)
func ToEntries[K comparable, V any](r map[K]V) Entries[K, V] {
	return G.ToEntries[map[K]V, Entries[K, V]](r)
}

// FromEntries creates a map from an array of key-value pairs
func FromEntries[K comparable, V any](fa Entries[K, V]) map[K]V {
	return G.FromEntries[map[K]V](fa)
}

// UpsertAt returns a function that inserts or updates a key-value pair in a map
func UpsertAt[K comparable, V any](k K, v V) Operator[K, V, V] {
	return G.UpsertAt[map[K]V](k, v)
}

// DeleteAt returns a function that removes a key from a map
func DeleteAt[K comparable, V any](k K) Operator[K, V, V] {
	return G.DeleteAt[map[K]V](k)
}

// Singleton creates a new map with a single entry
func Singleton[K comparable, V any](k K, v V) map[K]V {
	return G.Singleton[map[K]V](k, v)
}

// FilterMapWithIndex creates a new map with only the elements for which the transformation function creates a Some
func FilterMapWithIndex[K comparable, V1, V2 any](f func(K, V1) O.Option[V2]) Operator[K, V1, V2] {
	return G.FilterMapWithIndex[map[K]V1, map[K]V2](f)
}

// FilterMap creates a new map with only the elements for which the transformation function creates a Some
func FilterMap[K comparable, V1, V2 any](f func(V1) O.Option[V2]) Operator[K, V1, V2] {
	return G.FilterMap[map[K]V1, map[K]V2](f)
}

// Filter creates a new map with only the elements that match the predicate
func Filter[K comparable, V any](f func(K) bool) Operator[K, V, V] {
	return G.Filter[map[K]V](f)
}

// FilterWithIndex creates a new map with only the elements that match the predicate
func FilterWithIndex[K comparable, V any](f func(K, V) bool) Operator[K, V, V] {
	return G.FilterWithIndex[map[K]V](f)
}

// IsNil checks if the map is set to nil
func IsNil[K comparable, V any](m map[K]V) bool {
	return G.IsNil(m)
}

// IsNonNil checks if the map is set to nil
func IsNonNil[K comparable, V any](m map[K]V) bool {
	return G.IsNonNil(m)
}

// ConstNil return a nil map
func ConstNil[K comparable, V any]() map[K]V {
	return map[K]V(nil)
}

// MonadChainWithIndex chains a map transformation function that produces maps, combining results using the provided Monoid
func MonadChainWithIndex[V1 any, K comparable, V2 any](m Mo.Monoid[map[K]V2], r map[K]V1, f func(K, V1) map[K]V2) map[K]V2 {
	return G.MonadChainWithIndex(m, r, f)
}

// MonadChain chains a map transformation function that produces maps, combining results using the provided Monoid
func MonadChain[V1 any, K comparable, V2 any](m Mo.Monoid[map[K]V2], r map[K]V1, f func(V1) map[K]V2) map[K]V2 {
	return G.MonadChain(m, r, f)
}

// ChainWithIndex returns a function that chains a map transformation function that produces maps, combining results using the provided Monoid
func ChainWithIndex[V1 any, K comparable, V2 any](m Mo.Monoid[map[K]V2]) func(func(K, V1) map[K]V2) Operator[K, V1, V2] {
	return G.ChainWithIndex[map[K]V1](m)
}

// Chain returns a function that chains a map transformation function that produces maps, combining results using the provided Monoid
func Chain[V1 any, K comparable, V2 any](m Mo.Monoid[map[K]V2]) func(func(V1) map[K]V2) Operator[K, V1, V2] {
	return G.Chain[map[K]V1](m)
}

// Flatten converts a nested map into a regular map
func Flatten[K comparable, V any](m Mo.Monoid[map[K]V]) func(map[K]map[K]V) map[K]V {
	return G.Flatten[map[K]map[K]V](m)
}

// FilterChainWithIndex creates a new map with only the elements for which the transformation function creates a Some
func FilterChainWithIndex[V1 any, K comparable, V2 any](m Mo.Monoid[map[K]V2]) func(func(K, V1) O.Option[map[K]V2]) Operator[K, V1, V2] {
	return G.FilterChainWithIndex[map[K]V1](m)
}

// FilterChain creates a new map with only the elements for which the transformation function creates a Some
func FilterChain[V1 any, K comparable, V2 any](m Mo.Monoid[map[K]V2]) func(func(V1) O.Option[map[K]V2]) Operator[K, V1, V2] {
	return G.FilterChain[map[K]V1](m)
}

// FoldMap maps and folds a record. Map the record passing each value to the iterating function. Then fold the results using the provided Monoid.
func FoldMap[K comparable, A, B any](m Mo.Monoid[B]) func(func(A) B) func(map[K]A) B {
	return G.FoldMap[map[K]A](m)
}

// FoldMapWithIndex maps and folds a record. Map the record passing each value to the iterating function. Then fold the results using the provided Monoid.
func FoldMapWithIndex[K comparable, A, B any](m Mo.Monoid[B]) func(func(K, A) B) func(map[K]A) B {
	return G.FoldMapWithIndex[map[K]A](m)
}

// Fold folds the record using the provided Monoid.
func Fold[K comparable, A any](m Mo.Monoid[A]) func(map[K]A) A {
	return G.Fold[map[K]A](m)
}

// ReduceOrdWithIndex reduces a map into a single value via a reducer function making sure that the keys are passed to the reducer in the specified order
func ReduceOrdWithIndex[V, R any, K comparable](o ord.Ord[K]) func(func(K, R, V) R, R) func(map[K]V) R {
	return G.ReduceOrdWithIndex[map[K]V, K, V, R](o)
}

// ReduceOrd reduces a map into a single value via a reducer function making sure that the keys are passed to the reducer in the specified order
func ReduceOrd[V, R any, K comparable](o ord.Ord[K]) func(func(R, V) R, R) func(map[K]V) R {
	return G.ReduceOrd[map[K]V, K, V, R](o)
}

// FoldMap maps and folds a record. Map the record passing each value to the iterating function. Then fold the results using the provided Monoid and the items in the provided order
func FoldMapOrd[A, B any, K comparable](o ord.Ord[K]) func(m Mo.Monoid[B]) func(func(A) B) func(map[K]A) B {
	return G.FoldMapOrd[map[K]A, K, A, B](o)
}

// Fold folds the record using the provided Monoid with the items passed in the given order
func FoldOrd[A any, K comparable](o ord.Ord[K]) func(m Mo.Monoid[A]) func(map[K]A) A {
	return G.FoldOrd[map[K]A](o)
}

// FoldMapWithIndex maps and folds a record. Map the record passing each value to the iterating function. Then fold the results using the provided Monoid and the items in the provided order
func FoldMapOrdWithIndex[K comparable, A, B any](o ord.Ord[K]) func(m Mo.Monoid[B]) func(func(K, A) B) func(map[K]A) B {
	return G.FoldMapOrdWithIndex[map[K]A, K, A, B](o)
}

// KeysOrd returns the keys in the map in their given order
func KeysOrd[V any, K comparable](o ord.Ord[K]) func(r map[K]V) []K {
	return G.KeysOrd[map[K]V, []K](o)
}

// ValuesOrd returns the values in the map ordered by their keys in the given order
func ValuesOrd[V any, K comparable](o ord.Ord[K]) func(r map[K]V) []V {
	return G.ValuesOrd[map[K]V, []V](o)
}

// MonadFlap applies a value to a map of functions, producing a map of results
func MonadFlap[B any, K comparable, A any](fab map[K]func(A) B, a A) map[K]B {
	return G.MonadFlap[map[K]func(A) B, map[K]B](fab, a)
}

// Flap returns a function that applies a value to a map of functions, producing a map of results
func Flap[B any, K comparable, A any](a A) func(map[K]func(A) B) map[K]B {
	return G.Flap[map[K]func(A) B, map[K]B](a)
}

// Copy creates a shallow copy of the map
func Copy[K comparable, V any](m map[K]V) map[K]V {
	return G.Copy(m)
}

// Clone creates a deep copy of the map using the provided endomorphism to clone the values
func Clone[K comparable, V any](f EM.Endomorphism[V]) EM.Endomorphism[map[K]V] {
	return G.Clone[map[K]V](f)
}

// FromFoldableMap converts from a reducer to a map
// Duplicate keys are resolved by the provided [Mg.Magma]
func FromFoldableMap[
	FOLDABLE ~func(func(map[K]V, A) map[K]V, map[K]V) func(HKTA) map[K]V, // the reduce function
	A any,
	HKTA any,
	K comparable,
	V any](m Mg.Magma[V], red FOLDABLE) func(f func(A) Entry[K, V]) func(fa HKTA) map[K]V {
	return G.FromFoldableMap[func(A) Entry[K, V]](m, red)
}

// FromArrayMap converts from an array to a map
// Duplicate keys are resolved by the provided [Mg.Magma]
func FromArrayMap[
	A any,
	K comparable,
	V any](m Mg.Magma[V]) func(f func(A) Entry[K, V]) func(fa []A) map[K]V {
	return G.FromArrayMap[func(A) Entry[K, V], []A, map[K]V](m)
}

// FromFoldable converts from a reducer to a map
// Duplicate keys are resolved by the provided [Mg.Magma]
func FromFoldable[
	HKTA any,
	FOLDABLE ~func(func(map[K]V, Entry[K, V]) map[K]V, map[K]V) func(HKTA) map[K]V, // the reduce function
	K comparable,
	V any](m Mg.Magma[V], red FOLDABLE) func(fa HKTA) map[K]V {
	return G.FromFoldable(m, red)
}

// FromArray converts from an array to a map
// Duplicate keys are resolved by the provided [Mg.Magma]
func FromArray[
	K comparable,
	V any](m Mg.Magma[V]) func(fa Entries[K, V]) map[K]V {
	return G.FromArray[Entries[K, V], map[K]V](m)
}

// MonadAp applies a map of functions to a map of values, combining results using the provided Monoid
func MonadAp[A any, K comparable, B any](m Mo.Monoid[map[K]B], fab map[K]func(A) B, fa map[K]A) map[K]B {
	return G.MonadAp(m, fab, fa)
}

// Ap returns a function that applies a map of functions to a map of values, combining results using the provided Monoid
func Ap[A any, K comparable, B any](m Mo.Monoid[map[K]B]) func(fa map[K]A) func(map[K]func(A) B) map[K]B {
	return G.Ap[map[K]B, map[K]func(A) B, map[K]A](m)
}

// Of creates a map with a single key-value pair
func Of[K comparable, A any](k K, a A) map[K]A {
	return map[K]A{k: a}
}
