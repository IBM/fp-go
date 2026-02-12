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
	"github.com/IBM/fp-go/v2/pair"
)

// ZipWith applies a function to pairs of elements at the same index in two arrays,
// collecting the results in a new array. If one input array is shorter, excess elements
// of the longer array are discarded.
//
// Example:
//
//	names := []string{"Alice", "Bob", "Charlie"}
//	ages := []int{30, 25, 35}
//
//	result := array.ZipWith(names, ages, func(name string, age int) string {
//	    return fmt.Sprintf("%s is %d years old", name, age)
//	})
//	// Result: ["Alice is 30 years old", "Bob is 25 years old", "Charlie is 35 years old"]
//
//go:inline
func ZipWith[FCT ~func(A, B) C, A, B, C any](fa []A, fb []B, f FCT) []C {
	return G.ZipWith[[]A, []B, []C](fa, fb, f)
}

// Zip takes two arrays and returns an array of corresponding pairs (tuples).
// If one input array is shorter, excess elements of the longer array are discarded.
//
// Example:
//
//	names := []string{"Alice", "Bob", "Charlie"}
//	ages := []int{30, 25, 35}
//
//	pairs := array.Zip(ages)(names)
//	// Result: [(Alice, 30), (Bob, 25), (Charlie, 35)]
//
//	// With different lengths
//	pairs2 := array.Zip([]int{1, 2})([]string{"a", "b", "c"})
//	// Result: [(a, 1), (b, 2)]
//
//go:inline
func Zip[A, B any](fb []B) func([]A) []pair.Pair[A, B] {
	return G.Zip[[]A, []B, []pair.Pair[A, B]](fb)
}

// Unzip is the reverse of Zip. It takes an array of pairs (tuples) and returns
// two corresponding arrays, one containing all first elements and one containing all second elements.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/tuple"
//
//	pairs := []tuple.Tuple2[string, int]{
//	    tuple.MakeTuple2("Alice", 30),
//	    tuple.MakeTuple2("Bob", 25),
//	    tuple.MakeTuple2("Charlie", 35),
//	}
//
//	result := array.Unzip(pairs)
//	// Result: (["Alice", "Bob", "Charlie"], [30, 25, 35])
//	names := result.Head  // ["Alice", "Bob", "Charlie"]
//	ages := result.Tail   // [30, 25, 35]
//
//go:inline
func Unzip[A, B any](cs []pair.Pair[A, B]) pair.Pair[[]A, []B] {
	return G.Unzip[[]A, []B](cs)
}
