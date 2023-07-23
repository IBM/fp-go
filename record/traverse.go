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

package record

import (
	G "github.com/IBM/fp-go/internal/record"
)

func TraverseWithIndex[K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(map[K]B) HKTRB,
	fmap func(func(map[K]B) func(B) map[K]B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(K, A) HKTB) func(map[K]A) HKTRB {
	return G.TraverseWithIndex[map[K]A](fof, fmap, fap, f)
}

// HKTA = HKT<A>
// HKTB = HKT<B>
// HKTAB = HKT<func(A)B>
// HKTRB = HKT<map[K]B>
func Traverse[K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(map[K]B) HKTRB,
	fmap func(func(map[K]B) func(B) map[K]B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,
	f func(A) HKTB) func(map[K]A) HKTRB {
	return G.Traverse[map[K]A](fof, fmap, fap, f)
}

// HKTA = HKT[A]
// HKTAA = HKT[func(A)map[K]A]
// HKTRA = HKT[map[K]A]
func Sequence[K comparable, A, HKTA, HKTAA, HKTRA any](
	fof func(map[K]A) HKTRA,
	fmap func(func(map[K]A) func(A) map[K]A) func(HKTRA) HKTAA,
	fap func(HKTA) func(HKTAA) HKTRA,
	ma map[K]HKTA) HKTRA {
	return G.Sequence(fof, fmap, fap, ma)

}
