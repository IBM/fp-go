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

// Cycle creates a sequence that repeats the elements of the input sequence indefinitely.
//
// This function takes a finite sequence and creates an infinite sequence by cycling through
// all elements repeatedly. When the end of the input sequence is reached, it starts over
// from the beginning, continuing this pattern forever.
//
// RxJS Equivalent: [repeat] - https://rxjs.dev/api/operators/repeat
//
// WARNING: This creates an INFINITE sequence for non-empty inputs. It must be used with
// operations that limit the output (such as Take, First, or early termination in iteration)
// to avoid infinite loops.
//
// If the input sequence is empty, Cycle returns an empty sequence immediately. It does NOT
// loop indefinitely - the result is simply an empty sequence.
//
// The operation is lazy - elements are only generated as they are consumed. The input sequence
// is re-iterated each time the cycle completes, so any side effects in the source sequence
// will be repeated.
//
// Type Parameters:
//   - U: The type of elements in the sequence
//
// Parameters:
//   - ma: The input sequence to cycle through. Should be finite.
//
// Returns:
//   - An infinite sequence that repeats the elements of the input sequence
//
// Example - Basic cycling with Take:
//
//	seq := From(1, 2, 3)
//	cycled := Cycle(seq)
//	result := Take[int](7)(cycled)
//	// yields: 1, 2, 3, 1, 2, 3, 1
//
// Example - Cycling strings:
//
//	seq := From("A", "B", "C")
//	cycled := Cycle(seq)
//	result := Take[string](5)(cycled)
//	// yields: "A", "B", "C", "A", "B"
//
// Example - Using with First:
//
//	seq := From(10, 20, 30)
//	cycled := Cycle(seq)
//	first := First(cycled)
//	// returns: Some(10)
//
// Example - Combining with filter and take:
//
//	seq := From(1, 2, 3, 4, 5)
//	cycled := Cycle(seq)
//	evens := MonadFilter(cycled, func(x int) bool { return x%2 == 0 })
//	result := Take[int](5)(evens)
//	// yields: 2, 4, 2, 4, 2 (cycles through even numbers)
//
// Example - Empty sequence (returns empty, does not loop):
//
//	seq := Empty[int]()
//	cycled := Cycle(seq)
//	result := Take[int](10)(cycled)
//	// yields: nothing (empty sequence, terminates immediately)
func Cycle[U any](ma Seq[U]) Seq[U] {
	return func(yield func(U) bool) {
		for {
			isEmpty := true
			for u := range ma {
				if !yield(u) {
					return
				}
				isEmpty = false
			}
			if isEmpty {
				return
			}
		}
	}
}
