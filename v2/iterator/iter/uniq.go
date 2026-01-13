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

package iter

import (
	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
)

// Uniq returns an operator that filters a sequence to contain only unique elements,
// where uniqueness is determined by a key extraction function.
//
// This function takes a key extraction function and returns an operator that removes
// duplicate elements from a sequence. Two elements are considered duplicates if the
// key extraction function returns the same key for both. Only the first occurrence
// of each unique key is kept in the output sequence.
//
// The operation maintains a map of seen keys internally, so memory usage grows with
// the number of unique keys encountered. The operation is lazy - elements are processed
// and filtered as they are consumed.
//
// RxJS Equivalent: [distinct] - https://rxjs.dev/api/operators/distinct
//
// Type Parameters:
//   - A: The type of elements in the sequence
//   - K: The type of the key used for uniqueness comparison (must be comparable)
//
// Parameters:
//   - f: A function that extracts a comparable key from each element
//
// Returns:
//   - An Operator that filters the sequence to contain only unique elements based on the key
//
// Example - Remove duplicate integers:
//
//	seq := From(1, 2, 3, 2, 4, 1, 5)
//	unique := Uniq(reader.Ask[int]())
//	result := unique(seq)
//	// yields: 1, 2, 3, 4, 5
//
// Example - Unique by string length:
//
//	seq := From("a", "bb", "c", "dd", "eee")
//	uniqueByLength := Uniq(S.Size)
//	result := uniqueByLength(seq)
//	// yields: "a", "bb", "eee" (first occurrence of each length)
//
// Example - Unique structs by field:
//
//	type Person struct { ID int; Name string }
//	seq := From(
//	    Person{1, "Alice"},
//	    Person{2, "Bob"},
//	    Person{1, "Alice2"},  // duplicate ID
//	)
//	uniqueByID := Uniq(func(p Person) int { return p.ID })
//	result := uniqueByID(seq)
//	// yields: Person{1, "Alice"}, Person{2, "Bob"}
//
// Example - Case-insensitive unique strings:
//
//	seq := From("Hello", "world", "HELLO", "World", "test")
//	uniqueCaseInsensitive := Uniq(func(s string) string {
//	    return strings.ToLower(s)
//	})
//	result := uniqueCaseInsensitive(seq)
//	// yields: "Hello", "world", "test"
//
// Example - Empty sequence:
//
//	seq := Empty[int]()
//	unique := Uniq(reader.Ask[int]())
//	result := unique(seq)
//	// yields: nothing (empty sequence)
//
// Example - All duplicates:
//
//	seq := From(1, 1, 1, 1)
//	unique := Uniq(reader.Ask[int]())
//	result := unique(seq)
//	// yields: 1 (only first occurrence)
func Uniq[A any, K comparable](f func(A) K) Operator[A, A] {
	return func(s Seq[A]) Seq[A] {
		return func(yield func(A) bool) {
			items := make(map[K]Void)
			for a := range s {
				k := f(a)
				if _, ok := items[k]; !ok {
					items[k] = function.VOID
					if !yield(a) {
						return
					}
				}
			}
		}
	}
}

// StrictUniq filters a sequence to contain only unique elements using direct comparison.
//
// This is a convenience function that uses the identity function as the key extractor,
// meaning elements are compared directly for uniqueness. It's equivalent to calling
// Uniq with the identity function, but provides a simpler API when the elements
// themselves are comparable.
//
// The operation maintains a map of seen elements internally, so memory usage grows with
// the number of unique elements. Only the first occurrence of each unique element is kept.
//
// RxJS Equivalent: [distinct] - https://rxjs.dev/api/operators/distinct
//
// Type Parameters:
//   - A: The type of elements in the sequence (must be comparable)
//
// Parameters:
//   - as: The input sequence to filter for unique elements
//
// Returns:
//   - A sequence containing only the first occurrence of each unique element
//
// Example - Remove duplicate integers:
//
//	seq := From(1, 2, 3, 2, 4, 1, 5)
//	result := StrictUniq(seq)
//	// yields: 1, 2, 3, 4, 5
//
// Example - Remove duplicate strings:
//
//	seq := From("apple", "banana", "apple", "cherry", "banana")
//	result := StrictUniq(seq)
//	// yields: "apple", "banana", "cherry"
//
// Example - Single element:
//
//	seq := From(42)
//	result := StrictUniq(seq)
//	// yields: 42
//
// Example - All duplicates:
//
//	seq := From("x", "x", "x")
//	result := StrictUniq(seq)
//	// yields: "x" (only first occurrence)
//
// Example - Empty sequence:
//
//	seq := Empty[int]()
//	result := StrictUniq(seq)
//	// yields: nothing (empty sequence)
//
// Example - Already unique:
//
//	seq := From(1, 2, 3, 4, 5)
//	result := StrictUniq(seq)
//	// yields: 1, 2, 3, 4, 5 (no changes)
func StrictUniq[A comparable](as Seq[A]) Seq[A] {
	return Uniq(F.Identity[A])(as)
}
