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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/array"
	FC "github.com/IBM/fp-go/v2/internal/functor"
	M "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

// Of constructs a single element array
//
//go:inline
func Of[GA ~[]A, A any](value A) GA {
	return array.Of[GA](value)
}

func Reduce[GA ~[]A, A, B any](f func(B, A) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduce(as, f, initial)
	}
}

func ReduceWithIndex[GA ~[]A, A, B any](f func(int, B, A) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduceWithIndex(as, f, initial)
	}
}

func ReduceRight[GA ~[]A, A, B any](f func(A, B) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduceRight(as, f, initial)
	}
}

func ReduceRightWithIndex[GA ~[]A, A, B any](f func(int, A, B) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduceRightWithIndex(as, f, initial)
	}
}

func MonadReduce[GA ~[]A, A, B any](fa GA, f func(B, A) B, initial B) B {
	return array.Reduce(fa, f, initial)
}

func MonadReduceWithIndex[GA ~[]A, A, B any](fa GA, f func(int, B, A) B, initial B) B {
	return array.ReduceWithIndex(fa, f, initial)
}

func MonadReduceRight[GA ~[]A, A, B any](fa GA, f func(A, B) B, initial B) B {
	return array.ReduceRight(fa, f, initial)
}

func MonadReduceRightWithIndex[GA ~[]A, A, B any](fa GA, f func(int, A, B) B, initial B) B {
	return array.ReduceRightWithIndex(fa, f, initial)
}

// From constructs an array from a set of variadic arguments
func From[GA ~[]A, A any](data ...A) GA {
	return data
}

// MakeBy returns a `Array` of length `n` with element `i` initialized with `f(i)`.
func MakeBy[AS ~[]A, F ~func(int) A, A any](n int, f F) AS {
	// sanity check
	if n <= 0 {
		return Empty[AS]()
	}
	// run the generator function across the input
	as := make(AS, n)
	for i := range n {
		as[i] = f(i)
	}
	return as
}

func Replicate[AS ~[]A, A any](n int, a A) AS {
	return MakeBy[AS](n, F.Constant1[int](a))
}

func Lookup[GA ~[]A, A any](idx int) func(GA) O.Option[A] {
	none := O.None[A]()
	if idx < 0 {
		return F.Constant1[GA](none)
	}
	return func(as GA) O.Option[A] {
		if idx < len(as) {
			return O.Some(as[idx])
		}
		return none
	}
}

func Tail[GA ~[]A, A any](as GA) O.Option[GA] {
	if array.IsEmpty(as) {
		return O.None[GA]()
	}
	return O.Some(as[1:])
}

func Head[GA ~[]A, A any](as GA) O.Option[A] {
	if array.IsEmpty(as) {
		return O.None[A]()
	}
	return O.Some(as[0])
}

func First[GA ~[]A, A any](as GA) O.Option[A] {
	return Head(as)
}

func Last[GA ~[]A, A any](as GA) O.Option[A] {
	if array.IsEmpty(as) {
		return O.None[A]()
	}
	return O.Some(as[len(as)-1])
}

func Append[GA ~[]A, A any](as GA, a A) GA {
	return array.Append(as, a)
}

func Empty[GA ~[]A, A any]() GA {
	return array.Empty[GA]()
}

//go:inline
func UpsertAt[GA ~[]A, A any](a A) func(GA) GA {
	return array.UpsertAt[GA](a)
}

//go:inline
func MonadMap[GA ~[]A, GB ~[]B, A, B any](as GA, f func(a A) B) GB {
	return array.MonadMap[GA, GB](as, f)
}

//go:inline
func Map[GA ~[]A, GB ~[]B, A, B any](f func(a A) B) func(GA) GB {
	return array.Map[GA, GB](f)
}

//go:inline
func MonadMapWithIndex[GA ~[]A, GB ~[]B, A, B any](as GA, f func(int, A) B) GB {
	return array.MonadMapWithIndex[GA, GB](as, f)
}

//go:inline
func MapWithIndex[GA ~[]A, GB ~[]B, A, B any](f func(int, A) B) func(GA) GB {
	return F.Bind2nd(MonadMapWithIndex[GA, GB, A, B], f)
}

func Size[GA ~[]A, A any](as GA) int {
	return len(as)
}

func filterMap[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(A) O.Option[B]) GB {
	result := make(GB, 0, len(fa))
	for _, a := range fa {
		if b, ok := O.Unwrap(f(a)); ok {
			result = append(result, b)
		}
	}
	return result
}

func filterMapWithIndex[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(int, A) O.Option[B]) GB {
	result := make(GB, 0, len(fa))
	for i, a := range fa {
		if b, ok := O.Unwrap(f(i, a)); ok {
			result = append(result, b)
		}
	}
	return result
}

func MonadFilterMap[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(A) O.Option[B]) GB {
	return filterMap[GA, GB](fa, f)
}

func MonadFilterMapWithIndex[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(int, A) O.Option[B]) GB {
	return filterMapWithIndex[GA, GB](fa, f)
}

func filterWithIndex[AS ~[]A, PRED ~func(int, A) bool, A any](fa AS, pred PRED) AS {
	result := make(AS, 0, len(fa))
	for i, a := range fa {
		if pred(i, a) {
			result = append(result, a)
		}
	}
	return result
}

func FilterWithIndex[AS ~[]A, PRED ~func(int, A) bool, A any](pred PRED) func(AS) AS {
	return F.Bind2nd(filterWithIndex[AS, PRED, A], pred)
}

func Filter[AS ~[]A, PRED ~func(A) bool, A any](pred PRED) func(AS) AS {
	return FilterWithIndex[AS](F.Ignore1of2[int](pred))
}

func ChainOptionK[GA ~[]A, GB ~[]B, A, B any](f func(a A) O.Option[GB]) func(GA) GB {
	return F.Flow2(
		FilterMap[GA, []GB](f),
		Flatten[[]GB],
	)
}

func Flatten[GAA ~[]GA, GA ~[]A, A any](mma GAA) GA {
	return MonadChain(mma, F.Identity[GA])
}

func FilterMap[GA ~[]A, GB ~[]B, A, B any](f func(A) O.Option[B]) func(GA) GB {
	return F.Bind2nd(MonadFilterMap[GA, GB, A, B], f)
}

func FilterMapWithIndex[GA ~[]A, GB ~[]B, A, B any](f func(int, A) O.Option[B]) func(GA) GB {
	return F.Bind2nd(MonadFilterMapWithIndex[GA, GB, A, B], f)
}

func MonadPartition[GA ~[]A, A any](as GA, pred func(A) bool) pair.Pair[GA, GA] {
	left := Empty[GA]()
	right := Empty[GA]()
	array.Reduce(as, func(c bool, a A) bool {
		if pred(a) {
			right = append(right, a)
		} else {
			left = append(left, a)
		}
		return c
	}, true)
	// returns the partition
	return pair.MakePair(left, right)
}

func Partition[GA ~[]A, A any](pred func(A) bool) func(GA) pair.Pair[GA, GA] {
	return F.Bind2nd(MonadPartition[GA, A], pred)
}

func MonadChain[AS ~[]A, BS ~[]B, A, B any](fa AS, f func(a A) BS) BS {
	return array.Reduce(fa, func(bs BS, a A) BS {
		return append(bs, f(a)...)
	}, Empty[BS]())
}

func Chain[AS ~[]A, BS ~[]B, A, B any](f func(A) BS) func(AS) BS {
	return F.Bind2nd(MonadChain[AS, BS, A, B], f)
}

func MonadAp[BS ~[]B, ABS ~[]func(A) B, AS ~[]A, B, A any](fab ABS, fa AS) BS {
	return MonadChain(fab, F.Bind1st(MonadMap[AS, BS, A, B], fa))
}

func Ap[BS ~[]B, ABS ~[]func(A) B, AS ~[]A, B, A any](fa AS) func(ABS) BS {
	return F.Bind2nd(MonadAp[BS, ABS, AS], fa)
}

func IsEmpty[AS ~[]A, A any](as AS) bool {
	return array.IsEmpty(as)
}

func IsNil[GA ~[]A, A any](as GA) bool {
	return array.IsNil(as)
}

func IsNonNil[GA ~[]A, A any](as GA) bool {
	return array.IsNonNil(as)
}

func Match[AS ~[]A, A, B any](onEmpty func() B, onNonEmpty func(AS) B) func(AS) B {
	return func(as AS) B {
		if IsEmpty(as) {
			return onEmpty()
		}
		return onNonEmpty(as)
	}
}

func MatchLeft[AS ~[]A, A, B any](onEmpty func() B, onNonEmpty func(A, AS) B) func(AS) B {
	return func(as AS) B {
		if IsEmpty(as) {
			return onEmpty()
		}
		return onNonEmpty(as[0], as[1:])
	}
}

//go:inline
func Slice[AS ~[]A, A any](start, end int) func(AS) AS {
	return array.Slice[AS](start, end)
}

//go:inline
func SliceRight[AS ~[]A, A any](start int) func(AS) AS {
	return array.SliceRight[AS](start)
}

func Copy[AS ~[]A, A any](b AS) AS {
	buf := make(AS, len(b))
	copy(buf, b)
	return buf
}

func Clone[AS ~[]A, A any](f func(A) A) func(as AS) AS {
	// implementation assumes that map does not optimize for the empty array
	return Map[AS, AS](f)
}

func FoldMap[AS ~[]A, A, B any](m M.Monoid[B]) func(func(A) B) func(AS) B {
	empty := m.Empty()
	concat := m.Concat
	return func(f func(A) B) func(AS) B {
		return func(as AS) B {
			return array.Reduce(as, func(cur B, a A) B {
				return concat(cur, f(a))
			}, empty)
		}
	}
}

func FoldMapWithIndex[AS ~[]A, A, B any](m M.Monoid[B]) func(func(int, A) B) func(AS) B {
	empty := m.Empty()
	concat := m.Concat
	return func(f func(int, A) B) func(AS) B {
		return func(as AS) B {
			return array.ReduceWithIndex(as, func(idx int, cur B, a A) B {
				return concat(cur, f(idx, a))
			}, empty)
		}
	}
}

func Fold[AS ~[]A, A any](m M.Monoid[A]) func(AS) A {
	empty := m.Empty()
	concat := m.Concat
	return func(as AS) A {
		return array.Reduce(as, concat, empty)
	}
}

func Push[ENDO ~func(GA) GA, GA ~[]A, A any](a A) ENDO {
	return F.Bind2nd(array.Push[GA, A], a)
}

func MonadFlap[FAB ~func(A) B, GFAB ~[]FAB, GB ~[]B, A, B any](fab GFAB, a A) GB {
	return FC.MonadFlap(MonadMap[GFAB, GB], fab, a)
}

func Flap[FAB ~func(A) B, GFAB ~[]FAB, GB ~[]B, A, B any](a A) func(GFAB) GB {
	return FC.Flap(Map[GFAB, GB], a)
}

//go:inline
func Prepend[ENDO ~func(AS) AS, AS []A, A any](head A) ENDO {
	return array.Prepend[ENDO](head)
}

//go:inline
func Reverse[GT ~[]T, T any](as GT) GT {
	return array.Reverse(as)
}

// Extract returns the first element of an array, or a zero value if empty.
// This is the comonad extract operation for arrays.
//
// Extract is the dual of the monadic return/of operation. While Of wraps a value
// in a context, Extract unwraps a value from its context.
//
// Type Parameters:
//   - GA: The array type constraint
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
//	result := Extract([]int{1, 2, 3})
//	// result: 1
//
// Example with empty array:
//
//	result := Extract([]int{})
//	// result: 0 (zero value for int)
//
// Comonad laws:
//   - Extract ∘ Of == Identity (extracting from a singleton returns the value)
//   - Extract ∘ Extend(f) == f (extract after extend equals applying f)
//
//go:inline
func Extract[GA ~[]A, A any](as GA) A {
	if len(as) > 0 {
		return as[0]
	}
	var zero A
	return zero
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
//   - GA: The input array type constraint
//   - GB: The output array type constraint
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
//	sumSuffix := Extend[[]int, []int](func(as []int) int {
//	    return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
//	})
//	result := sumSuffix([]int{1, 2, 3, 4})
//	// result: []int{10, 9, 7, 4}
//	// Explanation: [1+2+3+4, 2+3+4, 3+4, 4]
//
// Example with length:
//
//	// Get remaining length at each position
//	lengths := Extend[[]int, []int](Size[[]int, int])
//	result := lengths([]int{10, 20, 30})
//	// result: []int{3, 2, 1}
//
// Comonad laws:
//   - Left identity: Extend(Extract) == Identity
//   - Right identity: Extract ∘ Extend(f) == f
//   - Associativity: Extend(f) ∘ Extend(g) == Extend(f ∘ Extend(g))
//
//go:inline
func Extend[GA ~[]A, GB ~[]B, A, B any](f func(GA) B) func(GA) GB {
	return func(as GA) GB {
		return MakeBy[GB](len(as), func(i int) B { return f(as[i:]) })
	}
}
