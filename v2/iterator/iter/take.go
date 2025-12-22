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

import F "github.com/IBM/fp-go/v2/function"

// Take returns an operator that limits the number of elements in a sequence to at most n elements.
//
// This function creates a transformation that takes the first n elements from a sequence
// and discards the rest. If n is less than or equal to 0, it returns an empty sequence.
// If the input sequence has fewer than n elements, all elements are returned.
//
// The operation is lazy and only consumes elements from the source sequence as needed.
// Once n elements have been yielded, iteration stops immediately without consuming
// the remaining elements from the source.
//
// RxJS Equivalent: [take] - https://rxjs.dev/api/operators/take
//
// Type Parameters:
//   - U: The type of elements in the sequence
//
// Parameters:
//   - n: The maximum number of elements to take from the sequence
//
// Returns:
//   - An Operator that transforms a Seq[U] by taking at most n elements
//
// Example - Take first 3 elements:
//
//	seq := From(1, 2, 3, 4, 5)
//	result := Take[int](3)(seq)
//	// yields: 1, 2, 3
//
// Example - Take more than available:
//
//	seq := From(1, 2)
//	result := Take[int](5)(seq)
//	// yields: 1, 2 (all available elements)
//
// Example - Take zero or negative:
//
//	seq := From(1, 2, 3)
//	result := Take[int](0)(seq)
//	// yields: nothing (empty sequence)
//
// Example - Chaining with other operations:
//
//	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	evens := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
//	result := Take[int](3)(evens)
//	// yields: 2, 4, 6 (first 3 even numbers)
func Take[U any](n int) Operator[U, U] {
	if n <= 0 {
		return F.Constant1[Seq[U]](Empty[U]())
	}
	return func(s Seq[U]) Seq[U] {
		return func(yield Predicate[U]) {
			i := 0
			for u := range s {
				if i >= n || !yield(u) {
					return
				}
				i += 1
			}
		}
	}
}
