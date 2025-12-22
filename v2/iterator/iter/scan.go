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

// Scan applies an accumulator function over a sequence, emitting each intermediate result.
//
// This function is similar to Reduce, but instead of returning only the final accumulated value,
// it returns a sequence containing all intermediate accumulated values. Each element in the
// output sequence is the result of applying the accumulator function to the previous accumulated
// value and the current input element.
//
// The operation is lazy - intermediate values are computed only as they are consumed.
//
// RxJS Equivalent: [scan] - https://rxjs.dev/api/operators/scan
//
// Scan is useful for:
//   - Computing running totals or cumulative sums
//   - Tracking state changes over a sequence
//   - Building up complex values incrementally
//   - Generating sequences based on previous values
//
// Type Parameters:
//   - FCT: The accumulator function type, must be ~func(V, U) V
//   - U: The type of elements in the input sequence
//   - V: The type of the accumulated value and elements in the output sequence
//
// Parameters:
//   - f: The accumulator function that takes the current accumulated value and the next
//     input element, returning the new accumulated value
//   - initial: The initial accumulated value (not included in the output sequence)
//
// Returns:
//   - An Operator that transforms a Seq[U] into a Seq[V] containing all intermediate
//     accumulated values
//
// Example - Running sum:
//
//	seq := From(1, 2, 3, 4, 5)
//	runningSum := Scan(func(acc, x int) int { return acc + x }, 0)
//	result := runningSum(seq)
//	// yields: 1, 3, 6, 10, 15
//
// Example - Running product:
//
//	seq := From(2, 3, 4)
//	runningProduct := Scan(func(acc, x int) int { return acc * x }, 1)
//	result := runningProduct(seq)
//	// yields: 2, 6, 24
//
// Example - Building strings:
//
//	seq := From("a", "b", "c")
//	concat := Scan(func(acc, x string) string { return acc + x }, "")
//	result := concat(seq)
//	// yields: "a", "ab", "abc"
//
// Example - Tracking maximum:
//
//	seq := From(3, 1, 4, 1, 5, 9, 2)
//	maxSoFar := Scan(func(acc, x int) int {
//	    if x > acc { return x }
//	    return acc
//	}, 0)
//	result := maxSoFar(seq)
//	// yields: 3, 3, 4, 4, 5, 9, 9
//
// Example - Empty sequence:
//
//	seq := Empty[int]()
//	runningSum := Scan(func(acc, x int) int { return acc + x }, 0)
//	result := runningSum(seq)
//	// yields: nothing (empty sequence)
//
// Example - Single element:
//
//	seq := From(42)
//	runningSum := Scan(func(acc, x int) int { return acc + x }, 10)
//	result := runningSum(seq)
//	// yields: 52
func Scan[FCT ~func(V, U) V, U, V any](f FCT, initial V) Operator[U, V] {
	return func(s Seq[U]) Seq[V] {
		return func(yield func(V) bool) {
			current := initial
			for u := range s {
				current = f(current, u)
				if !yield(current) {
					return
				}
			}
		}
	}
}
