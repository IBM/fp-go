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

package optiont

import (
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/apply"
	FC "github.com/IBM/fp-go/internal/functor"
	O "github.com/IBM/fp-go/option"
)

func Of[A, HKTA any](fof func(O.Option[A]) HKTA, a A) HKTA {
	return F.Pipe2(a, O.Of[A], fof)
}

func None[A, HKTA any](fof func(O.Option[A]) HKTA) HKTA {
	return F.Pipe1(O.None[A](), fof)
}

func OfF[A, HKTA, HKTEA any](fmap func(HKTA, func(A) O.Option[A]) HKTEA, fa HKTA) HKTEA {
	return fmap(fa, O.Of[A])
}

func MonadMap[A, B, HKTFA, HKTFB any](fmap func(HKTFA, func(O.Option[A]) O.Option[B]) HKTFB, fa HKTFA, f func(A) B) HKTFB {
	// HKTGA = Either[E, A]
	// HKTGB = Either[E, B]
	return FC.MonadMap(fmap, O.MonadMap[A, B], fa, f)
}

func MonadChain[A, B, HKTFA, HKTFB any](
	fchain func(HKTFA, func(O.Option[A]) HKTFB) HKTFB,
	fof func(O.Option[B]) HKTFB,
	ma HKTFA,
	f func(A) HKTFB) HKTFB {
	// dispatch to the even more generic implementation
	return fchain(ma, O.Fold(F.Nullary2(O.None[B], fof), f))
}

func MonadAp[A, B, HKTFAB, HKTFGAB, HKTFA, HKTFB any](
	fap func(HKTFGAB, HKTFA) HKTFB,
	fmap func(HKTFAB, func(O.Option[func(A) B]) func(O.Option[A]) O.Option[B]) HKTFGAB,
	fab HKTFAB,
	fa HKTFA) HKTFB {
	// HKTGA  = O.Option[A]
	// HKTGB  = O.Option[B]
	// HKTGAB = O.Option[func(a A) B]
	return apply.MonadAp(fap, fmap, O.MonadAp[B, A], fab, fa)
}

func MatchE[A, HKTEA, HKTB any](mchain func(HKTEA, func(O.Option[A]) HKTB) HKTB, onNone func() HKTB, onSome func(A) HKTB) func(HKTEA) HKTB {
	return F.Bind2nd(mchain, O.Fold(onNone, onSome))
}

func FromOptionK[A, B, HKTB any](
	fof func(O.Option[B]) HKTB,
	f func(A) O.Option[B]) func(A) HKTB {
	return F.Flow2(f, fof)
}

func MonadChainOptionK[A, B, HKTA, HKTB any](
	fchain func(HKTA, func(O.Option[A]) HKTB) HKTB,
	fof func(O.Option[B]) HKTB,
	ma HKTA,
	f func(A) O.Option[B],
) HKTB {
	return MonadChain(fchain, fof, ma, FromOptionK(fof, f))
}
