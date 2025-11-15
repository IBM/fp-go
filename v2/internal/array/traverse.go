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
	F "github.com/IBM/fp-go/v2/function"
)

func MonadSequenceSegment[HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	empty HKTRB,
	concat func(HKTRB, HKTRB) HKTRB,
	fbs []HKTB,
	start, end int,
) HKTRB {

	switch end - start {
	case 0:
		return empty
	case 1:
		return fof(fbs[start])
	default:
		mid := (start + end) / 2
		return concat(
			MonadSequenceSegment(fof, empty, concat, fbs, start, mid),
			MonadSequenceSegment(fof, empty, concat, fbs, mid, end),
		)
	}
}

func SequenceSegment[HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	empty HKTRB,
	concat func(HKTRB, HKTRB) HKTRB,
) func([]HKTB) HKTRB {

	concat_f := func(left, right func([]HKTB) HKTRB) func([]HKTB) HKTRB {
		return func(fbs []HKTB) HKTRB {
			return concat(left(fbs), right(fbs))
		}
	}
	empty_f := F.Constant1[[]HKTB](empty)
	at := func(idx int) func([]HKTB) HKTRB {
		return func(fbs []HKTB) HKTRB {
			return fof(fbs[idx])
		}
	}

	var divide func(start, end int) func([]HKTB) HKTRB
	divide = func(start, end int) func([]HKTB) HKTRB {
		switch end - start {
		case 0:
			return empty_f
		case 1:
			return at(start)
		default:
			mid := (start + end) / 2
			left := divide(start, mid)
			right := divide(mid, end)

			return concat_f(left, right)
		}
	}

	// TODO this could be cached by length
	get_divide := func(len int) func([]HKTB) HKTRB {
		return divide(0, len)
	}

	return func(fbs []HKTB) HKTRB {
		return get_divide(len(fbs))(fbs)
	}
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadTraverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,
	f func(A) HKTB) HKTRB {
	return MonadTraverseReduce(fof, fmap, fap, ta, f, Append[GB, B], Empty[GB]())
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadTraverseWithIndex[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,
	f func(int, A) HKTB) HKTRB {
	return MonadTraverseReduceWithIndex(fof, fmap, fap, ta, f, Append[GB, B], Empty[GB]())
}

func Traverse[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverse(fof, fmap, fap, ma, f)
	}
}

func TraverseWithIndex[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(int, A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverseWithIndex(fof, fmap, fap, ma, f)
	}
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadSequence[GA ~[]HKTA, HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	empty HKTRA,
	concat func(HKTRA, HKTRA) HKTRA,

	ta GA) HKTRA {
	return MonadSequenceSegment(fof, empty, concat, ta, 0, len(ta))
}

func Sequence[GA ~[]HKTA, HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	empty HKTRA,
	concat func(HKTRA, HKTRA) HKTRA,
) func(GA) HKTRA {

	return func(ma GA) HKTRA {
		return MonadSequence(fof, empty, concat, ma)
	}
}

func MonadTraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return Reduce(ta, func(r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(a)),
		)
	}, fof(initial))
}

func MonadTraverseReduceWithIndex[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(int, A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return ReduceWithIndex(ta, func(idx int, r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(idx, a)),
		)
	}, fof(initial))
}

func TraverseReduce[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) func(GA) HKTRB {
	return func(ta GA) HKTRB {
		return MonadTraverseReduce(fof, fmap, fap, ta, transform, reduce, initial)
	}
}

func TraverseReduceWithIndex[GA ~[]A, GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	transform func(int, A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) func(GA) HKTRB {
	return func(ta GA) HKTRB {
		return MonadTraverseReduceWithIndex(fof, fmap, fap, ta, transform, reduce, initial)
	}
}
