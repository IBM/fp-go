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

package array

import (
	G "github.com/IBM/fp-go/v2/array/generic"
	O "github.com/IBM/fp-go/v2/ord"
)

// Sort implements a stable sort on the array given the provided ordering.
// The sort is stable, meaning that elements that compare equal retain their original order.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/ord"
//
//	numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}
//	sorted := array.Sort(ord.FromStrictCompare[int]())(numbers)
//	// Result: [1, 1, 2, 3, 4, 5, 6, 9]
//
//go:inline
func Sort[T any](ord O.Ord[T]) Operator[T, T] {
	return G.Sort[[]T](ord)
}

// SortByKey implements a stable sort on the array given the provided ordering on an extracted key.
// This is useful when you want to sort complex types by a specific field.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/ord"
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	people := []Person{
//	    {"Alice", 30},
//	    {"Bob", 25},
//	    {"Charlie", 35},
//	}
//
//	sortByAge := array.SortByKey(
//	    ord.FromStrictCompare[int](),
//	    func(p Person) int { return p.Age },
//	)
//	sorted := sortByAge(people)
//	// Result: [{"Bob", 25}, {"Alice", 30}, {"Charlie", 35}]
//
//go:inline
func SortByKey[K, T any](ord O.Ord[K], f func(T) K) Operator[T, T] {
	return G.SortByKey[[]T](ord, f)
}

// SortBy implements a stable sort on the array using multiple ordering criteria.
// The orderings are applied in sequence: if two elements are equal according to the first
// ordering, the second ordering is used, and so on.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/ord"
//
//	type Person struct {
//	    LastName  string
//	    FirstName string
//	}
//
//	people := []Person{
//	    {"Smith", "John"},
//	    {"Smith", "Alice"},
//	    {"Jones", "Bob"},
//	}
//
//	sortByName := array.SortBy([]ord.Ord[Person]{
//	    ord.Contramap(func(p Person) string { return p.LastName })(ord.FromStrictCompare[string]()),
//	    ord.Contramap(func(p Person) string { return p.FirstName })(ord.FromStrictCompare[string]()),
//	})
//	sorted := sortByName(people)
//	// Result: [{"Jones", "Bob"}, {"Smith", "Alice"}, {"Smith", "John"}]
//
//go:inline
func SortBy[T any](ord []O.Ord[T]) Operator[T, T] {
	return G.SortBy[[]T](ord)
}
