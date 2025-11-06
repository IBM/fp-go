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
	G "github.com/IBM/fp-go/array/generic"
	EM "github.com/IBM/fp-go/endomorphism"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/array"
	M "github.com/IBM/fp-go/monoid"
	O "github.com/IBM/fp-go/option"
	"github.com/IBM/fp-go/tuple"
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

func MonadMap[A, B any](as []A, f func(a A) B) []B {
	return G.MonadMap[[]A, []B](as, f)
}

func MonadMapRef[A, B any](as []A, f func(a *A) B) []B {
	count := len(as)
	bs := make([]B, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(&as[i])
	}
	return bs
}

func MapWithIndex[A, B any](f func(int, A) B) func([]A) []B {
	return G.MapWithIndex[[]A, []B](f)
}

func Map[A, B any](f func(a A) B) func([]A) []B {
	return G.Map[[]A, []B, A, B](f)
}

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

func FilterRef[A any](pred func(*A) bool) EM.Endomorphism[[]A] {
	return F.Bind2nd(filterRef[A], pred)
}

func MonadFilterMap[A, B any](fa []A, f func(A) O.Option[B]) []B {
	return G.MonadFilterMap[[]A, []B](fa, f)
}

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

func MonadReduce[A, B any](fa []A, f func(B, A) B, initial B) B {
	return G.MonadReduce(fa, f, initial)
}

func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B {
	return G.Reduce[[]A](f, initial)
}

func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B {
	return G.ReduceWithIndex[[]A](f, initial)
}

func ReduceRight[A, B any](f func(A, B) B, initial B) func([]A) B {
	return G.ReduceRight[[]A](f, initial)
}

func ReduceRightWithIndex[A, B any](f func(int, A, B) B, initial B) func([]A) B {
	return G.ReduceRightWithIndex[[]A](f, initial)
}

func ReduceRef[A, B any](f func(B, *A) B, initial B) func([]A) B {
	return func(as []A) B {
		return reduceRef(as, f, initial)
	}
}

func Append[A any](as []A, a A) []A {
	return G.Append(as, a)
}

func IsEmpty[A any](as []A) bool {
	return G.IsEmpty(as)
}

func IsNonEmpty[A any](as []A) bool {
	return len(as) > 0
}

func Empty[A any]() []A {
	return G.Empty[[]A]()
}

func Zero[A any]() []A {
	return Empty[A]()
}

// Of constructs a single element array
func Of[A any](a A) []A {
	return G.Of[[]A](a)
}

func MonadChain[A, B any](fa []A, f func(a A) []B) []B {
	return G.MonadChain[[]A, []B](fa, f)
}

func Chain[A, B any](f func(A) []B) func([]A) []B {
	return G.Chain[[]A, []B](f)
}

func MonadAp[B, A any](fab []func(A) B, fa []A) []B {
	return G.MonadAp[[]B](fab, fa)
}

func Ap[B, A any](fa []A) func([]func(A) B) []B {
	return G.Ap[[]B, []func(A) B](fa)
}

func Match[A, B any](onEmpty func() B, onNonEmpty func([]A) B) func([]A) B {
	return G.Match[[]A](onEmpty, onNonEmpty)
}

func MatchLeft[A, B any](onEmpty func() B, onNonEmpty func(A, []A) B) func([]A) B {
	return G.MatchLeft[[]A](onEmpty, onNonEmpty)
}

func Tail[A any](as []A) O.Option[[]A] {
	return G.Tail(as)
}

func Head[A any](as []A) O.Option[A] {
	return G.Head(as)
}

func First[A any](as []A) O.Option[A] {
	return G.First(as)
}

func Last[A any](as []A) O.Option[A] {
	return G.Last(as)
}

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

func Intersperse[A any](middle A) EM.Endomorphism[[]A] {
	prepend := PrependAll(middle)
	return func(as []A) []A {
		if IsEmpty(as) {
			return as
		}
		return prepend(as)[1:]
	}
}

func Intercalate[A any](m M.Monoid[A]) func(A) func([]A) A {
	concatAll := ConcatAll[A](m)
	return func(middle A) func([]A) A {
		return Match(m.Empty, F.Flow2(Intersperse(middle), concatAll))
	}
}

func Flatten[A any](mma [][]A) []A {
	return G.Flatten(mma)
}

func Slice[A any](low, high int) func(as []A) []A {
	return array.Slice[[]A](low, high)
}

func Lookup[A any](idx int) func([]A) O.Option[A] {
	return G.Lookup[[]A](idx)
}

func UpsertAt[A any](a A) EM.Endomorphism[[]A] {
	return G.UpsertAt[[]A](a)
}

func Size[A any](as []A) int {
	return G.Size(as)
}

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

func Push[A any](a A) EM.Endomorphism[[]A] {
	return G.Push[EM.Endomorphism[[]A]](a)
}

func MonadFlap[B, A any](fab []func(A) B, a A) []B {
	return G.MonadFlap[func(A) B, []func(A) B, []B, A, B](fab, a)
}

func Flap[B, A any](a A) func([]func(A) B) []B {
	return G.Flap[func(A) B, []func(A) B, []B, A, B](a)
}

func Prepend[A any](head A) EM.Endomorphism[[]A] {
	return G.Prepend[EM.Endomorphism[[]A]](head)
}
