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
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	// Endomorphism represents a function from a type to itself (A -> A).
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Monoid represents a monoid structure with an associative binary operation and identity element.
	Monoid[A any] = monoid.Monoid[A]

	// Semigroup represents a semigroup structure with an associative binary operation.
	Semigroup[A any] = semigroup.Semigroup[A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Record represents a map with comparable keys and values of any type.
	// This is the primary data structure for the record package, providing
	// functional operations over Go's native map type.
	//
	// Example:
	//
	//	type UserRecord = Record[string, User]
	//	users := UserRecord{
	//	    "alice": User{Name: "Alice", Age: 30},
	//	    "bob":   User{Name: "Bob", Age: 25},
	//	}
	Record[K comparable, V any] = map[K]V

	// Predicate is a function that tests whether a key satisfies a condition.
	// Used in filtering operations to determine which entries to keep.
	//
	// Example:
	//
	//	isVowel := func(k string) bool {
	//	    return strings.ContainsAny(k, "aeiou")
	//	}
	Predicate[K any] = predicate.Predicate[K]

	// PredicateWithIndex is a function that tests whether a key-value pair satisfies a condition.
	// Used in filtering operations that need access to both key and value.
	//
	// Example:
	//
	//	isAdult := func(name string, user User) bool {
	//	    return user.Age >= 18
	//	}
	PredicateWithIndex[K comparable, V any] = func(K, V) bool

	// Operator transforms a record from one value type to another while preserving keys.
	// This is the fundamental transformation type for record operations.
	//
	// Example:
	//
	//	doubleValues := Map(func(x int) int { return x * 2 })
	//	result := doubleValues(Record[string, int]{"a": 1, "b": 2})
	//	// result: {"a": 2, "b": 4}
	Operator[K comparable, V1, V2 any] = func(Record[K, V1]) Record[K, V2]

	// OperatorWithIndex transforms a record using both key and value information.
	// Useful when the transformation depends on the key.
	//
	// Example:
	//
	//	prefixWithKey := MapWithIndex(func(k string, v string) string {
	//	    return k + ":" + v
	//	})
	OperatorWithIndex[K comparable, V1, V2 any] = func(func(K, V1) V2) Operator[K, V1, V2]

	// Kleisli represents a monadic function that transforms a value into a record.
	// Used in chain operations for composing record-producing functions.
	//
	// Example:
	//
	//	expand := func(x int) Record[string, int] {
	//	    return Record[string, int]{
	//	        "double": x * 2,
	//	        "triple": x * 3,
	//	    }
	//	}
	Kleisli[K comparable, V1, V2 any] = func(V1) Record[K, V2]

	// KleisliWithIndex is a monadic function that uses both key and value to produce a record.
	//
	// Example:
	//
	//	expandWithKey := func(k string, v int) Record[string, int] {
	//	    return Record[string, int]{
	//	        k + "_double": v * 2,
	//	        k + "_triple": v * 3,
	//	    }
	//	}
	KleisliWithIndex[K comparable, V1, V2 any] = func(K, V1) Record[K, V2]

	// Reducer accumulates values from a record into a single result.
	// The function receives the accumulator and current value, returning the new accumulator.
	//
	// Example:
	//
	//	sum := Reduce(func(acc int, v int) int {
	//	    return acc + v
	//	}, 0)
	Reducer[V, R any] = func(R, V) R

	// ReducerWithIndex accumulates values using both key and value information.
	//
	// Example:
	//
	//	weightedSum := ReduceWithIndex(func(k string, acc int, v int) int {
	//	    weight := len(k)
	//	    return acc + (v * weight)
	//	}, 0)
	ReducerWithIndex[K comparable, V, R any] = func(K, R, V) R

	// Collector transforms key-value pairs into a result type and collects them into an array.
	//
	// Example:
	//
	//	toStrings := Collect(func(k string, v int) string {
	//	    return fmt.Sprintf("%s=%d", k, v)
	//	})
	Collector[K comparable, V, R any] = func(K, V) R

	// Entry represents a single key-value pair from a record.
	// This is an alias for Tuple2 to provide semantic clarity.
	//
	// Example:
	//
	//	entries := ToEntries(record)
	//	for _, entry := range entries {
	//	    key := entry.F1
	//	    value := entry.F2
	//	}
	Entry[K comparable, V any] = pair.Pair[K, V]

	// Entries is a slice of key-value pairs.
	//
	// Example:
	//
	//	entries := Entries[string, int]{
	//	    T.MakeTuple2("a", 1),
	//	    T.MakeTuple2("b", 2),
	//	}
	//	record := FromEntries(entries)
	Entries[K comparable, V any] = []Entry[K, V]
)
