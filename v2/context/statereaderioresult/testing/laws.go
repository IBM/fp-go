// Copyright (c) 2024 - 2025 IBM Corp.
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

package testing

import (
	"context"
	"testing"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	ST "github.com/IBM/fp-go/v2/context/statereaderioresult"
	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
	P "github.com/IBM/fp-go/v2/pair"
	RES "github.com/IBM/fp-go/v2/result"
)

// AssertLaws asserts the monad laws for the StateReaderIOResult monad
func AssertLaws[S, A, B, C any](t *testing.T,
	eqs EQ.Eq[S],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,

	s S,
	ctx context.Context,
) func(a A) bool {

	eqra := RIORES.Eq(RES.Eq(P.Eq(eqs, eqa)))(ctx)
	eqrb := RIORES.Eq(RES.Eq(P.Eq(eqs, eqb)))(ctx)
	eqrc := RIORES.Eq(RES.Eq(P.Eq(eqs, eqc)))(ctx)

	fofc := ST.Pointed[S, C]()
	fofaa := ST.Pointed[S, func(A) A]()
	fofbc := ST.Pointed[S, func(B) C]()
	fofabb := ST.Pointed[S, func(func(A) B) B]()

	fmap := ST.Functor[S, func(B) C, func(func(A) B) func(A) C]()

	fapabb := ST.Applicative[S, func(A) B, B]()
	fapabac := ST.Applicative[S, func(A) B, func(A) C]()

	maa := ST.Monad[S, A, A]()
	mab := ST.Monad[S, A, B]()
	mac := ST.Monad[S, A, C]()
	mbc := ST.Monad[S, B, C]()

	return L.MonadAssertLaws(t,
		ST.Eq(eqra)(s),
		ST.Eq(eqrb)(s),
		ST.Eq(eqrc)(s),

		fofc,
		fofaa,
		fofbc,
		fofabb,

		fmap,

		fapabb,
		fapabac,

		maa,
		mab,
		mac,
		mbc,

		ab,
		bc,
	)

}
