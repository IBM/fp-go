// Copyright (c) 2024 IBM Corp.
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

// Package iter provides functional programming utilities for working with Go 1.23+ iterators.
// It offers operations for reducing, mapping, concatenating, and transforming iterator sequences
// in a functional style, compatible with the range-over-func pattern.
package iter

import (
	"slices"

	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
)

func From[A any](as ...A) Seq[A] {
	return slices.Values(as)
}

// MonadReduceWithIndex reduces an iterator sequence to a single value using a reducer function
// that receives the current index, accumulated value, and current element.
//
// The function iterates through all elements in the sequence, applying the reducer function
// at each step with the element's index. This is useful when the position of elements matters
// in the reduction logic.
//
// Parameters:
//   - fa: The iterator sequence to reduce
//   - f: The reducer function that takes (index, accumulator, element) and returns the new accumulator
//   - initial: The initial value for the accumulator
//
// Returns:
//   - The final accumulated value after processing all elements
//
// Example:
//
//	iter := func(yield func(int) bool) {
//	    yield(10)
//	    yield(20)
//	    yield(30)
//	}
//	// Sum with index multiplier: 0*10 + 1*20 + 2*30 = 80
//	result := MonadReduceWithIndex(iter, func(i, acc, val int) int {
//	    return acc + i*val
//	}, 0)
func MonadReduceWithIndex[GA ~func(yield func(A) bool), A, B any](fa GA, f func(int, B, A) B, initial B) B {
	current := initial
	var i int
	for a := range fa {
		current = f(i, current, a)
		i += 1
	}
	return current
}

// MonadReduce reduces an iterator sequence to a single value using a reducer function.
//
// This is similar to MonadReduceWithIndex but without index tracking, making it more
// efficient when the position of elements is not needed in the reduction logic.
//
// Parameters:
//   - fa: The iterator sequence to reduce
//   - f: The reducer function that takes (accumulator, element) and returns the new accumulator
//   - initial: The initial value for the accumulator
//
// Returns:
//   - The final accumulated value after processing all elements
//
// Example:
//
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	    yield(3)
//	}
//	sum := MonadReduce(iter, func(acc, val int) int {
//	    return acc + val
//	}, 0) // Returns: 6
func MonadReduce[GA ~func(yield func(A) bool), A, B any](fa GA, f func(B, A) B, initial B) B {
	current := initial
	for a := range fa {
		current = f(current, a)
	}
	return current
}

// Concat concatenates two iterator sequences, yielding all elements from left followed by all elements from right.
//
// The resulting iterator will first yield all elements from the left sequence, then all elements
// from the right sequence. If the consumer stops early (yield returns false), iteration stops
// immediately without processing remaining elements.
//
// Parameters:
//   - left: The first iterator sequence
//   - right: The second iterator sequence
//
// Returns:
//   - A new iterator that yields elements from both sequences in order
//
// Example:
//
//	left := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	}
//	right := func(yield func(int) bool) {
//	    yield(3)
//	    yield(4)
//	}
//	combined := Concat(left, right) // Yields: 1, 2, 3, 4
func Concat[GT ~func(yield func(T) bool), T any](left, right GT) GT {
	return func(yield func(T) bool) {
		for t := range left {
			if !yield(t) {
				return
			}
		}
		for t := range right {
			if !yield(t) {
				return
			}
		}
	}
}

// Of creates an iterator sequence containing a single element.
//
// This is the unit/return operation for the iterator monad, lifting a single value
// into the iterator context.
//
// Parameters:
//   - a: The element to wrap in an iterator
//
// Returns:
//   - An iterator that yields exactly one element
//
// Example:
//
//	iter := Of[func(yield func(int) bool)](42)
//	// Yields: 42
func Of[GA ~func(yield func(A) bool), A any](a A) GA {
	return func(yield func(A) bool) {
		yield(a)
	}
}

// MonadAppend appends a single element to the end of an iterator sequence.
//
// This creates a new iterator that yields all elements from the original sequence
// followed by the tail element.
//
// Parameters:
//   - f: The original iterator sequence
//   - tail: The element to append
//
// Returns:
//   - A new iterator with the tail element appended
//
// Example:
//
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	}
//	result := MonadAppend(iter, 3) // Yields: 1, 2, 3
func MonadAppend[GA ~func(yield func(A) bool), A any](f GA, tail A) GA {
	return Concat(f, Of[GA](tail))
}

// Append returns a function that appends a single element to the end of an iterator sequence.
//
// This is the curried version of MonadAppend, useful for partial application and composition.
//
// Parameters:
//   - tail: The element to append
//
// Returns:
//   - A function that takes an iterator and returns a new iterator with the tail element appended
//
// Example:
//
//	appendThree := Append[func(yield func(int) bool)](3)
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	}
//	result := appendThree(iter) // Yields: 1, 2, 3
func Append[GA ~func(yield func(A) bool), A any](tail A) func(GA) GA {
	return F.Bind2nd(Concat[GA], Of[GA](tail))
}

// Prepend returns a function that prepends a single element to the beginning of an iterator sequence.
//
// This is the curried version for prepending, useful for partial application and composition.
//
// Parameters:
//   - head: The element to prepend
//
// Returns:
//   - A function that takes an iterator and returns a new iterator with the head element prepended
//
// Example:
//
//	prependZero := Prepend[func(yield func(int) bool)](0)
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	}
//	result := prependZero(iter) // Yields: 0, 1, 2
func Prepend[GA ~func(yield func(A) bool), A any](head A) func(GA) GA {
	return F.Bind1st(Concat[GA], Of[GA](head))
}

// Empty creates an empty iterator sequence that yields no elements.
//
// This is the identity element for the Concat operation and represents an empty collection
// in the iterator context.
//
// Returns:
//   - An iterator that yields no elements
//
// Example:
//
//	iter := Empty[func(yield func(int) bool), int]()
//	// Yields nothing
func Empty[GA ~func(yield func(A) bool), A any]() GA {
	return func(_ func(A) bool) {}
}

// ToArray collects all elements from an iterator sequence into a slice.
//
// This eagerly evaluates the entire iterator sequence and materializes all elements
// into memory as a slice.
//
// Parameters:
//   - fa: The iterator sequence to collect
//
// Returns:
//   - A slice containing all elements from the iterator
//
// Example:
//
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	    yield(3)
//	}
//	arr := ToArray[func(yield func(int) bool), []int](iter) // Returns: []int{1, 2, 3}
func ToArray[GA ~func(yield func(A) bool), GB ~[]A, A any](fa GA) GB {
	bs := make(GB, 0)
	for a := range fa {
		bs = append(bs, a)
	}
	return bs
}

// MonadMapToArray maps each element of an iterator sequence through a function and collects the results into a slice.
//
// This combines mapping and collection into a single operation, eagerly evaluating the entire
// iterator sequence and materializing the transformed elements into memory.
//
// Parameters:
//   - fa: The iterator sequence to map and collect
//   - f: The mapping function to apply to each element
//
// Returns:
//   - A slice containing the mapped elements
//
// Example:
//
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	    yield(3)
//	}
//	doubled := MonadMapToArray[func(yield func(int) bool), []int](iter, func(x int) int {
//	    return x * 2
//	}) // Returns: []int{2, 4, 6}
func MonadMapToArray[GA ~func(yield func(A) bool), GB ~[]B, A, B any](fa GA, f func(A) B) GB {
	bs := make(GB, 0)
	for a := range fa {
		bs = append(bs, f(a))
	}
	return bs
}

// MapToArray returns a function that maps each element through a function and collects the results into a slice.
//
// This is the curried version of MonadMapToArray, useful for partial application and composition.
//
// Parameters:
//   - f: The mapping function to apply to each element
//
// Returns:
//   - A function that takes an iterator and returns a slice of mapped elements
//
// Example:
//
//	double := MapToArray[func(yield func(int) bool), []int](func(x int) int {
//	    return x * 2
//	})
//	iter := func(yield func(int) bool) {
//	    yield(1)
//	    yield(2)
//	}
//	result := double(iter) // Returns: []int{2, 4}
func MapToArray[GA ~func(yield func(A) bool), GB ~[]B, A, B any](f func(A) B) func(GA) GB {
	return F.Bind2nd(MonadMapToArray[GA, GB], f)
}

// MonadMapToArrayWithIndex maps each element of an iterator sequence through a function that receives
// the element's index, and collects the results into a slice.
//
// This is similar to MonadMapToArray but the mapping function also receives the zero-based index
// of each element, useful when the position matters in the transformation logic.
//
// Parameters:
//   - fa: The iterator sequence to map and collect
//   - f: The mapping function that takes (index, element) and returns the transformed element
//
// Returns:
//   - A slice containing the mapped elements
//
// Example:
//
//	iter := func(yield func(string) bool) {
//	    yield("a")
//	    yield("b")
//	    yield("c")
//	}
//	indexed := MonadMapToArrayWithIndex[func(yield func(string) bool), []string](iter,
//	    func(i int, s string) string {
//	        return fmt.Sprintf("%d:%s", i, s)
//	    }) // Returns: []string{"0:a", "1:b", "2:c"}
func MonadMapToArrayWithIndex[GA ~func(yield func(A) bool), GB ~[]B, A, B any](fa GA, f func(int, A) B) GB {
	bs := make(GB, 0)
	var i int
	for a := range fa {
		bs = append(bs, f(i, a))
		i += 1
	}
	return bs
}

// MapToArrayWithIndex returns a function that maps each element through an indexed function
// and collects the results into a slice.
//
// This is the curried version of MonadMapToArrayWithIndex, useful for partial application and composition.
//
// Parameters:
//   - f: The mapping function that takes (index, element) and returns the transformed element
//
// Returns:
//   - A function that takes an iterator and returns a slice of mapped elements
//
// Example:
//
//	addIndex := MapToArrayWithIndex[func(yield func(string) bool), []string](
//	    func(i int, s string) string {
//	        return fmt.Sprintf("%d:%s", i, s)
//	    })
//	iter := func(yield func(string) bool) {
//	    yield("a")
//	    yield("b")
//	}
//	result := addIndex(iter) // Returns: []string{"0:a", "1:b"}
func MapToArrayWithIndex[GA ~func(yield func(A) bool), GB ~[]B, A, B any](f func(int, A) B) func(GA) GB {
	return F.Bind2nd(MonadMapToArrayWithIndex[GA, GB], f)
}

// Monoid returns a Monoid instance for iterator sequences.
//
// The monoid uses Concat as the binary operation and Empty as the identity element,
// allowing iterator sequences to be combined in an associative way with a neutral element.
// This enables generic operations that work with any monoid, such as folding a collection
// of iterators into a single iterator.
//
// Returns:
//   - A Monoid instance with Concat and Empty operations
//
// Example:
//
//	m := Monoid[func(yield func(int) bool), int]()
//	iter1 := func(yield func(int) bool) { yield(1); yield(2) }
//	iter2 := func(yield func(int) bool) { yield(3); yield(4) }
//	combined := m.Concat(iter1, iter2) // Yields: 1, 2, 3, 4
//	empty := m.Empty() // Yields nothing
func Monoid[GA ~func(yield func(A) bool), A any]() M.Monoid[GA] {
	return M.MakeMonoid(Concat[GA], Empty[GA]())
}
