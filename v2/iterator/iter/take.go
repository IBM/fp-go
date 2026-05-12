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
// Marble Diagram:
//
//	Input:  --1--2--3--4--5--6--7--8-->
//	Take(3)
//	Output: --1--2--3|
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

func takeWhile[U any](p Predicate[U], inclusive bool) Operator[U, U] {
	return func(s Seq[U]) Seq[U] {
		return func(yield func(U) bool) {
			for u := range s {
				if !p(u) {
					if inclusive {
						yield(u)
					}
					return
				}
				if !yield(u) {
					return
				}
			}
		}
	}
}

// TakeWhile returns an operator that emits elements from a sequence while a predicate is satisfied.
//
// This function creates a transformation that yields elements from the source sequence
// as long as each element satisfies the provided predicate. Once an element fails the
// predicate test, the sequence terminates immediately, and no further elements are
// emitted, even if subsequent elements would satisfy the predicate.
//
// The operation is lazy and only consumes elements from the source sequence as needed.
// Once the predicate returns false, iteration stops immediately without consuming
// the remaining elements from the source.
//
// Marble Diagram:
//
//	Input:       --1--2--3--4--5--2--1-->
//	TakeWhile(x < 4)
//	Output:      --1--2--3|
//	                      (stops at 4)
//
// RxJS Equivalent: [takeWhile] - https://rxjs.dev/api/operators/takeWhile
//
// Type Parameters:
//   - U: The type of elements in the sequence
//
// Parameters:
//   - p: A predicate function that tests each element. Returns true to continue, false to stop
//
// Returns:
//   - An Operator that transforms a Seq[U] by taking elements while the predicate is satisfied
//
// Example - Take while less than threshold:
//
//	seq := From(1, 2, 3, 4, 5, 2, 1)
//	result := TakeWhile(func(x int) bool { return x < 4 })(seq)
//	// yields: 1, 2, 3 (stops at 4, doesn't continue to 2, 1)
//
// Example - Take while condition is met:
//
//	seq := From("a", "b", "c", "1", "d", "e")
//	isLetter := func(s string) bool { return s >= "a" && s <= "z" }
//	result := TakeWhile(isLetter)(seq)
//	// yields: "a", "b", "c" (stops at "1")
//
// Example - Take all when predicate always true:
//
//	seq := From(2, 4, 6, 8)
//	result := TakeWhile(func(x int) bool { return x%2 == 0 })(seq)
//	// yields: 2, 4, 6, 8 (all elements satisfy predicate)
//
// Example - Take none when first element fails:
//
//	seq := From(5, 1, 2, 3)
//	result := TakeWhile(func(x int) bool { return x < 5 })(seq)
//	// yields: nothing (first element fails predicate)
//
// Example - Chaining with other operations:
//
//	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	result := F.Pipe2(
//	    seq,
//	    MonadMap(seq, func(x int) int { return x * 2 }),
//	    TakeWhile(func(x int) bool { return x < 10 }),
//	)
//	// yields: 2, 4, 6, 8 (stops when doubled value reaches 10)
func TakeWhile[U any](p Predicate[U]) Operator[U, U] {
	return takeWhile(p, false)
}

// TakeWhileInclusive returns an operator that emits elements from a sequence while a predicate is satisfied, including the first element that fails the predicate.
//
// This function creates a transformation that yields elements from the source sequence
// as long as each element satisfies the provided predicate. Unlike TakeWhile, when an
// element fails the predicate test, that element is still emitted before the sequence
// terminates. This is useful when you want to include the boundary element that caused
// the condition to fail.
//
// The operation is lazy and only consumes elements from the source sequence as needed.
// Once the predicate returns false, the failing element is yielded, then iteration stops
// immediately without consuming the remaining elements from the source.
//
// Marble Diagram:
//
//	Input:                --1--2--3--4--5--2--1-->
//	TakeWhileInclusive(x < 4)
//	Output:               --1--2--3--4|
//	                               (includes 4, then stops)
//
// RxJS Equivalent: [takeWhile with inclusive option] - https://rxjs.dev/api/operators/takeWhile
//
// Type Parameters:
//   - U: The type of elements in the sequence
//
// Parameters:
//   - p: A predicate function that tests each element. Returns true to continue, false to stop after including the current element
//
// Returns:
//   - An Operator that transforms a Seq[U] by taking elements while the predicate is satisfied, plus the first failing element
//
// Example - Take while less than threshold, including boundary:
//
//	seq := From(1, 2, 3, 4, 5, 2, 1)
//	result := TakeWhileInclusive(func(x int) bool { return x < 4 })(seq)
//	// yields: 1, 2, 3, 4 (includes 4 which fails the predicate)
//
// Example - Take while condition is met, including first non-letter:
//
//	seq := From("a", "b", "c", "1", "d", "e")
//	isLetter := func(s string) bool { return s >= "a" && s <= "z" }
//	result := TakeWhileInclusive(isLetter)(seq)
//	// yields: "a", "b", "c", "1" (includes "1" which fails the predicate)
//
// Example - Take all when predicate always true:
//
//	seq := From(2, 4, 6, 8)
//	result := TakeWhileInclusive(func(x int) bool { return x%2 == 0 })(seq)
//	// yields: 2, 4, 6, 8 (all elements satisfy predicate, no failing element)
//
// Example - Take only first when it fails:
//
//	seq := From(5, 1, 2, 3)
//	result := TakeWhileInclusive(func(x int) bool { return x < 5 })(seq)
//	// yields: 5 (first element fails, but is included)
//
// Example - Chaining with other operations:
//
//	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	result := F.Pipe2(
//	    seq,
//	    MonadMap(seq, func(x int) int { return x * 2 }),
//	    TakeWhileInclusive(func(x int) bool { return x < 10 }),
//	)
//	// yields: 2, 4, 6, 8, 10 (includes 10 which fails x < 10)
//
// See Also:
//   - TakeWhile: Similar but excludes the first failing element
//   - Take: Takes a fixed number of elements
//   - SkipWhile: Skips elements while predicate is satisfied
func TakeWhileInclusive[U any](p Predicate[U]) Operator[U, U] {
	return takeWhile(p, true)
}

// SkipWhile returns an operator that skips elements from a sequence while a predicate is satisfied.
//
// This function creates a transformation that discards elements from the source sequence
// as long as each element satisfies the provided predicate. Once an element fails the
// predicate test, that element and all subsequent elements are yielded, regardless of
// whether they satisfy the predicate.
//
// The operation is lazy and only consumes elements from the source sequence as needed.
// Once the predicate returns false, all remaining elements are yielded without further
// predicate evaluation.
//
// Marble Diagram:
//
//	Input:        --1--2--3--4--5--2--1-->
//	SkipWhile(x < 4)
//	Output:       -----------4--5--2--1-->
//	                         (starts at 4, continues with all)
//
// RxJS Equivalent: [skipWhile] - https://rxjs.dev/api/operators/skipWhile
//
// Type Parameters:
//   - U: The type of elements in the sequence
//
// Parameters:
//   - p: A predicate function that tests each element. Returns true to skip, false to start yielding
//
// Returns:
//   - An Operator that transforms a Seq[U] by skipping elements while the predicate is satisfied
//
// Example - Skip while less than threshold:
//
//	seq := From(1, 2, 3, 4, 5, 2, 1)
//	result := SkipWhile(func(x int) bool { return x < 4 })(seq)
//	// yields: 4, 5, 2, 1 (starts at 4, continues with all remaining)
//
// Example - Skip while condition is met:
//
//	seq := From("a", "b", "c", "1", "d", "e")
//	isLetter := func(s string) bool { return s >= "a" && s <= "z" }
//	result := SkipWhile(isLetter)(seq)
//	// yields: "1", "d", "e" (starts at "1", continues with all remaining)
//
// Example - Skip none when first element fails:
//
//	seq := From(5, 1, 2, 3)
//	result := SkipWhile(func(x int) bool { return x < 5 })(seq)
//	// yields: 5, 1, 2, 3 (first element fails predicate, all yielded)
//
// Example - Skip all when predicate always true:
//
//	seq := From(2, 4, 6, 8)
//	result := SkipWhile(func(x int) bool { return x%2 == 0 })(seq)
//	// yields: nothing (all elements satisfy predicate)
//
// Example - Chaining with other operations:
//
//	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//	result := F.Pipe2(
//	    seq,
//	    SkipWhile(func(x int) bool { return x < 5 }),
//	    MonadMap(seq, func(x int) int { return x * 2 }),
//	)
//	// yields: 10, 12, 14, 16, 18, 20 (skip until 5, then double remaining)
func SkipWhile[U any](p Predicate[U]) Operator[U, U] {
	return func(s Seq[U]) Seq[U] {
		return func(yield func(U) bool) {
			skipping := true
			for u := range s {
				if skipping && p(u) {
					continue
				}
				skipping = false
				if !yield(u) {
					return
				}
			}
		}
	}
}
