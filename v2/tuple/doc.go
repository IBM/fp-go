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

// Package tuple provides type-safe heterogeneous tuple data structures and operations.
//
// Tuples are immutable data structures that can hold a fixed number of values of different types.
// Unlike arrays or slices which hold homogeneous data, tuples can contain values of different types
// at compile-time known positions.
//
// # Tuple Types
//
// The package provides tuple types from Tuple1 to Tuple15, where the number indicates how many
// elements the tuple contains. Each tuple type is generic over its element types:
//
//	Tuple1[T1]                    // Single element
//	Tuple2[T1, T2]                // Pair
//	Tuple3[T1, T2, T3]            // Triple
//	// ... up to Tuple15
//
// # Creating Tuples
//
// Tuples can be created using the MakeTupleN functions:
//
//	t2 := tuple.MakeTuple2("hello", 42)           // Tuple2[string, int]
//	t3 := tuple.MakeTuple3(1.5, true, "world")    // Tuple3[float64, bool, string]
//
// For single-element tuples, you can also use the Of function:
//
//	t1 := tuple.Of(42)  // Equivalent to MakeTuple1(42)
//
// # Accessing Elements
//
// Tuple elements are accessed via their F1, F2, F3, ... fields:
//
//	t := tuple.MakeTuple2("hello", 42)
//	s := t.F1  // "hello"
//	n := t.F2  // 42
//
// For Tuple2, convenience accessors are provided:
//
//	s := tuple.First(t)   // Same as t.F1
//	n := tuple.Second(t)  // Same as t.F2
//
// # Transforming Tuples
//
// The package provides several transformation functions:
//
// Map functions transform each element independently:
//
//	t := tuple.MakeTuple2(5, "hello")
//	mapper := tuple.Map2(
//	    func(n int) string { return fmt.Sprintf("%d", n) },
//	    S.Size,
//	)
//	result := mapper(t)  // Tuple2[string, int]{"5", 5}
//
// BiMap transforms both elements of a Tuple2:
//
//	mapper := tuple.BiMap(
//	    S.Size,
//	    func(n int) string { return fmt.Sprintf("%d", n*2) },
//	)
//
// Swap exchanges the elements of a Tuple2:
//
//	t := tuple.MakeTuple2("hello", 42)
//	swapped := tuple.Swap(t)  // Tuple2[int, string]{42, "hello"}
//
// # Function Conversion
//
// Tupled and Untupled functions convert between regular multi-parameter functions
// and functions that take/return tuples:
//
//	// Regular function
//	add := func(a, b int) int { return a + b }
//
//	// Convert to tuple-taking function
//	tupledAdd := tuple.Tupled2(add)
//	result := tupledAdd(tuple.MakeTuple2(3, 4))  // 7
//
//	// Convert back
//	untupledAdd := tuple.Untupled2(tupledAdd)
//	result = untupledAdd(3, 4)  // 7
//
// # Array Conversion
//
// Tuples can be converted to and from arrays using transformation functions:
//
//	t := tuple.MakeTuple3(1, 2, 3)
//	toArray := tuple.ToArray3(
//	    func(n int) int { return n },
//	    func(n int) int { return n },
//	    func(n int) int { return n },
//	)
//	arr := toArray(t)  // []int{1, 2, 3}
//
// # Algebraic Operations
//
// The package supports algebraic structures:
//
// Monoid operations for combining tuples:
//
//	import "github.com/IBM/fp-go/v2/monoid"
//	import "github.com/IBM/fp-go/v2/string"
//
//	m := tuple.Monoid2(string.Monoid, monoid.MonoidSum[int]())
//	t1 := tuple.MakeTuple2("hello", 5)
//	t2 := tuple.MakeTuple2(" world", 3)
//	result := m.Concat(t1, t2)  // Tuple2[string, int]{"hello world", 8}
//
// Ord operations for comparing tuples:
//
//	import "github.com/IBM/fp-go/v2/ord"
//
//	o := tuple.Ord2(ord.FromStrictCompare[string](), ord.FromStrictCompare[int]())
//	t1 := tuple.MakeTuple2("a", 1)
//	t2 := tuple.MakeTuple2("b", 2)
//	cmp := o.Compare(t1, t2)  // -1 (t1 < t2)
//
// # JSON Serialization
//
// Tuples support JSON marshaling and unmarshaling as arrays:
//
//	t := tuple.MakeTuple2("hello", 42)
//	data, _ := json.Marshal(t)  // ["hello", 42]
//
//	var t2 tuple.Tuple2[string, int]
//	json.Unmarshal(data, &t2)   // Reconstructs the tuple
//
// # Building Tuples Incrementally
//
// The Push functions allow building larger tuples from smaller ones:
//
//	t1 := tuple.MakeTuple1(42)
//	push := tuple.Push1[int, string]("hello")
//	t2 := push(t1)  // Tuple2[int, string]{42, "hello"}
//
// # Replication
//
// Create tuples with all elements set to the same value:
//
//	t := tuple.Replicate3(42)  // Tuple3[int, int, int]{42, 42, 42}
//
// # Use Cases
//
// Tuples are useful when you need to:
//   - Return multiple values from a function in a type-safe way
//   - Group related but differently-typed values together
//   - Work with fixed-size heterogeneous collections
//   - Implement functional programming patterns like bifunctors
//   - Avoid defining custom struct types for simple data groupings
//
// For homogeneous collections of unknown or variable size, consider using
// arrays or slices instead.
package tuple

//go:generate go run .. tuple --count 15 --filename gen.go
