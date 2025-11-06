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
	"github.com/IBM/fp-go/v2/tuple"
)

// Of constructs a single element array
func Of[GA ~[]A, A any](value A) GA {
	return GA{value}
}

func Reduce[GA ~[]A, A, B any](f func(B, A) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduce[GA](as, f, initial)
	}
}

func ReduceWithIndex[GA ~[]A, A, B any](f func(int, B, A) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduceWithIndex[GA](as, f, initial)
	}
}

func ReduceRight[GA ~[]A, A, B any](f func(A, B) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduceRight[GA](as, f, initial)
	}
}

func ReduceRightWithIndex[GA ~[]A, A, B any](f func(int, A, B) B, initial B) func(GA) B {
	return func(as GA) B {
		return MonadReduceRightWithIndex[GA](as, f, initial)
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
	for i := n - 1; i >= 0; i-- {
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

func UpsertAt[GA ~[]A, A any](a A) func(GA) GA {
	return array.UpsertAt[GA](a)
}

func MonadMap[GA ~[]A, GB ~[]B, A, B any](as GA, f func(a A) B) GB {
	return array.MonadMap[GA, GB](as, f)
}

func Map[GA ~[]A, GB ~[]B, A, B any](f func(a A) B) func(GA) GB {
	return array.Map[GA, GB](f)
}

func MonadMapWithIndex[GA ~[]A, GB ~[]B, A, B any](as GA, f func(int, A) B) GB {
	return array.MonadMapWithIndex[GA, GB](as, f)
}

func MapWithIndex[GA ~[]A, GB ~[]B, A, B any](f func(int, A) B) func(GA) GB {
	return F.Bind2nd(MonadMapWithIndex[GA, GB, A, B], f)
}

func Size[GA ~[]A, A any](as GA) int {
	return len(as)
}

func filterMap[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(A) O.Option[B]) GB {
	result := make(GB, 0, len(fa))
	for _, a := range fa {
		O.Map(func(b B) B {
			result = append(result, b)
			return b
		})(f(a))
	}
	return result
}

func filterMapWithIndex[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(int, A) O.Option[B]) GB {
	result := make(GB, 0, len(fa))
	for i, a := range fa {
		O.Map(func(b B) B {
			result = append(result, b)
			return b
		})(f(i, a))
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

func FilterChain[GA ~[]A, GB ~[]B, A, B any](f func(a A) O.Option[GB]) func(GA) GB {
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

func MonadPartition[GA ~[]A, A any](as GA, pred func(A) bool) tuple.Tuple2[GA, GA] {
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
	return tuple.MakeTuple2(left, right)
}

func Partition[GA ~[]A, A any](pred func(A) bool) func(GA) tuple.Tuple2[GA, GA] {
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

func Slice[AS ~[]A, A any](start int, end int) func(AS) AS {
	return func(a AS) AS {
		return a[start:end]
	}
}

func SliceRight[AS ~[]A, A any](start int) func(AS) AS {
	return func(a AS) AS {
		return a[start:]
	}
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

func Prepend[ENDO ~func(AS) AS, AS []A, A any](head A) ENDO {
	return array.Prepend[ENDO](head)
}
