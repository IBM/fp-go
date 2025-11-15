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
	"github.com/IBM/fp-go/v2/tuple"
)

// From constructs an array from a set of variadic arguments
//
//go:inline
func From[A any](data ...A) []A {
	return G.From[[]A](data...)
}

// MakeBy returns a `Array` of length `n` with element `i` initialized with `f(i)`.
//
//go:inline
func MakeBy[F ~func(int) A, A any](n int, f F) []A {
	return G.MakeBy[[]A](n, f)
}

// Replicate creates a `Array` containing a value repeated the specified number of times.
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
//	double := array.Map(func(x int) int { return x * 2 })
//	result := double([]int{1, 2, 3}) // [2, 4, 6]
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return G.Map[[]A, []B](f)
}

// MapRef applies a function to a pointer to each element of an array, returning a new array with the results.
// This is the curried version that returns a function.
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

// Filter returns a new array with all elements from the original array that match a predicate
//
//go:inline
func Filter[A any](pred func(A) bool) Operator[A, A] {
	return G.Filter[[]A](pred)
}

// FilterWithIndex returns a new array with all elements from the original array that match a predicate
//
//go:inline
func FilterWithIndex[A any](pred func(int, A) bool) Operator[A, A] {
	return G.FilterWithIndex[[]A](pred)
}

// FilterRef returns a new array with all elements from the original array that match a predicate operating on pointers.
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

// FilterMap maps an array with an iterating function that returns an [Option] and it keeps only the Some values discarding the Nones.
//
//go:inline
func FilterMap[A, B any](f option.Kleisli[A, B]) Operator[A, B] {
	return G.FilterMap[[]A, []B](f)
}

// FilterMapWithIndex maps an array with an iterating function that returns an [Option] and it keeps only the Some values discarding the Nones.
//
//go:inline
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B] {
	return G.FilterMapWithIndex[[]A, []B](f)
}

// FilterChain maps an array with an iterating function that returns an [Option] of an array. It keeps only the Some values discarding the Nones and then flattens the result.
//
//go:inline
func FilterChain[A, B any](f option.Kleisli[A, []B]) Operator[A, B] {
	return G.FilterChain[[]A](f)
}

// FilterMapRef filters an array using a predicate on pointers and maps the matching elements using a function on pointers.
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

//go:inline
func MonadReduce[A, B any](fa []A, f func(B, A) B, initial B) B {
	return G.MonadReduce(fa, f, initial)
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
func ReduceRef[A, B any](f func(B, *A) B, initial B) func([]A) B {
	return func(as []A) B {
		return reduceRef(as, f, initial)
	}
}

// Append adds an element to the end of an array, returning a new array.
//
//go:inline
func Append[A any](as []A, a A) []A {
	return G.Append(as, a)
}

// IsEmpty checks if an array has no elements.
//
//go:inline
func IsEmpty[A any](as []A) bool {
	return G.IsEmpty(as)
}

// IsNonEmpty checks if an array has at least one element.
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
// Example:
//
//	result := array.Intersperse(0)([]int{1, 2, 3}) // [1, 0, 2, 0, 3]
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
func Intercalate[A any](m M.Monoid[A]) func(A) func([]A) A {
	return func(middle A) func([]A) A {
		return Match(m.Empty, F.Flow2(Intersperse(middle), ConcatAll(m)))
	}
}

// Flatten converts a nested array into a flat array by concatenating all inner arrays.
//
// Example:
//
//	result := array.Flatten([][]int{{1, 2}, {3, 4}, {5}}) // [1, 2, 3, 4, 5]
//
//go:inline
func Flatten[A any](mma [][]A) []A {
	return G.Flatten(mma)
}

// Slice extracts a subarray from index low (inclusive) to high (exclusive).
func Slice[A any](low, high int) Operator[A, A] {
	return array.Slice[[]A](low, high)
}

// Lookup returns the element at the specified index, wrapped in an Option.
// Returns None if the index is out of bounds.
//
//go:inline
func Lookup[A any](idx int) func([]A) Option[A] {
	return G.Lookup[[]A](idx)
}

// UpsertAt returns a function that inserts or updates an element at a specific index.
// If the index is out of bounds, the element is appended.
//
//go:inline
func UpsertAt[A any](a A) Operator[A, A] {
	return G.UpsertAt[[]A](a)
}

// Size returns the number of elements in an array.
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
func MonadPartition[A any](as []A, pred func(A) bool) tuple.Tuple2[[]A, []A] {
	return G.MonadPartition(as, pred)
}

// Partition creates two new arrays out of one, the left result contains the elements
// for which the predicate returns false, the right one those for which the predicate returns true
//
//go:inline
func Partition[A any](pred func(A) bool) func([]A) tuple.Tuple2[[]A, []A] {
	return G.Partition[[]A](pred)
}

// IsNil checks if the array is set to nil
func IsNil[A any](as []A) bool {
	return array.IsNil(as)
}

// IsNonNil checks if the array is set to nil
func IsNonNil[A any](as []A) bool {
	return array.IsNonNil(as)
}

// ConstNil returns a nil array
func ConstNil[A any]() []A {
	return array.ConstNil[[]A]()
}

// SliceRight extracts a subarray from the specified start index to the end.
//
//go:inline
func SliceRight[A any](start int) Operator[A, A] {
	return G.SliceRight[[]A](start)
}

// Copy creates a shallow copy of the array
//
//go:inline
func Copy[A any](b []A) []A {
	return G.Copy(b)
}

// Clone creates a deep copy of the array using the provided endomorphism to clone the values
//
//go:inline
func Clone[A any](f func(A) A) Operator[A, A] {
	return G.Clone[[]A](f)
}

// FoldMap maps and folds an array. Map the Array passing each value to the iterating function. Then fold the results using the provided Monoid.
//
//go:inline
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func([]A) B {
	return G.FoldMap[[]A](m)
}

// FoldMapWithIndex maps and folds an array. Map the Array passing each value to the iterating function. Then fold the results using the provided Monoid.
//
//go:inline
func FoldMapWithIndex[A, B any](m M.Monoid[B]) func(func(int, A) B) func([]A) B {
	return G.FoldMapWithIndex[[]A](m)
}

// Fold folds the array using the provided Monoid.
//
//go:inline
func Fold[A any](m M.Monoid[A]) func([]A) A {
	return G.Fold[[]A](m)
}

// Push adds an element to the end of an array (alias for Append).
//
//go:inline
func Push[A any](a A) Operator[A, A] {
	return G.Push[Operator[A, A]](a)
}

// MonadFlap applies a value to an array of functions, producing an array of results.
// This is the monadic version that takes both parameters.
//
//go:inline
func MonadFlap[B, A any](fab []func(A) B, a A) []B {
	return G.MonadFlap[func(A) B, []func(A) B, []B](fab, a)
}

// Flap applies a value to an array of functions, producing an array of results.
// This is the curried version.
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return G.Flap[func(A) B, []func(A) B, []B](a)
}

// Prepend adds an element to the beginning of an array, returning a new array.
//
//go:inline
func Prepend[A any](head A) Operator[A, A] {
	return G.Prepend[Operator[A, A]](head)
}
