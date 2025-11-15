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

// Package array provides functional programming utilities for working with Go slices.
//
// This package treats Go slices as immutable arrays and provides a rich set of operations
// for transforming, filtering, folding, and combining arrays in a functional style.
// All operations return new arrays rather than modifying existing ones.
//
// # Core Concepts
//
// The array package implements several functional programming abstractions:
//   - Functor: Transform array elements with Map
//   - Applicative: Apply functions in arrays to values in arrays
//   - Monad: Chain operations that produce arrays with Chain/FlatMap
//   - Foldable: Reduce arrays to single values with Reduce/Fold
//   - Traversable: Transform arrays while preserving structure
//
// # Basic Operations
//
//	// Creating arrays
//	arr := array.From(1, 2, 3, 4, 5)
//	repeated := array.Replicate(3, "hello")
//	generated := array.MakeBy(5, func(i int) int { return i * 2 })
//
//	// Transforming arrays
//	doubled := array.Map(N.Mul(2))(arr)
//	filtered := array.Filter(func(x int) bool { return x > 2 })(arr)
//
//	// Combining arrays
//	combined := array.Flatten([][]int{{1, 2}, {3, 4}})
//	zipped := array.Zip([]string{"a", "b"})([]int{1, 2})
//
// # Mapping and Filtering
//
// Transform array elements with Map, or filter elements with Filter:
//
//	numbers := []int{1, 2, 3, 4, 5}
//
//	// Map transforms each element
//	doubled := array.Map(N.Mul(2))(numbers)
//	// Result: [2, 4, 6, 8, 10]
//
//	// Filter keeps elements matching a predicate
//	evens := array.Filter(func(x int) bool { return x%2 == 0 })(numbers)
//	// Result: [2, 4]
//
//	// FilterMap combines both operations
//	import "github.com/IBM/fp-go/v2/option"
//	result := array.FilterMap(func(x int) option.Option[int] {
//	    if x%2 == 0 {
//	        return option.Some(x * 2)
//	    }
//	    return option.None[int]()
//	})(numbers)
//	// Result: [4, 8]
//
// # Folding and Reducing
//
// Reduce arrays to single values:
//
//	numbers := []int{1, 2, 3, 4, 5}
//
//	// Sum all elements
//	sum := array.Reduce(func(acc, x int) int { return acc + x }, 0)(numbers)
//	// Result: 15
//
//	// Using a Monoid
//	import "github.com/IBM/fp-go/v2/monoid"
//	sum := array.Fold(monoid.MonoidSum[int]())(numbers)
//	// Result: 15
//
// # Chaining Operations
//
// Chain operations that produce arrays (also known as FlatMap):
//
//	numbers := []int{1, 2, 3}
//	result := array.Chain(func(x int) []int {
//	    return []int{x, x * 10}
//	})(numbers)
//	// Result: [1, 10, 2, 20, 3, 30]
//
// # Finding Elements
//
// Search for elements matching predicates:
//
//	numbers := []int{1, 2, 3, 4, 5}
//
//	// Find first element > 3
//	first := array.FindFirst(func(x int) bool { return x > 3 })(numbers)
//	// Result: Some(4)
//
//	// Find last element > 3
//	last := array.FindLast(func(x int) bool { return x > 3 })(numbers)
//	// Result: Some(5)
//
//	// Get head and tail
//	head := array.Head(numbers) // Some(1)
//	tail := array.Tail(numbers) // Some([2, 3, 4, 5])
//
// # Sorting
//
// Sort arrays using Ord instances:
//
//	import "github.com/IBM/fp-go/v2/ord"
//
//	numbers := []int{3, 1, 4, 1, 5}
//	sorted := array.Sort(ord.FromStrictCompare[int]())(numbers)
//	// Result: [1, 1, 3, 4, 5]
//
//	// Sort by extracted key
//	type Person struct { Name string; Age int }
//	people := []Person{{"Alice", 30}, {"Bob", 25}}
//	byAge := array.SortByKey(ord.FromStrictCompare[int](), func(p Person) int {
//	    return p.Age
//	})(people)
//
// # Uniqueness
//
// Remove duplicate elements:
//
//	numbers := []int{1, 2, 2, 3, 3, 3}
//	unique := array.StrictUniq(numbers)
//	// Result: [1, 2, 3]
//
//	// Unique by key
//	type Person struct { Name string; Age int }
//	people := []Person{{"Alice", 30}, {"Bob", 25}, {"Alice", 35}}
//	uniqueByName := array.Uniq(func(p Person) string { return p.Name })(people)
//	// Result: [{"Alice", 30}, {"Bob", 25}]
//
// # Zipping
//
// Combine multiple arrays:
//
//	names := []string{"Alice", "Bob", "Charlie"}
//	ages := []int{30, 25, 35}
//
//	// Zip into tuples
//	pairs := array.Zip(ages)(names)
//	// Result: [(Alice, 30), (Bob, 25), (Charlie, 35)]
//
//	// Zip with custom function
//	result := array.ZipWith(names, ages, func(name string, age int) string {
//	    return fmt.Sprintf("%s is %d", name, age)
//	})
//
// # Monadic Do Notation
//
// Build complex array computations using do-notation style:
//
//	result := array.Do(
//	    struct{ X, Y int }{},
//	)(
//	    array.Bind(
//	        func(x int) func(s struct{}) struct{ X int } {
//	            return func(s struct{}) struct{ X int } { return struct{ X int }{x} }
//	        },
//	        func(s struct{}) []int { return []int{1, 2, 3} },
//	    ),
//	    array.Bind(
//	        func(y int) func(s struct{ X int }) struct{ X, Y int } {
//	            return func(s struct{ X int }) struct{ X, Y int } {
//	                return struct{ X, Y int }{s.X, y}
//	            }
//	        },
//	        func(s struct{ X int }) []int { return []int{4, 5} },
//	    ),
//	)
//	// Produces all combinations: [{1,4}, {1,5}, {2,4}, {2,5}, {3,4}, {3,5}]
//
// # Sequence and Traverse
//
// Transform arrays of effects into effects of arrays:
//
//	import "github.com/IBM/fp-go/v2/option"
//
//	// Sequence: []Option[A] -> Option[[]A]
//	opts := []option.Option[int]{
//	    option.Some(1),
//	    option.Some(2),
//	    option.Some(3),
//	}
//	result := array.ArrayOption[int]()(opts)
//	// Result: Some([1, 2, 3])
//
//	// If any is None, result is None
//	opts2 := []option.Option[int]{
//	    option.Some(1),
//	    option.None[int](),
//	    option.Some(3),
//	}
//	result2 := array.ArrayOption[int]()(opts2)
//	// Result: None
//
// # Equality and Comparison
//
// Compare arrays for equality:
//
//	import "github.com/IBM/fp-go/v2/eq"
//
//	eq := array.Eq(eq.FromStrictEquals[int]())
//	equal := eq.Equals([]int{1, 2, 3}, []int{1, 2, 3})
//	// Result: true
//
// # Monoid Operations
//
// Combine arrays using monoid operations:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//
//	// Concatenate arrays
//	m := array.Monoid[int]()
//	result := m.Concat([]int{1, 2}, []int{3, 4})
//	// Result: [1, 2, 3, 4]
//
//	// Concatenate multiple arrays efficiently
//	result := array.ArrayConcatAll(
//	    []int{1, 2},
//	    []int{3, 4},
//	    []int{5, 6},
//	)
//	// Result: [1, 2, 3, 4, 5, 6]
//
// # Performance Considerations
//
// Most operations create new arrays rather than modifying existing ones. For performance-critical
// code, consider:
//   - Using Copy for shallow copies when needed
//   - Using Clone with a custom cloning function for deep copies
//   - Batching operations to minimize intermediate allocations
//   - Using ArrayConcatAll for efficient concatenation of multiple arrays
//
// # Subpackages
//
//   - array/generic: Generic implementations for custom array-like types
//   - array/nonempty: Operations for non-empty arrays with compile-time guarantees
//   - array/testing: Testing utilities for array laws and properties
package array
