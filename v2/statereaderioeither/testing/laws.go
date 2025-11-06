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
	"testing"

	ET "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
	P "github.com/IBM/fp-go/v2/pair"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	ST "github.com/IBM/fp-go/v2/statereaderioeither"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[S, E, R, A, B, C any](t *testing.T,
	eqs EQ.Eq[S],
	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,

	s S,
	r R,
) func(a A) bool {

	eqra := RIOE.Eq[R](ET.Eq(eqe, P.Eq(eqs, eqa)))(r)
	eqrb := RIOE.Eq[R](ET.Eq(eqe, P.Eq(eqs, eqb)))(r)
	eqrc := RIOE.Eq[R](ET.Eq(eqe, P.Eq(eqs, eqc)))(r)

	fofc := ST.Pointed[S, R, E, C]()
	fofaa := ST.Pointed[S, R, E, func(A) A]()
	fofbc := ST.Pointed[S, R, E, func(B) C]()
	fofabb := ST.Pointed[S, R, E, func(func(A) B) B]()

	fmap := ST.Functor[S, R, E, func(B) C, func(func(A) B) func(A) C]()

	fapabb := ST.Applicative[S, R, E, func(A) B, B]()
	fapabac := ST.Applicative[S, R, E, func(A) B, func(A) C]()

	maa := ST.Monad[S, R, E, A, A]()
	mab := ST.Monad[S, R, E, A, B]()
	mac := ST.Monad[S, R, E, A, C]()
	mbc := ST.Monad[S, R, E, B, C]()

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
