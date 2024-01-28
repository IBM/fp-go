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

package readert

import (
	F "github.com/IBM/fp-go/function"
	R "github.com/IBM/fp-go/reader/generic"
)

// here we implement the monadic operations using callbacks from
// higher kinded types, as good a golang allows use to do this

func MonadMap[GEA ~func(E) HKTA, GEB ~func(E) HKTB, E, A, B, HKTA, HKTB any](
	fmap func(HKTA, func(A) B) HKTB,
	fa GEA,
	f func(A) B,
) GEB {
	return R.MonadMap[GEA, GEB](fa, F.Bind2nd(fmap, f))
}

func Map[GEA ~func(E) HKTA, GEB ~func(E) HKTB, E, A, B, HKTA, HKTB any](
	fmap func(func(A) B) func(HKTA) HKTB,
	f func(A) B,
) func(GEA) GEB {
	return F.Pipe2(
		f,
		fmap,
		R.Map[GEA, GEB, E, HKTA, HKTB],
	)
}

func MonadChain[GEA ~func(E) HKTA, GEB ~func(E) HKTB, A, E, HKTA, HKTB any](fchain func(HKTA, func(A) HKTB) HKTB, ma GEA, f func(A) GEB) GEB {
	return R.MakeReader(func(r E) HKTB {
		return fchain(ma(r), func(a A) HKTB {
			return f(a)(r)
		})
	})
}

func MonadOf[GEA ~func(E) HKTA, E, A, HKTA any](fof func(A) HKTA, a A) GEA {
	return R.MakeReader(func(_ E) HKTA {
		return fof(a)
	})
}

// HKTFAB = HKT[func(A)B]
func MonadAp[GEA ~func(E) HKTA, GEB ~func(E) HKTB, GEFAB ~func(E) HKTFAB, E, A, HKTA, HKTB, HKTFAB any](fap func(HKTFAB, HKTA) HKTB, fab GEFAB, fa GEA) GEB {
	return R.MakeReader(func(r E) HKTB {
		return fap(fab(r), fa(r))
	})
}

func MonadFromReader[GA ~func(E) A, GEA ~func(E) HKTA, E, A, HKTA any](
	fof func(A) HKTA, ma GA) GEA {
	return R.MakeReader(F.Flow2(ma, fof))
}

func FromReader[GA ~func(E) A, GEA ~func(E) HKTA, E, A, HKTA any](
	fof func(A) HKTA) func(ma GA) GEA {
	return F.Bind1st(MonadFromReader[GA, GEA, E, A, HKTA], fof)
}
