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

import (
	F "github.com/IBM/fp-go/v2/function"
	INTA "github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
)

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadTraverse[GA ~func(yield func(A) bool), GB ~func(yield func(B) bool), A, B, HKT_B, HKT_GB_GB, HKT_GB any](
	fmap_b func(HKT_B, func(B) GB) HKT_GB,

	fof_gb func(GB) HKT_GB,
	fmap_gb func(HKT_GB, func(GB) func(GB) GB) HKT_GB_GB,
	fap_gb func(HKT_GB_GB, HKT_GB) HKT_GB,

	ta GA,
	f func(A) HKT_B) HKT_GB {

	fof := F.Bind2nd(fmap_b, Of[GB])

	empty := fof_gb(Empty[GB]())

	cb := F.Curry2(Concat[GB])
	concat_gb := F.Bind2nd(fmap_gb, cb)
	concat := func(first HKT_GB, second HKT_GB) HKT_GB {
		return fap_gb(concat_gb(first), second)
	}

	// convert to an array
	hktb := MonadMapToArray[GA, []HKT_B](ta, f)
	return INTA.MonadSequenceSegment(fof, empty, concat, hktb, 0, len(hktb))
}

func Traverse[GA ~func(yield func(A) bool), GB ~func(yield func(B) bool), A, B, HKT_B, HKT_GB_GB, HKT_GB any](
	fmap_b func(func(B) GB) func(HKT_B) HKT_GB,

	fof_gb func(GB) HKT_GB,
	fmap_gb func(func(GB) func(GB) GB) func(HKT_GB) HKT_GB_GB,
	fap_gb func(HKT_GB_GB, HKT_GB) HKT_GB,

	f func(A) HKT_B) func(GA) HKT_GB {

	fof := fmap_b(Of[GB])
	empty := fof_gb(Empty[GB]())
	cb := F.Curry2(Concat[GB])
	concat_gb := fmap_gb(cb)

	concat := func(first, second HKT_GB) HKT_GB {
		return fap_gb(concat_gb(first), second)
	}

	return func(ma GA) HKT_GB {
		// return INTA.SequenceSegment(fof, empty, concat)(MapToArray[GA, []HKT_B](f)(ma))
		hktb := MonadMapToArray[GA, []HKT_B](ma, f)
		return INTA.MonadSequenceSegment(fof, empty, concat, hktb, 0, len(hktb))
	}
}

func MonadSequence[GA ~func(yield func(HKTA) bool), HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	m M.Monoid[HKTRA],

	ta GA) HKTRA {

	// convert to an array
	hktb := ToArray[GA, []HKTA](ta)
	return INTA.MonadSequenceSegment(fof, m.Empty(), m.Concat, hktb, 0, len(hktb))
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<GB>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func MonadTraverseWithIndex[GA ~func(yield func(A) bool), A, HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	m M.Monoid[HKTRB],

	ta GA,
	f func(int, A) HKTB) HKTRB {

	// convert to an array
	hktb := MonadMapToArrayWithIndex[GA, []HKTB](ta, f)
	return INTA.MonadSequenceSegment(fof, m.Empty(), m.Concat, hktb, 0, len(hktb))
}

func Sequence[GA ~func(yield func(HKTA) bool), HKTA, HKTRA any](
	fof func(HKTA) HKTRA,
	m M.Monoid[HKTRA]) func(GA) HKTRA {

	return func(ma GA) HKTRA {
		return MonadSequence(fof, m, ma)
	}
}

func TraverseWithIndex[GA ~func(yield func(A) bool), A, HKTB, HKTRB any](
	fof func(HKTB) HKTRB,
	m M.Monoid[HKTRB],

	f func(int, A) HKTB) func(GA) HKTRB {

	return func(ma GA) HKTRB {
		return MonadTraverseWithIndex(fof, m, ma, f)
	}
}

func MonadTraverseReduce[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return MonadReduce(ta, func(r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(a)),
		)
	}, fof(initial))
}

func MonadTraverseReduceWithIndex[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta GA,

	transform func(int, A) HKTB,
	reduce func(GB, B) GB,
	initial GB,
) HKTRB {
	mmap := fmap(F.Curry2(reduce))

	return MonadReduceWithIndex(ta, func(idx int, r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap,
			fap(transform(idx, a)),
		)
	}, fof(initial))
}

func TraverseReduce[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
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

func TraverseReduceWithIndex[GA ~func(yield func(A) bool), GB, A, B, HKTB, HKTAB, HKTRB any](
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
