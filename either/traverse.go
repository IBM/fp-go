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

package either

import (
	F "github.com/IBM/fp-go/function"
)

/*
*
We need to pass the members of the applicative explicitly, because golang does neither support higher kinded types nor template methods on structs or interfaces

HKTRB = HKT<Either[B]>
HKTA = HKT<A>
HKTB = HKT<B>
*/
func traverse[E, A, B, HKTB, HKTRB any](
	mof func(Either[E, B]) HKTRB,
	mmap func(func(B) Either[E, B]) func(HKTB) HKTRB,
) func(Either[E, A], func(A) HKTB) HKTRB {

	left := F.Flow2(Left[B, E], mof)
	right := mmap(Right[E, B])

	return func(ta Either[E, A], f func(A) HKTB) HKTRB {
		return MonadFold(ta,
			left,
			F.Flow2(f, right),
		)
	}
}

// Traverse converts an [Either] of some higher kinded type into the higher kinded type of an [Either]
func Traverse[A, E, B, HKTB, HKTRB any](
	mof func(Either[E, B]) HKTRB,
	mmap func(func(B) Either[E, B]) func(HKTB) HKTRB,
) func(func(A) HKTB) func(Either[E, A]) HKTRB {
	delegate := traverse[E, A, B](mof, mmap)
	return func(f func(A) HKTB) func(Either[E, A]) HKTRB {
		return F.Bind2nd(delegate, f)
	}
}

// Sequence converts an [Either] of some higher kinded type into the higher kinded type of an [Either]
func Sequence[E, A, HKTA, HKTRA any](
	mof func(Either[E, A]) HKTRA,
	mmap func(func(A) Either[E, A]) func(HKTA) HKTRA,
) func(Either[E, HKTA]) HKTRA {
	return Fold(F.Flow2(Left[A, E], mof), mmap(Right[E, A]))
}
