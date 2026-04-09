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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

// From constructs an array from a set of variadic arguments.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - data: Variadic arguments to include in the array
//
// # Returns
//
//   - A new array containing all provided arguments
//
// # Example
//
//	arr := array.From(1, 2, 3, 4, 5)
//	// arr: []int{1, 2, 3, 4, 5}
//
//go:inline
func From[A any](data ...A) []A {
	return G.From[[]A](data...)
}

// MakeBy returns an array of length n with element i initialized with f(i).
//
// # Type Parameters
//
//   - F: Function type that takes an int and returns A
//   - A: The type of elements in the array
//
// # Parameters
//
//   - n: The length of the array to create
//   - f: Function to generate each element based on its index
//
// # Returns
//
//   - A new array where element at index i equals f(i)
//
// # Example
//
//	squares := array.MakeBy(5, func(i int) int { return i * i })
//	// squares: []int{0, 1, 4, 9, 16}
//
//go:inline
func MakeBy[F ~func(int) A, A any](n int, f F) []A {
	return G.MakeBy[[]A](n, f)
}

// Replicate creates an array containing a value repeated the specified number of times.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - n: The number of times to repeat the value
//   - a: The value to repeat
//
// # Returns
//
//   - A new array containing n copies of a
//
// # Example
//
//	zeros := array.Replicate(5, 0)
//	// zeros: []int{0, 0, 0, 0, 0}
//
//go:inline
func Replicate[A any](n int, a A) []A {
	return G.Replicate[[]A](n, a)
}

// MonadMap applies a function to each element of an array, returning a new array with the results.
// This is the monadic version of Map that takes the array as the first parameter.
//
//go:inline
func MonadMap[A, B any](as []A, f func(A) B) []B {
	return G.MonadMap[[]A, []B](as, f)
}

// MonadMapRef applies a function to a pointer to each element of an array, returning a new array with the results.
// This is useful when you need to access elements by reference without copying.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// # Parameters
//
//   - as: The input array
//   - f: Function that takes a pointer to an element and returns a transformed value
//
// # Returns
//
//   - A new array with transformed elements
//
// # Example
//
//	type Point struct { X, Y int }
//	points := []Point{{1, 2}, {3, 4}}
//	xs := array.MonadMapRef(points, func(p *Point) int { return p.X })
//	// xs: []int{1, 3}
func MonadMapRef[A, B any](as []A, f func(*A) B) []B {
	count := len(as)
	bs := make([]B, count)
	for i := range count {
		bs[i] = f(&as[i])
	}
	return bs
}

// MapWithIndex applies a function to each element and its index in an array, returning a new array with the results.
//
//go:inline
func MapWithIndex[A, B any](f func(int, A) B) Operator[A, B] {
	return G.MapWithIndex[[]A, []B](f)
}

// Map applies a function to each element of an array, returning a new array with the results.
// This is the curried version that returns a function.
//
// Example:
//
//	double := array.Map(N.Mul(2))
//	result := double([]int{1, 2, 3}) // [2, 4, 6]
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return G.Map[[]A, []B](f)
}

// MapRef applies a function to a pointer to each element of an array, returning a new array with the results.
// This is the curried version that returns a function.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// # Parameters
//
//   - f: Function that takes a pointer to an element and returns a transformed value
//
// # Returns
//
//   - A function that transforms an array of A into an array of B
//
// # Example
//
//	type Point struct { X, Y int }
//	extractX := array.MapRef(func(p *Point) int { return p.X })
//	points := []Point{{1, 2}, {3, 4}}
//	xs := extractX(points)
//	// xs: []int{1, 3}
func MapRef[A, B any](f func(*A) B) Operator[A, B] {
	return F.Bind2nd(MonadMapRef[A, B], f)
}

func filterRef[A any](fa []A, pred func(*A) bool) []A {
	count := len(fa)
	var result []A = make([]A, 0, count)
	for i := range count {
		a := &fa[i]
		if pred(a) {
			result = append(result, *a)
		}
	}
	return result
}

func filterMapRef[A, B any](fa []A, pred func(*A) bool, f func(*A) B) []B {
	count := len(fa)
	var result []B = make([]B, 0, count)
	for i := range count {
		a := &fa[i]
		if pred(a) {
			result = append(result, f(a))
		}
	}
	return result
}

// Filter returns a new array with all elements from the original array that match a predicate.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - pred: Predicate function to test each element
//
// # Returns
//
//   - A function that filters an array based on the predicate
//
// # Example
//
//	isEven := array.Filter(func(x int) bool { return x%2 == 0 })
//	result := isEven([]int{1, 2, 3, 4, 5, 6})
//	// result: []int{2, 4, 6}
//
//go:inline
func Filter[A any](pred func(A) bool) Operator[A, A] {
	return G.Filter[[]A](pred)
}

// FilterWithIndex returns a new array with all elements from the original array that match a predicate.
// The predicate receives both the index and the element.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - pred: Predicate function that takes an index and element
//
// # Returns
//
//   - A function that filters an array based on the predicate
//
// # Example
//
//	filterOddIndices := array.FilterWithIndex(func(i int, _ int) bool { return i%2 == 1 })
//	result := filterOddIndices([]int{10, 20, 30, 40, 50})
//	// result: []int{20, 40}
//
//go:inline
func FilterWithIndex[A any](pred func(int, A) bool) Operator[A, A] {
	return G.FilterWithIndex[[]A](pred)
}

// FilterRef returns a new array with all elements from the original array that match a predicate operating on pointers.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - pred: Predicate function that takes a pointer to an element
//
// # Returns
//
//   - A function that filters an array based on the predicate
//
// # Example
//
//	type Point struct { X, Y int }
//	filterPositiveX := array.FilterRef(func(p *Point) bool { return p.X > 0 })
//	points := []Point{{-1, 2}, {3, 4}, {-5, 6}}
//	result := filterPositiveX(points)
//	// result: []Point{{3, 4}}
func FilterRef[A any](pred func(*A) bool) Operator[A, A] {
	return F.Bind2nd(filterRef[A], pred)
}

// MonadFilterMap maps an array with a function that returns an Option and keeps only the Some values.
// This is the monadic version that takes the array as the first parameter.
//
//go:inline
func MonadFilterMap[A, B any](fa []A, f option.Kleisli[A, B]) []B {
	return G.MonadFilterMap[[]A, []B](fa, f)
}

// MonadFilterMapWithIndex maps an array with a function that takes an index and returns an Option,
// keeping only the Some values. This is the monadic version that takes the array as the first parameter.
//
//go:inline
func MonadFilterMapWithIndex[A, B any](fa []A, f func(int, A) Option[B]) []B {
	return G.MonadFilterMapWithIndex[[]A, []B](fa, f)
}

// FilterMap maps an array with an iterating function that returns an Option and keeps only the Some values discarding the Nones.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// # Parameters
//
//   - f: Function that maps elements to Option values
//
// # Returns
//
//   - A function that transforms and filters an array
//
// # Example
//
//	parseInt := array.FilterMap(func(s string) option.Option[int] {
//	    if n, err := strconv.Atoi(s); err == nil {
//	        return option.Some(n)
//	    }
//	    return option.None[int]()
//	})
//	result := parseInt([]string{"1", "bad", "3", "4"})
//	// result: []int{1, 3, 4}
//
//go:inline
func FilterMap[A, B any](f option.Kleisli[A, B]) Operator[A, B] {
	return G.FilterMap[[]A, []B](f)
}

// FilterMapWithIndex maps an array with an iterating function that returns an Option and keeps only the Some values discarding the Nones.
// The function receives both the index and the element.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// # Parameters
//
//   - f: Function that takes an index and element and returns an Option
//
// # Returns
//
//   - A function that transforms and filters an array
//
//go:inline
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B] {
	return G.FilterMapWithIndex[[]A, []B](f)
}

// ChainOptionK maps an array with an iterating function that returns an Option of an array.
// It keeps only the Some values discarding the Nones and then flattens the result.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// # Parameters
//
//   - f: Function that maps elements to Option of arrays
//
// # Returns
//
//   - A function that transforms, filters, and flattens an array
//
//go:inline
func ChainOptionK[A, B any](f option.Kleisli[A, []B]) Operator[A, B] {
	return G.ChainOptionK[[]A](f)
}

// FilterMapRef filters an array using a predicate on pointers and maps the matching elements using a function on pointers.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// # Parameters
//
//   - pred: Predicate function that takes a pointer to an element
//   - f: Function that transforms a pointer to an element
//
// # Returns
//
//   - A function that filters and transforms an array
func FilterMapRef[A, B any](pred func(a *A) bool, f func(*A) B) Operator[A, B] {
	return func(fa []A) []B {
		return filterMapRef(fa, pred, f)
	}
}

func reduceRef[A, B any](fa []A, f func(B, *A) B, initial B) B {
	current := initial
	for i := range len(fa) {
		current = f(current, &fa[i])
	}
	return current
}

// MonadReduce folds an array from left to right, applying a function to accumulate a result.
// This is the monadic version that takes the array as the first parameter.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//   - B: The type of the accumulated result
//
// # Parameters
//
//   - fa: The input array
//   - f: Function that combines the accumulator with each element
//   - initial: The initial accumulator value
//
// # Returns
//
//   - The final accumulated result
//
//go:inline
func MonadReduce[A, B any](fa []A, f func(B, A) B, initial B) B {
	return G.MonadReduce(fa, f, initial)
}

// MonadReduceWithIndex folds an array from left to right with access to the index,
// applying a function to accumulate a result.
// This is the monadic version that takes the array as the first parameter.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//   - B: The type of the accumulated result
//
// # Parameters
//
//   - fa: The input array
//   - f: Function that combines the index, accumulator, and element
//   - initial: The initial accumulator value
//
// # Returns
//
//   - The final accumulated result
//
//go:inline
func MonadReduceWithIndex[A, B any](fa []A, f func(int, B, A) B, initial B) B {
	return G.MonadReduceWithIndex(fa, f, initial)
}

// Reduce folds an array from left to right, applying a function to accumulate a result.
//
// Example:
//
//	sum := array.Reduce(func(acc, x int) int { return acc + x }, 0)
//	result := sum([]int{1, 2, 3, 4, 5}) // 15
//
//go:inline
func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B {
	return G.Reduce[[]A](f, initial)
}

// ReduceWithIndex folds an array from left to right with access to the index,
// applying a function to accumulate a result.
//
//go:inline
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B {
	return G.ReduceWithIndex[[]A](f, initial)
}

// ReduceRight folds an array from right to left, applying a function to accumulate a result.
//
//go:inline
func ReduceRight[A, B any](f func(A, B) B, initial B) func([]A) B {
	return G.ReduceRight[[]A](f, initial)
}

// ReduceRightWithIndex folds an array from right to left with access to the index,
// applying a function to accumulate a result.
//
//go:inline
func ReduceRightWithIndex[A, B any](f func(int, A, B) B, initial B) func([]A) B {
	return G.ReduceRightWithIndex[[]A](f, initial)
}

// ReduceRef folds an array from left to right using pointers to elements,
// applying a function to accumulate a result.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//   - B: The type of the accumulated result
//
// # Parameters
//
//   - f: Function that combines the accumulator with a pointer to each element
//   - initial: The initial accumulator value
//
// # Returns
//
//   - A function that reduces an array to a single value
func ReduceRef[A, B any](f func(B, *A) B, initial B) func([]A) B {
	return func(as []A) B {
		return reduceRef(as, f, initial)
	}
}

// Append adds an element to the end of an array, returning a new array.
// This is a non-curried version that takes both the array and element as parameters.
//
// Example:
//
//	arr := []int{1, 2, 3}
//	result := array.Append(arr, 4)
//	// result: []int{1, 2, 3, 4}
//	// arr: []int{1, 2, 3} (unchanged)
//
// For a curried version, see Push.
//
//go:inline
func Append[A any](as []A, a A) []A {
	return G.Append(as, a)
}

// IsEmpty checks if an array has no elements.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - as: The array to check
//
// # Returns
//
//   - true if the array is empty, false otherwise
//
//go:inline
func IsEmpty[A any](as []A) bool {
	return G.IsEmpty(as)
}

// IsNonEmpty checks if an array has at least one element.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - as: The array to check
//
// # Returns
//
//   - true if the array has at least one element, false otherwise
func IsNonEmpty[A any](as []A) bool {
	return len(as) > 0
}

// Empty returns an empty array of type A.
//
//go:inline
func Empty[A any]() []A {
	return G.Empty[[]A]()
}

// Zero returns an empty array of type A (alias for Empty).
//
//go:inline
func Zero[A any]() []A {
	return Empty[A]()
}

// Of constructs a single element array
//
//go:inline
func Of[A any](a A) []A {
	return G.Of[[]A](a)
}

// MonadChain applies a function that returns an array to each element and flattens the results.
// This is the monadic version that takes the array as the first parameter (also known as FlatMap).
//
//go:inline
func MonadChain[A, B any](fa []A, f Kleisli[A, B]) []B {
	return G.MonadChain(fa, f)
}

// Chain applies a function that returns an array to each element and flattens the results.
// This is the curried version (also known as FlatMap).
//
// Example:
//
//	duplicate := array.Chain(func(x int) []int { return []int{x, x} })
//	result := duplicate([]int{1, 2, 3}) // [1, 1, 2, 2, 3, 3]
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return G.Chain[[]A](f)
}

// MonadAp applies an array of functions to an array of values, producing all combinations.
// This is the monadic version that takes both arrays as parameters.
//
//go:inline
func MonadAp[B, A any](fab []func(A) B, fa []A) []B {
	return G.MonadAp[[]B](fab, fa)
}

// Ap applies an array of functions to an array of values, producing all combinations.
// This is the curried version.
//
//go:inline
func Ap[B, A any](fa []A) Operator[func(A) B, B] {
	return G.Ap[[]B, []func(A) B](fa)
}

// Match performs pattern matching on an array, calling onEmpty if empty or onNonEmpty if not.
//
//go:inline
func Match[A, B any](onEmpty func() B, onNonEmpty func([]A) B) func([]A) B {
	return G.Match(onEmpty, onNonEmpty)
}

// MatchLeft performs pattern matching on an array, calling onEmpty if empty or onNonEmpty with head and tail if not.
//
//go:inline
func MatchLeft[A, B any](onEmpty func() B, onNonEmpty func(A, []A) B) func([]A) B {
	return G.MatchLeft(onEmpty, onNonEmpty)
}

// Tail returns all elements except the first, wrapped in an Option.
// Returns None if the array is empty.
//
//go:inline
func Tail[A any](as []A) Option[[]A] {
	return G.Tail(as)
}

// Head returns the first element of an array, wrapped in an Option.
// Returns None if the array is empty.
//
//go:inline
func Head[A any](as []A) Option[A] {
	return G.Head(as)
}

// First returns the first element of an array, wrapped in an Option (alias for Head).
// Returns None if the array is empty.
//
//go:inline
func First[A any](as []A) Option[A] {
	return G.First(as)
}

// Last returns the last element of an array, wrapped in an Option.
// Returns None if the array is empty.
//
//go:inline
func Last[A any](as []A) Option[A] {
	return G.Last(as)
}

// PrependAll inserts a separator before each element of an array.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - middle: The separator to insert before each element
//
// # Returns
//
//   - A function that transforms an array by prepending the separator to each element
//
// # Example
//
//	result := array.PrependAll(0)([]int{1, 2, 3})
//	// result: []int{0, 1, 0, 2, 0, 3}
func PrependAll[A any](middle A) Operator[A, A] {
	return func(as []A) []A {
		count := len(as)
		dst := count * 2
		result := make([]A, dst)
		for i := count - 1; i >= 0; i-- {
			dst--
			result[dst] = as[i]
			dst--
			result[dst] = middle
		}
		return result
	}
}

// Intersperse inserts a separator between each element of an array.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - middle: The separator to insert between elements
//
// # Returns
//
//   - A function that transforms an array by inserting the separator between elements
//
// # Example
//
//	result := array.Intersperse(0)([]int{1, 2, 3})
//	// result: []int{1, 0, 2, 0, 3}
func Intersperse[A any](middle A) Operator[A, A] {
	prepend := PrependAll(middle)
	return func(as []A) []A {
		if IsEmpty(as) {
			return as
		}
		return prepend(as)[1:]
	}
}

// Intercalate inserts a separator between elements and concatenates them using a Monoid.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - m: The Monoid to use for concatenation
//
// # Returns
//
//   - A curried function that takes a separator and returns a function that reduces an array
func Intercalate[A any](m M.Monoid[A]) func(A) func([]A) A {
	return func(middle A) func([]A) A {
		return Match(m.Empty, F.Flow2(Intersperse(middle), ConcatAll(m)))
	}
}

// Flatten converts a nested array into a flat array by concatenating all inner arrays.
//
// # Type Parameters
//
//   - A: The type of elements in the inner arrays
//
// # Parameters
//
//   - mma: A nested array (array of arrays)
//
// # Returns
//
//   - A flat array containing all elements from all inner arrays
//
// # Example
//
//	result := array.Flatten([][]int{{1, 2}, {3, 4}, {5}})
//	// result: []int{1, 2, 3, 4, 5}
//
//go:inline
func Flatten[A any](mma [][]A) []A {
	return G.Flatten(mma)
}

// Slice extracts a subarray from index low (inclusive) to high (exclusive).
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - low: The starting index (inclusive)
//   - high: The ending index (exclusive)
//
// # Returns
//
//   - A function that extracts a subarray
//
// # Example
//
//	middle := array.Slice[int](2, 5)
//	result := middle([]int{0, 1, 2, 3, 4, 5, 6})
//	// result: []int{2, 3, 4}
func Slice[A any](low, high int) Operator[A, A] {
	return array.Slice[[]A](low, high)
}

// Lookup returns the element at the specified index, wrapped in an Option.
// Returns None if the index is out of bounds.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - idx: The index to look up
//
// # Returns
//
//   - A function that retrieves an element at the given index, wrapped in an Option
//
// # Example
//
//	getSecond := array.Lookup[int](1)
//	result := getSecond([]int{10, 20, 30})
//	// result: option.Some(20)
//
//go:inline
func Lookup[A any](idx int) func([]A) Option[A] {
	return G.Lookup[[]A](idx)
}

// UpsertAt returns a function that inserts or updates an element at a specific index.
// If the index is out of bounds, the element is appended.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - a: The element to insert or update
//
// # Returns
//
//   - A function that takes an index and returns a function that upserts at that index
//
//go:inline
func UpsertAt[A any](a A) Operator[A, A] {
	return G.UpsertAt[[]A](a)
}

// Size returns the number of elements in an array.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - as: The array to measure
//
// # Returns
//
//   - The number of elements in the array
//
//go:inline
func Size[A any](as []A) int {
	return G.Size(as)
}

// MonadPartition splits an array into two arrays based on a predicate.
// The first array contains elements for which the predicate returns false,
// the second contains elements for which it returns true.
//
//go:inline
func MonadPartition[A any](as []A, pred func(A) bool) pair.Pair[[]A, []A] {
	return G.MonadPartition(as, pred)
}

// Partition creates two new arrays out of one. The left result contains the elements
// for which the predicate returns false, the right one those for which the predicate returns true.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - pred: Predicate function to test each element
//
// # Returns
//
//   - A function that partitions an array into a pair of arrays
//
// # Example
//
//	isEven := array.Partition(func(x int) bool { return x%2 == 0 })
//	result := isEven([]int{1, 2, 3, 4, 5, 6})
//	// result: pair.Pair{Left: []int{1, 3, 5}, Right: []int{2, 4, 6}}
//
//go:inline
func Partition[A any](pred func(A) bool) func([]A) pair.Pair[[]A, []A] {
	return G.Partition[[]A](pred)
}

// IsNil checks if the array is set to nil.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - as: The array to check
//
// # Returns
//
//   - true if the array is nil, false otherwise
func IsNil[A any](as []A) bool {
	return array.IsNil(as)
}

// IsNonNil checks if the array is not nil.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - as: The array to check
//
// # Returns
//
//   - true if the array is not nil, false otherwise
func IsNonNil[A any](as []A) bool {
	return array.IsNonNil(as)
}

// ConstNil returns a nil array.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Returns
//
//   - A nil array of type A
func ConstNil[A any]() []A {
	return array.ConstNil[[]A]()
}

// SliceRight extracts a subarray from the specified start index to the end.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - start: The starting index (inclusive)
//
// # Returns
//
//   - A function that extracts a subarray from start to end
//
// # Example
//
//	fromThird := array.SliceRight[int](2)
//	result := fromThird([]int{0, 1, 2, 3, 4, 5})
//	// result: []int{2, 3, 4, 5}
//
//go:inline
func SliceRight[A any](start int) Operator[A, A] {
	return G.SliceRight[[]A](start)
}

// Copy creates a shallow copy of the array.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - b: The array to copy
//
// # Returns
//
//   - A new array with the same elements
//
//go:inline
func Copy[A any](b []A) []A {
	return G.Copy(b)
}

// Clone creates a deep copy of the array using the provided endomorphism to clone the values.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - f: Function to clone each element
//
// # Returns
//
//   - A function that creates a deep copy of an array
//
//go:inline
func Clone[A any](f func(A) A) Operator[A, A] {
	return G.Clone[[]A](f)
}

// FoldMap maps and folds an array. Maps each value using the iterating function,
// then folds the results using the provided Monoid.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements after mapping
//
// # Parameters
//
//   - m: The Monoid to use for folding
//
// # Returns
//
//   - A curried function that takes a mapping function and returns a function that folds an array
//
//go:inline
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func([]A) B {
	return G.FoldMap[[]A](m)
}

// FoldMapWithIndex maps and folds an array with access to indices. Maps each value using the iterating function,
// then folds the results using the provided Monoid.
//
// # Type Parameters
//
//   - A: The type of elements in the input array
//   - B: The type of elements after mapping
//
// # Parameters
//
//   - m: The Monoid to use for folding
//
// # Returns
//
//   - A curried function that takes a mapping function and returns a function that folds an array
//
//go:inline
func FoldMapWithIndex[A, B any](m M.Monoid[B]) func(func(int, A) B) func([]A) B {
	return G.FoldMapWithIndex[[]A](m)
}

// Fold folds the array using the provided Monoid.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - m: The Monoid to use for folding
//
// # Returns
//
//   - A function that folds an array to a single value
//
//go:inline
func Fold[A any](m M.Monoid[A]) func([]A) A {
	return G.Fold[[]A](m)
}

// Push adds an element to the end of an array (curried version of Append).
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - a: The element to add
//
// # Returns
//
//   - A function that appends the element to an array
//
// # Example
//
//	addFive := array.Push(5)
//	result := addFive([]int{1, 2, 3})
//	// result: []int{1, 2, 3, 5}
//
// # See Also
//
//   - Append: Non-curried version
//
//go:inline
func Push[A any](a A) Operator[A, A] {
	return G.Push[Operator[A, A]](a)
}

// Concat concatenates two arrays by appending a suffix array to a base array.
//
// This is a curried function that takes a suffix array and returns a function
// that takes a base array and produces a new array with the suffix appended.
// It follows the "data last" pattern, where the data to be operated on (base array)
// is provided last, making it ideal for use in functional pipelines.
//
// Semantic: Concat(suffix)(base) produces [base... suffix...]
//
// The function creates a new array containing all elements from the base array
// followed by all elements from the suffix array. Neither input array is modified.
//
// Type Parameters:
//
//   - A: The type of elements in the arrays
//
// Parameters:
//
//   - suffix: The array to append to the end of the base array
//
// Returns:
//
//   - A function that takes a base array and returns [base... suffix...]
//
// Behavior:
//
//   - Creates a new array with length equal to len(base) + len(suffix)
//   - Copies all elements from the base array first
//   - Appends all elements from the suffix array at the end
//   - Returns the base array unchanged if suffix is empty
//   - Returns suffix unchanged if the base array is empty
//   - Does not modify either input array
//   - Preserves element order within each array
//
// Example - Basic concatenation:
//
//	base := []int{1, 2, 3}
//	suffix := []int{4, 5, 6}
//	concat := array.Concat(suffix)
//	result := concat(base)
//	// result: []int{1, 2, 3, 4, 5, 6}
//	// base: []int{1, 2, 3} (unchanged)
//	// suffix: []int{4, 5, 6} (unchanged)
//
// Example - Direct application:
//
//	result := array.Concat([]int{4, 5, 6})([]int{1, 2, 3})
//	// result: []int{1, 2, 3, 4, 5, 6}
//	// Demonstrates: Concat(b)(a) = [a... b...]
//
// Example - Empty arrays:
//
//	base := []int{1, 2, 3}
//	empty := []int{}
//	result := array.Concat(empty)(base)
//	// result: []int{1, 2, 3}
//
// Example - Strings:
//
//	words1 := []string{"hello", "world"}
//	words2 := []string{"foo", "bar"}
//	result := array.Concat(words2)(words1)
//	// result: []string{"hello", "world", "foo", "bar"}
//
// Example - Functional composition:
//
//	numbers := []int{1, 2, 3}
//	result := F.Pipe2(
//	    numbers,
//	    array.Map(N.Mul(2)),           // [2, 4, 6]
//	    array.Concat([]int{10, 20}),   // [2, 4, 6, 10, 20]
//	)
//
// Example - Multiple concatenations:
//
//	result := F.Pipe2(
//	    []int{1},
//	    array.Concat([]int{2, 3}),     // [1, 2, 3]
//	    array.Concat([]int{4, 5}),     // [1, 2, 3, 4, 5]
//	)
//
// Example - Building arrays incrementally:
//
//	header := []string{"Name", "Age"}
//	data := []string{"Alice", "30"}
//	footer := []string{"Total: 1"}
//	result := F.Pipe2(
//	    header,
//	    array.Concat(data),
//	    array.Concat(footer),
//	)
//	// result: []string{"Name", "Age", "Alice", "30", "Total: 1"}
//
// Use cases:
//
//   - Combining multiple arrays into one
//   - Building arrays incrementally in pipelines
//   - Implementing array-based data structures (queues, buffers)
//   - Merging results from multiple operations
//   - Creating array transformation pipelines
//   - Appending batches of elements
//
// Mathematical properties:
//
//   - Associativity: Concat(c)(Concat(b)(a)) == Concat(Concat(c)(b))(a)
//   - Identity: Concat([])(a) == a and Concat(a)([]) == a
//   - Length: len(Concat(b)(a)) == len(a) + len(b)
//
// Performance:
//
//   - Time complexity: O(n + m) where n and m are the lengths of the arrays
//   - Space complexity: O(n + m) for the new array
//   - Optimized to avoid allocation when one array is empty
//
// Note: This function is immutable - it creates a new array rather than modifying
// the input arrays. For appending a single element, consider using Append or Push.
//
// See Also:
//
//   - Append: For appending a single element
//   - Push: Curried version of Append
//   - Flatten: For flattening nested arrays
//
//go:inline
func Concat[A any](suffix []A) Operator[A, A] {
	return F.Bind2nd(array.Concat[[]A, A], suffix)
}

// MonadFlap applies a value to an array of functions, producing an array of results.
// This is the monadic version that takes both parameters.
//
// # Type Parameters
//
//   - B: The type of results
//   - A: The type of the input value
//
// # Parameters
//
//   - fab: Array of functions to apply
//   - a: The value to apply to each function
//
// # Returns
//
//   - An array of results from applying the value to each function
//
//go:inline
func MonadFlap[B, A any](fab []func(A) B, a A) []B {
	return G.MonadFlap[func(A) B, []func(A) B, []B](fab, a)
}

// Flap applies a value to an array of functions, producing an array of results.
// This is the curried version.
//
// # Type Parameters
//
//   - B: The type of results
//   - A: The type of the input value
//
// # Parameters
//
//   - a: The value to apply to each function
//
// # Returns
//
//   - A function that applies the value to an array of functions
//
// # Example
//
//	fns := []func(int) int{
//	    func(x int) int { return x * 2 },
//	    func(x int) int { return x + 10 },
//	    func(x int) int { return x * x },
//	}
//	applyFive := array.Flap[int](5)
//	result := applyFive(fns)
//	// result: []int{10, 15, 25}
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return G.Flap[func(A) B, []func(A) B, []B](a)
}

// Prepend adds an element to the beginning of an array, returning a new array.
//
// # Type Parameters
//
//   - A: The type of elements in the array
//
// # Parameters
//
//   - head: The element to add at the beginning
//
// # Returns
//
//   - A function that prepends the element to an array
//
// # Example
//
//	addZero := array.Prepend(0)
//	result := addZero([]int{1, 2, 3})
//	// result: []int{0, 1, 2, 3}
//
//go:inline
func Prepend[A any](head A) Operator[A, A] {
	return G.Prepend[Operator[A, A]](head)
}

// Reverse returns a new slice with elements in reverse order.
// This function creates a new slice containing all elements from the input slice
// in reverse order, without modifying the original slice.
//
// Type Parameters:
//   - A: The type of elements in the slice
//
// Parameters:
//   - as: The input slice to reverse
//
// Returns:
//   - A new slice with elements in reverse order
//
// Behavior:
//   - Creates a new slice with the same length as the input
//   - Copies elements from the input slice in reverse order
//   - Does not modify the original slice
//   - Returns an empty slice if the input is empty
//   - Returns a single-element slice unchanged if input has one element
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	reversed := array.Reverse(numbers)
//	// reversed: []int{5, 4, 3, 2, 1}
//	// numbers: []int{1, 2, 3, 4, 5} (unchanged)
//
// Example with strings:
//
//	words := []string{"hello", "world", "foo", "bar"}
//	reversed := array.Reverse(words)
//	// reversed: []string{"bar", "foo", "world", "hello"}
//
// Example with empty slice:
//
//	empty := []int{}
//	reversed := array.Reverse(empty)
//	// reversed: []int{} (empty slice)
//
// Example with single element:
//
//	single := []string{"only"}
//	reversed := array.Reverse(single)
//	// reversed: []string{"only"}
//
// Use cases:
//   - Reversing the order of elements for display or processing
//   - Implementing stack-like behavior (LIFO)
//   - Processing data in reverse chronological order
//   - Reversing transformation pipelines
//   - Creating palindrome checks
//   - Implementing undo/redo functionality
//
// Example with processing in reverse:
//
//	events := []string{"start", "middle", "end"}
//	reversed := array.Reverse(events)
//	// Process events in reverse order
//	for _, event := range reversed {
//	    fmt.Println(event) // Prints: "end", "middle", "start"
//	}
//
// Example with functional composition:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	result := F.Pipe2(
//	    numbers,
//	    array.Map(N.Mul(2)),
//	    array.Reverse,
//	)
//	// result: []int{10, 8, 6, 4, 2}
//
// Performance:
//   - Time complexity: O(n) where n is the length of the slice
//   - Space complexity: O(n) for the new slice
//   - Does not allocate if the input slice is empty
//
// Note: This function is immutable - it does not modify the original slice.
// If you need to reverse a slice in-place, consider using a different approach
// or modifying the slice directly.
//
//go:inline
func Reverse[A any](as []A) []A {
	return G.Reverse(as)
}

// Extend applies a function to every suffix of an array, creating a new array of results.
// This is the comonad extend operation for arrays.
//
// The function f is applied to progressively smaller suffixes of the input array:
//   - f(as[0:]) for the first element
//   - f(as[1:]) for the second element
//   - f(as[2:]) for the third element
//   - and so on...
//
// Type Parameters:
//   - A: The type of elements in the input array
//   - B: The type of elements in the output array
//
// Parameters:
//   - f: A function that takes an array suffix and returns a value
//
// Returns:
//   - A function that transforms an array of A into an array of B
//
// Behavior:
//   - Creates a new array with the same length as the input
//   - For each position i, applies f to the suffix starting at i
//   - Returns an empty array if the input is empty
//
// Example:
//
//	// Sum all elements from current position to end
//	sumSuffix := array.Extend(func(as []int) int {
//	    return array.Reduce(func(acc, x int) int { return acc + x }, 0)(as)
//	})
//	result := sumSuffix([]int{1, 2, 3, 4})
//	// result: []int{10, 9, 7, 4}
//	// Explanation: [1+2+3+4, 2+3+4, 3+4, 4]
//
// Example with length:
//
//	// Get remaining length at each position
//	lengths := array.Extend(array.Size[int])
//	result := lengths([]int{10, 20, 30})
//	// result: []int{3, 2, 1}
//
// Example with head:
//
//	// Duplicate each element (extract head of each suffix)
//	duplicate := array.Extend(func(as []int) int {
//	    return F.Pipe1(as, array.Head[int], O.GetOrElse(F.Constant(0)))
//	})
//	result := duplicate([]int{1, 2, 3})
//	// result: []int{1, 2, 3}
//
// Use cases:
//   - Computing cumulative or rolling operations
//   - Implementing sliding window algorithms
//   - Creating context-aware transformations
//   - Building comonadic computations
//
// Comonad laws:
//   - Left identity: Extend(Extract) == Identity
//   - Right identity: Extract ∘ Extend(f) == f
//   - Associativity: Extend(f) ∘ Extend(g) == Extend(f ∘ Extend(g))
//
//go:inline
func Extend[A, B any](f func([]A) B) Operator[A, B] {
	return func(as []A) []B {
		return G.MakeBy[[]B](len(as), func(i int) B { return f(as[i:]) })
	}
}

// Extract returns the first element of an array, or a zero value if empty.
// This is the comonad extract operation for arrays.
//
// Extract is the dual of the monadic return/of operation. While Of wraps a value
// in a context, Extract unwraps a value from its context.
//
// Type Parameters:
//   - A: The type of elements in the array
//
// Parameters:
//   - as: The input array
//
// Returns:
//   - The first element if the array is non-empty, otherwise the zero value of type A
//
// Behavior:
//   - Returns as[0] if the array has at least one element
//   - Returns the zero value of A if the array is empty
//   - Does not modify the input array
//
// Example:
//
//	result := array.Extract([]int{1, 2, 3})
//	// result: 1
//
// Example with empty array:
//
//	result := array.Extract([]int{})
//	// result: 0 (zero value for int)
//
// Example with strings:
//
//	result := array.Extract([]string{"hello", "world"})
//	// result: "hello"
//
// Example with empty string array:
//
//	result := array.Extract([]string{})
//	// result: "" (zero value for string)
//
// Use cases:
//   - Extracting the current focus from a comonadic context
//   - Getting the head element with a default zero value
//   - Implementing comonad-based computations
//
// Comonad laws:
//   - Extract ∘ Of == Identity (extracting from a singleton returns the value)
//   - Extract ∘ Extend(f) == f (extract after extend equals applying f)
//
// Note: For a safer alternative that handles empty arrays explicitly,
// consider using Head which returns an Option[A].
//
//go:inline
func Extract[A any](as []A) A {
	return G.Extract(as)
}
