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

package record

import (
	F "github.com/IBM/fp-go/v2/function"
)

// createEmpty creates a new empty, read-write map
// this is different to Empty which creates a new read-only empty map
func createEmpty[N ~map[K]A, K comparable, A any]() N {
	return make(N)
}

// inserts the key/value pair into a read-write map for performance
// order of parameters is adjusted to be curryable
func addKey[N ~map[K]A, K comparable, A any](key K, m N, value A) N {
	m[key] = value
	return m
}

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<MB>
HKTA = HKT<A>
HKTB = HKT<B>
HKTAB = HKT<func(A)B>
*/
func traverseWithIndex[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(MB) HKTRB,
	fmap func(func(MB) func(B) MB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta MA, f func(K, A) HKTB) HKTRB {
	// this function inserts a value into a map with a given key
	mmap := F.Flow2(F.Curry3(addKey[MB, K, B]), fmap)

	return ReduceWithIndex(ta, func(k K, r HKTRB, a A) HKTRB {
		return F.Pipe2(
			r,
			mmap(k),
			fap(f(k, a)),
		)
	}, fof(createEmpty[MB]()))
}

func MonadTraverse[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(MB) HKTRB,
	fmap func(func(MB) func(B) MB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	r MA, f func(A) HKTB) HKTRB {
	return traverseWithIndex(fof, fmap, fap, r, F.Ignore1of2[K](f))
}

func MonadTraverseWithIndex[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(MB) HKTRB,
	fmap func(func(MB) func(B) MB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	r MA, f func(K, A) HKTB) HKTRB {
	return traverseWithIndex(fof, fmap, fap, r, f)
}

func TraverseWithIndex[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(MB) HKTRB,
	fmap func(func(MB) func(B) MB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(K, A) HKTB) func(MA) HKTRB {

	return func(ma MA) HKTRB {
		return traverseWithIndex(fof, fmap, fap, ma, f)
	}
}

// HKTA = HKT<A>
// HKTB = HKT<B>
// HKTAB = HKT<func(A)B>
// HKTRB = HKT<MB>
func Traverse[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(MB) HKTRB,
	fmap func(func(MB) func(B) MB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(A) HKTB) func(MA) HKTRB {

	return func(ma MA) HKTRB {
		return MonadTraverse(fof, fmap, fap, ma, f)
	}
}

// HKTA = HKT[A]
// HKTAA = HKT[func(A)MA]
// HKTRA = HKT[MA]
func Sequence[MA ~map[K]A, MKTA ~map[K]HKTA, K comparable, A, HKTA, HKTAA, HKTRA any](
	fof func(MA) HKTRA,
	fmap func(func(MA) func(A) MA) func(HKTRA) HKTAA,
	fap func(HKTA) func(HKTAA) HKTRA,

	ma MKTA) HKTRA {
	return MonadTraverse(fof, fmap, fap, ma, F.Identity[HKTA])
}
