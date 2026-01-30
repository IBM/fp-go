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
	Mg "github.com/IBM/fp-go/v2/magma"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	G "github.com/IBM/fp-go/v2/record/generic"
)

// IsEmpty tests if a record is empty (contains no entries).
//
// Returns true if the record has no key-value pairs, false otherwise.
//
// Example:
//
//	empty := Record[string, int]{}
//	IsEmpty(empty) // true
//
//	nonEmpty := Record[string, int]{"a": 1}
//	IsEmpty(nonEmpty) // false
func IsEmpty[K comparable, V any](r Record[K, V]) bool {
	return G.IsEmpty(r)
}

// IsNonEmpty tests if a record is not empty (contains at least one entry).
//
// Returns true if the record has at least one key-value pair, false otherwise.
// This is the logical negation of IsEmpty.
//
// Example:
//
//	record := Record[string, int]{"a": 1}
//	IsNonEmpty(record) // true
func IsNonEmpty[K comparable, V any](r Record[K, V]) bool {
	return G.IsNonEmpty(r)
}

// Keys returns all keys from a record as a slice.
//
// The order of keys is non-deterministic due to Go's map iteration behavior.
// Use KeysOrd if you need keys in a specific order.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	keys := Keys(record) // ["a", "b", "c"] in any order
func Keys[K comparable, V any](r Record[K, V]) []K {
	return G.Keys[Record[K, V], []K](r)
}

// Values returns all values from a record as a slice.
//
// The order of values is non-deterministic due to Go's map iteration behavior.
// Use ValuesOrd if you need values ordered by their keys.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	values := Values(record) // [1, 2, 3] in any order
func Values[K comparable, V any](r Record[K, V]) []V {
	return G.Values[Record[K, V], []V](r)
}

// Collect transforms each key-value pair in a record using a collector function
// and returns the results as a slice.
//
// The collector function receives both the key and value, allowing transformations
// that depend on both. The order of results is non-deterministic.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	toStrings := Collect(func(k string, v int) string {
//	    return fmt.Sprintf("%s=%d", k, v)
//	})
//	result := toStrings(record) // ["a=1", "b=2"] in any order
func Collect[K comparable, V, R any](f func(K, V) R) func(Record[K, V]) []R {
	return G.Collect[Record[K, V], []R](f)
}

// CollectOrd transforms each key-value pair in a record using a collector function
// and returns the results as a slice in the order specified by the Ord instance.
//
// Unlike Collect, this function guarantees the order of results based on key ordering.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	toStrings := CollectOrd(S.Ord)(func(k string, v int) string {
//	    return fmt.Sprintf("%s=%d", k, v)
//	})
//	result := toStrings(record) // ["a=1", "b=2", "c=3"] (ordered by key)
func CollectOrd[V, R any, K comparable](o ord.Ord[K]) func(func(K, V) R) func(Record[K, V]) []R {
	return G.CollectOrd[Record[K, V], []R](o)
}

// Reduce reduces a record to a single value by applying a reducer function to each value.
//
// The reducer function receives the accumulated result and the current value,
// returning the new accumulated result. The iteration order is non-deterministic.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	sum := Reduce(func(acc int, v int) int {
//	    return acc + v
//	}, 0)
//	result := sum(record) // 6
func Reduce[K comparable, V, R any](f func(R, V) R, initial R) func(Record[K, V]) R {
	return G.Reduce[Record[K, V]](f, initial)
}

// ReduceWithIndex reduces a record to a single value by applying a reducer function
// to each key-value pair.
//
// The reducer function receives the key, accumulated result, and current value,
// allowing reductions that depend on the key. The iteration order is non-deterministic.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	weightedSum := ReduceWithIndex(func(k string, acc int, v int) int {
//	    weight := len(k)
//	    return acc + (v * weight)
//	}, 0)
//	result := weightedSum(record) // 1 + 2 + 3 = 6
func ReduceWithIndex[K comparable, V, R any](f func(K, R, V) R, initial R) func(Record[K, V]) R {
	return G.ReduceWithIndex[Record[K, V]](f, initial)
}

// ReduceRef reduces a record to a single value by applying a reducer function
// to each value reference.
//
// Similar to Reduce, but passes value pointers instead of values, which can be
// more efficient for large value types and allows mutation if needed.
//
// Example:
//
//	record := Record[string, LargeStruct]{...}
//	result := ReduceRef(func(acc int, v *LargeStruct) int {
//	    return acc + v.Size
//	}, 0)(record)
func ReduceRef[K comparable, V, R any](f func(R, *V) R, initial R) func(Record[K, V]) R {
	return G.ReduceRef[Record[K, V]](f, initial)
}

// ReduceRefWithIndex reduces a record to a single value by applying a reducer function
// to each key-value pair with value references.
//
// Combines the benefits of ReduceWithIndex and ReduceRef, providing both key access
// and value references for efficient processing of large value types.
//
// Example:
//
//	record := Record[string, LargeStruct]{...}
//	result := ReduceRefWithIndex(func(k string, acc int, v *LargeStruct) int {
//	    return acc + len(k) * v.Size
//	}, 0)(record)
func ReduceRefWithIndex[K comparable, V, R any](f func(K, R, *V) R, initial R) func(Record[K, V]) R {
	return G.ReduceRefWithIndex[Record[K, V]](f, initial)
}

// MonadMap transforms each value in a record using the provided function.
//
// This is the monadic version of Map, taking the record as the first parameter.
// Useful for method chaining or when the record is already available.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	result := MonadMap(record, func(v int) int { return v * 2 })
//	// result: {"a": 2, "b": 4}
func MonadMap[K comparable, V, R any](r Record[K, V], f func(V) R) Record[K, R] {
	return G.MonadMap[Record[K, V], Record[K, R]](r, f)
}

// MonadMapWithIndex transforms each key-value pair in a record using the provided function.
//
// This is the monadic version of MapWithIndex, taking the record as the first parameter.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	result := MonadMapWithIndex(record, func(k string, v int) string {
//	    return fmt.Sprintf("%s=%d", k, v)
//	})
//	// result: {"a": "a=1", "b": "b=2"}
func MonadMapWithIndex[K comparable, V, R any](r Record[K, V], f func(K, V) R) Record[K, R] {
	return G.MonadMapWithIndex[Record[K, V], Record[K, R]](r, f)
}

// MonadMapRefWithIndex transforms each key-value pair in a record using the provided
// function with value references.
//
// Combines MonadMapWithIndex with reference passing for efficient processing of large values.
//
// Example:
//
//	record := Record[string, LargeStruct]{...}
//	result := MonadMapRefWithIndex(record, func(k string, v *LargeStruct) int {
//	    return len(k) + v.Size
//	})
func MonadMapRefWithIndex[K comparable, V, R any](r Record[K, V], f func(K, *V) R) Record[K, R] {
	return G.MonadMapRefWithIndex[Record[K, V], Record[K, R]](r, f)
}

// MonadMapRef transforms each value in a record using the provided function with value references.
//
// This is the monadic version of MapRef, useful for efficient processing of large value types.
//
// Example:
//
//	record := Record[string, LargeStruct]{...}
//	result := MonadMapRef(record, func(v *LargeStruct) int {
//	    return v.Size
//	})
func MonadMapRef[K comparable, V, R any](r Record[K, V], f func(*V) R) Record[K, R] {
	return G.MonadMapRef[Record[K, V], Record[K, R]](r, f)
}

// Map returns a function that transforms each value in a record using the provided function.
//
// This is the Functor map operation for records, creating a new record with transformed values
// while preserving all keys. This is one of the most commonly used operations.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	double := Map(func(v int) int { return v * 2 })
//	result := double(record) // {"a": 2, "b": 4, "c": 6}
func Map[K comparable, V, R any](f func(V) R) Operator[K, V, R] {
	return G.Map[Record[K, V], Record[K, R]](f)
}

// MapRef returns a function that transforms each value in a record using the provided
// function with value references.
//
// More efficient than Map for large value types as it avoids copying values.
//
// Example:
//
//	record := Record[string, LargeStruct]{...}
//	extractSize := MapRef(func(v *LargeStruct) int { return v.Size })
//	result := extractSize(record)
func MapRef[K comparable, V, R any](f func(*V) R) Operator[K, V, R] {
	return G.MapRef[Record[K, V], Record[K, R]](f)
}

// MapWithIndex returns a function that transforms each key-value pair in a record
// using the provided function.
//
// Useful when the transformation depends on both the key and value.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	prefixWithKey := MapWithIndex(func(k string, v int) string {
//	    return fmt.Sprintf("%s:%d", k, v)
//	})
//	result := prefixWithKey(record) // {"a": "a:1", "b": "b:2"}
func MapWithIndex[K comparable, V, R any](f func(K, V) R) Operator[K, V, R] {
	return G.MapWithIndex[Record[K, V], Record[K, R]](f)
}

// MapRefWithIndex returns a function that transforms each key-value pair in a record
// using the provided function with value references.
//
// Combines the benefits of MapWithIndex and MapRef for efficient key-aware transformations
// of large value types.
//
// Example:
//
//	record := Record[string, LargeStruct]{...}
//	transform := MapRefWithIndex(func(k string, v *LargeStruct) int {
//	    return len(k) + v.Size
//	})
//	result := transform(record)
func MapRefWithIndex[K comparable, V, R any](f func(K, *V) R) Operator[K, V, R] {
	return G.MapRefWithIndex[Record[K, V], Record[K, R]](f)
}

// Lookup returns a function that retrieves the value for a key in a record if it exists.
//
// Returns Some(value) if the key exists, None otherwise. This is the curried version
// that returns a Kleisli arrow, useful for composition.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	lookupA := Lookup[int]("a")
//	result := lookupA(record) // Some(1)
//
//	lookupC := Lookup[int]("c")
//	result2 := lookupC(record) // None
func Lookup[V any, K comparable](k K) option.Kleisli[Record[K, V], V] {
	return G.Lookup[Record[K, V]](k)
}

// MonadLookup retrieves the value for a key in a record if it exists.
//
// This is the monadic version of Lookup, taking the record as the first parameter.
// Returns Some(value) if the key exists, None otherwise.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	result := MonadLookup(record, "a") // Some(1)
//	result2 := MonadLookup(record, "c") // None
func MonadLookup[V any, K comparable](m Record[K, V], k K) Option[V] {
	return G.MonadLookup(m, k)
}

// Has tests if a key exists in a record.
//
// Returns true if the key is present in the record, false otherwise.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	Has("a", record) // true
//	Has("c", record) // false
func Has[K comparable, V any](k K, r Record[K, V]) bool {
	return G.Has(k, r)
}

// Union combines two records using the provided Magma to resolve conflicts for duplicate keys.
//
// The Magma defines how to combine values when the same key exists in both records.
// This is useful for custom merge strategies beyond simple replacement.
//
// Example:
//
//	// Sum values for duplicate keys
//	sumMagma := Mg.MakeMagma(func(a, b int) int { return a + b })
//	unionSum := Union(sumMagma)
//
//	r1 := Record[string, int]{"a": 1, "b": 2}
//	r2 := Record[string, int]{"b": 3, "c": 4}
//	result := unionSum(r1)(r2) // {"a": 1, "b": 5, "c": 4}
func Union[K comparable, V any](m Mg.Magma[V]) func(Record[K, V]) Operator[K, V, V] {
	return G.Union[Record[K, V]](m)
}

// Merge combines two records, giving precedence to values in the right record for duplicate keys.
//
// This is a simpler alternative to Union that always takes the right value when keys conflict.
// Also see MergeMonoid for the monoid version.
//
// Example:
//
//	r1 := Record[string, int]{"a": 1, "b": 2}
//	r2 := Record[string, int]{"b": 3, "c": 4}
//	result := Merge(r2)(r1) // {"a": 1, "b": 3, "c": 4}
func Merge[K comparable, V any](right Record[K, V]) Operator[K, V, V] {
	return G.Merge(right)
}

// Empty creates an empty record with no entries.
//
// This is useful as an identity element for record operations or as a starting point
// for building records incrementally.
//
// Example:
//
//	empty := Empty[string, int]()
//	IsEmpty(empty) // true
func Empty[K comparable, V any]() Record[K, V] {
	return G.Empty[Record[K, V]]()
}

// Size returns the number of key-value pairs in a record.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	Size(record) // 3
//
//	empty := Empty[string, int]()
//	Size(empty) // 0
func Size[K comparable, V any](r Record[K, V]) int {
	return G.Size(r)
}

// ToArray converts a record to a slice of key-value pairs (entries).
//
// The order of entries is non-deterministic due to Go's map iteration behavior.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	entries := ToArray(record)
//	// entries: []Entry[string, int]{{F1: "a", F2: 1}, {F1: "b", F2: 2}} in any order
func ToArray[K comparable, V any](r Record[K, V]) Entries[K, V] {
	return G.ToArray[Record[K, V], Entries[K, V]](r)
}

// ToEntries converts a record to a slice of key-value pairs (entries).
//
// This is an alias for ToArray, providing semantic clarity when working with entries.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	entries := ToEntries(record)
func ToEntries[K comparable, V any](r Record[K, V]) Entries[K, V] {
	return G.ToEntries[Record[K, V], Entries[K, V]](r)
}

// FromEntries creates a record from a slice of key-value pairs.
//
// If duplicate keys exist in the slice, the last occurrence wins.
//
// Example:
//
//	entries := Entries[string, int]{
//	    P.MakePair("a", 1),
//	    P.MakePair("b", 2),
//	}
//	record := FromEntries(entries) // {"a": 1, "b": 2}
func FromEntries[K comparable, V any](fa Entries[K, V]) Record[K, V] {
	return G.FromEntries[Record[K, V]](fa)
}

// UpsertAt returns a function that inserts or updates a key-value pair in a record.
//
// If the key exists, its value is updated. If the key doesn't exist, it's added.
// The original record is not modified; a new record is returned.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	addC := UpsertAt("c", 3)
//	result := addC(record) // {"a": 1, "b": 2, "c": 3}
//
//	updateA := UpsertAt("a", 10)
//	result2 := updateA(record) // {"a": 10, "b": 2}
func UpsertAt[K comparable, V any](k K, v V) Operator[K, V, V] {
	return G.UpsertAt[Record[K, V]](k, v)
}

// DeleteAt returns a function that removes a key from a record.
//
// If the key doesn't exist, the record is returned unchanged.
// The original record is not modified; a new record is returned.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	removeB := DeleteAt[string, int]("b")
//	result := removeB(record) // {"a": 1, "c": 3}
func DeleteAt[K comparable, V any](k K) Operator[K, V, V] {
	return G.DeleteAt[Record[K, V]](k)
}

// Singleton creates a new record with a single key-value pair.
//
// This is useful for creating records with one entry or as a building block
// for more complex record operations.
//
// Example:
//
//	record := Singleton("key", 42)
//	// record: {"key": 42}
func Singleton[K comparable, V any](k K, v V) Record[K, V] {
	return G.Singleton[Record[K, V]](k, v)
}

// FilterMapWithIndex filters and transforms a record simultaneously.
//
// The transformation function returns Some(value) to keep and transform an entry,
// or None to exclude it. This combines filtering and mapping in a single pass.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	evenDoubled := FilterMapWithIndex(func(k string, v int) Option[int] {
//	    if v%2 == 0 {
//	        return O.Some(v * 2)
//	    }
//	    return O.None[int]()
//	})
//	result := evenDoubled(record) // {"b": 4}
func FilterMapWithIndex[K comparable, V1, V2 any](f func(K, V1) Option[V2]) Operator[K, V1, V2] {
	return G.FilterMapWithIndex[Record[K, V1], Record[K, V2]](f)
}

// FilterMap filters and transforms a record based on values only.
//
// Similar to FilterMapWithIndex but the transformation function only receives values.
// Returns Some(value) to keep and transform, None to exclude.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	evenDoubled := FilterMap(func(v int) Option[int] {
//	    if v%2 == 0 {
//	        return O.Some(v * 2)
//	    }
//	    return O.None[int]()
//	})
//	result := evenDoubled(record) // {"b": 4}
func FilterMap[K comparable, V1, V2 any](f option.Kleisli[V1, V2]) Operator[K, V1, V2] {
	return G.FilterMap[Record[K, V1], Record[K, V2]](f)
}

// Filter creates a new record with only the entries whose keys match the predicate.
//
// The predicate tests keys only, not values. Use FilterWithIndex to test both.
//
// Example:
//
//	record := Record[string, int]{"apple": 1, "banana": 2, "cherry": 3}
//	startsWithA := Filter[string, int](func(k string) bool {
//	    return strings.HasPrefix(k, "a")
//	})
//	result := startsWithA(record) // {"apple": 1}
func Filter[K comparable, V any](f Predicate[K]) Operator[K, V, V] {
	return G.Filter[Record[K, V]](f)
}

// FilterWithIndex creates a new record with only the entries that match the predicate.
//
// The predicate receives both key and value, allowing filtering based on both.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	evenValues := FilterWithIndex(func(k string, v int) bool {
//	    return v%2 == 0
//	})
//	result := evenValues(record) // {"b": 2}
func FilterWithIndex[K comparable, V any](f PredicateWithIndex[K, V]) Operator[K, V, V] {
	return G.FilterWithIndex[Record[K, V]](f)
}

// IsNil checks if the record is nil (not initialized).
//
// Note: This checks for nil, not empty. An empty record {} is not nil.
//
// Example:
//
//	var record Record[string, int]
//	IsNil(record) // true
//
//	record = Record[string, int]{}
//	IsNil(record) // false
func IsNil[K comparable, V any](m Record[K, V]) bool {
	return G.IsNil(m)
}

// IsNonNil checks if the record is not nil (is initialized).
//
// This is the logical negation of IsNil.
//
// Example:
//
//	record := Record[string, int]{"a": 1}
//	IsNonNil(record) // true
func IsNonNil[K comparable, V any](m Record[K, V]) bool {
	return G.IsNonNil(m)
}

// ConstNil returns a nil record.
//
// This is useful as a constant function that always returns nil,
// particularly in functional composition scenarios.
//
// Example:
//
//	nilRecord := ConstNil[string, int]()
//	IsNil(nilRecord) // true
func ConstNil[K comparable, V any]() Record[K, V] {
	return Record[K, V](nil)
}

// MonadChainWithIndex chains a record transformation that produces records, combining results using a Monoid.
//
// This is the monadic bind operation for records. Each value is transformed into a record,
// and all resulting records are combined using the provided Monoid. The transformation
// function receives both key and value.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	expand := func(k string, v int) Record[string, int] {
//	    return Record[string, int]{
//	        k + "_double": v * 2,
//	        k + "_triple": v * 3,
//	    }
//	}
//	result := MonadChainWithIndex(MergeMonoid[string, int](), record, expand)
//	// result: {"a_double": 2, "a_triple": 3, "b_double": 4, "b_triple": 6}
func MonadChainWithIndex[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]], r Record[K, V1], f KleisliWithIndex[K, V1, V2]) Record[K, V2] {
	return G.MonadChainWithIndex(m, r, f)
}

// MonadChain chains a record transformation that produces records, combining results using a Monoid.
//
// Similar to MonadChainWithIndex but the transformation function only receives values.
// This is the monadic bind (flatMap) operation for records.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	expand := func(v int) Record[string, int] {
//	    return Record[string, int]{
//	        "double": v * 2,
//	        "triple": v * 3,
//	    }
//	}
//	result := MonadChain(MergeMonoid[string, int](), record, expand)
func MonadChain[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]], r Record[K, V1], f Kleisli[K, V1, V2]) Record[K, V2] {
	return G.MonadChain(m, r, f)
}

// ChainWithIndex returns a function that chains record transformations with key-value access.
//
// This is the curried version of MonadChainWithIndex, useful for composition.
//
// Example:
//
//	expand := ChainWithIndex[int, string, string](MergeMonoid[string, string]())(
//	    func(k string, v int) Record[string, string] {
//	        return Record[string, string]{k + "_str": fmt.Sprint(v)}
//	    },
//	)
//	result := expand(Record[string, int]{"a": 1})
func ChainWithIndex[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(KleisliWithIndex[K, V1, V2]) Operator[K, V1, V2] {
	return G.ChainWithIndex[Record[K, V1]](m)
}

// Chain returns a function that chains record transformations.
//
// This is the curried version of MonadChain, useful for composition.
// The monadic bind operation for records.
//
// Example:
//
//	expand := Chain[int, string, string](MergeMonoid[string, string]())(
//	    func(v int) Record[string, string] {
//	        return Record[string, string]{"result": fmt.Sprint(v)}
//	    },
//	)
//	result := expand(Record[string, int]{"a": 1})
func Chain[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(Kleisli[K, V1, V2]) Operator[K, V1, V2] {
	return G.Chain[Record[K, V1]](m)
}

// Flatten converts a nested record (record of records) into a flat record.
//
// When keys conflict between nested records, the Monoid determines how to combine values.
// This is the monadic join operation.
//
// Example:
//
//	nested := Record[string, Record[string, int]]{
//	    "group1": {"a": 1, "b": 2},
//	    "group2": {"c": 3, "d": 4},
//	}
//	flat := Flatten(MergeMonoid[string, int]())(nested)
//	// flat: {"a": 1, "b": 2, "c": 3, "d": 4}
func Flatten[K comparable, V any](m Monoid[Record[K, V]]) func(Record[K, Record[K, V]]) Record[K, V] {
	return G.Flatten[Record[K, Record[K, V]]](m)
}

// FilterChainWithIndex filters and chains transformations that produce records.
//
// Combines filtering with chaining: the transformation returns Some(record) to include
// the result, or None to exclude it. Results are combined using the Monoid.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	expandEven := FilterChainWithIndex[int, string, int](MergeMonoid[string, int]())(
//	    func(k string, v int) Option[Record[string, int]] {
//	        if v%2 == 0 {
//	            return O.Some(Record[string, int]{k + "_doubled": v * 2})
//	        }
//	        return O.None[Record[string, int]]()
//	    },
//	)
//	result := expandEven(record) // {"b_doubled": 4}
func FilterChainWithIndex[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(func(K, V1) Option[Record[K, V2]]) Operator[K, V1, V2] {
	return G.FilterChainWithIndex[Record[K, V1]](m)
}

// FilterChain filters and chains transformations that produce records.
//
// Similar to FilterChainWithIndex but without key access in the transformation function.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	expandEven := FilterChain[int, string, int](MergeMonoid[string, int]())(
//	    func(v int) Option[Record[string, int]] {
//	        if v%2 == 0 {
//	            return O.Some(Record[string, int]{"doubled": v * 2})
//	        }
//	        return O.None[Record[string, int]]()
//	    },
//	)
//	result := expandEven(record)
func FilterChain[V1 any, K comparable, V2 any](m Monoid[Record[K, V2]]) func(option.Kleisli[V1, Record[K, V2]]) Operator[K, V1, V2] {
	return G.FilterChain[Record[K, V1]](m)
}

// FoldMap maps each value in a record and folds the results using a Monoid.
//
// This is a two-step operation: first map each value using the provided function,
// then combine all results using the Monoid's operation.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	sumDoubled := FoldMap(N.MonoidSum[int]())(func(v int) int {
//	    return v * 2
//	})
//	result := sumDoubled(record) // 12 (2 + 4 + 6)
func FoldMap[K comparable, A, B any](m Monoid[B]) func(func(A) B) func(Record[K, A]) B {
	return G.FoldMap[Record[K, A]](m)
}

// FoldMapWithIndex maps each key-value pair in a record and folds the results using a Monoid.
//
// Similar to FoldMap but the mapping function receives both key and value.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2}
//	weightedSum := FoldMapWithIndex(N.MonoidSum[int]())(func(k string, v int) int {
//	    return len(k) * v
//	})
//	result := weightedSum(record) // 1*1 + 1*2 = 3
func FoldMapWithIndex[K comparable, A, B any](m Monoid[B]) func(func(K, A) B) func(Record[K, A]) B {
	return G.FoldMapWithIndex[Record[K, A]](m)
}

// Fold combines all values in a record using a Monoid.
//
// This is useful when the record values are already of the target type and you
// just need to combine them.
//
// Example:
//
//	record := Record[string, int]{"a": 1, "b": 2, "c": 3}
//	sum := Fold(N.MonoidSum[int]())
//	result := sum(record) // 6
func Fold[K comparable, A any](m Monoid[A]) func(Record[K, A]) A {
	return G.Fold[Record[K, A]](m)
}

// ReduceOrdWithIndex reduces a record to a single value with keys processed in order.
//
// Unlike ReduceWithIndex, this guarantees the order in which keys are processed
// based on the provided Ord instance.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	concat := ReduceOrdWithIndex(S.Ord)(func(k string, acc string, v int) string {
//	    return acc + k + fmt.Sprint(v)
//	}, "")
//	result := concat(record) // "a1b2c3" (alphabetical order)
func ReduceOrdWithIndex[V, R any, K comparable](o ord.Ord[K]) func(func(K, R, V) R, R) func(Record[K, V]) R {
	return G.ReduceOrdWithIndex[Record[K, V], K, V, R](o)
}

// ReduceOrd reduces a record to a single value with keys processed in order.
//
// Similar to ReduceOrdWithIndex but the reducer function doesn't receive keys.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	sum := ReduceOrd(S.Ord)(func(acc int, v int) int {
//	    return acc + v
//	}, 0)
//	result := sum(record) // 6 (order doesn't affect sum)
func ReduceOrd[V, R any, K comparable](o ord.Ord[K]) func(func(R, V) R, R) func(Record[K, V]) R {
	return G.ReduceOrd[Record[K, V], K, V, R](o)
}

// FoldMapOrd maps and folds a record with keys processed in order.
//
// Similar to FoldMap but guarantees the order in which entries are processed
// based on the provided Ord instance.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	concat := FoldMapOrd(S.Ord)(S.MonoidConcat)(func(v int) string {
//	    return fmt.Sprint(v)
//	})
//	result := concat(record) // "123" (alphabetical key order)
func FoldMapOrd[A, B any, K comparable](o ord.Ord[K]) func(m Monoid[B]) func(func(A) B) func(Record[K, A]) B {
	return G.FoldMapOrd[Record[K, A], K, A, B](o)
}

// FoldOrd combines all values in a record using a Monoid with keys processed in order.
//
// Similar to Fold but guarantees the order based on the provided Ord instance.
//
// Example:
//
//	record := Record[string, string]{"c": "3", "a": "1", "b": "2"}
//	concat := FoldOrd(S.Ord)(S.MonoidConcat)
//	result := concat(record) // "123" (alphabetical key order)
func FoldOrd[A any, K comparable](o ord.Ord[K]) func(m Monoid[A]) func(Record[K, A]) A {
	return G.FoldOrd[Record[K, A]](o)
}

// FoldMapOrdWithIndex maps and folds a record with key-value access and ordered processing.
//
// Combines FoldMapWithIndex with ordered key processing based on the Ord instance.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	concat := FoldMapOrdWithIndex(S.Ord)(S.MonoidConcat)(func(k string, v int) string {
//	    return k + fmt.Sprint(v)
//	})
//	result := concat(record) // "a1b2c3" (alphabetical key order)
func FoldMapOrdWithIndex[K comparable, A, B any](o ord.Ord[K]) func(m Monoid[B]) func(func(K, A) B) func(Record[K, A]) B {
	return G.FoldMapOrdWithIndex[Record[K, A], K, A, B](o)
}

// KeysOrd returns the keys from a record in the order specified by the Ord instance.
//
// Unlike Keys, this guarantees a specific order for the returned keys.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	getKeys := KeysOrd(S.Ord)
//	keys := getKeys(record) // ["a", "b", "c"]
func KeysOrd[V any, K comparable](o ord.Ord[K]) func(r Record[K, V]) []K {
	return G.KeysOrd[Record[K, V], []K](o)
}

// ValuesOrd returns the values from a record ordered by their keys.
//
// The values are returned in the order determined by sorting the keys using the Ord instance.
//
// Example:
//
//	record := Record[string, int]{"c": 3, "a": 1, "b": 2}
//	getValues := ValuesOrd(S.Ord)
//	values := getValues(record) // [1, 2, 3] (ordered by key: a, b, c)
func ValuesOrd[V any, K comparable](o ord.Ord[K]) func(r Record[K, V]) []V {
	return G.ValuesOrd[Record[K, V], []V](o)
}

// MonadFlap applies a value to a record of functions, producing a record of results.
//
// This is the monadic version of Flap. Each function in the record is applied to the value,
// preserving the keys.
//
// Example:
//
//	funcs := Record[string, func(int) int]{
//	    "double": func(x int) int { return x * 2 },
//	    "triple": func(x int) int { return x * 3 },
//	}
//	result := MonadFlap(funcs, 5)
//	// result: {"double": 10, "triple": 15}
func MonadFlap[B any, K comparable, A any](fab Record[K, func(A) B], a A) Record[K, B] {
	return G.MonadFlap[Record[K, func(A) B], Record[K, B]](fab, a)
}

// Flap returns a function that applies a value to a record of functions.
//
// This is the curried version of MonadFlap, useful for composition.
// It's the "flipped" version of Ap where the value is fixed and functions vary.
//
// Example:
//
//	funcs := Record[string, func(int) int]{
//	    "double": func(x int) int { return x * 2 },
//	    "triple": func(x int) int { return x * 3 },
//	}
//	applyFive := Flap[int, string, int](5)
//	result := applyFive(funcs) // {"double": 10, "triple": 15}
func Flap[B any, K comparable, A any](a A) Operator[K, func(A) B, B] {
	return G.Flap[Record[K, func(A) B], Record[K, B]](a)
}

// Copy creates a shallow copy of a record.
//
// The keys and values are copied, but if values are pointers or contain pointers,
// they will point to the same underlying data. Use Clone for deep copying.
//
// Example:
//
//	original := Record[string, int]{"a": 1, "b": 2}
//	copy := Copy(original)
//	// Modifying copy doesn't affect original
func Copy[K comparable, V any](m Record[K, V]) Record[K, V] {
	return G.Copy(m)
}

// Clone creates a deep copy of a record using the provided endomorphism to clone values.
//
// The endomorphism is applied to each value to create a deep copy. This is useful
// when values contain pointers or other references that need to be duplicated.
//
// Example:
//
//	type Data struct { Value int }
//	cloneData := func(d Data) Data { return Data{Value: d.Value} }
//
//	original := Record[string, Data]{"a": {Value: 1}}
//	deepCopy := Clone(cloneData)(original)
func Clone[K comparable, V any](f Endomorphism[V]) Endomorphism[Record[K, V]] {
	return G.Clone[Record[K, V]](f)
}

// FromFoldableMap converts a foldable structure to a record by mapping elements to entries.
//
// The mapping function transforms each element into a key-value entry. When duplicate keys
// occur, the Magma determines how to combine their values. This is useful for building
// records from custom data structures.
//
// Example:
//
//	type Person struct { ID string; Score int }
//	people := []Person{{"alice", 10}, {"bob", 20}, {"alice", 15}}
//
//	sumMagma := Mg.MakeMagma(func(a, b int) int { return a + b })
//	toRecord := FromArrayMap[Person, string, int](sumMagma)(func(p Person) Entry[string, int] {
//	    return P.MakePair(p.ID, p.Score)
//	})
//	result := toRecord(people) // {"alice": 25, "bob": 20}
func FromFoldableMap[
	FOLDABLE ~func(func(Record[K, V], A) Record[K, V], Record[K, V]) func(HKTA) Record[K, V], // the reduce function
	A any,
	HKTA any,
	K comparable,
	V any](m Mg.Magma[V], red FOLDABLE) func(f func(A) Entry[K, V]) Kleisli[K, HKTA, V] {
	return G.FromFoldableMap[func(A) Entry[K, V]](m, red)
}

// FromArrayMap converts an array to a record by mapping elements to entries.
//
// Each element is transformed into a key-value entry. When duplicate keys occur,
// the Magma determines how to combine their values.
//
// Example:
//
//	type Item struct { Name string; Count int }
//	items := []Item{{"apple", 5}, {"banana", 3}, {"apple", 2}}
//
//	sumMagma := Mg.MakeMagma(func(a, b int) int { return a + b })
//	toRecord := FromArrayMap[Item, string, int](sumMagma)(func(item Item) Entry[string, int] {
//	    return P.MakePair(item.Name, item.Count)
//	})
//	result := toRecord(items) // {"apple": 7, "banana": 3}
func FromArrayMap[
	A any,
	K comparable,
	V any](m Mg.Magma[V]) func(f func(A) Entry[K, V]) Kleisli[K, []A, V] {
	return G.FromArrayMap[func(A) Entry[K, V], []A, Record[K, V]](m)
}

// FromFoldable converts a foldable structure of entries to a record.
//
// The foldable structure should contain Entry[K, V] elements. When duplicate keys
// occur, the Magma determines how to combine their values.
//
// Example:
//
//	entries := []Entry[string, int]{
//	    P.MakePair("a", 1),
//	    P.MakePair("b", 2),
//	    P.MakePair("a", 3),
//	}
//	sumMagma := Mg.MakeMagma(func(a, b int) int { return a + b })
//	toRecord := FromArray(sumMagma)
//	result := toRecord(entries) // {"a": 4, "b": 2}
func FromFoldable[
	HKTA any,
	FOLDABLE ~func(func(Record[K, V], Entry[K, V]) Record[K, V], Record[K, V]) func(HKTA) Record[K, V], // the reduce function
	K comparable,
	V any](m Mg.Magma[V], red FOLDABLE) Kleisli[K, HKTA, V] {
	return G.FromFoldable(m, red)
}

// FromArray converts an array of entries to a record.
//
// When duplicate keys occur, the Magma determines how to combine their values.
// This is useful for aggregating data with the same keys.
//
// Example:
//
//	entries := Entries[string, int]{
//	    P.MakePair("a", 1),
//	    P.MakePair("b", 2),
//	    P.MakePair("a", 3),
//	}
//	sumMagma := Mg.MakeMagma(func(a, b int) int { return a + b })
//	toRecord := FromArray(sumMagma)
//	result := toRecord(entries) // {"a": 4, "b": 2}
func FromArray[
	K comparable,
	V any](m Mg.Magma[V]) Kleisli[K, Entries[K, V], V] {
	return G.FromArray[Entries[K, V], Record[K, V]](m)
}

// MonadAp applies a record of functions to a record of values, producing a record of results.
//
// This is the applicative apply operation for records. For each matching key in both records,
// the function is applied to the value. When keys exist in both records, results are combined
// using the provided Monoid.
//
// Example:
//
//	funcs := Record[string, func(int) int]{
//	    "double": func(x int) int { return x * 2 },
//	    "triple": func(x int) int { return x * 3 },
//	}
//	values := Record[string, int]{"double": 5, "triple": 7}
//	result := MonadAp(MergeMonoid[string, int](), funcs, values)
//	// result: {"double": 10, "triple": 21}
func MonadAp[A any, K comparable, B any](m Monoid[Record[K, B]], fab Record[K, func(A) B], fa Record[K, A]) Record[K, B] {
	return G.MonadAp(m, fab, fa)
}

// Ap returns a function that applies a record of functions to a record of values.
//
// This is the curried version of MonadAp, useful for composition.
// The applicative apply operation for records.
//
// Example:
//
//	funcs := Record[string, func(int) int]{
//	    "double": func(x int) int { return x * 2 },
//	}
//	applyFuncs := Ap[int, string, int](MergeMonoid[string, int]())
//	values := Record[string, int]{"double": 5}
//	result := applyFuncs(values)(funcs) // {"double": 10}
func Ap[A any, K comparable, B any](m Monoid[Record[K, B]]) func(fa Record[K, A]) Operator[K, func(A) B, B] {
	return G.Ap[Record[K, B], Record[K, func(A) B], Record[K, A]](m)
}

// Of creates a record with a single key-value pair.
//
// This is the pointed functor operation for records, lifting a value into the record context.
// It's an alias for Singleton but follows the standard functional programming naming convention.
//
// Example:
//
//	record := Of("key", 42)
//	// record: {"key": 42}
func Of[K comparable, A any](k K, a A) Record[K, A] {
	return Record[K, A]{k: a}
}
