// Copyright (c) 2023 IBM Corp.
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
	EM "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/tuple"
)

// From constructs an array from a set of variadic arguments
func From[A any](data ...A) []A {
	return G.From[[]A](data...)
}

// MakeBy returns a `Array` of length `n` with element `i` initialized with `f(i)`.
func MakeBy[F ~func(int) A, A any](n int, f F) []A {
	return G.MakeBy[[]A](n, f)
}

// Replicate creates a `Array` containing a value repeated the specified number of times.
func Replicate[A any](n int, a A) []A {
	return G.Replicate[[]A](n, a)
}

// MonadMap applies a function to each element of an array, returning a new array with the results.
// This is the monadic version of Map that takes the array as the first parameter.
func MonadMap[A, B any](as []A, f func(a A) B) []B {
	return G.MonadMap[[]A, []B](as, f)
}

// MonadMapRef applies a function to a pointer to each element of an array, returning a new array with the results.
// This is useful when you need to access elements by reference without copying.
func MonadMapRef[A, B any](as []A, f func(a *A) B) []B {
	count := len(as)
	bs := make([]B, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(&as[i])
	}
	return bs
}

// MapWithIndex applies a function to each element and its index in an array, returning a new array with the results.
func MapWithIndex[A, B any](f func(int, A) B) func([]A) []B {
	return G.MapWithIndex[[]A, []B](f)
}

// Map applies a function to each element of an array, returning a new array with the results.
// This is the curried version that returns a function.
//
// Example:
//
//	double := array.Map(func(x int) int { return x * 2 })
//	result := double([]int{1, 2, 3}) // [2, 4, 6]
func Map[A, B any](f func(a A) B) func([]A) []B {
	return G.Map[[]A, []B, A, B](f)
}

// MapRef applies a function to a pointer to each element of an array, returning a new array with the results.
// This is the curried version that returns a function.
func MapRef[A, B any](f func(a *A) B) func([]A) []B {
	return F.Bind2nd(MonadMapRef[A, B], f)
}

func filterRef[A any](fa []A, pred func(a *A) bool) []A {
	var result []A
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(&a) {
			result = append(result, a)
		}
	}
	return result
}

func filterMapRef[A, B any](fa []A, pred func(a *A) bool, f func(a *A) B) []B {
	var result []B
	count := len(fa)
	for i := 0; i < count; i++ {
		a := fa[i]
		if pred(&a) {
			result = append(result, f(&a))
		}
	}
	return result
}

// Filter returns a new array with all elements from the original array that match a predicate
func Filter[A any](pred func(A) bool) EM.Endomorphism[[]A] {
	return G.Filter[[]A](pred)
}

// FilterWithIndex returns a new array with all elements from the original array that match a predicate
func FilterWithIndex[A any](pred func(int, A) bool) EM.Endomorphism[[]A] {
	return G.FilterWithIndex[[]A](pred)
}

// FilterRef returns a new array with all elements from the original array that match a predicate operating on pointers.
func FilterRef[A any](pred func(*A) bool) EM.Endomorphism[[]A] {
	return F.Bind2nd(filterRef[A], pred)
}

// MonadFilterMap maps an array with a function that returns an Option and keeps only the Some values.
// This is the monadic version that takes the array as the first parameter.
func MonadFilterMap[A, B any](fa []A, f func(A) O.Option[B]) []B {
	return G.MonadFilterMap[[]A, []B](fa, f)
}

// MonadFilterMapWithIndex maps an array with a function that takes an index and returns an Option,
// keeping only the Some values. This is the monadic version that takes the array as the first parameter.
func MonadFilterMapWithIndex[A, B any](fa []A, f func(int, A) O.Option[B]) []B {
	return G.MonadFilterMapWithIndex[[]A, []B](fa, f)
}

// FilterMap maps an array with an iterating function that returns an [O.Option] and it keeps only the Some values discarding the Nones.
func FilterMap[A, B any](f func(A) O.Option[B]) func([]A) []B {
	return G.FilterMap[[]A, []B](f)
}

// FilterMapWithIndex maps an array with an iterating function that returns an [O.Option] and it keeps only the Some values discarding the Nones.
func FilterMapWithIndex[A, B any](f func(int, A) O.Option[B]) func([]A) []B {
	return G.FilterMapWithIndex[[]A, []B](f)
}

// FilterChain maps an array with an iterating function that returns an [O.Option] of an array. It keeps only the Some values discarding the Nones and then flattens the result.
func FilterChain[A, B any](f func(A) O.Option[[]B]) func([]A) []B {
	return G.FilterChain[[]A](f)
}

// FilterMapRef filters an array using a predicate on pointers and maps the matching elements using a function on pointers.
func FilterMapRef[A, B any](pred func(a *A) bool, f func(a *A) B) func([]A) []B {
	return func(fa []A) []B {
		return filterMapRef(fa, pred, f)
	}
}

func reduceRef[A, B any](fa []A, f func(B, *A) B, initial B) B {
	current := initial
	count := len(fa)
	for i := 0; i < count; i++ {
		current = f(current, &fa[i])
	}
	return current
}

// Reduce folds an array from left to right, applying a function to accumulate a result.
//
// Example:
//
//	sum := array.Reduce(func(acc, x int) int { return acc + x }, 0)
//	result := sum([]int{1, 2, 3, 4, 5}) // 15
func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B {
	return G.Reduce[[]A](f, initial)
}

// ReduceWithIndex folds an array from left to right with access to the index,
// applying a function to accumulate a result.
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B {
	return G.ReduceWithIndex[[]A](f, initial)
}

// ReduceRight folds an array from right to left, applying a function to accumulate a result.
func ReduceRight[A, B any](f func(A, B) B, initial B) func([]A) B {
	return G.ReduceRight[[]A](f, initial)
}

// ReduceRightWithIndex folds an array from right to left with access to the index,
// applying a function to accumulate a result.
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
func Append[A any](as []A, a A) []A {
	return G.Append(as, a)
}

// IsEmpty checks if an array has no elements.
func IsEmpty[A any](as []A) bool {
	return G.IsEmpty(as)
}

// IsNonEmpty checks if an array has at least one element.
func IsNonEmpty[A any](as []A) bool {
	return len(as) > 0
}

// Empty returns an empty array of type A.
func Empty[A any]() []A {
	return G.Empty[[]A]()
}

// Zero returns an empty array of type A (alias for Empty).
func Zero[A any]() []A {
	return Empty[A]()
}

// Of constructs a single element array
func Of[A any](a A) []A {
	return G.Of[[]A](a)
}

// MonadChain applies a function that returns an array to each element and flattens the results.
// This is the monadic version that takes the array as the first parameter (also known as FlatMap).
func MonadChain[A, B any](fa []A, f func(a A) []B) []B {
	return G.MonadChain[[]A, []B](fa, f)
}

// Chain applies a function that returns an array to each element and flattens the results.
// This is the curried version (also known as FlatMap).
//
// Example:
//
//	duplicate := array.Chain(func(x int) []int { return []int{x, x} })
//	result := duplicate([]int{1, 2, 3}) // [1, 1, 2, 2, 3, 3]
func Chain[A, B any](f func(A) []B) func([]A) []B {
	return G.Chain[[]A, []B](f)
}

// MonadAp applies an array of functions to an array of values, producing all combinations.
// This is the monadic version that takes both arrays as parameters.
func MonadAp[B, A any](fab []func(A) B, fa []A) []B {
	return G.MonadAp[[]B](fab, fa)
}

// Ap applies an array of functions to an array of values, producing all combinations.
// This is the curried version.
func Ap[B, A any](fa []A) func([]func(A) B) []B {
	return G.Ap[[]B, []func(A) B](fa)
}

// Match performs pattern matching on an array, calling onEmpty if empty or onNonEmpty if not.
func Match[A, B any](onEmpty func() B, onNonEmpty func([]A) B) func([]A) B {
	return G.Match[[]A](onEmpty, onNonEmpty)
}

// MatchLeft performs pattern matching on an array, calling onEmpty if empty or onNonEmpty with head and tail if not.
func MatchLeft[A, B any](onEmpty func() B, onNonEmpty func(A, []A) B) func([]A) B {
	return G.MatchLeft[[]A](onEmpty, onNonEmpty)
}

// Tail returns all elements except the first, wrapped in an Option.
// Returns None if the array is empty.
func Tail[A any](as []A) O.Option[[]A] {
	return G.Tail(as)
}

// Head returns the first element of an array, wrapped in an Option.
// Returns None if the array is empty.
func Head[A any](as []A) O.Option[A] {
	return G.Head(as)
}

// First returns the first element of an array, wrapped in an Option (alias for Head).
// Returns None if the array is empty.
func First[A any](as []A) O.Option[A] {
	return G.First(as)
}

// Last returns the last element of an array, wrapped in an Option.
// Returns None if the array is empty.
func Last[A any](as []A) O.Option[A] {
	return G.Last(as)
}

// PrependAll inserts a separator before each element of an array.
func PrependAll[A any](middle A) EM.Endomorphism[[]A] {
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
func Intersperse[A any](middle A) EM.Endomorphism[[]A] {
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
		return Match(m.Empty, F.Flow2(Intersperse(middle), ConcatAll[A](m)))
	}
}

// Flatten converts a nested array into a flat array by concatenating all inner arrays.
//
// Example:
//
//	result := array.Flatten([][]int{{1, 2}, {3, 4}, {5}}) // [1, 2, 3, 4, 5]
func Flatten[A any](mma [][]A) []A {
	return G.Flatten(mma)
}

// Slice extracts a subarray from index low (inclusive) to high (exclusive).
func Slice[A any](low, high int) func(as []A) []A {
	return array.Slice[[]A](low, high)
}

// Lookup returns the element at the specified index, wrapped in an Option.
// Returns None if the index is out of bounds.
func Lookup[A any](idx int) func([]A) O.Option[A] {
	return G.Lookup[[]A](idx)
}

// UpsertAt returns a function that inserts or updates an element at a specific index.
// If the index is out of bounds, the element is appended.
func UpsertAt[A any](a A) EM.Endomorphism[[]A] {
	return G.UpsertAt[[]A](a)
}

// Size returns the number of elements in an array.
func Size[A any](as []A) int {
	return G.Size(as)
}

// MonadPartition splits an array into two arrays based on a predicate.
// The first array contains elements for which the predicate returns false,
// the second contains elements for which it returns true.
func MonadPartition[A any](as []A, pred func(A) bool) tuple.Tuple2[[]A, []A] {
	return G.MonadPartition(as, pred)
}

// Partition creates two new arrays out of one, the left result contains the elements
// for which the predicate returns false, the right one those for which the predicate returns true
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
func SliceRight[A any](start int) EM.Endomorphism[[]A] {
	return G.SliceRight[[]A](start)
}

// Copy creates a shallow copy of the array
func Copy[A any](b []A) []A {
	return G.Copy(b)
}

// Clone creates a deep copy of the array using the provided endomorphism to clone the values
func Clone[A any](f func(A) A) func(as []A) []A {
	return G.Clone[[]A](f)
}

// FoldMap maps and folds an array. Map the Array passing each value to the iterating function. Then fold the results using the provided Monoid.
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func([]A) B {
	return G.FoldMap[[]A](m)
}

// FoldMapWithIndex maps and folds an array. Map the Array passing each value to the iterating function. Then fold the results using the provided Monoid.
func FoldMapWithIndex[A, B any](m M.Monoid[B]) func(func(int, A) B) func([]A) B {
	return G.FoldMapWithIndex[[]A](m)
}

// Fold folds the array using the provided Monoid.
func Fold[A any](m M.Monoid[A]) func([]A) A {
	return G.Fold[[]A](m)
}

// Push adds an element to the end of an array (alias for Append).
func Push[A any](a A) EM.Endomorphism[[]A] {
	return G.Push[EM.Endomorphism[[]A]](a)
}

// MonadFlap applies a value to an array of functions, producing an array of results.
// This is the monadic version that takes both parameters.
func MonadFlap[B, A any](fab []func(A) B, a A) []B {
	return G.MonadFlap[func(A) B, []func(A) B, []B, A, B](fab, a)
}

// Flap applies a value to an array of functions, producing an array of results.
// This is the curried version.
func Flap[B, A any](a A) func([]func(A) B) []B {
	return G.Flap[func(A) B, []func(A) B, []B, A, B](a)
}

// Prepend adds an element to the beginning of an array, returning a new array.
func Prepend[A any](head A) EM.Endomorphism[[]A] {
	return G.Prepend[EM.Endomorphism[[]A]](head)
}
